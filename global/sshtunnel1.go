package global

import (
	// "fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"sshtunnelweb/myorm"
	"sync"
)

// type Sshinfo struct {
// 	// User   []string
// 	// Remote string
// 	St *net.Listener
// }

var (
	// {"80": *net.Listener,}
	GlobalSshtunnelInfo map[string]*net.Listener = make(map[string]*net.Listener, 0)

	ST *ssh.Client
	Lk sync.RWMutex

	// {"port1": {"ip1": "1", "ip2": "1",}, "port2": {"ip1": "1", "ip3": "1",}} map好判断是否存在
	// 当授权关联普通账户后，把关联账户的IP，自动加到关联的port的map中
	// 取消关联时，把关联账户的IP，从port的map中删除掉
	// 检测权限时，探测该port中是否有该IP，有这通过，无则拒绝
	LocalPortAndUserIP map[string]map[string]string = make(map[string]map[string]string, 0)
)

func StartST() error {
	Lk.RLock()
	if ST == nil {
		tmpres := myorm.Sshinfo{}

		err := DB.Model(&myorm.Sshinfo{}).First(&tmpres).Error
		if err != nil {
			return err
		}
		Lk.RUnlock()
		username := tmpres.Username
		passwd := tmpres.Passwd
		host := tmpres.Host
		port := tmpres.Port
		st, err := StartSshClient(username, passwd, host, port)
		if err != nil {
			return err
		}
		Lk.Lock()
		// global.SshClient = true
		ST = st
		Lk.Unlock()
		Logger.Info("启动连接中转机成功")
	}
	return nil
}
