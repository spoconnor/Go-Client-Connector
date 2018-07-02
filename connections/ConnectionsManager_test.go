package connections

import (
	"log"
	"net"
	"pharos/client-connector/gopool"
	"testing"
	"time"
)

//--------------------------------------------------------
// Implements net.Addr interface
type Addr struct {
}

func (a Addr) Network() string {
	return "tcp"
}
func (a Addr) String() string {
	return "127.0.0.1"
}

//--------------------------------------------------------
// Implements net.Conn interface
type TestConn struct {
}

// Read reads data from the connection.
func (c *TestConn) Read(b []byte) (n int, err error) {
	log.Println("[TestConn.Read]")
	return 0, nil
}

// Write writes data to the connection.
func (c *TestConn) Write(b []byte) (n int, err error) {
	log.Println("[TestConn.Write]")
	return 0, nil
}

// Close closes the connection.
func (c *TestConn) Close() error {
	log.Println("[TestConn.Close]")
	return nil
}

// LocalAddr returns the local network address.
func (c *TestConn) LocalAddr() net.Addr {
	log.Println("[TestConn.LocalAddr]")
	return Addr{}
}

// RemoteAddr returns the remote network address.
func (c *TestConn) RemoteAddr() net.Addr {
	log.Println("[TestConn.RemoteAddr]")
	return Addr{}
}

func (c *TestConn) SetDeadline(t time.Time) error {
	log.Println("[TestConn.ReadDeadline]")
	return nil
}
func (c *TestConn) SetReadDeadline(t time.Time) error {
	log.Println("[TestConn.SetReadDeadline]")
	return nil
}
func (c *TestConn) SetWriteDeadline(t time.Time) error {
	log.Println("[TestConn.SetWriteDeadline]")
	return nil
}

//--------------------------------------------------------

func TestSetKey(t *testing.T) {
	workers, queue := 2, 1
	pool := gopool.NewPool(workers, queue, 1)
	conns := NewConnectionsManager(pool)

	conn := &TestConn{}
	connection := conns.Register(conn)

	if !conns.HaveConnectionKey(connection.Key) {
		t.Errorf("Do not have key %v", connection.Key)
	}

	conns.SetConnectionKey(connection, "NewKey")

	if !conns.HaveConnectionKey("NewKey") {
		t.Errorf("Do not have key %v", "NewKey")
	}
}

//func NewConnectionsManager(pool *gopool.Pool) *ConnectionsManager {

// Register registers new connection as a Connection.
//func (c *ConnectionsManager) Register(conn net.Conn) *Connection {

//func (c *ConnectionsManager) SetConnectionKey(connection *Connection, key string) {

// Remove removes connection from connectionsManager
//func (c *ConnectionsManager) Remove(connection *Connection) {

//func (c *ConnectionsManager) Broadcast(method string, params contracts.RpcParams) error {
//func (c *ConnectionsManager) SendToClient(key string, method string, params contracts.RpcParams) error {
