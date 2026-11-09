package command

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"io"
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
)

func TestSSHKeyFiles(t *testing.T) {
	pub, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	sshPub, err := ssh.NewPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	block, err := ssh.MarshalPrivateKey(private, "test")
	if err != nil {
		t.Fatal(err)
	}
	ageIdentity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	publicFile := filepath.Join(dir, "recipients")
	privateFile := filepath.Join(dir, "identity")
	data := "# recipients\n\n" + ageIdentity.Recipient().String() + "\n" + string(ssh.MarshalAuthorizedKey(sshPub))
	if err := os.WriteFile(publicFile, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(privateFile, pem.EncodeToMemory(block), 0600); err != nil {
		t.Fatal(err)
	}
	recipients, err := parseRecipientFile(publicFile)
	if err != nil || len(recipients) != 2 {
		t.Fatalf("recipients: %v, %v", recipients, err)
	}
	identities, err := parseIdentityFile(&cli.Context{}, privateFile)
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
	r, err := age.Decrypt(&encrypted, identities...)
	if err != nil {
		t.Fatal(err)
	}
	plain, err := io.ReadAll(r)
	if err != nil || string(plain) != "secret" {
		t.Fatalf("plaintext: %q, %v", plain, err)
	}
}

func TestSSHEncryptedIdentityRequiresTerminal(t *testing.T) {
	_, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	block, err := ssh.MarshalPrivateKeyWithPassphrase(key, "test", []byte("passphrase"))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "identity")
	if err := os.WriteFile(path, pem.EncodeToMemory(block), 0600); err != nil {
		t.Fatal(err)
	}
	ctx := &cli.Context{App: &cli.App{Reader: bytes.NewReader(nil), Writer: io.Discard}}
	if _, err := parseIdentityFile(ctx, path); err == nil {
		t.Fatal("expected an error without a terminal")
	}
}

func TestDecryptReportsInvalidData(t *testing.T) {
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}
	w := Decrypt(nopWriteCloser{io.Discard}, []age.Identity{id})
	_, _ = w.Write([]byte("invalid"))
	if err := w.Close(); err == nil {
		t.Fatal("expected decryption error")
	}
}

type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }
