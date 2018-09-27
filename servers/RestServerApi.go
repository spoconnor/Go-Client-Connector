//go:generate swagger generate spec

package servers

import (
	"encoding/json"
	"log"

	contracts "github.com/spoconnor/Go-Client-Connector/contracts"

	"github.com/go-ozzo/ozzo-routing"
)

//-------------------------------------------------

// @Title hello
// @Description Test Rest service is running
// @Success 200 {string} Reponse message
// @Router /hello [get]
func (r *RestServer) hello(c *routing.Context) error {
	log.Println("[RestServer.hello]")
	c.Write("Hello")
	return nil
}

// @Title listConnections
// @Description List all client connections
// @Success 200 {array} Client connections
// @Router /listConnections [get]
func (r *RestServer) listConnections(c *routing.Context) error {
	log.Println("[RestServer.listConnections]")
	var res = r.connectionsManager.ListConnections()
	c.Write(res)
	return nil
}

// @Title ping
// @Description Test connection to a specified client
// @Success 200 {string} Reponse message
// @Router /ping/key/{key} [get]
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

// @Title jsonRpc
// @Description Send a json rpc message to a specified client
// @Accept json
// @Param key path string true "Client Id"
// @Param req body contracts.RpcRequest true "Rpc request"
// @Success 200 {string} Reponse message
// @Router /jsonRpc/key/{key} [get]
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
