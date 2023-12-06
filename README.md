# ytpl: YouTube PlayList manager

## Usage: retrieve your playlist items

1. generate OAuth 2.0 credentials in Google Cloud console.
2. deploy the client\_secret.json downloaded from Google Cloud console.
3. execute `go run ./cmd/get_playlist_items -i [PlaylistID]`
4. display the generated `playlist_items.csv`
