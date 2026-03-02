package ws

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/olahol/melody"
)

type Hub struct {
	melody *melody.Melody
	mu     sync.RWMutex
	subs   map[string][]*melody.Session
}

type Notification struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func NewHub() *Hub {
	h := &Hub{
		melody: melody.New(),
		subs:   make(map[string][]*melody.Session),
	}

	h.melody.HandleConnect(func(s *melody.Session) {
		// connected
	})

	h.melody.HandleMessage(func(s *melody.Session, msg []byte) {

		var orderID string
		json.Unmarshal(msg, &orderID)

		h.mu.Lock()
		h.subs[orderID] = append(h.subs[orderID], s)
		h.mu.Unlock()
	})

	return h
}

func (h *Hub) HandleRequest(w http.ResponseWriter, r *http.Request) {
	h.melody.HandleRequest(w, r)
}

func (h *Hub) Notify(orderID, status string) {

	h.mu.RLock()
	sessions := h.subs[orderID]
	h.mu.RUnlock()

	notif := Notification{
		OrderID: orderID,
		Status:  status,
	}

	data, _ := json.Marshal(notif)

	for _, s := range sessions {
		s.Write(data)
	}
}
