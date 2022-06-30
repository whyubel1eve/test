package main

import (
	"context"
	"fmt"
	testinstance "github.com/ipfs/go-bitswap/testinstance"
	tn "github.com/ipfs/go-bitswap/testnet"
	"github.com/ipfs/go-cid"
	blocksutil "github.com/ipfs/go-ipfs-blocksutil"
	delay "github.com/ipfs/go-ipfs-delay"
	mockrouting "github.com/ipfs/go-ipfs-routing/mock"
	tu "github.com/libp2p/go-libp2p-testing/etc"
	"time"
)

const kNetworkDelay = 0 * time.Millisecond

func main() {
	net := tn.VirtualNetwork(mockrouting.NewServer(), delay.Fixed(kNetworkDelay))
	ig := testinstance.NewTestInstanceGenerator(net, nil, nil)
	defer ig.Close()
	bg := blocksutil.NewBlockGenerator()

	instances := ig.Instances(3)
	blocks := bg.Blocks(10)
	var keys []cid.Cid
	for _, b := range blocks {
		keys = append(keys, b.Cid())
	}

	// First peer has block
	err := instances[0].Exchange.HasBlock(context.Background(), blocks[0])
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Second peer broadcasts want for block CID
	// (Received by first and third peers)
	blk, err := instances[1].Exchange.GetBlock(ctx, blocks[0].Cid())
	if err != nil {
		panic(err)
	}
	fmt.Println("ins[2]'s wantlist for ins[1]: ", instances[2].Exchange.WantlistForPeer(instances[1].Peer))

	// When second peer receives block, it should send out a cancel, so third
	// peer should no longer keep second peer's want
	if err = tu.WaitFor(ctx, func() error {
		fmt.Println("ins[2]'s wantlist for ins[1]: ", instances[2].Exchange.WantlistForPeer(instances[1].Peer))
		fmt.Println("ins[1]'s wantlist: ", instances[1].Exchange.GetWantlist())
		return nil
	}); err != nil {
		panic(err)
	}

	fmt.Println("block data: ", blk.RawData())

	ctx = context.Background()
	_, err = instances[1].Exchange.GetBlocks(ctx, keys[1:10])
	if err != nil {
		panic(err)
	}
	fmt.Println("ins[1]'s wantHaves: ", instances[1].Exchange.GetWantHaves())

	for _, inst := range instances {
		err := inst.Exchange.Close()
		if err != nil {
			panic(err)
		}
	}
}
