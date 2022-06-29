package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func newGithubClient(ctx context.Context, config *configuration) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func main() {
	config, err := newConfiguration()
	if err != nil {
		fmt.Printf("Error getting config: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf(
		"Analyzing %s/%s for prefix \"%s\" with separator \"%s\"\n",
		config.repoOwner,
		config.repoName,
		config.labelPrefix,
		config.labelSeparator,
	)

	ctx := context.Background()
	githubClient := newGithubClient(ctx, config)

	if err := processIssues(ctx, githubClient, config); err != nil {
		fmt.Printf("Error processing issues: %v\n", err)
		os.Exit(1)
	}
}
