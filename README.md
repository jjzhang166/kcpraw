# kcptun(KT)
TCP流转换为KCP+UDP流，:snowflake:[下载地址](https://github.com/xtaci/kcptun/releases/latest):snowflake:工作示意图：     
```
                +---------------------------------------+
                |                                       |
                |                KCPTUN                 |
                |                                       |
+--------+      |  +------------+       +------------+  |      +--------+
|        |      |  |            |       |            |  |      |        |
| Client | +--> |  | KCP Client | +---> | KCP Server |  | +--> | Server |
|        | TCP  |  |            |  UDP  |            |  | TCP  |        |
+--------+      |  +------------+       +------------+  |      +--------+
                |                                       |
                |                                       |
                +---------------------------------------+
```
***kcptun是[kcp](https://github.com/skywind3000/kcp)协议的一个简单应用，可以用于任意tcp网络程序的传输承载，以提高网络流畅度，降低掉线情况。***   

<img src="kitty.jpg" style="width: 300px;"/>

### 使用の方法 :hash:
```bat
D:\>client_windows_amd64.exe -h
NAME:
   kcptun - kcptun client

USAGE:
   client_windows_amd64.exe [global options] command [command options] [arguments...]

VERSION:
   20160517

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --localaddr, -l ":12948"	local listen address
   --remoteaddr, -r "vps:29900"	kcp server address
   --key "it's a secrect"	key for communcation, must be the same as kcptun server [$KCPTUN_KEY]
   --mode "fast"		mode for communication: fast2, fast, normal, default
   --tuncrypt			enable tunnel encryption, adds extra secrecy for data transfer
   --mtu "1350"			set MTU of UDP packets, suggest 'tracepath' to discover path mtu
   --sndwnd "128"		set send window size(num of packets)
   --rcvwnd "1024"		set receive window size(num of packets)
   --help, -h			show help
   --version, -v		print the version

D:\>server_windows_amd64.exe -h
NAME:
   kcptun - kcptun server

USAGE:
   server_windows_amd64.exe [global options] command [command options] [arguments...]

VERSION:
   20160517

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --listen, -l ":29900"		kcp server listen address
   --target, -t "127.0.0.1:12948"	target server address
   --key "it's a secrect"		key for communcation, must be the same as kcptun client [$KCPTUN_KEY]
   --mode "fast"			mode for communication: fast2, fast, normal, default
   --tuncrypt				enable tunnel encryption, adds extra secrecy for data transfer
   --mtu "1350"				set MTU of UDP packets, suggest 'tracepath' to discover path mtu
   --sndwnd "1024"			set send window size(num of packets)
   --rcvwnd "1024"			set receive window size(num of packets)
   --help, -h				show help
   --version, -v			print the version
```
### 适用范围限定:hash:     
1. 实时网络游戏的数据传输        
2. 跨运营商的流量传输               
3. 其他高丢包通信链路的TCP承载      

### 性能对比:hash:
```
root@vultr:~# iperf -s
------------------------------------------------------------
Server listening on TCP port 5001
TCP window size: 4.00 MByte (default)
------------------------------------------------------------
[  4] local 172.7.7.1 port 5001 connected with 172.7.7.2 port 55453
[ ID] Interval       Transfer     Bandwidth
[  4]  0.0-18.0 sec  5.50 MBytes  2.56 Mbits/sec     <-- connection via kcptun
[  5] local 45.32.xxx.xxx port 5001 connected with 218.88.xxx.xxx port 17220
[  5]  0.0-17.9 sec  2.12 MBytes   997 Kbits/sec     <-- direct connnection via tcp
```

# 免责申明
用户以各种方式使用本软件（包括但不限于修改使用、直接使用、通过第三方使用）的过程中，不得以任何方式利用本软件直接或间接从事违反中国法律、以及社会公德的行为。软件的使用者需对自身行为负责，因使用软件引发的一切纠纷，由使用者承担全部法律及连带责任。作者不承担任何法律及连带责任。       

对免责声明的解释、修改及更新权均属于作者本人所有。
