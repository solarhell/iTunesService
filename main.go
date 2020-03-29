package main

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/imroc/req"
	"github.com/smallnest/rpcx/server"
	"github.com/solarhell/iTunesService/public"
)

func GetArtistPictureImageUrl(ctx context.Context, args *public.Args, reply *public.Reply) error {
	if args.Name == "" {
		return errors.New("歌手名不能为空")
	}
	u, err := url.Parse("http://ax.itunes.apple.com/WebObjects/MZStoreServices.woa/wa/wsSearch")
	if err != nil {
		return err
	}

	q := u.Query()
	q.Add("country", "hk")
	q.Add("entity", "allArtist")
	q.Add("term", args.Name)
	q.Add("limit", "1")

	u.RawQuery = q.Encode()

	r, err := req.Get(u.String())
	if err != nil {
		return err
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
		return err
	}

	if len(artist.Results) == 0 {
		return err
	}

	r, err = req.Get(artist.Results[0].ArtistLinkURL)
	if err != nil {
		return err
	}

	str, err := r.ToString()
	if err != nil {
		return err
	}

	doc, err := htmlquery.Parse(strings.NewReader(str))
	if err != nil {
		return err
	}

	imageEle := htmlquery.FindOne(doc, "//meta[@property='og:image']/@content")
	imageStr := htmlquery.InnerText(imageEle)
	imageProcessed := strings.Replace(imageStr, "1200x630cw.png", "99999x99999.png", 1)
	reply.URL = imageProcessed
	return nil
}

func main() {
	s := server.NewServer()
	var err error
	err = s.RegisterFunctionName("iTunes", "GetArtistPictureImageUrl", GetArtistPictureImageUrl, "")
	if err != nil {
		panic(err)
	}
	err = s.Serve("tcp", ":8972")
	if err != nil {
		panic(err)
	}
}
