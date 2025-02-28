//go:build !solution

package pubsub

import (
	"context"
	"errors"
	"sync"
)

var _ Subscription = (*MySubscription)(nil)
var _ PubSub = (*MyPubSub)(nil)

type MySubscription struct {
	pubSub         *MyPubSub
	obj            string
	handler        MsgHandler
	messageChannel chan interface{}
	waitGroup      *sync.WaitGroup
}

type MyPubSub struct {
	subscribers map[string][]*MySubscription
	mutex       sync.RWMutex
	waitGroup   sync.WaitGroup
	closed      bool
}

func NewPubSub() PubSub {
	return &MyPubSub{
		subscribers: make(map[string][]*MySubscription),
	}
}

func (pubSub *MyPubSub) Subscribe(obj string, handler MsgHandler) (Subscription, error) {
	if err := pubSub.checkIfClosed(); err != nil {
		return nil, err
	}
	pubSub.mutex.Lock()
	defer pubSub.mutex.Unlock()

	sub := pubSub.createAndStartSubscription(obj, handler)
	pubSub.addSubscription(obj, sub)

	return sub, nil
}

func (pubSub *MyPubSub) Publish(obj string, message interface{}) error {
	if err := pubSub.checkIfClosed(); err != nil {
		return err
	}
	pubSub.mutex.RLock()
	defer pubSub.mutex.RUnlock()

	pubSub.sendToSubscribers(obj, message)
	return nil
}

func (pubSub *MyPubSub) Close(context context.Context) error {
	if err := pubSub.markAsClosed(); err != nil {
		return err
	}

	pubSub.closeAllSubscriptionChannels()

	return pubSub.waitForCompletion(context)
}

func (pubSub *MyPubSub) checkIfClosed() error {
	if pubSub.closed {
		return errors.New("PubSub is closed")
	}
	return nil
}

func (pubSub *MyPubSub) markAsClosed() error {
	pubSub.mutex.Lock()
	defer pubSub.mutex.Unlock()

	if pubSub.closed {
		return errors.New("PubSub already closed")
	}
	pubSub.closed = true
	return nil
}

func (pubSub *MyPubSub) createAndStartSubscription(obj string, handler MsgHandler) *MySubscription {
	sub := pubSub.createSubscription(obj, handler)
	sub.startListening()
	return sub
}

func (pubSub *MyPubSub) createSubscription(obj string, handler MsgHandler) *MySubscription {
	return &MySubscription{
		pubSub:         pubSub,
		obj:            obj,
		handler:        handler,
		messageChannel: make(chan interface{}, 100),
		waitGroup:      &pubSub.waitGroup,
	}
}

func (pubSub *MyPubSub) addSubscription(obj string, mySubscription *MySubscription) {
	pubSub.subscribers[obj] = append(pubSub.subscribers[obj], mySubscription)
}

func (pubSub *MyPubSub) sendToSubscribers(obj string, message interface{}) {
	for _, sub := range pubSub.getSubscribers(obj) {
		pubSub.sendMessageToSubscriber(sub, message)
	}
}

func (pubSub *MyPubSub) getSubscribers(obj string) []*MySubscription {
	if subs, found := pubSub.subscribers[obj]; found {
		return subs
	}
	return nil
}

func (pubSub *MyPubSub) sendMessageToSubscriber(sub *MySubscription, message interface{}) {
	pubSub.waitGroup.Add(1)
	sub.messageChannel <- message
}

func (pubSub *MyPubSub) closeAllSubscriptionChannels() {
	for _, subs := range pubSub.subscribers {
		for _, sub := range subs {
			close(sub.messageChannel)
		}
	}
}

func (pubSub *MyPubSub) waitForCompletion(context context.Context) error {
	doneCh := make(chan struct{})
	go pubSub.waitGroupCompletion(doneCh)

	select {
	case <-doneCh:
		return nil
	case <-context.Done():
		return context.Err()
	}
}

func (pubSub *MyPubSub) waitGroupCompletion(doneChannel chan struct{}) {
	pubSub.waitGroup.Wait()
	close(doneChannel)
}

func (pubSub *MyPubSub) removeSubscription(obj string, mySub *MySubscription) {
	pubSub.mutex.Lock()
	defer pubSub.mutex.Unlock()

	if subs, found := pubSub.subscribers[obj]; found {
		for i, mySubscription := range subs {
			if mySubscription == mySub {
				pubSub.subscribers[obj] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
	}
}

func (mySubscription *MySubscription) Unsubscribe() {
	mySubscription.pubSub.removeSubscription(mySubscription.obj, mySubscription)
	close(mySubscription.messageChannel)
}

func (mySubscription *MySubscription) startListening() {
	go mySubscription.listenForMessages()
}

func (mySubscription *MySubscription) listenForMessages() {
	for msg := range mySubscription.messageChannel {
		mySubscription.handleMessage(msg)
	}
}

func (mySubscription *MySubscription) handleMessage(msg interface{}) {
	mySubscription.handler(msg)
	mySubscription.waitGroup.Done()
}
