package nats_package

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsHandler struct {
	Conn *nats.Conn
}

func (h *NatsHandler) OnHelloWorld(msg *nats.Msg) {
	fmt.Println("==DATA : ", string(msg.Data))
	resp := new(Data)
	err := json.Unmarshal(msg.Data, resp)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}
	log.Println("Received message:", resp.Data)
	msg.Respond([]byte("Hello " + resp.Data))

	if err := h.PublishOnDemand(pubSubject); err != nil {
		log.Printf("nitesh publishing message: %v", err)
		return
	}
}

func (h *NatsHandler) PublishOnDemand(subject string) error {
	var i int = 0
	for {
		time.Sleep(3 * time.Second)
		data := Data{Data: fmt.Sprintf("OnDemand, %d", i)}
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		if err := h.Conn.Publish(subject, dataBytes); err != nil {
			log.Printf("Error publishing message: %v", err)
			return err
		}

		i += 1
		log.Printf("Published message: %s", data.Data)
	}
}

func NewNatsHandler(conn *nats.Conn) *NatsHandler {
	return &NatsHandler{
		Conn: conn,
	}
}
