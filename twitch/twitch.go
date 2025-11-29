package twitchmsg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/adeithe/go-twitch/api"
	"github.com/gempir/go-twitch-irc/v4"
)

var (
	CLIENT_ID     = os.Getenv("TWITCH_CLIENT_ID")
	CLIENT_SECRET = os.Getenv("TWITCH_CLIENT_SECRET")
)

const (
	messagesToFetch = 5
	streamsToFetch  = 3
)

// FetchTwitchMessages fetches chat messages from live Twitch streams.
// It finds the top `streamsToFetch` live streams and listens to their chat channels,
// collecting `messagesToFetch` messages from each stream.
// It returns an array of collected messages.
func FetchTwitchMessages() [messagesToFetch * streamsToFetch]string {
	ctx := context.Background()
	token := fetchOauthToken()
	client := api.New(CLIENT_ID, api.WithDefaultBearerToken(token))

	streams, err := client.Streams.List().First(streamsToFetch).Do(ctx)
	if err != nil {
		panic(err)
	}

	var all_messages [messagesToFetch * streamsToFetch]string
	var wg sync.WaitGroup

	for streamIdx, stream := range streams.Data {
		wg.Go(func() {
			msg := listenToChat(stream.UserLogin)

			for i, m := range msg {
				currentIdx := streamIdx*messagesToFetch + i
				all_messages[currentIdx] = m
			}
		})
	}

	wg.Wait()

	return all_messages
}

// listenToChat connects to a Twitch chat channel and listens for messages.
// It collects `messagesToFetch` messages and then disconnects.
// It returns an array of collected messages.
func listenToChat(channel string) [messagesToFetch]string {
	var messages [messagesToFetch]string
	var count int = 0

	client := twitch.NewAnonymousClient()
	client.Join(channel)
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		// fmt.Println(message.Message)
		messages[count] = message.Message
		count++
		if count >= messagesToFetch {
			client.Depart(channel)
			client.Disconnect()
		}
	})
	err := client.Connect()
	if err != nil {
		if err == twitch.ErrClientDisconnected {
			// Normal disconnection
			return messages
		}
		panic(err)
	}

	fmt.Println("Finished listening to chat for channel:", channel)

	return messages
}

// Response structure from Twitch
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// fetchOauthToken fetches an OAuth token from Twitch using client credentials.
func fetchOauthToken() string {
	// Twitch OAuth token endpoint
	tokenURL := "https://id.twitch.tv/oauth2/token"

	data := url.Values{}
	data.Set("client_id", CLIENT_ID)
	data.Set("client_secret", CLIENT_SECRET)
	data.Set("grant_type", "client_credentials")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Failed to get token: %s", body))
	}

	var token TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		panic(err)
	}

	return token.AccessToken
}
