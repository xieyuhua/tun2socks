# tun2socks

需求：网卡A和网卡B是不同的网络，网卡A本地默认网卡，网卡B虚拟网卡（代理到其他网络）
我们需要访问特定的ip地址使用网卡B，其他正常访问走网卡A   
- route add 192.168.0.0 mask 255.255.0.0 192.168.123.1 metric 6 -p    
```
release:
	GOOS=linux go build -o tun2socks_linux main.go
	GOOS=linux GOARCH=arm go build -o tun2socks_linux_arm main.go
	GOOS=darwin go build -o tun2socks_darwin main.go
	GOOS=windows GOARCH=amd64 go build -o tun2socks_windows_64.exe main.go
	GOOS=windows GOARCH=386 go build -o tun2socks_windows_32.exe main.go
```

# use 

```
//socks5 auth 
//这里的0x05表示使用的是SOCKS版本5，0x01表示一个方法的数量，，0x02表示一个通常称为"USERNAME/PASSWORD"，  0x00表示选择的是无需认证的方法。
socksConn.Write([]byte{0x05, 0x01, 0x00})
authBack := make([]byte, 2)
_, err := io.ReadFull(socksConn, authBack)
if err != nil {
	log.Println(err)
	return err
}
```


```
网卡A的网关为192.168.5.1，网卡B的网关为192.168.123.1

netsh int ip reset （重置ip设置）
netsh winsock reset （充值网络设置）
netsh winhttp reset proxy （重置代理设置）
ipconfig /flushdns （刷新dns缓存）

// METRIC 1 越小优先级越高
route delete 0.0.0.0 mask 0.0.0.0 192.168.5.1
route add 0.0.0.0 mask 0.0.0.0 192.168.5.1 metric 100 -p 
route print

//创建设备 xieyuhua
tun2socks-windows-amd64.exe -device xieyuhua -proxy socks5://192.168.9.21:1080
go-tun2socks.exe -proxy 192.168.9.21:1080 -addr 192.168.123.2 -gate 192.168.123.1 -mask 255.255.255.0
# 设置地址和子网掩码
netsh interface ipv4 set address name="xieyuhua" source=static addr=192.168.123.2 mask=255.255.255.0
#设置 dns
netsh interface ipv4 set dnsservers name="xieyuhua" static address=114.114.114.114 register=none validate=no
#指定ip走 192.168.123.1 这个网卡网关
route add 192.168.9.26 mask 255.255.255.255 192.168.123.1 metric 6 -p 
#设置网关 不走默认节点，就设置大一点，不然影响本地网络
netsh interface ipv4 add route 0.0.0.0/0 "xieyuhua" 192.168.123.1 metric=200
```


# thank

  github.com/google/netstack
  
  github.com/miekg/dns



