package main

import (
	"fmt"
	"testing"

	"github.com/google/go-github/github"
)

func issuesContainIssue(issues []*github.Issue, issue *github.Issue) bool {
	for _, sourceIssue := range issues {
		if issue.GetTitle() == sourceIssue.GetTitle() {
			return true
		}
	}

	return false
}

func issuesAreEqual(source, destination []*github.Issue) error {
	sourceLength := len(source)
	destinationLength := len(destination)
	if sourceLength != destinationLength {
		return fmt.Errorf(
			"source (%v) length %d is not equal to destination (%v) length %d",
			source,
			sourceLength,
			destination,
			destinationLength,
		)
	}

	for _, sourceIssue := range source {
		if !issuesContainIssue(destination, sourceIssue) {
			return fmt.Errorf(
				"destination issues does not contain issue '%s'",
				sourceIssue.GetTitle(),
			)
		}
	}

	return nil
}

func TestIssuesToModify(t *testing.T) {
	issue1Title := "issue1"
	issue2Title := "issue2"
	labelPrefix := "prefix"
	separator := "/"
	labelNameMatching := fmt.Sprintf("%s%s%s", labelPrefix, separator, "suffix")

	testCases := []struct {
		name     string
		config   *configuration
		input    []*github.Issue
		expected []*github.Issue
	}{
		{
			name: "single_issue_missing_label",
			config: &configuration{
				onlyMilestone:  false,
				labelPrefix:    labelPrefix,
				labelSeparator: separator,
			},
			input: []*github.Issue{
				{
					Title: &issue1Title,
				},
			},
			expected: []*github.Issue{
				{
					Title: &issue1Title,
				},
			},
		},
		{
			name: "single_issue_contains_label",
			config: &configuration{
				onlyMilestone:  false,
				labelPrefix:    labelPrefix,
				labelSeparator: separator,
			},
			input: []*github.Issue{
				{
					Title: &issue1Title,
					Labels: []github.Label{
						{
							Name: &labelNameMatching,
						},
					},
				},
			},
			expected: []*github.Issue{},
		},
		{
			name: "single_pull_request",
			config: &configuration{
				onlyMilestone:  false,
				labelPrefix:    labelPrefix,
				labelSeparator: separator,
			},
			input: []*github.Issue{
				{
					Title:            &issue1Title,
					PullRequestLinks: &github.PullRequestLinks{},
				},
			},
			expected: []*github.Issue{},
		},
		{
			name: "multiple_issues",
			config: &configuration{
				onlyMilestone:  false,
				labelPrefix:    labelPrefix,
				labelSeparator: separator,
			},
			input: []*github.Issue{
				{
					Title: &issue1Title,
					Labels: []github.Label{
						{
							Name: &labelNameMatching,
						},
					},
				},
				{
					Title: &issue2Title,
				},
			},
			expected: []*github.Issue{
				{
					Title: &issue2Title,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := issuesToModify(testCase.input, testCase.config)
			if err := issuesAreEqual(testCase.expected, actual); err != nil {
				t.Errorf("%v", err)
			}
		})
	}
}
