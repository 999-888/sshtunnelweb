package resps

import "time"

type Conn struct {
	ID uint
	// port
	Local string
	// ip:port
	Svcname   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Remote struct {
	ID        uint
	Remote    string
	Svcname   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RemoteSelect struct {
	ID      uint
	Svcname string
}

type User struct {
	ID        uint
	Username  string
	Passwd    string
	Ip        string
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LocalPortUser struct {
	ID       uint
	Username string
}

type Sshinfo struct {
	ID        uint
	Username  string
	Passwd    string
	Port      string
	Host      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Workflow struct {
	ID        uint
	Username  string
	Localport string
	Svcname   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
