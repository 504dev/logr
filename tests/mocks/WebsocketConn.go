package mocks

type WebsocketConn struct {
	journal chan string
}

func NewWebsocketConn() *WebsocketConn {
	return &WebsocketConn{
		journal: make(chan string, 100),
	}
}

func (w WebsocketConn) Drain() []string {
	var result []string
	for {
		select {
		case v := <-w.journal:
			result = append(result, v)
		default:
			return result
		}
	}
}

func (w WebsocketConn) Write(p []byte) (int, error) {
	w.journal <- string(p)
	return 0, nil
}

func (w WebsocketConn) Close() error {
	close(w.journal)
	return nil
}
