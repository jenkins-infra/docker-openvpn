package helpers

import (
	"fmt"
	"log"
	"os"

	"go.mozilla.org/sops/v3/decrypt"
)

// DecryptPrivateDir decrypt ca private key and ca private key password
func DecryptPrivateDir() {
	dirname := "cert/pki/private/"
	files := []string{"ca.key"}

	for _, file := range files {
		out, err := decrypt.File(dirname+file+".enc", "txt")
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(dirname + file)
		if err != nil {
			panic(err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				panic(err)
			}
		}()
		_, err = fmt.Fprintf(file, "%s", string(out))
		if err != nil {
			panic(err)
		}
	}
}

// CleanPrivateDir delete decrypte ca private key and the file containing his password
func CleanPrivateDir() {
	dirname := "cert/pki/private/"
	files := []string{"ca.key"}

	for _, file := range files {
		err := os.Remove(dirname + file)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// GetUsernameFile returns a list of all CN file configuration for a specific extension: .req or .crt
func GetUsernameFile(path, extension string) []string {
	dir, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	files, err := dir.Readdirnames(-1)
	if err != nil {
		log.Fatal(err)
	}
	return files
}
