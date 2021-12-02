package utility

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func GetHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func DivideBodyIntoParagraphs(body string) []string {
	var result []string
	var paragraph []byte

	base64Body, err := base64.StdEncoding.DecodeString(body)
	CheckErr(err)

	for i := 0; i < len(base64Body); i++ {
		paragraph = append(paragraph, base64Body[i])

		if base64Body[i] == 13 {
			result = append(result, string(paragraph))
			i = i + 2
			paragraph = make([]byte, 0)
		}

		if len(paragraph) != 0 && i == len(base64Body)-1 {
			result = append(result, string(paragraph))
		}
	}

	return result
}

func SavePostImage(r *http.Request) string {
	var path string

	in, header, err := r.FormFile("myImage")
	imageData := strings.Split(header.Filename, ".")
	if err != nil {
		log.Println(err)
	}
	defer in.Close()

	randBytes := make([]byte, 16)
	rand.Read(randBytes)

	path = "images/" + hex.EncodeToString(randBytes) + "." + imageData[1]

	out, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}
	defer out.Close()
	io.Copy(out, in)

	return path
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
