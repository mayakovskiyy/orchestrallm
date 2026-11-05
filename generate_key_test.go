package main

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateKeyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "keys", "master.key")
	key := gkFile(path)
	if len(key) != 32 {
		t.Fatal("expected 32-byte key")
	}
	if !bytes.Equal(key, gkFile(path)) {
		t.Fatal("existing key was replaced")
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatal("expected file permissions 600")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(bytes.TrimSpace(data)))
	if err != nil || !bytes.Equal(decoded, key) {
		t.Fatal("saved key does not match returned key")
	}
}

func TestInvalidKeyFile(t *testing.T) {
	for _, value := range []string{"!", "", base64.StdEncoding.EncodeToString(make([]byte, 16))} {
		path := filepath.Join(t.TempDir(), "master.key")
		if err := os.WriteFile(path, []byte(value), 0600); err != nil {
			t.Fatal(err)
		}
		if gkFile(path) != nil {
			t.Fatal("accepted invalid key")
		}
		data, err := os.ReadFile(path)
		if err != nil || string(data) != value {
			t.Fatal("invalid existing key was overwritten")
		}
	}
}
