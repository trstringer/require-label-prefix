package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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

func labelsContainPrefix(labels []github.Label, prefix, separator string) bool {
	for _, label := range labels {
		if strings.HasPrefix(label.GetName(), fmt.Sprintf("%s%s", prefix, separator)) {
			return true
		}
	}

	return false
}

func processIssues(ctx context.Context, githubClient *github.Client, config *configuration) error {
	issues, _, err := githubClient.Issues.ListByRepo(
		ctx,
		config.repoOwner,
		config.repoName,
		&github.IssueListByRepoOptions{
			State: "open",
		})
	if err != nil {
		return fmt.Errorf("error getting issues: %v", err)
	}

	for _, issue := range issues {
		// Filter out pull requests, which are returned from the issues API.
		if issue.GetPullRequestLinks() != nil {
			continue
		}

		if config.onlyMilestone && issue.GetMilestone() == nil {
			continue
		}

		if !labelsContainPrefix(issue.Labels, config.labelPrefix, config.labelSeparator) {
			fmt.Printf(
				"Issue #%d does not have the required label prefix: \"%s%s\"\n",
				issue.GetNumber(),
				config.labelPrefix,
				config.labelSeparator,
			)

			var comment string
			if config.addLabel {
				githubClient.Issues.AddLabelsToIssue(
					ctx,
					config.repoOwner,
					config.repoName,
					issue.GetNumber(),
					[]string{config.defaultLabel},
				)
				comment = fmt.Sprintf("Added default label `%s`. Please consider re-labeling this issue appropriately.", config.defaultLabel)
			} else {
				comment = fmt.Sprintf(
					"No label with prefix \"%s%s\" found. Please add the appropriate label.",
					config.labelPrefix,
					config.labelSeparator,
				)
			}
			_, _, err := githubClient.Issues.CreateComment(
				ctx,
				config.repoOwner,
				config.repoName,
				issue.GetNumber(),
				&github.IssueComment{Body: &comment},
			)
			if err != nil {
				return fmt.Errorf("error adding comment: %v", err)
			}
		}
	}

	return nil
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
