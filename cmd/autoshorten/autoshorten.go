package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

var (
	content   string
	shortlink string
)

func main() {
	fmt.Println("Scanning for copied URLs.")
	for {
		content, err := clipboard.ReadAll()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if strings.HasPrefix(content, "http") && len(content) > 30 && content != shortlink && !strings.Contains(content, "\n") {
			fmt.Println("URL Copied: " + content)
			shortlink, err = shorten(content)
			if err != nil {
				shortlink = ""
				fmt.Println(err)
			}
			fmt.Println("Shortlink produced: " + shortlink)
			err = clipboard.WriteAll(shortlink)
			if err != nil {
				shortlink = ""
				fmt.Println(err)
			}
		}
		time.Sleep(time.Second)
	}
}

func shorten(link string) (string, error) {
	response, err := http.Get("https://tinyurl.com/api-create.php?url=" + link)

	if err != nil {
		return "", err
	}

	output, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	response.Body.Close()
	return string(output), err
}
