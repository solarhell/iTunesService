package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/antchfx/htmlquery"
	"github.com/imroc/req"
	pb "github.com/solarhell/iTunesService/proto"
	"google.golang.org/grpc"
)

type checkServer struct{}

func (*checkServer) GetArtistPicture(_ context.Context, input *pb.CheckRequest) (*pb.CheckReply, error) {
	name := input.ArtistName

	if name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "歌手名不能为空")
	}

	u, err := url.Parse("http://ax.itunes.apple.com/WebObjects/MZStoreServices.woa/wa/wsSearch")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("country", "hk")
	q.Add("entity", "allArtist")
	q.Add("term", name)
	q.Add("limit", "1")

	u.RawQuery = q.Encode()

	r, err := req.Get(u.String())
	if err != nil {
		return nil, err
	}

	type artistResponse struct {
		ResultCount int `json:"resultCount"`
		Results     []struct {
			WrapperType      string `json:"wrapperType"`
			ArtistType       string `json:"artistType"`
			ArtistName       string `json:"artistName"`
			ArtistLinkURL    string `json:"artistLinkUrl"`
			ArtistID         int    `json:"artistId"`
			PrimaryGenreName string `json:"primaryGenreName"`
			PrimaryGenreID   int    `json:"primaryGenreId"`
		} `json:"results"`
	}

	artist := artistResponse{}
	err = r.ToJSON(&artist)
	if err != nil {
		return nil, err
	}

	if len(artist.Results) == 0 {
		return nil, status.Errorf(codes.NotFound, "没有查找到该歌手")
	}

	r, err = req.Get(artist.Results[0].ArtistLinkURL)
	if err != nil {
		return nil, err
	}

	str, err := r.ToString()
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(str))
	if err != nil {
		return nil, err
	}

	imageEle := htmlquery.FindOne(doc, "//meta[@property='og:image']/@content")
	imageStr := htmlquery.InnerText(imageEle)
	imageProcessed := strings.Replace(imageStr, "1200x630cw.png", "99999x99999.png", 1)
	return &pb.CheckReply{Picture: imageProcessed}, nil
}

const GRPC_PORT = 50051

func main() {
	grpcServer := grpc.NewServer()

	pb.RegisterMusicServer(grpcServer, &checkServer{})
	reflection.Register(grpcServer)

	go func() {
		log.Printf("grpc server is on 127.0.0.1:%d\n", GRPC_PORT)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", GRPC_PORT))
		if err != nil {
			log.Fatalf("start grpc server error: %v", err)
		}
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("start grpc server error: %v", err)
		}
	}()

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("exiting:", <-termChan)
	grpcServer.GracefulStop()
	log.Println("exited")
}
