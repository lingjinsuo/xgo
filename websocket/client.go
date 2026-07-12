package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.chuangzhen-sh.net/golang/xgo/logging/applogger"
	"sync"
	"time"
)

const (
	TimerIntervalSecond = 5
	ReconnectWaitSecond = 60
)

type ConnectedHandler func(c *Connect)
type MessageHandler func(c *Connect, message string) (interface{}, error)
type ResponseHandler func(c *Connect, response interface{})

type Connect struct {
	host              string
	path              string
	conn              *websocket.Conn
	connectedHandler  ConnectedHandler
	messageHandler    MessageHandler
	responseHandler   ResponseHandler
	stopReadChannel   chan int
	stopTickerChannel chan int
	ticker            *time.Ticker
	lastReceivedTime  time.Time
	sendMutex         *sync.Mutex
}

func NewClient(host, path string) *Connect {
	return &Connect{
		host:              host,
		path:              path,
		stopReadChannel:   make(chan int, 1),
		stopTickerChannel: make(chan int, 1),
		sendMutex:         &sync.Mutex{},
		connectedHandler:  func(c *Connect) {},
		messageHandler:    func(c *Connect, message string) (interface{}, error) { return nil, nil },
		responseHandler:   func(c *Connect, response interface{}) {},
	}
}

func (c *Connect) SetHandler(connHandler ConnectedHandler, msgHandler MessageHandler, repHandler ResponseHandler) {
	if connHandler != nil {
		c.connectedHandler = connHandler
	}

	if msgHandler != nil {
		c.messageHandler = msgHandler
	}

	if repHandler != nil {
		c.responseHandler = repHandler
	}
}

func (c *Connect) Connect(autoReconnect bool) {
	c.connectWebSocket()

	if autoReconnect {
		c.startTicker()
	}
}

func (c *Connect) SendString(data []byte) error {
	if c.conn == nil {
		return errors.New("WebSocket sent error: no connection available")
	}

	c.sendMutex.Lock()
	err := c.conn.WriteMessage(websocket.TextMessage, data)
	c.sendMutex.Unlock()
	return err
}

func (c *Connect) SendJson(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.SendString(b)
}

func (c *Connect) Close() {
	c.stopTicker()
	c.disconnectWebSocket()
}

func (c *Connect) connectWebSocket() {
	var err error
	url := fmt.Sprintf("ws://%s%s", c.host, c.path)
	applogger.Debug("WebSocket connecting...")
	c.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		applogger.Error("WebSocket connected error: %s", err)
		return
	}
	applogger.Info("WebSocket connected")
	c.lastReceivedTime = time.Now()

	// start loop to read and handle message
	c.startReadLoop()

	c.connectedHandler(c)
}

// disconnect with server
func (c *Connect) disconnectWebSocket() {
	if c.conn == nil {
		return
	}

	// start a new goroutine to send stop signal
	go c.stopReadLoop()

	applogger.Debug("WebSocket disconnecting...")
	err := c.conn.Close()
	if err != nil {
		applogger.Error("WebSocket disconnect error: %s", err)
		return
	}

	applogger.Info("WebSocket disconnected")
}

func (c *Connect) startTicker() {
	c.ticker = time.NewTicker(TimerIntervalSecond * time.Second)

	go c.tickerLoop()
}

func (c *Connect) stopTicker() {
	c.ticker.Stop()
	c.stopTickerChannel <- 1
}

func (c *Connect) tickerLoop() {
	applogger.Debug("tickerLoop started")
	for {
		select {
		// Receive data from stopChannel
		case <-c.stopTickerChannel:
			applogger.Debug("tickerLoop stopped")
			return

		// Receive tick from tickChannel
		case <-c.ticker.C:
			elapsedSecond := time.Now().Sub(c.lastReceivedTime).Seconds()
			applogger.Debug("WebSocket received data %f sec ago", elapsedSecond)

			if elapsedSecond > ReconnectWaitSecond {
				applogger.Info("WebSocket reconnect...")
				c.disconnectWebSocket()
				c.connectWebSocket()
			}
		}
	}
}

func (c *Connect) startReadLoop() {
	go c.readLoop()
}

func (c *Connect) stopReadLoop() {
	c.stopReadChannel <- 1
}

func (c *Connect) readLoop() {
	applogger.Debug("readLoop started")
	for {
		select {
		// Receive data from stopChannel
		case <-c.stopReadChannel:
			applogger.Debug("readLoop stopped")
			return

		default:
			if c.conn == nil {
				applogger.Error("Read error: no connection available")
				time.Sleep(TimerIntervalSecond * time.Second)
				continue
			}

			msgType, buf, err := c.conn.ReadMessage()
			if err != nil {
				applogger.Error("Read error: %s", err)
				time.Sleep(TimerIntervalSecond * time.Second)
				continue
			}

			c.lastReceivedTime = time.Now()

			if msgType == websocket.PingMessage {
				applogger.Debug("Read PingMessage")
				_ = c.conn.WriteMessage(websocket.PongMessage, []byte("thisisappdata"))
			}

			var message string
			if msgType == websocket.BinaryMessage {
				message = string(buf)
			} else {
				message = string(buf)
			}

			result, err := c.messageHandler(c, message)
			if err != nil {
				applogger.Error("Handle message error: %s", err)
				continue
			}
			c.responseHandler(c, result)
		}
	}
}
