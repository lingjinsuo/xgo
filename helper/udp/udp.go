package udp

import (
	"fmt"
	"gitlab.chuangzhen-sh.net/golang/xgo/logging/applogger"
	"go.uber.org/zap"
	"net"
)

type Datagrams struct {
	Data []byte
}

type Client struct {
	conn   *net.UDPConn
	ch     chan *Datagrams
	bufsiz int
	cb     func(datagrams *Datagrams)
}

func ListenMulticast(address string, bufsiz int, cb func(datagrams *Datagrams)) (*Client, error) {
	c := &Client{
		ch:     make(chan *Datagrams, 1),
		bufsiz: 1500,
		cb:     func(datagrams *Datagrams) {},
	}

	if bufsiz > 0 {
		c.bufsiz = bufsiz
	}

	if cb != nil {
		c.cb = cb
	}

	applogger.Info(fmt.Sprintf("udp listen %s", address))

	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		applogger.Error("udp listen failed, error", zap.Any("err", err))
		return nil, err
	}
	c.conn, err = net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		applogger.Error("udp listen failed, error", zap.Any("err", err))
		return nil, err
	}

	go func() {
		for datum := range c.ch {
			c.cb(datum)
		}
	}()

	go c.read()
	return c, nil
}

func (c *Client) read() {
	defer c.conn.Close()
	data := make([]byte, c.bufsiz)
	for {
		n, remoteAddr, err := c.conn.ReadFromUDP(data)
		if err != nil {
			applogger.Error("from addr error during read:", zap.Any("addr", remoteAddr), zap.Any("error", err))
		}

		if err != nil {
			applogger.Error("from addr error during read:", zap.Any("addr", remoteAddr), zap.Any("error", err))
		} else {
			c.ch <- &Datagrams{Data: data[:n]}
		}
	}
}
