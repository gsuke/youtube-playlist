package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"youtubeapi/oauth2"

	"github.com/go-gota/gota/dataframe"
	"google.golang.org/api/option"
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

	client := oauth2.GetClient(youtube.YoutubeReadonlyScope)

	service, err := youtube.NewService(context.Background(), option.WithHTTPClient(client))
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

	file, err := os.Create("playlist_items.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	if err := df.WriteCSV(file); err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}
