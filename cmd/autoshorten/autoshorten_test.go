package main

import (
	"strings"
	"testing"
)

func TestRequest(t *testing.T) {
	testURL := "https://dev.to/chand1012/transisioning-from-python-to-golang-and-why-python-programmers-should-consider-it-4pc9"
	shortURL, err := shorten(testURL)

	if err != nil {
		t.Errorf("There was an error testing the shorten function: %v", err)
	}

	if !strings.Contains(shortURL, "tinyurl.com") {
		t.Errorf("Was expecting the string to contain 'tinyurl.com', but it didn't. Returned string: %s", shortURL)
	}
}

func TestGitIORequest(t *testing.T) {
	testURL := "https://github.com/chand1012/chand1012/discussions"
	shortURL, err := gitIOShorten(testURL)

	if err != nil {
		t.Errorf("There was an error testing the shorten function: %v", err)
	}

	if !strings.Contains(shortURL, "git.io") {
		t.Errorf("Was expecting the string to contain 'tinyurl.com', but it didn't. Returned string: %s", shortURL)
	}
}
