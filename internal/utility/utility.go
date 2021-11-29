package utility

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"mime/multipart"

	"golang.org/x/crypto/bcrypt"
)

func GetHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func CreateImageContainer(file *multipart.FileHeader) string {
	imageByteContainer := make([]byte, (1024 * 1024 * 2))
	fileContent, err := file.Open()

	imageByteContainer, err = ioutil.ReadAll(fileContent)
	if err != nil {
		panic(err)
	}

	fileContent.Close()

	return base64.StdEncoding.EncodeToString(imageByteContainer)
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

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
