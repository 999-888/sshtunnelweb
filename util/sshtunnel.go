package util

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func GetOnePort() string {
	var port string
	for {
		rand.Seed(time.Now().UnixNano())
		port = fmt.Sprintf("%d%d%d%d%d", rand.Intn(3)+3, rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10))
		if _, err := net.Dial("tcp", "localhost:"+port); err == nil {
			continue
		}
		break
	}
	return port
}
