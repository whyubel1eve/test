package Legs

import (
	"bytes"
	"fmt"
	blocks "github.com/ipfs/go-block-format"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"io"
)

// MkLinkSystem Make a linksystem with blockstore
func MkLinkSystem(bs blockstore.Blockstore) ipld.LinkSystem {
	ls := cidlink.DefaultLinkSystem()
	ls.StorageReadOpener = func(lnkCtx ipld.LinkContext, lnk ipld.Link) (io.Reader, error) {
		asCidLink, ok := lnk.(cidlink.Link)
		if !ok {
			return nil, fmt.Errorf("unsupported link types")
		}
		block, err := bs.Get(lnkCtx.Ctx, asCidLink.Cid)
		if err != nil {
			return nil, err
		}
		return bytes.NewBuffer(block.RawData()), nil
	}
	ls.StorageWriteOpener = func(linkContext linking.LinkContext) (io.Writer, linking.BlockWriteCommitter, error) {
		buf := bytes.NewBuffer(nil)
		return buf, func(lnk ipld.Link) error {
			c := lnk.(cidlink.Link).Cid
			buffer := buf.Bytes()
			block, err := blocks.NewBlockWithCid(buffer, c)
			if err != nil {
				return err
			}
			err = bs.Put(linkContext.Ctx, block)
			if err != nil {
				return err
			}
			return nil
		}, nil
	}
	return ls
}
