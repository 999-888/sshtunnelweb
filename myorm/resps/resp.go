package resps

type Conn struct {
	ID uint
	// port
	Local string
	// ip:port
	Svcname string
}

type Remote struct {
	ID      uint
	Remote  string
	Svcname string
}

type RemoteSelect struct {
	ID      uint
	Svcname string
}

type User struct {
	ID       uint
	Username string
	Passwd   string
	Ip       string
	IsAdmin  bool
}

type LocalPortUser struct {
	ID       uint
	Username string
}

type Sshinfo struct {
	ID       uint
	Username string
	Passwd   string
	Port     string
	Host     string
}

type Workflow struct {
	ID        uint
	Username  string
	Localport string
	Svcname   string
}
