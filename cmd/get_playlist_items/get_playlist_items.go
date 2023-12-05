package main

import (
	"context"
	"flag"
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

	// Flag
	playlistId := flag.String("i", "", "Playlist ID")
	flag.Parse()
	if *playlistId == "" {
		log.Fatalln("Specify the playlist ID. e.g. go run ./cmd/get_playlist_items -i abcdefghijklmnopqrstuvwxyz")
	}

	client := oauth2.GetClient(youtube.YoutubeReadonlyScope)

	service, err := youtube.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	// Get all playlist items
	videos := []Video{}
	nextPageToken := ""
	retrievedItemCount := 0
	for {

		call := service.PlaylistItems.List([]string{"snippet,status"}).
			PlaylistId(*playlistId).
			MaxResults(50) // Probably MaxResults is max 50
		if nextPageToken != "" {
			call = call.PageToken(nextPageToken)
		}

		response, err := call.Do()
		handleError(err, "")

		nextPageToken = response.NextPageToken
		retrievedItemCount += len(response.Items)
		log.Printf("Page loaded. (%d/%d)\n", retrievedItemCount, response.PageInfo.TotalResults)

		for _, item := range response.Items {
			videos = append(videos, Video{
				VideoId:      item.Snippet.ResourceId.VideoId,
				VideoTitle:   item.Snippet.Title,
				ChannelId:    item.Snippet.VideoOwnerChannelId,
				ChannelTitle: item.Snippet.VideoOwnerChannelTitle,
				IsPublic:     item.Status.PrivacyStatus == "public",
			})
		}

		if response.NextPageToken == "" {
			break
		}
	}

	// Write to CSV
	file, err := os.Create("playlist_items.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	if err := dataframe.LoadStructs(videos).WriteCSV(file); err != nil {
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
