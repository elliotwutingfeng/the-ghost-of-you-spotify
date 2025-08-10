package src

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/toqueteos/webbrowser"
)

const (
	SCOPE              string = "user-library-read user-library-modify"
	AUTHORIZE_ENDPOINT string = "https://accounts.spotify.com/authorize"
	API_TOKEN_ENDPOINT string = "https://accounts.spotify.com/api/token"
)

type accessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Opens up web browser for user to sign in and grant app access to read/modify their Spotify account's track library.
// Returns an authorization code that can be exchanged for an access token and a refresh token.
func signInToSpotify(clientID string, redirectHost string) (authorizationCode string) {
	// Randomly generated state value for maintaining state between the request and callback.
	// This provides protection against cross-site request forgery attacks.
	// See https://datatracker.ietf.org/doc/html/rfc6749#section-4.1
	state := rand.Text()

	u, _ := url.Parse(AUTHORIZE_ENDPOINT)
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", clientID)
	q.Set("scope", SCOPE)
	q.Set("redirect_uri", fmt.Sprintf("http://%s/callback", redirectHost))
	q.Set("state", state)
	u.RawQuery = q.Encode()

	fmt.Println("Opening Spotify sign-in page...")
	fmt.Println()
	fmt.Println("If your browser did not open automatically, copy and paste the following link into your web browser's address bar.")
	fmt.Println()
	fmt.Println(u)
	fmt.Println()
	webbrowser.Open(u.String())

	// Start local HTTP server to catch Spotify API redirect.
	listener, err := net.Listen("tcp", redirectHost)
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}

	done := make(chan struct{})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		receivedState := r.FormValue("state")
		if receivedState != state {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			log.Fatalf("State mismatch: expected %s, got %s", state, receivedState)
		} else {
			w.Write([]byte("Authenticated with Spotify. You may now close this window."))
		}
		authorizationCode = r.FormValue("code")
		go func() {
			done <- struct{}{}
		}()
	})

	go func() {
		http.Serve(listener, nil)
	}()

	<-done
	return
}

func requestTokensWithCode(authorizationCode string, clientID string, clientSecret string, redirectHost string) *accessTokenResponse {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authorizationCode)
	data.Set("redirect_uri", "http://"+redirectHost+"/callback")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", API_TOKEN_ENDPOINT, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	res := accessTokenResponse{}
	json.Unmarshal(body, &res)
	return &res
}

func requestTokensWithRefreshToken(refreshToken string, clientID string, clientSecret string) *accessTokenResponse {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", API_TOKEN_ENDPOINT, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	res := accessTokenResponse{}
	json.Unmarshal(body, &res)

	if res.RefreshToken == "" {
		// If no refresh token provided in response, reuse old refresh token.
		res.RefreshToken = refreshToken
	}
	return &res
}

func GetTokens(clientID string, clientSecret string, redirectHost string, refreshToken string, pathToEnv string) (string, string) {
	var code string
	var res *accessTokenResponse

	// Try using refresh token to obtain a new access token (1 hour expiry).
	res = requestTokensWithRefreshToken(refreshToken, clientID, clientSecret)

	// No access token means refresh token is invalid; user has to sign in to get a new refresh token.
	if res.AccessToken == "" {
		code = signInToSpotify(clientID, redirectHost)
		res = requestTokensWithCode(code, clientID, clientSecret, redirectHost)
	}
	if res.AccessToken == "" {
		log.Fatal("Unable to obtain new access token.")
	}

	return res.AccessToken, res.RefreshToken
}
