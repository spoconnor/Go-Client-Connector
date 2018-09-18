//go:generate swagger generate spec

package servers

import (
	"log"
	"net/http"

	"github.com/spoconnor/Go-Client-Connector/connections"

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
