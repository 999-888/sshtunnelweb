package myorm

func GetAll() []interface{} {
	data := make([]interface{}, 0)
	data = append(data, &User{})
	data = append(data, &Conn{})
	data = append(data, &Sshinfo{})
	data = append(data, &Workflow{})
	return data
}
