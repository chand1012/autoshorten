package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/chand1012/autoshorten/icon"
	"github.com/getlantern/systray"
	"github.com/pkg/browser"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {

	_, err := lockFileCreate()
	if err != nil {
		panic(err)
	}

	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("Autoshorten")
	systray.SetTooltip("Autoshorten")
	openBrowser := systray.AddMenuItem("Homepage", "Go to application homepage.")
	tinyURL := systray.AddMenuItem("TinyURL", "TinyURL Homepage.")
	gitIO := systray.AddMenuItem("Git.io", "Got to Git.io.")
	systray.AddSeparator()
	trayQuit := systray.AddMenuItem("Quit", "Quit the application.")

	go quitRoutine(trayQuit)
	go mainRoutine(trayQuit)
	go homepageRoutine(openBrowser)
	go tinyURLRoutine(tinyURL)
	go gitIORoutine(gitIO)
}

func onExit() {
	clipboard.WriteAll("")
	os.Remove("./thread.lock")
}

func quitRoutine(quit *systray.MenuItem) {
	<-quit.ClickedCh
	systray.Quit()
}

func homepageRoutine(button *systray.MenuItem) {
	for {
		<-button.ClickedCh
		fmt.Println("Opening Homepage.")
		browser.OpenURL("https://github.com/chand1012/autoshorten")
	}
}

func tinyURLRoutine(button *systray.MenuItem) {
	for {
		<-button.ClickedCh
		fmt.Println("Opening TinyURL.")
		browser.OpenURL("https://tinyurl.com/")
	}
}

func gitIORoutine(button *systray.MenuItem) {
	for {
		<-button.ClickedCh
		fmt.Println("Opening Git.io.")
		browser.OpenURL("https://git.io/")
	}
}

func mainRoutine(quit *systray.MenuItem) {
	fmt.Println("Scanning for copied URLs.")
	var content string
	var shortlink string
	var err error
	for {
		fmt.Println("Scanning....")
		content, err = clipboard.ReadAll()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if strings.HasPrefix(content, "http") && len(content) > 30 && content != shortlink && !strings.Contains(content, "\n") {
			var shortlink string
			fmt.Println("URL Copied: " + content)
			if strings.Contains(content, "github.com") || strings.Contains(content, "github.io") || strings.Contains(content, "githubusercontent.com") {
				shortlink, err = gitIOShorten(content)
			} else {
				shortlink, err = shorten(content)
			}
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
		time.Sleep(time.Second / 2)
		if !lockFileExists() {
			break
		}
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

func gitIOShorten(link string) (string, error) {
	data := url.Values{
		"url": []string{link},
	}

	resp, err := http.PostForm("https://git.io/create", data)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	endpoint := string(body)

	return "https://git.io/" + endpoint, nil

}

func lockFileExists() bool {
	info, err := os.Stat("./thread.lock")
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func lockFileCreate() ([]byte, error) {
	fileData := make([]byte, 8)
	rand.Read(fileData)
	err := ioutil.WriteFile("./thread.lock", fileData, 0644)
	return fileData, err
}
