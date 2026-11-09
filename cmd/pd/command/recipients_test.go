package command

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"filippo.io/age"
)

func TestHybridRecipients(t *testing.T) {
	identity, err := age.GenerateHybridIdentity()
	if err != nil {
		t.Fatal(err)
	}
	key := identity.Recipient().String()
	path := filepath.Join(t.TempDir(), "recipients")
	if err := os.WriteFile(path, []byte("# comment\n\n"+key+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	for name, parse := range map[string]func() ([]age.Recipient, error){
		"argument": func() ([]age.Recipient, error) { return parseRecipients([]string{key}) },
		"file":     func() ([]age.Recipient, error) { return parseRecipientFile(path) },
	} {
		t.Run(name, func(t *testing.T) {
			recipients, err := parse()
			if err != nil {
				t.Fatal(err)
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
			r, err := age.Decrypt(&encrypted, identity)
			if err != nil {
				t.Fatal(err)
			}
			plain, err := io.ReadAll(r)
			if err != nil || string(plain) != "secret" {
				t.Fatalf("plaintext: %q, %v", plain, err)
			}
		})
	}
}

func TestRecipientFileErrors(t *testing.T) {
	for _, tc := range []struct {
		name, contents, want string
	}{
		{"empty", "# comment\n\n", "no recipients found"},
		{"invalid key", "# comment\ninvalid\n", "error at line 2"},
		{"invalid UTF-8", "\xff\n", "not valid UTF-8"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "recipients")
			if err := os.WriteFile(path, []byte(tc.contents), 0600); err != nil {
				t.Fatal(err)
			}
			_, err := parseRecipientFile(path)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %q, got %v", tc.want, err)
			}
		})
	}
}

func TestHybridCannotMixWithOtherRecipients(t *testing.T) {
	hybrid, err := age.GenerateHybridIdentity()
	if err != nil {
		t.Fatal(err)
	}
	classic, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	recipients, err := parseRecipients([]string{hybrid.Recipient().String(), classic.Recipient().String()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := age.Encrypt(io.Discard, recipients...); err == nil {
		t.Fatal("expected incompatible recipients to be rejected")
	}
}
