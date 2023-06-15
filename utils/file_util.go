package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func File_get_contents(filepath string) string {
	fileContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string
	text := string(fileContent)
	return text
}
func File_put_contents(filepath string, content string) {
	data := []byte(content)

	err := ioutil.WriteFile(filepath, data, 0)

	if err != nil {
		log.Fatal(err)
	}
}
func Filesize(filename string) int64 {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	return fileInfo.Size()
}
func File_exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
func Md5_file(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
