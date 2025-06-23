package nats_package

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats-server/server"
	"github.com/nats-io/nats.go"
)

type Data struct {
	Data string `json:"data"`
}

func NatsIO() {
	nc, ns, err := RunEmbeddedServer()
	if err != nil {
		log.Fatal(err)
	}

	natsHandler := NewNatsHandler(nc)

	defer nc.Close()

	nc.Subscribe(subSubject, natsHandler.OnHelloWorld)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Shutdown the server
	ns.Shutdown()
}

func RunEmbeddedServer() (*nats.Conn, *server.Server, error) {
	ns := server.New(&server.Options{
		Debug:   false,
		Trace:   false,
		Logtime: false,
	})

	go ns.Start()

	if !ns.ReadyForConnections(5 * time.Second) {
		return nil, nil, fmt.Errorf("NATS server timeout")
	}

	ns.ConfigureLogger()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, nil, fmt.Errorf("NATS connection error: %v", err)
	}

	fmt.Println("NATS server is running at", ns.Addr())
	return nc, ns, nil
}
