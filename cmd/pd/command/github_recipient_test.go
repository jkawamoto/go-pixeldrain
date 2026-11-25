package command

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"filippo.io/age"
	"filippo.io/age/agessh"
	"golang.org/x/crypto/ssh"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestResolveGitHubRecipients(t *testing.T) {
	firstPublic, firstPrivate, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	secondPublic, secondPrivate, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	firstSSH, err := ssh.NewPublicKey(firstPublic)
	if err != nil {
		t.Fatal(err)
	}
	secondSSH, err := ssh.NewPublicKey(secondPublic)
	if err != nil {
		t.Fatal(err)
	}
	keys := string(ssh.MarshalAuthorizedKey(firstSSH)) + string(ssh.MarshalAuthorizedKey(secondSSH))
	client := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodGet || request.URL.String() != "https://github.com/octocat.keys" {
			t.Errorf("unexpected request: %s %s", request.Method, request.URL)
		}
		return githubResponse(http.StatusOK, keys), nil
	})}
	recipients, err := parseUploadRecipients(context.Background(), []string{"github:octocat"}, client, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(recipients) != 2 {
		t.Fatalf("expected 2 recipients, got %d", len(recipients))
	}

	var encrypted bytes.Buffer
	w, err := age.Encrypt(&encrypted, recipients...)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte("secret")); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	for _, privateKey := range []ed25519.PrivateKey{firstPrivate, secondPrivate} {
		identity, err := agessh.NewEd25519Identity(privateKey)
		if err != nil {
			t.Fatal(err)
		}
		r, err := age.Decrypt(bytes.NewReader(encrypted.Bytes()), identity)
		if err != nil {
			t.Fatal(err)
		}
		plain, err := io.ReadAll(r)
		if err != nil || string(plain) != "secret" {
			t.Fatalf("decrypted plaintext: %q, %v", plain, err)
		}
	}
}

func TestResolveGitHubRecipientsSkipsUnsupportedKey(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	unsupported, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := ssh.NewPublicKey(privateKey.Public())
	if err != nil {
		t.Fatal(err)
	}
	client := githubClient(http.StatusOK, string(ssh.MarshalAuthorizedKey(unsupported))+string(ssh.MarshalAuthorizedKey(publicKey)), nil)
	var warnings bytes.Buffer
	recipients, err := resolveGitHubRecipients(context.Background(), "octocat", client, &warnings)
	if err != nil {
		t.Fatal(err)
	}
	if len(recipients) != 1 || !strings.Contains(warnings.String(), "unsupported SSH key type") {
		t.Fatalf("recipients: %d, warnings: %q", len(recipients), warnings.String())
	}
}

func TestResolveGitHubRecipientsErrors(t *testing.T) {
	for _, tc := range []struct {
		name     string
		username string
		client   *http.Client
		want     string
	}{
		{name: "invalid username", username: "bad/name", want: "invalid GitHub username"},
		{name: "empty response", username: "octocat", client: githubClient(http.StatusOK, "\n", nil), want: "no supported SSH keys"},
		{name: "invalid key", username: "octocat", client: githubClient(http.StatusOK, "not-a-key\n", nil), want: "invalid SSH key"},
		{name: "HTTP error", username: "octocat", client: githubClient(http.StatusNotFound, "not found", nil), want: "HTTP status 404"},
		{name: "oversized response", username: "octocat", client: githubClient(http.StatusOK, strings.Repeat("a", githubKeysMaxLength+1), nil), want: "too large"},
		{name: "network error", username: "octocat", client: githubClient(0, "", errors.New("connection failed")), want: "connection failed"},
		{name: "redirect", username: "octocat", client: githubClient(http.StatusFound, "", nil), want: "HTTP status 302"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := resolveGitHubRecipients(context.Background(), tc.username, tc.client, nil)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected error containing %q, got %v", tc.want, err)
			}
		})
	}
}

func githubClient(status int, body string, requestError error) *http.Client {
	return &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		if requestError != nil {
			return nil, requestError
		}
		response := githubResponse(status, body)
		if status == http.StatusFound {
			response.Header.Set("Location", "https://attacker.example/keys")
		}
		return response, nil
	})}
}

func githubResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
