package command

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"filippo.io/age"
	"github.com/jkawamoto/go-pixeldrain/cmd/pd/status"
	"github.com/urfave/cli/v2"
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
	if err := validateRecipientMix([]age.Recipient{hybrid.Recipient(), classic.Recipient()}); err == nil {
		t.Fatal("expected incompatible recipients to be rejected")
	}
}

func TestCmdUploadRejectsHybridRecipientMixBeforeUpload(t *testing.T) {
	hybrid, err := age.GenerateHybridIdentity()
	if err != nil {
		t.Fatal(err)
	}
	classic, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	for _, recipientFile := range []bool{false, true} {
		t.Run(map[bool]string{false: "arguments", true: "recipient file"}[recipientFile], func(t *testing.T) {
			args := []string{"--" + FlagRecipient, hybrid.Recipient().String()}
			if recipientFile {
				path := filepath.Join(t.TempDir(), "recipients")
				if err := os.WriteFile(path, []byte(classic.Recipient().String()), 0600); err != nil {
					t.Fatal(err)
				}
				args = append(args, "--"+FlagRecipientFile, path)
			} else {
				args = append(args, "--"+FlagRecipient, classic.Recipient().String())
			}
			args = append(args, "missing-file")
			flagSet := flag.NewFlagSet("upload", flag.PanicOnError)
			flagSet.Var(&cli.StringSlice{}, FlagRecipient, "")
			flagSet.String(FlagRecipientFile, "", "")
			if err := flagSet.Parse(args); err != nil {
				t.Fatal(err)
			}
			ctx := cli.NewContext(cli.NewApp(), flagSet, nil)
			ctx.App.Reader = strings.NewReader("")
			ctx.App.Writer = io.Discard
			ctx.Context = context.Background()
			err := CmdUpload(ctx)
			var exit cli.ExitCoder
			if !errors.As(err, &exit) || exit.ExitCode() != status.InvalidArgument || !strings.Contains(err.Error(), "hybrid age keys") {
				t.Fatalf("expected pre-upload invalid-argument error, got %v", err)
			}
		})
	}
}
