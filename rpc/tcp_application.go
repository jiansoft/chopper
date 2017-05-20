package rpc

import (
	"fmt"
	"log"
	"net"
)

//處理連線的 發送、接收、連入事件管理

type tcpApplication struct {
	bufferManager *BufferManager
	listener      *net.TCPListener
	config        ApplicationConfig
	serverEvent   IServerEvent
}

type ApplicationConfig struct {
	Host               string
	Port               int
	ConnectionPoolSize int
	BufferSize         int
}

type IServerEvent interface {
	CreateApplication()
	CreateSession(ISession) IClientEvent
}

type IClientEvent interface {
	OnDisconnect()
	OnRecvComplete([]byte)
}

//
//var once sync.Once
//var _self *tcpApplication
//
//const (
//    VERSION = "1.0"
//)

//func FetchInstance() *tcpApplication {
//	if _self == nil {
//		once.Do(func() {
//			if nil == _self {
//				_self = new(tcpApplication)
//			}
//		})
//	}
//	return _self
//}

func NewApplication(config ApplicationConfig, serverEvent IServerEvent) *tcpApplication {
	//s := FetchInstance()
	s := new(tcpApplication)
	s.serverEvent = serverEvent
	s.config = config
	s.bufferManager = NewBufferManager(config.ConnectionPoolSize, config.BufferSize)
	s.serverEvent.CreateApplication()
	return s
}

func (t *tcpApplication) Start() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", t.config.Host, t.config.Port))
	if err != nil {
		log.Printf("ResolveTCPAddr error => %v", err)
		return
	}
	t.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Printf("ListenTCP error => %v", err)
		return
	}

	go func() {
		for {
			conn, err := t.listener.AcceptTCP()
			if err != nil {
				log.Printf("AcceptTCP error => %v", err)
				continue
			}
			session := newClientSession(t, conn)
			session.clientEvent = t.serverEvent.CreateSession(session)
			session.Start()
		}
	}()
}

func (t *tcpApplication) OnDisconnect(c *ClientSession) {
	if c.bufferOffst >= 0 {
		t.bufferManager.FreeBuffer(c)
	}
}