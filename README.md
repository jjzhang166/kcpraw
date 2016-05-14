# kcptun(KT)
TCP流转换为KCP+UDP流，>>>[下载地址](https://github.com/xtaci/kcptun/releases/latest)<<< 用于***高丢包***环境中的数据传输，工作示意图:      
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
***kcptun是[kcp](https://github.com/skywind3000/kcp)协议的一个简单应用，可以用于任意tcp网络程序的传输承载，以提高软件网络流畅度(如浏览器，telnet等)，降低掉线情况。***   

<img src="kitty.jpg" style="width: 300px;"/>

### Docker
```
docker pull xtaci/kcptun
```

### 使用方法
```
D:\>client_windows_amd64.exe -h
NAME:
   kcptun - kcptun client

USAGE:
   client_windows_amd64.exe [global options] command [command options] [arguments...]

VERSION:
   20160507

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --localaddr, -l ":12948"     local listen addr:
   --remoteaddr, -r "vps:29900" kcp server addr
   --key "it's a secrect"       key for communcation, must be the same as kcptun server [$KCPTUN_KEY]
   --mode "fast"                mode for communication: fast, normal, default
   --tuncrypt                   enable tunnel encryption, adds extra secrecy for data transfer
   --help, -h                   show help
   --version, -v                print the version

D:\>server_windows_amd64.exe -h
NAME:
   kcptun - kcptun server

USAGE:
   server_windows_amd64.exe [global options] command [command options] [arguments...]

VERSION:
   20160507

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --listen, -l ":29900"                kcp server listen addr:
   --target, -t "127.0.0.1:12948"       target server addr
   --key "it's a secrect"               key for communcation, must be the same as kcptun client [$KCPTUN_KEY]
   --mode "fast"                        mode for communication: fast, normal, default
   --tuncrypt                           enable tunnel encryption, adds extra secrecy for data transfer
   --help, -h                           show help
   --version, -v                        print the version
```
### 适用范围（包括但不限于）:           
1. 网络游戏的数据传输        
2. 跨运营商的流量传输               
3. 其他高丢包，高干扰通信环境的TCP数据传输      

# 免责申明
用户以各种方式使用本软件（包括但不限于修改使用、直接使用、通过第三方使用）的过程中，不得以任何方式利用本软件直接或间接从事违反中国法律、以及社会公德的行为。软件的使用者需对自身行为负责，因使用软件引发的一切纠纷，由使用者承担全部法律及连带责任。作者不承担任何法律及连带责任。       

对免责声明的解释、修改及更新权均属于作者本人所有。
