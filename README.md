## kcptun-raw

运行在伪造的 TCP 协议之上的 kcptun, 主要目的是避免 ISP 对 UDP 协议可能的 QOS.  
kcptun 的具体参数与使用方法参见 [kcptun](https://github.com/xtaci/kcptun)  

### 基本用法  

服务端  
```
./server_linux_amd64 -t "TARGET_IP:8388" -l "KCP_SERVER_IP:4000"  
```  
客户端  
```
./client_darwin_amd64 -r "KCP_SERVER_IP:4000" -l ":8388"  
```
为了避免内核返回的 RST 报文影响连接的建立，需要在服务端设置相应的 iptables 规则  
```
iptables -A OUTPUT -p tcp --sport 4000 --tcp-flags RST RST -j DROP
```

为了使用原始套接字，服务端与客户端都需要 root 权限  
