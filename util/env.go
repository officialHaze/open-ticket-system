package util

import (
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	return godotenv.Load(UseEnv())
}

func UseEnv() string {
	filepath := path.Join("envfile.txt")

	f, err := os.Open(filepath)
	if err != nil {
		log.Printf("error opening %s: %v", filepath, err)
		return ""
	}

	data, err := io.ReadAll(f)
	if err != nil {
		log.Printf("error reading %s: %v", filepath, err)
	}

	str := string(data)
	return strings.TrimSpace(str) // trim any trailing white spaces (including new line)
}
