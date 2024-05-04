package types

func NewServiceInfo(name, host, port string) ServiceInfo {
	return ServiceInfo{Name: name, Host: host, Port: port}
}
