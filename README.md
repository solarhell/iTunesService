# iTunesService

## usage

```go
package main

import (
	"context"
	"errors"
	"github.com/smallnest/rpcx/client"
    "github.com/solarhell/iTunesService/public"
	"log"
)

func main()  {
	d := client.NewPeer2PeerDiscovery("tcp@localhost:8972", "")
	xclient := client.NewXClient("iTunes", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &public.Args{
		Name: "李志",
	}

	reply := &public.Reply{}
	err := xclient.Call(context.Background(), "GetArtistPictureImageUrl", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Println(reply.URL)
}

```
