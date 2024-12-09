package weatherfeedreceiver

import (
	"log"
)

type IFeedReceiver interface {
	ReceiveMessages() error
}

type FeedReceiver struct {
	receiver IFeedReceiver
}

func NewFeedReceiver(receiver IFeedReceiver) *FeedReceiver {
	return &FeedReceiver{receiver}
}

func (r FeedReceiver) HandleReceiveMessages() {
	if err := r.receiver.ReceiveMessages(); err != nil {
		log.Printf("ERROR: Failed to receive weather feed messages due to following error: %v", err)
	}
}
