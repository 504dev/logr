package sockmap

type SockMessage struct {
	Action  string      `json:"action,omitempty"`
	Path    string      `json:"path"`
	Payload interface{} `json:"payload"`
}
