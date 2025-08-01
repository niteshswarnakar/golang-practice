package nats_package

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type MsgQueue chan *nats.Msg

// var subQueue = make(MsgQueue, 100)

type NatsHandler struct {
	Conn          *nats.Conn
	subscriptions map[string]*nats.Subscription
	mu            sync.RWMutex
}

func (h *NatsHandler) OnHelloWorld(msg *nats.Msg) {
	fmt.Println("\n$Started")
	startTime := time.Now()
	resp := new(Data)
	err := json.Unmarshal(msg.Data, resp)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
	}
	log.Println("Received message:", resp.Data)
	time.Sleep(10 * time.Second)

	log.Printf("Responded to message on subject: %s", msg.Subject)
	endTime := time.Now()
	fmt.Println("Elapsed time: ", endTime.Sub(startTime))

	err = msg.Respond([]byte("\nHello User" + endTime.Sub(startTime).String()))
	if err != nil {
		log.Printf("Error responding to message: %v", err)
	}
	// if err := h.PublishOnDemand(pubSubject); err != nil {
	// 	log.Printf("nitesh publishing message: %v", err)
	// 	return
	// }
}

func (h *NatsHandler) PublishOnDemand(subject string) error {
	var i int = 0
	for {
		time.Sleep(5 * time.Second)
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

func (h *NatsHandler) SubscribeToSubject(subject string, handler nats.MsgHandler) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	workerCount := 10
	msgChan := make(chan *nats.Msg, workerCount*2)

	// sub, err := h.Conn.Subscribe(subject, handler)
	sub, err := h.Conn.ChanSubscribe(subject, msgChan)
	if err != nil {
		return err
	}

	h.subscriptions[subject] = sub
	log.Printf("Subscribed to subject: %s", subject)
	return nil
}

func (h *NatsHandler) GetSubscribedSubjects() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	subjects := make([]string, 0, len(h.subscriptions))
	for subject := range h.subscriptions {
		subjects = append(subjects, subject)
	}
	return subjects
}

func (h *NatsHandler) PrintSubscribedSubjects() {
	subjects := h.GetSubscribedSubjects()
	fmt.Println("Currently subscribed subjects:")
	for i, subject := range subjects {
		fmt.Printf("%d. %s\n", i+1, subject)
	}
}

func (h *NatsHandler) UnsubscribeFromSubject(subject string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if sub, exists := h.subscriptions[subject]; exists {
		err := sub.Unsubscribe()
		if err != nil {
			return err
		}
		delete(h.subscriptions, subject)
		log.Printf("Unsubscribed from subject: %s", subject)
	}
	return nil
}

func NewNatsHandler(conn *nats.Conn) *NatsHandler {
	return &NatsHandler{
		Conn:          conn,
		subscriptions: make(map[string]*nats.Subscription),
	}
}
