package weatherfeedreceiver

import (
	"log"
	"sync"
)

type IFeedReceiver interface {
	ReceiveMessages(feedCounter *int, mu *sync.Mutex) error
}

type FeedReceiver struct {
	receiver IFeedReceiver
}

func NewFeedReceiver(receiver IFeedReceiver) *FeedReceiver {
	return &FeedReceiver{receiver}
}

func (r FeedReceiver) HandleReceiveMessages(feedCounter *int, mu *sync.Mutex) {
	if err := r.receiver.ReceiveMessages(feedCounter, mu); err != nil {
		log.Printf("ERROR: Failed to receive weather feed messages due to following error: %v", err)
	}
}
