package twitchmsg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/adeithe/go-twitch/api"
	"github.com/gempir/go-twitch-irc/v4"
)

const (
	ClientID        = ""
	ClientSecret    = ""
	messagesToFetch = 2
	streamsToFetch  = 2
)

func DoTwitchRequest() [messagesToFetch * streamsToFetch]string {
	ctx := context.Background()
	token := fetch_token()
	client := api.New(ClientID, api.WithDefaultBearerToken(token))

	streams, err := client.Streams.List().First(streamsToFetch).Do(ctx)
	if err != nil {
		panic(err)
	}

	var all_messages [messagesToFetch * streamsToFetch]string

	for streamIdx, stream := range streams.Data {
		// fmt.Printf("%s is streaming %s to %d viewers\n",
		// 	stream.UserLogin,
		// 	stream.GameName,
		// 	stream.ViewerCount,
		// )
		msg := listen_to_chat(stream.UserLogin)

		for i, m := range msg {
			currentIdx := streamIdx*messagesToFetch + i
			all_messages[currentIdx] = m
		}
	}

	return all_messages
}

func listen_to_chat(channel string) [messagesToFetch]string {
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

func fetch_token() string {
	// Twitch OAuth token endpoint
	tokenURL := "https://id.twitch.tv/oauth2/token"

	// Build the form data
	data := url.Values{}
	data.Set("client_id", ClientID)
	data.Set("client_secret", ClientSecret)
	data.Set("grant_type", "client_credentials")

	// Make the POST request
	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Check for non-200 status
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Failed to get token: %s", body))
	}

	// Parse JSON
	var token TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		panic(err)
	}

	return token.AccessToken
}
