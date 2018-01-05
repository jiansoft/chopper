package chopper

import (
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

type PeerSession struct {
	identifyId      string
	tcpConn         *net.TCPConn
	chanSend        chan []byte
	isClose         int32
	tcpBinaryHander *TcpBinaryHander
	tcpApp          *tcpApplication
	buffers         []byte
	bufferOffst     int64
	clientEvent     IClientEvent
}

func (c *PeerSession) init(tcpApp *tcpApplication, conn *net.TCPConn) *PeerSession {
	c.identifyId = fmt.Sprintf("%p-%d", &c, time.Now().Unix())
	c.tcpConn = conn
	c.isClose = 0
	c.chanSend = make(chan []byte, 12)
	c.tcpBinaryHander = NewTcpBinaryHander(c.onSocketReceived)
	c.tcpApp = tcpApp
	c.tcpApp.bufferManager.SetBuffer(c)
	return c
}

func newClientSession(tcpApp *tcpApplication, conn *net.TCPConn) *PeerSession {
	return new(PeerSession).init(tcpApp, conn)
}

func (c *PeerSession) Start() {
	go c.doReceive()
	go c.doWrite()
}

func (c *PeerSession) Send(buf []byte) {
	if c.isClose != 0 {
		return
	}
	c.chanSend <- buf
}

func (c *PeerSession) Disconnect() {
	if atomic.SwapInt32(&c.isClose, 1) != 0 {
		return
	}

	c.clientEvent.OnDisconnect()
	c.tcpApp.OnDisconnect(c)
	c.tcpConn.Close()
	c.tcpBinaryHander.Dispose()
	c.buffers = nil
	c.clientEvent = nil
	c.tcpConn = nil
	c.tcpApp = nil
	close(c.chanSend)
}

func (c PeerSession) RemoteAddr() net.Addr {
	return c.tcpConn.RemoteAddr()
}

func (c *PeerSession) onSocketReceived(data []byte) {
	c.clientEvent.OnRecvComplete(data)
}

func (c *PeerSession) doReceive() {
	for c.isClose == 0 {
		n, err := c.tcpConn.Read(c.buffers)
		if err != nil {
			c.Disconnect()
			break
		}
		c.tcpBinaryHander.Parse(c.buffers, n)
	}
}

func (c *PeerSession) doWrite() {
	for c.isClose == 0 {
		select {
		case buf := <-c.chanSend:
			buffer := c.tcpBinaryHander.Wrap(buf)
			c.tcpConn.Write(buffer)
		}
	}
}
