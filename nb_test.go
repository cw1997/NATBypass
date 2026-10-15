package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestCLIHelperProcess(t *testing.T) {
	if os.Getenv("NATBYPASS_TEST_CLI") != "1" {
		return
	}
	for i, arg := range os.Args {
		if arg == "--" {
			os.Args = append([]string{"nb"}, os.Args[i+1:]...)
			main()
			os.Exit(0)
		}
	}
}

func TestTranWithLog(t *testing.T) {
	dir, err := ioutil.TempDir("", "natbypass log-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	target, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer target.Close()
	request := bytes.Repeat([]byte("request\x00"), 32768)
	response := bytes.Repeat([]byte("response\x00"), 32768)
	done := make(chan error, 1)
	go func() {
		conn, err := target.Accept()
		if err != nil {
			done <- err
			return
		}
		defer conn.Close()
		conn.SetDeadline(time.Now().Add(10 * time.Second))
		got := make([]byte, len(request))
		if _, err = io.ReadFull(conn, got); err == nil && !bytes.Equal(got, request) {
			err = io.ErrUnexpectedEOF
		}
		if err == nil {
			_, err = conn.Write(response)
		}
		done <- err
	}()
	portListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := portListener.Addr().(*net.TCPAddr).Port
	portListener.Close()
	address := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	cmd := exec.Command(os.Args[0], "-test.run=^TestCLIHelperProcess$", "--", "-tran", strconv.Itoa(port), target.Addr().String(), "-log", dir)
	cmd.Env = append(os.Environ(), "NATBYPASS_TEST_CLI=1")
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
		t.Log(output.String())
	}()
	var client net.Conn
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		client, err = net.DialTimeout("tcp", address, 100*time.Millisecond)
		if err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	client.SetDeadline(time.Now().Add(10 * time.Second))
	if _, err := client.Write(request); err != nil {
		t.Fatal(err)
	}
	got := make([]byte, len(response))
	if _, err := io.ReadFull(client, got); err != nil {
		t.Fatalf("read forwarded response: %v", err)
	}
	if !bytes.Equal(got, response) {
		t.Fatal("forwarded response differs")
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	deadline = time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		files, err := filepath.Glob(filepath.Join(dir, "*.log"))
		if err != nil {
			t.Fatal(err)
		}
		var haveRequest, haveResponse bool
		for _, file := range files {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			haveRequest = haveRequest || bytes.Equal(data, request)
			haveResponse = haveResponse || bytes.Equal(data, response)
		}
		if haveRequest && haveResponse {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("log files do not contain both complete byte streams")
}
