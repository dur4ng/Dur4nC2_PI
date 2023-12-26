package state

import (
	"Dur4nC2/misc/protobuf/implantpb"
	"github.com/gofrs/uuid"
	"sync"
	"time"
)

// ImplantConnection - Abstract connection to an implant
type ImplantConnection struct {
	ID               string
	Send             chan *implantpb.Envelope
	RespMutex        *sync.RWMutex
	LastMessageMutex *sync.RWMutex
	Resp             map[int64]chan *implantpb.Envelope
	Transport        string
	RemoteAddress    string
	LastMessage      time.Time
	Cleanup          func()
}

// NewImplantConnection - Creates a new implant connection
func NewImplantConnection(transport string, remoteAddress string) *ImplantConnection {
	return &ImplantConnection{
		ID:               generateImplantConnectionID(),
		Send:             make(chan *implantpb.Envelope),
		RespMutex:        &sync.RWMutex{},
		LastMessageMutex: &sync.RWMutex{},
		Resp:             map[int64]chan *implantpb.Envelope{},
		Transport:        transport,
		RemoteAddress:    remoteAddress,
		Cleanup:          func() {},
	}
}

// GetLastMessage - Retrieves the last message time
func (c *ImplantConnection) GetLastMessage() time.Time {
	c.LastMessageMutex.RLock()
	defer c.LastMessageMutex.RUnlock()

	return c.LastMessage
}

// UpdateLastMessage - Updates the last message time
func (c *ImplantConnection) UpdateLastMessage() {
	c.LastMessageMutex.Lock()
	defer c.LastMessageMutex.Unlock()

	c.LastMessage = time.Now()
}

func generateImplantConnectionID() string {
	id, _ := uuid.NewV4()
	return id.String()
}
