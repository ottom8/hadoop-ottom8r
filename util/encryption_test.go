package util

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"testing"
)

func EncryptDecryptTest(t *testing.T) {

	if len(os.Args) != 2 {
		fmt.Println("You must drag and drop a file!")
		waitForKey()
		os.Exit(1)
	} else if !isFile(os.Args[1]) {
		fmt.Println("File does not exist!")
		waitForKey()
		os.Exit(1)
	}

	file := os.Args[1]
	key := "testtesttesttest"

	if strings.ToLower(filepath.Ext(file)) != ".enc" {
		content, err := readFromFile(file)
		if err != nil {
			fmt.Println(err)
			waitForKey()
			os.Exit(1)
		}
		encrypted := Encrypt(string(content), key)
		writeToFile(encrypted, file+".enc")
	} else {
		content, err := readFromFile(file)
		if err != nil {
			fmt.Println(err)
			waitForKey()
			os.Exit(1)
		}
		decrypted := Decrypt(string(content), key)
		writeToFile(decrypted, file[:len(file)-4])
	}
}

