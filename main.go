package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"time"
	//	_ "net/http/pprof"  // TODO
	"pharos/client-connector/servers"
)

var (
	addr = flag.String("listen", ":8080", "address to bind to")
	//debug     = flag.String("pprof", "", "address for pprof http")
	workers   = flag.Int("workers", 128, "max workers count")
	queue     = flag.Int("queue", 1, "workers task queue size")
	ioTimeout = flag.Duration("io_timeout", time.Millisecond*100, "i/o operations timeout")
)

func main() {
	SetupLog("client-connector", "info", "output.txt", true)
	log.Println("Starting client-connector")

	flag.Parse()

	ws := servers.NetWebSocketServer(*addr, *ioTimeout, *workers, *queue)
	go ws.Start()

	rs := servers.NewRestServer(ws.ConnectionsManager)
	go rs.Start()

	log.Println("Press any key to exit")
	reader := bufio.NewReader(os.Stdin)
	_, _, _ = reader.ReadRune()
	log.Println("Done")
}
