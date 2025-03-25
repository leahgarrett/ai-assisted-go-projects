package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v50/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set in the environment")
	}

	// Authenticate with GitHub API
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// Fetch and display authenticated user's profile
	user, _, err := client.Users.Get(oauth2.NoContext, "")
	if err != nil {
		log.Fatalf("Error fetching user profile: %v", err)
	}

	fmt.Printf("Authenticated as: %s\n", *user.Login)
	fmt.Printf("Name: %s\n", *user.Name)
	fmt.Printf("Bio: %s\n", *user.Bio)
}
