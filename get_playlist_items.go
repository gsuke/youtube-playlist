package main

import (
	"fmt"
	"log"

	"google.golang.org/api/youtube/v3"
)

func main() {

	const playlistId = ""

	client := getClient(youtube.YoutubeReadonlyScope)

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	response, err := service.PlaylistItems.List([]string{"snippet,status"}).
		PlaylistId(playlistId).
		MaxResults(3).
		Do()
	handleError(err, "")

	for _, item := range response.Items {
		fmt.Printf(
			"動画: %v(%v), チャンネル: %v(%v), %v\n",
			item.Snippet.Title,
			item.Snippet.ResourceId.VideoId,
			item.Snippet.VideoOwnerChannelTitle,
			item.Snippet.VideoOwnerChannelId,
			item.Status.PrivacyStatus,
		)
	}

}
