// github_recipient.go
//
// Copyright (c) 2018-2026 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"filippo.io/age"
	"filippo.io/age/agessh"
	"golang.org/x/crypto/ssh"
)

const (
	githubKeysURL       = "https://github.com/"
	githubKeysTimeout   = 10 * time.Second
	githubKeysMaxLength = 1 << 20
)

func resolveGitHubRecipients(
	ctx context.Context,
	username string,
	client *http.Client,
	warnings io.Writer,
) ([]age.Recipient, error) {
	if !validGitHubUsername(username) {
		return nil, fmt.Errorf("invalid GitHub username %q", username)
	}
	if client == nil {
		client = &http.Client{}
	}
	clientCopy := *client
	clientCopy.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
	ctx, cancel := context.WithTimeout(ctx, githubKeysTimeout)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, githubKeysURL+username+".keys", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "text/plain")
	response, err := clientCopy.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve SSH keys for GitHub user %q: %w", username, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub returned HTTP status %s for user %q", response.Status, username)
	}

	data, err := io.ReadAll(io.LimitReader(response.Body, githubKeysMaxLength+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH keys for GitHub user %q: %w", username, err)
	}
	if len(data) > githubKeysMaxLength {
		return nil, fmt.Errorf("SSH keys response for GitHub user %q is too large", username)
	}

	var recipients []age.Recipient
	for lineNumber, line := range strings.Split(string(data), "\n") {
		lineNumber++
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !utf8.ValidString(line) {
			return nil, fmt.Errorf("invalid SSH key at line %d: not valid UTF-8", lineNumber)
		}
		key, _, _, rest, err := ssh.ParseAuthorizedKey([]byte(line))
		if err != nil || len(strings.TrimSpace(string(rest))) != 0 {
			if err == nil {
				err = errors.New("unexpected data after SSH key")
			}
			return nil, fmt.Errorf("invalid SSH key at line %d: %w", lineNumber, err)
		}
		if key.Type() != ssh.KeyAlgoED25519 && key.Type() != ssh.KeyAlgoRSA {
			if warnings != nil {
				_, _ = fmt.Fprintf(warnings, "warning: skipping unsupported SSH key type %q for GitHub user %q\n", key.Type(), username)
			}
			continue
		}
		recipient, err := agessh.ParseRecipient(line)
		if err != nil {
			return nil, fmt.Errorf("invalid SSH key at line %d: %w", lineNumber, err)
		}
		recipients = append(recipients, recipient)
	}
	if len(recipients) == 0 {
		return nil, fmt.Errorf("GitHub user %q has no supported SSH keys", username)
	}
	return recipients, nil
}

func validGitHubUsername(username string) bool {
	if len(username) == 0 || len(username) > 39 || username[0] == '-' || username[len(username)-1] == '-' || strings.Contains(username, "--") {
		return false
	}
	for _, char := range username {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < '0' || char > '9') && char != '-' {
			return false
		}
	}
	return true
}
