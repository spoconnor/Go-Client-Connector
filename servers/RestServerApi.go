//go:generate swagger generate spec

package servers

import (
	"encoding/json"
	"log"

	contracts "github.com/spoconnor/Go-Client-Connector/contracts"

	"github.com/go-ozzo/ozzo-routing"
)

//-------------------------------------------------

func (r *RestServer) hello(c *routing.Context) error {
	log.Println("[RestServer.hello]")
	c.Write("Hello")
	return nil
}

func (r *RestServer) listConnections(c *routing.Context) error {
	log.Println("[RestServer.listConnections]")
	var res = r.connectionsManager.ListConnections()
	c.Write(res)
	return nil
}

func (r *RestServer) ping(c *routing.Context) error {
	log.Println("[RestServer.ping]")
	key := c.Param("key")
	//message := c.Query("message")
	res, err := r.connectionsManager.SendToClient(key, "Ping", nil, true)
	if err != nil {
		c.Write(res)
	}
	return err
}

func (r *RestServer) jsonRpc(c *routing.Context) error {
	log.Println("[RestServer.jsonRpc]")
	key := c.Param("key")
	var req contracts.RpcRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("[RestServer.jsonRpc] received '%s' for '%s'", req.Method, key)

	res, err := r.connectionsManager.SendToClient(key, req.Method, req.Params, true)
	if err != nil {
		log.Printf("[RestServer.jsonRpc] Error '%v'", err)
	} else {
		json, _ := json.Marshal(res.Result)
		log.Printf("[RestServer.jsonRpc] sending '%s'", json)
		c.Write(string(json))
	}
	return err
}
