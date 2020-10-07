package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
)

var (
	content   string
	shortlink string
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	iconData, err := ioutil.ReadFile("icon.ico")
	if err != nil {
		panic(err)
	}

	_, err = lockFileCreate()
	if err != nil {
		panic(err)
	}

	systray.SetTemplateIcon(iconData, iconData)
	systray.SetTitle("Autoshorten")
	systray.SetTooltip("Autoshorten Settings")
	trayQuit := systray.AddMenuItem("Quit", "Quit the application.")
	go quitRoutine(trayQuit)
	go mainRoutine(trayQuit)
}

func onExit() {
	clipboard.WriteAll("")
	os.Remove("./thread.lock")
}

func quitRoutine(quit *systray.MenuItem) {
	<-quit.ClickedCh
	systray.Quit()
}

func mainRoutine(quit *systray.MenuItem) {
	fmt.Println("Scanning for copied URLs.")
	for {
		fmt.Println("Scanning....")
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

func lockFileEqu(input []byte) (bool, error) {
	data, err := ioutil.ReadFile("./thread.lock")
	if err != nil {
		return false, err
	}
	if bytes.Compare(input, data) == 0 {
		return true, nil
	}
	return false, nil
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
