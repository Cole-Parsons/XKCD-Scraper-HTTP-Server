package main

import (
	"os"
	"testing"
)

// test sanitation
func TestSanitizeTitle(t *testing.T) {
	input := `Hello /:*?"<>| World`
	expected := "Hello_World"

	result := sanitizeTitle(input)

	if result != expected {
		t.Errorf("sanitizeTitle failed: gave %s, expected %s", result, expected)
	}
}

func TestGetComic(t *testing.T) {
	comic, err := getComic(614)
	if err != nil {
		t.Fatalf("getComic failed: %v", err)
	}
	if comic.Num != 614 {
		t.Errorf("Expected comic number 614, got %d", comic.Num)
	}
	if comic.Title == "" {
		t.Errorf("Expected comic title, got empty string")
	}
}

func TestGetLastComicNum(t *testing.T) {
	num, err := getLastComicNum()
	if err != nil {
		t.Fatalf("getLastComicNum failed: %v", err)
	}
	if num < 614 {
		t.Errorf("Expected latest comic >= 614, got %d", num)
	}
}

func TestDownloadImage(t *testing.T) {
	url := "https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg"
	filename := "test_image.png"

	err := downloadImage(url, filename)
	if err != nil {
		t.Errorf("File was not created: %s", filename)
	}

	os.Remove(filename)
}
