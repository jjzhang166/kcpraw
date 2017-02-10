## kcptun-raw

运行在伪造的 TCP 协议之上的 kcptun, 主要目的是避免 ISP 对 UDP 协议可能的 QOS.  
在三次握手后会进行 HTTP 握手, 将流量伪装成 HTTP 流量.  
kcptun 的具体参数与使用方法参见 [kcptun](https://github.com/xtaci/kcptun)  

### 基本用法  

为了使用原始套接字，服务端与客户端都需要 root 权限  
  
服务端  
```
./server_linux_amd64 -t "TARGET_IP:TARGET_PORT" -l "KCP_SERVER_IP:KCP_SERVER_PORT"
```  
客户端  
```
./client_darwin_amd64 -r "KCP_SERVER_IP:KCP_SERVER_PORT" -l ":LOCAL_PORT"
```

### 注意事项
~~为了避免内核返回的 RST 报文影响连接的建立，需要添加相应的 iptables 规则~~  
现在 linux 下的客户端和服务端会自动添加和清理 iptables 规则  
windows 客户端参考 [windows firewall port exceptions](https://www.veritas.com/support/en_US/article.000085856) 链接中的方法为 LOCAL_PORT 设置防火墙规则  
macos 下使用客户端可以参考 [enable steath mode](http://osxdaily.com/2015/11/18/enable-stealth-mode-mac-os-x-firewall/) 打开静默模式即可  
由于使用 windows 或者 macos 服务器运行 kcptun 服务的情况并不多见, 并没有进行相应的测试, 所以并不能保证服务端在非 linux 环境下能够正常使用  

伪装的 Host 可以通过选项 --host <name> 进行设置  
如果不希望伪装为 HTTP 流量可以通过设置选项 --nohttp 关闭此功能, 注意客户端和服务端在这一选项上必须保持一致

### 构建

```
go get github.com/ccsexyz/kcptun/client  
go get github.com/ccsexyz/kcptun/server  
```

windows 下依赖 [winpcap](http://www.winpcap.org/install/) 和 gcc 请自行解决环境问题    
  
对 windows10 使用者你可能需要参考这个链接 http://stackoverflow.com/questions/38047858/compile-gopacket-on-windows-64bit