//go:generate swagger generate spec

package servers

import (
	"encoding/json"
	"log"
	"net/http"
	"go-client-connector/connections"
	contracts "go-client-connector/contracts"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
)

type RestServer struct {
	connectionsManager *connections.ConnectionsManager
}

func NewRestServer(c *connections.ConnectionsManager) *RestServer {
	r := &RestServer{
		connectionsManager: c,
	}
	return r
}

func (r *RestServer) Start() {
	log.Println("Starting Rest Server at http://localhost:9000/ClientConnector...")
	router := routing.New()

	router.Use(
		// all these handlers are shared by every route
		access.Logger(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
	)

	// serve RESTful APIs
	api := router.Group("/ClientConnector")
	api.Use(
		// these handlers are shared by the routes in the api group only
		content.TypeNegotiator(content.JSON),
	)

	api.Get("/Hello", r.hello)

	api.Get("/ListConnections", r.listConnections)

	// TODO customerid regex - note ` character
	api.Get(`/Key/<key>/Ping`, r.ping)
	api.Post("/Key/<key>/JsonRpc", r.jsonRpc)

	/*
		// serve index file
		router.Get("/", file.Content("ui/index.html"))
		// serve files under the "ui" subdirectory
		router.Get("/*", file.Server(file.PathMap{
			"/": "/ui/",
		}))
	*/

	http.Handle("/", router)
	http.ListenAndServe(":9000", nil)
	log.Println("Stopping Rest Server...")
}

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

	//var req json.RawMessage
	//if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
	//return err
	//}
	//cMap := make(map[string]string)
	//e := json.Unmarshal(req, &cMap)
	//if e != nil {
	//panic(e)
	//}

	var req contracts.RpcRequest
	//if err := json.Unmarshal([]byte(body), &req); err != nil {
	//if err := c.Read(&req); err != nil {
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
