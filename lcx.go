package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println("Welcome to use the port transmit tool.")
	fmt.Println("Code by cw1997 at 2017-10-19 03:59:51")
	fmt.Println("If you have some problem when you use the tool,")
	fmt.Println("please submit a new issue at : https://github.com/cw1997/lcx .")
	fmt.Println()
	// sleep one second because the fmt is not thread-safety.
	// if not to do this, fmt.Print will print after the log.Print.
	time.Sleep(time.Second)
	args := os.Args
	switch args[1] {
	case "-listen":
		port1 := checkPort(args[2])
		port2 := checkPort(args[3])
		log.Println("start to listen port:", port1, "and port:", port2)
		port2port(port1, port2)
		break
	case "-tran":
		port := checkPort(args[2])
		var remoteAddress string
		if checkIp(args[3]) {
			remoteAddress = args[3]
		}
		split := strings.SplitN(remoteAddress, ":", 2)
		log.Println("start to transmit address:", remoteAddress, "to address:", split[0]+":"+port)
		port2host(port, remoteAddress)
		break
	case "-slave":
		var address1, address2 string
		checkIp(args[2])
		if checkIp(args[2]) {
			address1 = args[3]
		}
		checkIp(args[3])
		if checkIp(args[3]) {
			address2 = args[3]
		}
		log.Println("start to connect address:", address1, "and address:", address2)
		host2host(address1, address2)
		break
	}
}

func checkPort(port string) string {
	PortNum, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalln("port should be a number")
	}
	if PortNum < 1 && PortNum > 65535 {
		log.Fatalln("port should be a number and the range is [1,65536)")
	}
	return port
}

func checkIp(address string) bool {
	pattern := `(\d|[1-9]\d|1\d{2}|2[0-5][0-5])\.(\d|[1-9]\d|1\d{2}|2[0-5][0-5])\.(\d|[1-9]\d|1\d{2}|2[0-5][0-5])\.(\d|[1-9]\d|1\d{2}|2[0-5][0-5]):([0-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-5]{2}[0-3][0-5])`
	ok, err := regexp.MatchString(pattern, address)
	if err != nil || !ok {
		log.Fatalln("ip address error. should be a string like [ip:port]. ")
	}
	return ok
}

func port2port(port1 string, port2 string) {
	listen1 := start_server("127.0.0.1:" + port1)
	listen2 := start_server("127.0.0.1:" + port2)
	log.Println("listen port:", port1, "and", port2, "success. waiting for client...")
	for {
		conn1 := accept(listen1)
		conn2 := accept(listen2)
		if conn1 == nil || conn2 == nil {
			continue
		}
		forward(conn1, conn2)
	}
}

func port2host(allowPort string, targetAddress string) {
	server := start_server("0.0.0.0:" + allowPort)
	for {
		conn := accept(server)
		if conn == nil {
			continue
		}
		//println(targetAddress)
		go func(targetAddress string) {
			target, err := net.Dial("tcp", targetAddress)
			if err != nil {
				// temporarily unavailable, don't use fatal.
				log.Println("connect target address [" + targetAddress + "] faild.")
				return
			}
			log.Println("connect target address [" + targetAddress + "] success.")
			forward(target, conn)
		}(targetAddress)
	}
}

func host2host(localAddress string, targetAddress string) {
	target, err := net.Dial("tcp", targetAddress)
	if err != nil {
		log.Fatalln("connect target address [" + localAddress + "] faild.")
	}
	local, err := net.Dial("tcp", localAddress)
	if err != nil {
		log.Fatalln("connect user's host [" + localAddress + "] faild.")
	}
	log.Println("connect target address [" + localAddress + "] and user's host [" + localAddress + "] success.")
	forward(target, local)
}

func start_server(address string) net.Listener {
	server, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("listen address [" + address + "] faild.")
	}
	log.Println("start listen at address:[" + address + "]")
	return server
	/*defer server.Close()

	for {
		conn, err := server.Accept()
		log.Println("accept a new client. remote address:[" + conn.RemoteAddr().String() +
			"], local address:[" + conn.LocalAddr().String() + "]")
		if err != nil {
			log.Println("accept a new client faild.", err.Error())
			continue
		}
		//go recvConnMsg(conn)
	}*/
}

func accept(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		log.Println("accept connect ["+conn.RemoteAddr().String()+"] faild.", err.Error())
		return nil
	}
	log.Println("accept a new client. remote address:[" + conn.RemoteAddr().String() + "], local address:[" + conn.LocalAddr().String() + "]")
	return conn
}

func forward(conn1 net.Conn, conn2 net.Conn) {
	var wg sync.WaitGroup
	// wait tow goroutines
	wg.Add(2)
	go connCopy(conn1, conn2, &wg)
	go connCopy(conn2, conn1, &wg)
	//blocking when the wg is locked
	wg.Wait()
}

func connCopy(conn1 net.Conn, conn2 net.Conn, wg *sync.WaitGroup) {
	//TODO:log, record the data from conn1 and conn2.
	io.Copy(conn1, conn2)
	conn1.Close()
	log.Println("close the connect at local:[" + conn1.LocalAddr().String() + "] and remote:[" + conn1.RemoteAddr().String() + "]")
	wg.Done()
}
