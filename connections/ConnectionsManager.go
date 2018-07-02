package connections

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"sort"
	"sync"
	"time"

	"pharos/client-connector/contracts"
	"pharos/client-connector/gopool"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

const (
	Forbidden = "Forbidden"
	KeyPlease = "RoutingKeyPlease"
	KeyPrefix = "RoutingKey:"
)

// ConnectionsManager contains logic of connection interaction.
type ConnectionsManager struct {
	mu  sync.RWMutex
	seq uint
	us  []*Connection
	ns  map[string]*Connection

	pool *gopool.Pool
	out  chan []byte
}

// TODO - make parameter an interface type for testing
func NewConnectionsManager(pool *gopool.Pool) *ConnectionsManager {
	log.Printf("[NewConnectionsManager] Creating ConnectionsManager")
	connections := &ConnectionsManager{
		pool: pool,
		ns:   make(map[string]*Connection),
		out:  make(chan []byte, 1),
	}

	go connections.writer()

	return connections
}

// Register registers new connection as a Connection.
func (c *ConnectionsManager) Register(conn net.Conn) *Connection {
	connection := &Connection{
		connectionsManager: c,
		conn:               conn,
		awaitingReply:      make(map[int](chan *contracts.RpcResponse)),
	}

	connection.SendKeyRequest()
	//c.Broadcast("greet", websockets.Params{
	//	"name": connection.name,
	//	"time": timestamp(),
	//})

	return connection
}

func (c *ConnectionsManager) ListConnections() string {
	var res string
	for k, _ := range c.ns {
		res = res + k + ","
	}
	return res
}

func (c *ConnectionsManager) HaveConnectionKey(key string) bool {
	c.mu.Lock()
	_, has := c.ns[key]
	c.mu.Unlock()
	return has
}

func (c *ConnectionsManager) SetConnectionKey(connection *Connection, key string) {
	c.mu.Lock()
	{
		connection.id = c.seq
		connection.Key = key
		log.Printf("Adding connection %s", connection.Key)
		c.us = append(c.us, connection)
		c.ns[connection.Key] = connection

		c.seq++
	}
	c.mu.Unlock()

	/*
		c.mu.Lock()
		{
			var prev string
			if _, has := c.ns[key]; !has {
				prev, connection.Key = connection.Key, key
				log.Printf("Renamed connection %s to %s", prev, key)
				delete(c.ns, prev)
				c.ns[key] = connection
			}
		}
		c.mu.Unlock()
	*/
}

// Remove removes connection from connectionsManager
func (c *ConnectionsManager) Remove(connection *Connection) {
	c.mu.Lock()
	removed := c.remove(connection)
	c.mu.Unlock()

	if !removed {
		return
	}

	c.Broadcast("goodbye", contracts.RpcParams{
		"name": connection.name,
		"time": timestamp(),
	})
}

// Broadcast sends message to all alive connections.
func (c *ConnectionsManager) Broadcast(method string, params contracts.RpcParams) error {
	var buf bytes.Buffer

	w := wsutil.NewWriter(&buf, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(w)

	r := contracts.RpcRequest{Method: method, Params: params}
	if err := encoder.Encode(r); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	c.out <- buf.Bytes()

	return nil
}

func (c *ConnectionsManager) SendToClient(
	key string, method string,
	params contracts.RpcParams, waitForReply bool) (*contracts.RpcResponse, error) {

	log.Printf("[SendToClient] Sending %s to %s", method, key)
	var buf bytes.Buffer

	w := wsutil.NewWriter(&buf, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(w)

	r := contracts.RpcRequest{Method: method, Params: params}
	if err := encoder.Encode(r); err != nil {
		return nil, err
	}
	if err := w.Flush(); err != nil {
		return nil, err
	}
	log.Printf("[SendToClient] Sending raw: %s", buf.String())

	c.mu.Lock()
	defer c.mu.Unlock()
	if u, ok := c.ns[key]; ok {
		u := u              // for closure
		data := buf.Bytes() // for closure

		if waitForReply {
			err := u.writeRaw(data)
			if err != nil {
				log.Printf("[SendToClient] Error: %s", err)
			}
			res, err := u.AwaitReply(r.ID)
			return res, err

		} else {

			c.pool.Schedule(func() {
				log.Printf("[SendToClient] Writing...")
				err := u.writeRaw(data)
				if err != nil {
					log.Printf("[SendToClient] Error: %s", err)
				}
			})
		}

	} else {
		log.Printf("[SendToClient] '%s' not found", key)
		// TODO return error
	}

	return nil, nil // TODO
}

//-----------------------------------------------------------

// writer writes broadcast messages from out channel.
func (c *ConnectionsManager) writer() {
	for bts := range c.out {
		c.mu.RLock()
		us := c.us
		c.mu.RUnlock()

		for _, u := range us {
			u := u // For closure.
			c.pool.Schedule(func() {
				log.Printf("[ConnectionsManager.writer] Writing broadcast...")
				u.writeRaw(bts)
			})
		}
	}
}

// mutex must be held.
func (c *ConnectionsManager) remove(connection *Connection) bool {
	if _, has := c.ns[connection.name]; !has {
		return false
	}

	delete(c.ns, connection.name)

	i := sort.Search(len(c.us), func(i int) bool {
		return c.us[i].id >= connection.id
	})
	if i >= len(c.us) {
		panic("ConnectionsManager: inconsistent state")
	}

	without := make([]*Connection, len(c.us)-1)
	copy(without[:i], c.us[:i])
	copy(without[i:], c.us[i+1:])
	c.us = without

	return true
}

func timestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
