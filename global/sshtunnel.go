package global

import (
	"fmt"
	"io"
	"net"
	"sshtunnelweb/myorm"
	// "sshtunnelweb/util"
	"strings"

	"golang.org/x/crypto/ssh"
)

// const (
// 	host        string = "10.0.0.50"
// 	port        string = "22"
// 	tunnel_user string = "root"
// 	tunnel_pwd  string = "123"
// )

type LinkInfo struct {
	Local  string
	Remote string
}

// 连接中转机
func StartSshClient(tunnel_user string, tunnel_pwd string, host string, port string) (*ssh.Client, error) {
	scc := &ssh.ClientConfig{
		User: tunnel_user,
		Auth: []ssh.AuthMethod{
			ssh.Password(tunnel_pwd),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr := fmt.Sprintf("%s:%s", host, port)
	var f bool = false
	var (
		sc  *ssh.Client
		err error
	)
	for i := 0; i < 5; i++ {
		if sc, err = ssh.Dial("tcp", addr, scc); err != nil {
			Logger.Error(fmt.Sprintf("ssh 第%d次拨号 中转机失败\n%s", i+1, err.Error()))
			f = true
			continue
		} else {
			f = false
			break
		}
	}
	if f {
		return nil, err
	}
	return sc, nil
}

func CloseSshClient(st *ssh.Client) error {
	return st.Close()
}

func StartLocalListen(Local string) (net.Listener, error) {
	local_listen, err := net.Listen("tcp", Local)
	if err != nil {
		Logger.Error(fmt.Sprintf("打开本地监听端口失败: %s ", err.Error()))
		return nil, err
	}
	return local_listen, nil
}

func CloseLocalListen(local_listen net.Listener) {
	local_listen.Close()
}

// 先开启本地监听 获取net.listener
// 然后在打开本地accept的conn和远程的conn，用于交互数据传输
func StartTunnel(local_listen net.Listener, Remote string, st *ssh.Client) {

	for {
		defer func() {
			err := recover()
			if err != nil {
				Logger.Error(fmt.Printf("recover receive a err: %+v \n", err))

			}
		}()
		//等待交互信息传输，才会进行下一步，不然就一直等待，for循环就卡在这里
		// 这个等待才是关键，不然for一直循环，内存很快撑爆
		// 没法复用accept返回的net.conn
		local, err := local_listen.Accept()

		if err != nil {
			Logger.Error("本地accept失败 :" + err.Error())
			f := false
			for _, v := range GlobalSshtunnelInfo {
				if v == &local_listen {
					f = true
				}
			}
			if f {
				continue
			} else {
				break
			}
		}
		f := false // 判定是否有权限可以继续访问
		// fmt.Println("accept net.conn addr: ", local.LocalAddr().String(), strings.Split(local.RemoteAddr().String(), ":")[0])
		// 用户客户端 local.RemoteAddr().String() 用户的IP端口  local.LocalAddr().String() 本机IP端口
		// accept net.conn addr:  172.16.1.201:57791 192.168.1.60:25755
		tmpIP := strings.Split(local.RemoteAddr().String(), ":")
		// len(tmpIP)
		requestIP := tmpIP[0]                                          // 研发请求过来的IP,ipv4 OK, ipv6失败
		localPort := strings.Split(local.LocalAddr().String(), ":")[1] // 转发服务监听的本地端口
		findUser := myorm.User{}
		if err := DB.Model(&myorm.User{}).Where(myorm.User{Ip: requestIP}).First(&findUser).Error; err != nil {
			Logger.Error(requestIP + "查找关联用户失败: " + err.Error())
			continue
		}
		if findUser.IsAdmin { // admin可以使用所有的转发
			f = true
		} else {
			// 根据IP 查找用户关联的ssh转发信息，无论是没找到用户，还是没找到转发信息，
			// 都直接忽略本次数据请求，进入下一次本地端口监听等待中
			tmpuser := myorm.User{}
			if err := DB.Model(&myorm.User{}).Where(myorm.User{Ip: requestIP}).Preload("Conn").First(&tmpuser).Error; err != nil {
				Logger.Error(requestIP + "查找db失败: " + err.Error())
				continue
			}

			for _, k := range tmpuser.Conn {
				if k.Local == localPort {
					f = true
				}
			}
		}
		if f {

			// 没法复用dail返回的net.conn
			remote, err := st.Dial("tcp", Remote)
			// fmt.Println("remote addr ", &remote)
			if err != nil {
				fmt.Printf("连接远程目标%s失败 ", Remote)
				fmt.Println(err)
				local.Close()
				remote.Close()
				break
			}

			// fmt.Println("dail net.conn addr: ", remote.LocalAddr().String(), remote.RemoteAddr().String())
			// dail net.conn addr:  0.0.0.0:0 0.0.0.0:0

			// fmt.Println("连接远程目标主机端口成功")
			// 每次交互完信息，把local remote自动close
			go transfer(&local, &remote)
		} else {
			continue
		}
	}
}

func transfer(local *net.Conn, remote *net.Conn) {
	// 每次交互完信息，把local remote自动close
	defer (*local).Close()
	defer (*remote).Close()
	go func() {
		io.Copy((*remote), (*local))
	}()

	io.Copy((*local), (*remote))
}
