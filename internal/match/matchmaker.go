package match

type ServerInfo struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

func FindServer() ServerInfo {
	return ServerInfo{
		IP:   "127.0.0.1",
		Port: 7777,
	}
}
