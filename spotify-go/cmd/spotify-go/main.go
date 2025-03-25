package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

const (
	redirectURI = "http://localhost:8080/callback"
)

var (
	spotifyConfig *oauth2.Config
	state        = "random-state" // In production, use a proper random state
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	log.Printf("ClientID: %s", clientID)
	log.Printf("ClientSecret: %s", clientSecret)
	log.Printf("RedirectURI: %s", redirectURI)

	if clientID == "" || clientSecret == "" {
		log.Fatal("Missing Spotify credentials in .env file")
	}

	httpClient := &http.Client{
		Transport: &loggingTransport{http.DefaultTransport},
	}

	spotifyConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes: []string{
			"user-read-private",
			"user-read-email",
			"user-top-read",
		},
		Endpoint: spotify.Endpoint,
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	spotifyConfig.Client(ctx, nil)
}

type loggingTransport struct {
	rt http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("Making request to %s", req.URL)
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(body))
		log.Printf("Request body: %s", string(body))
	}
	resp, err := t.rt.RoundTrip(req)
	if err != nil {
		log.Printf("Request error: %v", err)
		return resp, err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
		log.Printf("Response status: %d, body: %s", resp.StatusCode, string(body))
	}
	return resp, err
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	html := `
		<html>
			<head>
				<title>Top 50 Spotify Artists</title>
				<style>
					body {
						font-family: Arial, sans-serif;
						max-width: 1000px;
						margin: 0 auto;
						padding: 2rem;
						text-align: center;
						background-color: #121212;
						color: white;
					}
					.login-button {
						display: inline-block;
						background-color: #1DB954;
						color: white;
						padding: 1rem 2rem;
						text-decoration: none;
						border-radius: 500px;
						font-weight: bold;
						margin-top: 2rem;
					}
					.login-button:hover {
						background-color: #1ed760;
						transform: scale(1.05);
						transition: all 0.2s;
					}
					h1 {
						color: #1DB954;
						font-size: 2.5rem;
						margin-bottom: 1rem;
					}
					p {
						color: #b3b3b3;
						font-size: 1.2rem;
					}
				</style>
			</head>
			<body>
				<h1>Top 50 Spotify Artists</h1>
				<p>Discover the most popular artists on Spotify right now</p>
				<a href="/login" class="login-button">View Top Artists</a>
			</body>
		</html>
	`
	fmt.Fprint(w, html)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := spotifyConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

type Artist struct {
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
	Followers  struct {
		Total int `json:"total"`
	} `json:"followers"`
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
	Genres []string `json:"genres"`
}

type SearchResponse struct {
	Artists struct {
		Items []Artist `json:"items"`
	} `json:"artists"`
}

func formatNumber(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received callback with state: %s", r.FormValue("state"))
	if r.FormValue("state") != state {
		log.Printf("State mismatch. Expected: %s, Got: %s", state, r.FormValue("state"))
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		log.Printf("No code received in callback")
		http.Error(w, "No code received", http.StatusBadRequest)
		return
	}
	log.Printf("Received authorization code: %s", code)

	ctx := context.WithValue(r.Context(), oauth2.HTTPClient, &http.Client{
		Transport: &loggingTransport{http.DefaultTransport},
	})

	token, err := spotifyConfig.Exchange(ctx, code)
	if err != nil {
		log.Printf("Token exchange error: %v", err)
		if oauthErr, ok := err.(*oauth2.RetrieveError); ok {
			log.Printf("OAuth2 error details - Status: %d, Body: %s", oauthErr.Response.StatusCode, string(oauthErr.Body))
		}
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := spotifyConfig.Client(ctx, token)

	// Search for popular artists
	searchResp, err := client.Get("https://api.spotify.com/v1/search?type=artist&limit=50&q=genre:pop genre:rock genre:hip-hop genre:rap genre:latin")
	if err != nil {
		http.Error(w, "Failed to search artists", http.StatusInternalServerError)
		return
	}
	defer searchResp.Body.Close()

	var response SearchResponse
	if err := json.NewDecoder(searchResp.Body).Decode(&response); err != nil {
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}

	// Sort artists by popularity (most popular first)
	artists := response.Artists.Items

	// Generate HTML response
	html := `
		<html>
			<head>
				<title>Top 50 Spotify Artists</title>
				<style>
					body {
						font-family: Arial, sans-serif;
						max-width: 1200px;
						margin: 0 auto;
						padding: 2rem;
						background-color: #121212;
						color: white;
					}
					h1 {
						text-align: center;
						color: #1DB954;
						font-size: 2.5rem;
						margin-bottom: 2rem;
					}
					.artists-grid {
						display: grid;
						grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
						gap: 2rem;
						padding: 1rem;
					}
					.artist-card {
						background: #282828;
						border-radius: 8px;
						padding: 1.5rem;
						transition: all 0.3s ease;
					}
					.artist-card:hover {
						transform: translateY(-5px);
						box-shadow: 0 10px 20px rgba(0,0,0,0.2);
					}
					.artist-image {
						width: 100%;
						height: 200px;
						object-fit: cover;
						border-radius: 4px;
						margin-bottom: 1rem;
					}
					.artist-name {
						font-size: 1.4rem;
						font-weight: bold;
						margin: 0.5rem 0;
						color: white;
					}
					.artist-stats {
						color: #b3b3b3;
						font-size: 0.9rem;
					}
					.genres {
						margin-top: 0.5rem;
						display: flex;
						flex-wrap: wrap;
						gap: 0.5rem;
					}
					.genre-tag {
						background: #1DB954;
						color: white;
						padding: 0.2rem 0.6rem;
						border-radius: 50px;
						font-size: 0.8rem;
					}
					.popularity-bar {
						width: 100%;
						height: 4px;
						background: #404040;
						border-radius: 2px;
						margin-top: 0.5rem;
					}
					.popularity-fill {
						height: 100%;
						background: #1DB954;
						border-radius: 2px;
					}
				</style>
			</head>
			<body>
				<h1>Top 50 Most Popular Artists on Spotify</h1>
				<div class="artists-grid">
	`

	for _, artist := range artists {
		imageURL := "/default-artist.jpg"
		if len(artist.Images) > 0 {
			imageURL = artist.Images[0].URL
		}

		genres := ""
		if len(artist.Genres) > 0 {
			genres = "<div class=\"genres\">"
			for i, genre := range artist.Genres {
				if i < 3 { // Show only first 3 genres
					genres += fmt.Sprintf("<span class=\"genre-tag\">%s</span>", genre)
				}
			}
			genres += "</div>"
		}

		html += fmt.Sprintf(`
			<div class="artist-card">
				<img src="%s" alt="%s" class="artist-image">
				<h2 class="artist-name">%s</h2>
				<div class="artist-stats">
					<div>Followers: %s</div>
					<div>Popularity: %d%%</div>
					<div class="popularity-bar">
						<div class="popularity-fill" style="width: %d%%"></div>
					</div>
				</div>
				%s
			</div>
		`, imageURL, artist.Name, artist.Name, formatNumber(artist.Followers.Total), artist.Popularity, artist.Popularity, genres)
	}

	html += `
				</div>
			</body>
		</html>
	`

	fmt.Fprint(w, html)
}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)

	log.Printf("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
