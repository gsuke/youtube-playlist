package main

import (
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
	"google.golang.org/api/youtube/v3"
)

type Video struct {
	VideoId      string
	VideoTitle   string
	ChannelId    string
	ChannelTitle string
	IsPublic     bool
}

func main() {

	const playlistId = ""

	client := getClient(youtube.YoutubeReadonlyScope)

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	response, err := service.PlaylistItems.List([]string{"snippet,status"}).
		PlaylistId(playlistId).
		MaxResults(3).
		Do()
	handleError(err, "")

	videos := []Video{}
	for _, item := range response.Items {
		videos = append(videos, Video{
			VideoId:      item.Snippet.ResourceId.VideoId,
			VideoTitle:   item.Snippet.Title,
			ChannelId:    item.Snippet.VideoOwnerChannelId,
			ChannelTitle: item.Snippet.VideoOwnerChannelTitle,
			IsPublic:     item.Status.PrivacyStatus == "public",
		})
	}
	df := dataframe.LoadStructs(videos)

	if err := df.WriteCSV(os.Stdout); err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}
}
