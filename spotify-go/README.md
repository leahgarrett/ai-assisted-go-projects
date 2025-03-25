# Spotify Go

A Go application that uses the Spotify API to display your profile and top tracks.

## Prompts used to generate this code

Code generated with Windsurf

Started with:  
'Can we create a new project that uses the spotify API using Go'

This ran fine so I asked for:
'Can you change this to display the top 50 most popular artists on spotify?'

The code then did not work. Windsurf continued to add debugging statements.



## Setup

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new application
3. Add `http://localhost:8080/callback` to the Redirect URIs in your app settings
4. Create a `.env` file in the root directory with the following content:
   ```
   SPOTIFY_CLIENT_ID=your_client_id_here
   SPOTIFY_CLIENT_SECRET=your_client_secret_here
   ```

## Running the Application

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the application:
   ```bash
   go run cmd/spotify-go/main.go
   ```

3. Open http://localhost:8080 in your browser

## Features

- OAuth2 authentication with Spotify
- View your Spotify profile
- See your top 10 tracks from the last month
