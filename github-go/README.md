# GitHub Go

A Go application that uses the GitHub API to explore repositories, users, and organizations.

## Code generation
generted by copilot

## Setup

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Create a new personal access token with the required scopes (e.g., `repo`, `read:org`, etc.)
3. Create a `.env` file in the root directory with the following content:
   ```
   GITHUB_TOKEN=your_personal_access_token_here
   ```

## Running the Application

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the application:
   ```bash
   go run cmd/github-go/main.go
   ```

3. Follow the instructions in the terminal to explore GitHub data.

## Features

- Authenticate with GitHub using a personal access token
- Fetch and display user profile information
- List repositories for a user or organization
- Search for repositories by keyword
