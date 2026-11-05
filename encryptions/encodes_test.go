package encryptions

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	key := bytes.Repeat([]byte{1}, 32)
	for _, plaintext := range []string{"", "sk-example", "ключ 🔑"} {
		encrypted, err := Encrypt(plaintext, key)
		if err != nil {
			t.Fatal(err)
		}
		decrypted, err := Decrypt(encrypted, key)
		if err != nil || decrypted != plaintext {
			t.Fatalf("round trip failed: %v", err)
		}
		second, err := Encrypt(plaintext, key)
		if err != nil {
			t.Fatal(err)
		}
		if encrypted == second {
			t.Fatal("nonce reused")
		}
	}
}

func TestRejectInvalidInput(t *testing.T) {
	key := bytes.Repeat([]byte{1}, 32)
	for _, size := range []int{0, 16, 24, 31, 33} {
		if _, err := Encrypt("secret", make([]byte, size)); err == nil {
			t.Fatalf("accepted key length %d", size)
		}
		if _, err := Decrypt("", make([]byte, size)); err == nil {
			t.Fatalf("accepted key length %d", size)
		}
	}
	encrypted, err := Encrypt("secret", key)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Decrypt(encrypted, bytes.Repeat([]byte{2}, 32)); err == nil {
		t.Fatal("accepted wrong key")
	}
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		t.Fatal(err)
	}
	data[len(data)-1] ^= 1
	for _, input := range []string{"!", "", base64.StdEncoding.EncodeToString(data[:3]), base64.StdEncoding.EncodeToString(data)} {
		if _, err := Decrypt(input, key); err == nil {
			t.Fatal("accepted damaged ciphertext")
		}
	}
}
