package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func gk() []byte {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return gkFile(filepath.Join(homeDir, ".config", "orchestrallm", "master.key"))
}

func gkFile(path string) []byte {
	key, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println(err)
			return nil
		}

		key, err = exec.Command("openssl", "rand", "-base64", "32").Output()
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			fmt.Println(err)
			return nil
		}

		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			if os.IsExist(err) {
				return gkFile(path)
			}
			fmt.Println(err)
			return nil
		}

		_, err = file.Write(key)
		closeErr := file.Close()
		if err != nil || closeErr != nil {
			os.Remove(path)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(closeErr)
			}
			return nil
		}
	}

	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(key)))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if len(decoded) != 32 {
		fmt.Println("master key must contain exactly 32 bytes")
		return nil
	}

	return decoded
}
