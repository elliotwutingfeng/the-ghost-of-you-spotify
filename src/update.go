package src

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	mathRand "math/rand"
	"net/http"
	"net/url"
	"slices"
	"time"
)

type SpotifySearchResponse struct {
	Tracks struct {
		Items []struct {
			ID               string   `json:"id"`
			AvailableMarkets []string `json:"available_markets"`
			IsLocal          bool     `json:"is_local"`
		} `json:"items"`
	} `json:"tracks"`
}

type SpotifyClient struct {
	Client      *http.Client
	AccessToken string
}

func (sc *SpotifyClient) newRequest(method, rawURL string, body any) (*http.Request, error) {
	var buf *bytes.Buffer
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	} else {
		buf = &bytes.Buffer{}
	}

	req, err := http.NewRequest(method, rawURL, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sc.AccessToken))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (sc *SpotifyClient) getJSON(req *http.Request, v any) error {
	resp, err := sc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Unexpected Spotify API response status code: %d", resp.StatusCode)
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

func Update(accessToken string, market string) error {
	sc := &SpotifyClient{
		Client:      &http.Client{Timeout: 30 * time.Second},
		AccessToken: accessToken,
	}

	var trackID string
	var attemptCount int
	const maxAttempts int = 30 // Avoid spamming Spotify API perpetually.

	fmt.Println("🔍 Looking for a suitable ghost track...")

	for trackID == "" && attemptCount <= maxAttempts {
		attemptCount++

		// Use a random keyword to search for tracks from Spotify catalogue.
		var n int64
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(SearchTerms))))
		if err != nil {
			// Fallback to deterministic random number generator.
			n = int64(mathRand.Intn(len(SearchTerms)))
		}
		n = nBig.Int64()
		searchTerm := SearchTerms[n]

		query := url.Values{}
		query.Set("q", searchTerm)
		query.Set("type", "track")
		searchURL := fmt.Sprintf("https://api.spotify.com/v1/search?%s", query.Encode())
		req, _ := sc.newRequest("GET", searchURL, nil)

		// Find a suitable track from search results.
		// The track should not be from a local file and must be available in the user's market.
		var searchResponse SpotifySearchResponse
		if err := sc.getJSON(req, &searchResponse); err != nil {
			return err
		}
		for _, item := range searchResponse.Tracks.Items {
			if item.IsLocal {
				continue
			}
			if !slices.Contains(item.AvailableMarkets, market) {
				continue
			}

			// Accept the track only if it does not exist in the user's "Liked Songs" playlist.
			query = url.Values{}
			query.Set("ids", item.ID)
			checkURL := fmt.Sprintf("https://api.spotify.com/v1/me/tracks/contains?%s", query.Encode())
			req, _ = sc.newRequest("GET", checkURL, nil)

			var existsInLibraryResponse []bool
			if err := sc.getJSON(req, &existsInLibraryResponse); err != nil {
				continue
			}
			if len(existsInLibraryResponse) > 0 && !existsInLibraryResponse[0] {
				trackID = item.ID
				fmt.Println("🎯 Found track   | ID:", item.ID)
				break
			}
		}
	}

	if trackID == "" {
		log.Fatalf("Exceeded maximum limit of %d attempts. No suitable track found.", maxAttempts)
	}

	tracksURL := "https://api.spotify.com/v1/me/tracks"
	tracksBody := map[string][]string{"ids": {trackID}}

	// Add track
	req, _ := sc.newRequest("PUT", tracksURL, tracksBody)
	if err := sc.getJSON(req, nil); err != nil {
		log.Fatalf("Failed to add track: %v", err)
	}
	fmt.Println("📝 Added track   | ID:", trackID)

	// Wait a short while
	time.Sleep(4 * time.Second)

	// Remove track
	req, _ = sc.newRequest("DELETE", tracksURL, tracksBody)
	if err := sc.getJSON(req, nil); err != nil {
		log.Fatalf("Failed to remove track: %v", err)
	}
	fmt.Println("❌ Removed track | ID:", trackID)
	fmt.Println("👻 Boo! Your \"Liked Songs\" playlist should be synced up now across all devices.")

	return nil
}
