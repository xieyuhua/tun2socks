package socks

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"net/url"  
	"strconv"
	"strings"
)

/*to socks5*/
func SocksCmds(socksConn net.Conn, cmd uint8, hosts string) error {
	//socks5 auth 
	log.Println(hosts)
	//这里的0x05表示使用的是SOCKS版本5，0x01表示一个方法的数量，，0x02表示一个通常称为"USERNAME/PASSWORD"，  0x00表示选择的是无需认证的方法。
	//socksConn.Write([]byte{0x05, 0x01, 0x00})
	socksConn.Write([]byte{0x05, 0x02,0x00, 0x02})
	// 0x05: SOCKS 版本 5, 0x02: 2个方法, 0x00: 无需认证, 0x02: USERNAME/PASSWORD  
	authBack := make([]byte, 2)
	_, err := io.ReadFull(socksConn, authBack)
	if err != nil {
		log.Println(err)
		return err
	}
	
	parsedURL, err := url.Parse(hosts)  
	if err != nil {  
		log.Println(err)
		return  err
	}
	
	// 提取主机和端口  
	host := parsedURL.Hostname()  
	port := parsedURL.Port()  
	if port == "" {  
	    port = "1080"
	}
	
	if authBack[1] == 0x02 { // 如果服务器选择了 USERNAME/PASSWORD 认证  
    	// 提取用户名和密码  
    	userinfo := parsedURL.User  
		username := userinfo.Username()  
		password, _ := userinfo.Password() // 注意：Password() 方法不会返回错误  
    	
		// 发送用户名和密码  
		userLen := uint8(len(username))  
		passLen := uint8(len(password))  
		// 下面的代码将发送完整的用户名和密码  
		var authReqFull bytes.Buffer  
		authReqFull.WriteByte(0x01) // 子协议版本  
		authReqFull.WriteByte(userLen)  
		authReqFull.WriteString(username)  
		authReqFull.WriteByte(passLen)  
		authReqFull.WriteString(password)  
		socksConn.Write(authReqFull.Bytes())  
  
		// 读取认证结果  
		authResult := make([]byte, 2)  
		if _, err := io.ReadFull(socksConn, authResult); err != nil {  
			log.Println(err)  
			return err  
		}  
		if authResult[1] != 0x00 { // 认证失败  
			log.Println(err)
			return err  
		}  
	} else if authBack[1] != 0x00 {  
		log.Println(err)
		return err  
	}
	
	//connect head
	rAddr := net.ParseIP(host)
	_port, _ := strconv.Atoi(port)
	msg := []byte{0x05, cmd, 0x00, 0x01}
	buffer := bytes.NewBuffer(msg)
	//ip
	binary.Write(buffer, binary.BigEndian, rAddr.To4())
	//port
	binary.Write(buffer, binary.BigEndian, uint16(_port))
	socksConn.Write(buffer.Bytes())
	conectBack := make([]byte, 10)
	_, err = io.ReadFull(socksConn, conectBack)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}


/*to socks5*/
func SocksCmd(socksConn net.Conn, cmd uint8, host string) error {
	//socks5 auth 
	//这里的0x05表示使用的是SOCKS版本5，0x01表示一个方法的数量，，0x02表示一个通常称为"USERNAME/PASSWORD"，  0x00表示选择的是无需认证的方法。
	socksConn.Write([]byte{0x05, 0x01, 0x00})
	authBack := make([]byte, 2)
	_, err := io.ReadFull(socksConn, authBack)
	if err != nil {
		log.Println(err)
		return err
	}
	//connect head
	hosts := strings.Split(host, ":")
	rAddr := net.ParseIP(hosts[0])
	_port, _ := strconv.Atoi(hosts[1])
	msg := []byte{0x05, cmd, 0x00, 0x01}
	buffer := bytes.NewBuffer(msg)
	//ip
	binary.Write(buffer, binary.BigEndian, rAddr.To4())
	//port
	binary.Write(buffer, binary.BigEndian, uint16(_port))
	socksConn.Write(buffer.Bytes())
	conectBack := make([]byte, 10)
	_, err = io.ReadFull(socksConn, conectBack)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
