package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-pubsub"
)

const ChatRoomBufSize = 128

type ChatRoom struct {
	Messages chan *ChatMessage

	ctx   context.Context
	ps    *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription

	roomName string
	self     peer.ID
	nick     string
}

type ChatMessage struct {
	Message    string
	SenderID   peer.ID
	SenderNick string
}

func JoinChatRoom(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID, nickname string, roomName string) (*ChatRoom, error) {
	// join the pubsub topic
	topic, err := ps.Join(roomName)
	if err != nil {
		return nil, err
	}

	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	cr := &ChatRoom{
		ctx:      ctx,
		ps:       ps,
		topic:    topic,
		sub:      sub,
		self:     selfID,
		nick:     nickname,
		roomName: roomName,
		Messages: make(chan *ChatMessage, ChatRoomBufSize),
	}

	// start reading messages from the subscription in a loop
	go cr.readLoop()
	return cr, nil
}

func (cr *ChatRoom) Publish(message string) error {
	m := ChatMessage{
		Message:    message,
		SenderID:   cr.self,
		SenderNick: cr.nick,
	}
	msg, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return cr.topic.Publish(cr.ctx, msg)
}

func (cr *ChatRoom) ListPeers() []peer.ID {
	return cr.ps.ListPeers(cr.roomName)
}

func (cr *ChatRoom) readLoop() {
	for {
		msg, err := cr.sub.Next(cr.ctx)
		if err != nil {
			close(cr.Messages)
			return
		}
		// only forward messages delivered by others
		if msg.ReceivedFrom == cr.self {
			continue
		}
		cm := new(ChatMessage)
		err = json.Unmarshal(msg.Data, cm)
		if err != nil {
			continue
		}
		fmt.Printf("\x1b[32m%s(ID:%s): \x1B[0;35m%s> ",
			cm.SenderNick, cm.SenderID.String()[:6], cm.Message)
		// send valid messages onto the Messages channel
		if len(cr.Messages) < 128 {
			cr.Messages <- cm
		} else {
			<-cr.Messages
			cr.Messages <- cm
		}
	}
}
