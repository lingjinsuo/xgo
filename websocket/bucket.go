package websocket

import (
	"sync"
)

type bucket struct {
	sessRWMutex sync.RWMutex
	sessions    map[*Session]struct{}

	register   chan *Session
	unregister chan *Session

	broadcast chan *envelope
	exit      chan *envelope

	open bool
}

func (b *bucket) init() *bucket {
	b.sessRWMutex = sync.RWMutex{}
	b.sessions = make(map[*Session]struct{})
	b.register = make(chan *Session)
	b.unregister = make(chan *Session)
	b.broadcast = make(chan *envelope)
	b.exit = make(chan *envelope)
	b.open = true
	return b
}

func (b *bucket) run() {
loop:
	for {
		select {
		case s := <-b.register:
			b.doRegister(s)
		case s := <-b.unregister:
			b.doUnregister(s)
		case m := <-b.broadcast:
			b.doBroadcast(m)
		case m := <-b.exit:
			b.doExit(m)
			break loop
		}
	}
}

func (b *bucket) closed() bool {
	b.sessRWMutex.RLock()
	defer b.sessRWMutex.RUnlock()
	return !b.open
}

func (b *bucket) len() int {
	b.sessRWMutex.RLock()
	defer b.sessRWMutex.RUnlock()

	return len(b.sessions)
}

func (b *bucket) doRegister(session *Session) {
	b.sessRWMutex.Lock()
	b.sessions[session] = struct{}{}
	b.sessRWMutex.Unlock()
}

func (b *bucket) doUnregister(session *Session) {
	if _, ok := b.sessions[session]; ok {
		b.sessRWMutex.Lock()
		delete(b.sessions, session)
		b.sessRWMutex.Unlock()
	}
}

func (b *bucket) setTags(alias string, tags ...string) {
	b.sessRWMutex.Lock()
	for s := range b.sessions {
		if s.Alias == alias {
			for _, tag := range tags {
				s.Tags[tag] = true
			}
		}
	}
	b.sessRWMutex.Unlock()
}

func (b *bucket) delTags(alias string, tags ...string) {
	b.sessRWMutex.Lock()
	for s := range b.sessions {
		if s.Alias == alias {
			for _, tag := range tags {
				delete(s.Tags, tag)
			}
		}
	}
	b.sessRWMutex.Unlock()
}

func (b *bucket) doBroadcast(msg *envelope) {
	b.sessRWMutex.RLock()
	for s := range b.sessions {
		if msg.filter != nil {
			if msg.filter(s) {
				s.writeMessage(msg)
			}
		} else {
			s.writeMessage(msg)
		}
	}
	b.sessRWMutex.RUnlock()
}

func (b *bucket) doExit(msg *envelope) {
	b.sessRWMutex.Lock()
	for s := range b.sessions {
		s.writeMessage(msg)
		delete(b.sessions, s)
		s.Close()
	}
	b.open = false
	b.sessRWMutex.Unlock()
}
