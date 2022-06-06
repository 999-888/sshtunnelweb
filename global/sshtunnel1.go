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
	// {"80": {"info": sshinfo },}
	// {"80": 本地80端口的地址, "800": 本地800端口的地址,}
	// GlobalSshtunnelInfo map[string]map[string]Sshinfo = make(map[string]map[string]Sshinfo, 0)
	GlobalSshtunnelInfo map[string]*net.Listener = make(map[string]*net.Listener, 0)
	// SshClient           bool                     = false // 判定是否对中转机做了拨号
	ST *ssh.Client
	Lk sync.RWMutex
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
