package connections

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"github.com/spoconnor/Go-Client-Connector/contracts"
	"strings"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

// Connection represents user connection.
// It contains logic of receiving and sending messages.
// That is, there are no active reader or writer. Some other layer of the
// application should call Receive() to read user's incoming message.
type Connection struct {
	io   sync.Mutex
	conn io.ReadWriteCloser

	id                 uint
	name               string              // TODO - remove
	connectionsManager *ConnectionsManager // TODO - remove?

	ConnectionId string //Guid
	Key          string
	DateTimeUtc  time.Time

	Properties map[string]string

	awaitingReply map[int](chan *contracts.RpcResponse)
}

func (u *Connection) SendKeyRequest() error {
	log.Printf("Sending Key request")
	var buf bytes.Buffer
	w := wsutil.NewWriter(&buf, ws.StateServerSide, ws.OpBinary)
	w.Write([]byte(KeyPlease))
	if err := w.Flush(); err != nil {
		log.Printf("[SendKeyRequest] Error: %s", err)
		u.conn.Close()
		return err
	}
	log.Printf("[SendKeyRequest] Writing key request...")
	return u.writeRaw(buf.Bytes())
}

// Receive reads next message from user's underlying connection.
// It blocks until full message received.
func (u *Connection) Receive() error {
	res, err := u.readResponse()
	if err != nil {
		log.Printf("[Receive] Error: %s", err)
		u.conn.Close()
		return err
	}
	if res == nil {
		// Handled some control message.
		return nil
	}
	log.Printf("[Connection.Receive] Received response %v", res.Result)
	/*
		switch req.Method {
		case "rename":
			name, ok := req.Params["name"].(string)
			if !ok {
				return u.writeErrorTo(req, websockets.Params{
					"error": "bad params",
				})
			}
			prev, ok := u.ConnectionsManager.Rename(u, name)
			if !ok {
				return u.writeErrorTo(req, websockets.Params{
					"error": "already exists",
				})
			}
			u.ConnectionsManager.Broadcast("rename", websockets.Params{
				"prev": prev,
				"name": name,
				"time": timestamp(),
			})
			return u.writeResultTo(req, nil)
		case "publish":
			req.Params["author"] = u.name
			req.Params["time"] = timestamp()
			u.ConnectionsManager.Broadcast("publish", req.Params)
		default:
			return u.writeErrorTo(req, websockets.Params{
				"error": "not implemented",
			})
		}
	*/
	return nil
}

//------------------------------------------------------------------

// readResponse reads json-rpc response from connection.
func (u *Connection) readResponse() (*contracts.RpcResponse, error) {
	u.io.Lock()
	defer u.io.Unlock()

	h, r, err := wsutil.NextReader(u.conn, ws.StateServerSide)
	if err != nil {
		return nil, err
	}
	if h.OpCode.IsControl() {
		return nil, wsutil.ControlHandler(u.conn, ws.StateServerSide)(h, r)
	}

	if h.OpCode == ws.OpBinary {

		log.Printf("Received binary")
		reader := bufio.NewReader(r)
		message, err := reader.ReadString('\n')
		log.Printf("Received %s", message)
		challenge, err := reader.ReadString('\n')
		log.Printf("Received %s", challenge)

		key := strings.TrimSpace(strings.TrimPrefix(message, KeyPrefix))
		u.connectionsManager.SetConnectionKey(u, key)

		//var toks = Encoding.UTF8.GetString(e.RawData).Split('\n');
		//var message = toks.FirstOrDefault();
		//var challenge = toks.Length > 1 ? toks[1] : null;

		//if (message == Constants.KeyPlease)
		//{
		//sessionId = Constants.KeyPrefix + Key;
		//auth = AuthorizationProvider?.GetAuthorization(challenge ?? sessionId);
		//routingInfo = PrepareRoutingInfo();
		//reply = string.Join("\n", new string[] { sessionId, auth, routingInfo });

		//		webSocketClient.Post(Encoding.UTF8.GetBytes(reply));
		//		Logger.Info("[WebSocketExchangeClient.webSocketClient_OnMessage] Connected to cloud with routing key: " + RoutingKey);
		if err == io.EOF {
			err = nil
		}
		return nil, err
	}
	if h.OpCode == ws.OpText {
		log.Printf("Received text")
		res := &contracts.RpcResponse{}
		decoder := json.NewDecoder(r)
		if err := decoder.Decode(res); err != nil {
			log.Printf("[Receive] Json decode response error: %s", err)
			return nil, err
		}

		c, has := u.awaitingReply[res.ID]
		if has {
			log.Println("[Connection.readReasponse] notifying of response")

			c <- res
		}
		delete(u.awaitingReply, res.ID)
		return res, nil
	}

	log.Printf("[Receive] Unhandled OpCode received %d", h.OpCode)
	return nil, nil // TODO - error
}

func (u *Connection) AwaitReply(id int) (*contracts.RpcResponse, error) {
	log.Println("[Connection.AwaitReply] waiting")
	c := make(chan *contracts.RpcResponse)
	u.awaitingReply[id] = c
	res := <-c
	log.Printf("[Connection.AwaitReply] returning response %v", res)
	return res, nil
}

func (u *Connection) writeErrorTo(req *contracts.RpcRequest, err contracts.RpcParams) error {
	log.Println("[writeError]")
	return u.write(contracts.RpcError{
		ID:    req.ID,
		Error: err,
	})
}

func (u *Connection) writeResultTo(req *contracts.RpcRequest, result contracts.RpcParams) error {
	log.Println("[writeResultTo]")
	return u.write(contracts.RpcResponse{
		ID:     req.ID,
		Result: result,
	})
}

func (u *Connection) writeNotice(method string, params contracts.RpcParams) error {
	log.Println("[writeNotice]")
	return u.write(contracts.RpcRequest{
		Method: method,
		Params: params,
	})
}

func (u *Connection) write(x interface{}) error {
	log.Println("[write] Writing interface")
	w := wsutil.NewWriter(u.conn, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(w)

	u.io.Lock()
	defer u.io.Unlock()

	if err := encoder.Encode(x); err != nil {
		return err
	}

	return w.Flush()
}

func (u *Connection) writeRaw(p []byte) error {
	log.Printf("[writeRaw] Writing %d raw bytes", len(p))
	u.io.Lock()
	defer u.io.Unlock()

	count, err := u.conn.Write(p)
	log.Printf("[writeRaw] Wrote %d bytes", count)

	return err
}
