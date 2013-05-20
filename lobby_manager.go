package golem

type lobbyRequest struct {
	name string
	conn *Connection
}

type lobbyManager struct {
	lobbies    map[string]*lobby
	register   chan *lobbyRequest
	unregister chan *lobbyRequest
	remove     chan string
}

func newLobbyManager() *lobbyManager {
	return &lobbyManager{
		lobbies:    make(map[string]*lobby),
		register:   make(chan *lobbyRequest),
		unregister: make(chan *lobbyRequest),
		remove:     make(chan string),
	}
}

func (lm *lobbyManager) run() {
	for {
		select {
		case req := <-lm.register:
			l, ok := lm.lobbies[req.name]
			if !ok {
				l := newLobby(lm, req.name)
				lm.lobbies[req.name] = l
				l.subscribe <- req.conn
				go l.run()
			} else {
				l.subscribe <- req.conn
			}
			req.conn.lobbies.add <- l
		case req := <-lm.unregister:
			l, ok := lm.lobbies[req.name]
			if ok {
				l.unsubscribe <- req.conn
				req.conn.lobbies.remove <- l
			}
		case ln := <-lm.remove:
			delete(lm.lobbies, ln)
		}
	}
}
