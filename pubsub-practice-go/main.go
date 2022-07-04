package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"os"
	"pubsubTest/chat"
)

func main() {
	nickname := flag.String("nick", "", "nickname in chatRoom")
	roomName := flag.String("roomName", "awesome-chatRoom", "ChatRoom name")
	flag.Parse()
	
	ctx := context.Background()
	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
	)
	if err != nil {
		panic(err)
	}

	gossipSub, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}

	// setup local mDNS discovery
	if err := chat.SetupDiscovery(h); err != nil {
		panic(err)
	}

	chatRoom, err := chat.JoinChatRoom(ctx, gossipSub, h.ID(), *nickname, *roomName)
	if err != nil {
		panic(err)
	}

	stdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		err = chatRoom.Publish(sendData)
		if err != nil {
			panic(err)
		}
	}


}