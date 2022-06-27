package main

import (
	"fmt"
	"os"
	"strings"
)

type configuration struct {
	repoOwner      string
	repoName       string
	token          string
	labelPrefix    string
	labelSeparator string
	addLabel       bool
	defaultLabel   string
	onlyMilestone  bool
}

func newConfiguration() (*configuration, error) {
	repoFullName := os.Getenv(envVarRepoFullName)
	if repoFullName == "" {
		return nil, fmt.Errorf("%s is unset for the repo name", envVarRepoFullName)
	}

	labelPrefix := os.Getenv(envVarLabelPrefix)
	if labelPrefix == "" {
		return nil, fmt.Errorf("%s is unset for input prefix", envVarLabelPrefix)
	}

	labelSeparator := os.Getenv(envVarLabelSeparator)
	if labelSeparator == "" {
		labelSeparator = "/"
	}

	addLabel := os.Getenv(envVarAddLabel) == "true"
	defaultLabel := os.Getenv(envVarDefaultLabel)
	if addLabel && defaultLabel == "" {
		return nil, fmt.Errorf("if add label is set, you must specify a default label")
	}

	repoParts := strings.Split(repoFullName, "/")
	if len(repoParts) != 2 {
		return nil, fmt.Errorf("unexpected repository format")
	}
	repoOwner := repoParts[0]
	repoName := repoParts[1]

	return &configuration{
		repoOwner:      repoOwner,
		repoName:       repoName,
		token:          os.Getenv(envVarToken),
		labelPrefix:    labelPrefix,
		labelSeparator: labelSeparator,
		addLabel:       addLabel,
		defaultLabel:   defaultLabel,
		onlyMilestone:  os.Getenv(envVarOnlyMilestone) == "true",
	}, nil
}
