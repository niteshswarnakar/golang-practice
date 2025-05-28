package nats_package

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats-server/server"
	"github.com/nats-io/nats.go"
)

func NatsIO() {
	nc, ns, err := RunEmbeddedServer(true)
	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	fmt.Println("Subscribing to NATS server...")
	nc.Subscribe("hello.world", func(msg *nats.Msg) {
		msg.Respond([]byte("Hello, this is data : " + string(msg.Data)))
		fmt.Println("Received message:", string(msg.Data))
	})

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Shutdown the server
	ns.Shutdown()
}

func RunEmbeddedServer(enableLogging bool) (*nats.Conn, *server.Server, error) {
	leafUrl, err := url.Parse("nats-leaf://connect.ngs.global")
	if err != nil {
		return nil, nil, fmt.Errorf("URL parse error: %v", err)
	}

	fmt.Println("Leaf URL:", leafUrl)

	ns := server.New(&server.Options{
		Debug:   enableLogging,
		Trace:   false,
		Logtime: enableLogging,
	})

	go ns.Start()

	if !ns.ReadyForConnections(5 * time.Second) {
		return nil, nil, fmt.Errorf("NATS server timeout")
	}

	if enableLogging {
		ns.ConfigureLogger()
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, nil, fmt.Errorf("NATS connection error: %v", err)
	}

	return nc, ns, nil
}
