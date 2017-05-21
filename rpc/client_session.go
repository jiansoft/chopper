package rpc

import (
	//    "log"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

type ISession interface {
	Start()
	Disconnect()
	Send(buf []byte)
	RemoteAddr() net.Addr
}

type ClientSession struct {
	identifyId      string
	tcpConn         *net.TCPConn
	chanSend        chan []byte
	chanClose       chan bool
	isClose         int32
	tcpBinaryHander *TcpBinaryHander
	tcpApp          *tcpApplication
	buffers         []byte
	bufferOffst     int64
	clientEvent     IClientEvent
}

func (c *ClientSession) init(tcpApp *tcpApplication, conn *net.TCPConn) *ClientSession {
	c.identifyId = fmt.Sprintf("%p-%d", &c, time.Now().Unix())
	c.tcpConn = conn
	c.isClose = 0
	c.chanSend = make(chan []byte, 12)
	c.chanClose = make(chan bool)
	c.tcpBinaryHander = NewTcpBinaryHander(c.onSocketReceived)
	c.tcpApp = tcpApp
	c.tcpApp.bufferManager.SetBuffer(c)
	return c
}

func newClientSession(tcpApp *tcpApplication, conn *net.TCPConn) *ClientSession {
	return new(ClientSession).init(tcpApp, conn)
}

func (c *ClientSession) Start() {
	go c.ioCompleted()
	go c.doReceive()
}

func (c *ClientSession) Send(buf []byte) {
	if c.isClose != 0 {
		return
	}
	c.chanSend <- buf
}

func (c *ClientSession) Disconnect() {
	if atomic.SwapInt32(&c.isClose, 1) != 0 {
		return
	}
	close(c.chanSend)
	close(c.chanClose)
	c.tcpConn.Close()

	c.clientEvent.OnDisconnect()
	c.tcpApp.OnDisconnect(c)
	c.tcpBinaryHander.Dispose()
	c.buffers = nil
	c.clientEvent = nil
	c.tcpConn = nil
	c.tcpApp = nil
}

func (c ClientSession) RemoteAddr() net.Addr {
	return c.tcpConn.RemoteAddr()
}

func (c *ClientSession) onSocketReceived(data []byte) {
	c.clientEvent.OnRecvComplete(data)
}

func (c *ClientSession) ioCompleted() {
	for c.isClose == 0 {
		select {
		case buf := <-c.chanSend:
			c.tcpConn.Write(c.tcpBinaryHander.Wrap(buf))
		case _ = <-c.chanClose:
			c.Disconnect()
			break
		}
	}
}

func (c *ClientSession) doReceive() {
	for c.isClose == 0 {
		n, err := c.tcpConn.Read(c.buffers)
		if err != nil {
			c.chanClose <- true
			break
		}
		c.tcpBinaryHander.Parse(c.buffers, n)
	}
}
