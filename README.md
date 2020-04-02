# iTunesService

## usage

```go
package main

import (
	"context"
	"log"
	"time"

	pb "github.com/solarhell/iTunesService/applemusic"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		"103.121.209.132:50051",
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		panic(err)
	}

	appleMusic := pb.NewMusicClient(conn)
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	reply, err := appleMusic.GetArtistPicture(ctx, &pb.CheckRequest{ArtistName: "李志"})
	if err != nil {
		panic(err)
	}

	log.Printf("获取到歌手照片: %s\n", reply.Picture)
}
```
