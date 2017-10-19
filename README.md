# ptt (Port Transmit Tool)
一款lcx在golang下的实现

# build
`go build ptt.go`

如果未出现任何消息则表示编译成功。（Linux哲学：没有消息就是最好的消息）

执行之后Linux或Mac OS在当前目录执行./ptt，Windows在当前目录执行ptt.exe，如果出现欢迎信息则表示一切正常。

如果编译出现错误，请检查当前系统上的golang是否被破坏，建议重装后再尝试或者直接下载已经编译好的二进制文件。

如果运行出现错误，请检查您所输入的参数是否有错，或者相应的端口被占用，请更换端口重试。如果未看见欢迎消息，请给编译好的可执行文件设置权限为777。如果无法写日志，请检查当前用户是否有权限在日志文件路径的读写权限。

# platform
Windows 7 7601 + go1.7.5(windows/amd64) 编译与测试通过
Ubuntu 16.04.1 + go1.6.2(linux/amd64) 编译与测试通过
Windows 2003 SP2 3790 + go1.9.1(windows/386) 编译与测试通过

# usage
- -listen port1 port2 

### 说明
同时监听port1端口和port2端口，当两个客户端主动连接上这两个监听端口之后，ptt负责这两个端口间的数据转发。

### 示例
`ptt -listen 1997 2017`

---

- -tran port1 ip:port2 

### 说明
本地开始监听port1端口，当port1端口上接收到来自客户端的主动连接之后，ptt将主动连接ip:port2，并且负责port1端口和ip:port2之间的数据转发。

### 示例
`ptt -tran 1997 192.168.1.2:338`

---

- -slave ip1:port1 ip2:port2

### 说明
本地开始主动连接ip1:port1主机和ip2:port2主机，当连接成功之后，ptt负责这两个主机之间的数据转发。

### 示例
`ptt -slave 127.0.0.1:3389 8.8.8.8:1997`

---

- log filepath

### 示例
`ptt -listen 1997 2017 -log D:/ptt`

`ptt -tran 1997 192.168.1.2:338 -log D:/ptt`

`ptt -slave 127.0.0.1:3389 8.8.8.8:1997 -log D:/ptt`

### 说明
`-log`为一个可选开关。如果在前面任意一个必选开关的末尾加上该开关，那么所有转发数据将会被记录到`D:/ptt/Y_m_d_H_i_s-agrs1-args2-args3.log`文件中，其中`YmdHis`以及`args`均会被替换为实际执行时的时间和参数。如果有特殊需求，可根据时间顺序，以及相关参数进行合并，以得到连续的转发数据日志记录。（由于转发数据可能并非文本文件，建议使用UltraEdit等支持二进制查看的编辑器打开）

警告：不要使用包含空格以及各种特殊字符的文件路径，比如说`C:\Documents and Settings\Administrator\桌面\go\bin`这个文件路径就是无效文件路径，因为其包含空格。

注意：由于日志流记录是即时的，建议将日志文件存储在机械硬盘分区中，而不要放在包括固态硬盘，U盘，SD卡等设备，防止大量小文件写入影响这些设备的寿命。

技巧：可使用Linux下的`tail -f`命令将转发数据实时显示出来。

# TODO
- UDP协议的转发支持
