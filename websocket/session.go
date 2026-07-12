package websocket

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.chuangzhen-sh.net/golang/xgo/logging/applogger"
	"sync"
	"time"
)

type Session struct {
	GinContext *gin.Context
	Alias      string
	Tags       map[string]interface{}
	server     *Server
	conn       *websocket.Conn
	send       chan *envelope
	once       sync.Once
	open       bool
	rwmutex    *sync.RWMutex
}

func (s *Session) Response(message string) {
	s.writeMessage(&envelope{
		t:      websocket.TextMessage,
		msg:    []byte(message),
		filter: nil,
	})
}

func (s *Session) ResponseBinary(message []byte) {
	s.writeMessage(&envelope{
		t:      websocket.BinaryMessage,
		msg:    message,
		filter: nil,
	})
}

func (s *Session) writeMessage(message *envelope) {
	if s.closed() {
		s.server.errorHandler(s, errors.New("tried to write to closed a session"))
		return
	}

	select {
	case s.send <- message:
		applogger.Debug("writeMessage t:%d msg:%s", message.t, string(message.msg))
	default:
		s.server.errorHandler(s, errors.New("session message buffer is full"))
	}
}

func (s *Session) writeRaw(message *envelope) error {
	if s.closed() {
		return errors.New("tried to write to a closed session")
	}

	s.conn.SetWriteDeadline(time.Now().Add(s.server.Config.WriteWait))
	err := s.conn.WriteMessage(message.t, message.msg)

	if err != nil {
		return err
	}

	return nil
}

func (s *Session) closed() bool {
	s.rwmutex.RLock()
	defer s.rwmutex.RUnlock()
	return !s.open
}

func (s *Session) Close() {
	if !s.closed() {
		s.rwmutex.Lock()
		s.open = false
		s.conn.Close()
		close(s.send)
		s.rwmutex.Unlock()
	}
}

func (s *Session) ping() {
	s.writeRaw(&envelope{t: websocket.PingMessage, msg: []byte{}})
}

func (s *Session) writePump() {
	ticker := time.NewTicker(s.server.Config.PingPeriod)
	defer ticker.Stop()

loop:
	for {
		select {
		case msg, ok := <-s.send:
			if !ok {
				break loop
			}

			err := s.writeRaw(msg)

			if err != nil {
				s.server.errorHandler(s, err)
				break loop
			}

			if msg.t == websocket.CloseMessage {
				break loop
			}

			if msg.t == websocket.TextMessage {
				s.server.messageSentHandler(s, msg.msg)
			}

			if msg.t == websocket.BinaryMessage {
				s.server.messageSentHandlerBinary(s, msg.msg)
			}
		case <-ticker.C:
			s.ping()
		}
	}
}

func (s *Session) readPump() {
	s.conn.SetReadLimit(s.server.Config.MaxMessageSize)
	s.conn.SetReadDeadline(time.Now().Add(s.server.Config.PongWait))

	s.conn.SetPongHandler(func(appData string) error {
		applogger.Debug("pong %s", appData)
		s.conn.SetReadDeadline(time.Now().Add(s.server.Config.PongWait))
		s.server.pongHandler(s)
		return nil
	})

	if s.server.closeHandler != nil {
		s.conn.SetCloseHandler(func(code int, text string) error {
			return s.server.closeHandler(s, code, text)
		})
	}

	for {
		t, message, err := s.conn.ReadMessage()

		if err != nil {
			s.server.errorHandler(s, err)
			break
		}

		if t == websocket.TextMessage {
			s.server.messageHandler(s, message)
		}

		if t == websocket.BinaryMessage {
			s.server.messageHandlerBinary(s, message)
		}
	}
}
