# ptt (Port Transmit Tool)
一款lcx在golang下的实现

# platform
Windows 7 7601 + go1.7.5(windows/amd64) 测试通过

# usage
- -listen port1 port2 

## 说明
同时监听port1端口和port2端口，当两个客户端主动连接上这两个监听端口之后，ptt负责这两个端口间的数据转发。

## 示例
`ptt -listen 1997 2017`

- -tran port1 ip:port2 

## 说明
本地开始监听port1端口，当port1端口上接收到来自客户端的主动连接之后，ptt将主动连接ip:port2，并且负责port1端口和ip:port2之间的数据转发。

## 示例
`ptt -tran 1997 192.168.1.2:338`

- -slave ip1:port1 ip2:port2

## 说明
本地开始主动连接ip1:port1主机和ip2:port2主机，当连接成功之后，ptt负责这两个主机之间的数据转发。

## 示例
`ptt -slave 127.0.0.1:3389 8.8.8.8:1997`

- log filepath example

## 示例
`ptt -listen 1997 2017 -log D:\ptt_log.txt`
`ptt -tran 1997 192.168.1.2:338 -log D:\ptt_log.txt`
`ptt -slave 127.0.0.1:3389 8.8.8.8:1997 -log D:\ptt_log.txt`

## 说明
`-log`为一个可选开关。如果在前面任意一个必选开关的末尾加上该开关，那么所有转发数据将会被记录到`D:\ptt_log.txt`文件中。（由于转发数据可能并非文本文件，建议使用UltraEdit等支持二进制查看的编辑器打开）


