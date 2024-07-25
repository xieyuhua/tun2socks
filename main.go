package main

import (
	"go-tun2socks/dns"
	"go-tun2socks/tun2socks"
	"flag"
)

var tunDevice = flag.String("dev", "xieyuhua", "tunDevice name")
var tunAddr   = flag.String("addr", "192.168.123.2", "tunAddr 192.168.123.2")
var mask      = flag.String("mask", "255.255.255.0", "mask")
var tunGW     = flag.String("gate", "192.168.123.1", "gate")
var socksAddr = flag.String("proxy", "socks5://xieyuhua:xieyuhua@192.168.9.21:1080", "socksAddr")
var tunDns    = flag.String("dns", "114.114.114.114:53", "tunDns")

func main() {
    flag.Parse()
    
	//start local dns server (doh)
	dns.StartDns(*tunDns)
    //tunDevice string, tunAddr , tunMask , tunGW ,
    //mtu int, _sock5Addr string, _tunDNS string
    // netsh interface ipv4 add route 192.168.9.26 "xieyuhua" 192.166.2.2 metric=1
	tun2socks.StartTunDevice(*tunDevice, *tunAddr, *mask, *tunGW, 1500, *socksAddr, *tunDns)
}
