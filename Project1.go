package main

import (
	"encoding/json" //parse JSON data
	"fmt"           //print messages to console
	"io"            // reading and writing streams of data
	"net/http"      //make http requests
	"os"            //to interact with file system
	"strings"       //for string manipulation
)

func main() {

	folder := "comics"
	fmt.Println("Downloading all comic .pngs up to the most recent")

	err := os.MkdirAll(folder, os.ModePerm) //creates new folder for comics
	if err != nil {
		fmt.Println("Error making folder")
		return
	}

	lastNum, err := getLastComicNum()
	if err != nil {
		fmt.Println("Error getting latest comic")
		return
	}

	for i := 1; i <= lastNum; i++ {

		Comic, err := getComic(i) //calls get comic ; gets comic JSON
		if err != nil {
			fmt.Println("Skipping comic ", i, ":", err)
			continue
		}

		//prints all info about the comic
		// fmt.Println("Comic #: ", Comic.Num)
		// fmt.Println("Comic Title: ", Comic.Title)
		// fmt.Println("Comic url: ", Comic.Img)
		// fmt.Println("Alt text: ", Comic.Alt)

		//dynamically creates file name for individual comic
		safeTitle := sanitizeTitle(Comic.Title)
		filename := fmt.Sprintf("%s/%d-%s.png", folder, Comic.Num, safeTitle)

		//checks if file already exists
		if _, err := os.Stat(filename); err == nil {
			fmt.Println("File already exists, skipping: ", filename)
			continue
		}

		//downloads comic to disk using dynamic file name
		err = downloadImage(Comic.Img, filename)

		//confirms if image saved correctly
		if err != nil {
			fmt.Println("Error downloading Comic: ", err)
			return
		}
		//fmt.Println("Comic saved Successfully: ", filename)
	}
	fmt.Println("All comics downloaded successfully")
}

// collection of related data grouped together
type Comic struct {
	Num   int    `json:"num"`
	Title string `json:"title"`
	Img   string `json:"img"` //img url
	Alt   string `json:"alt"`
}

func getComic(num int) (*Comic, error) {
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", num) //dynamically storing the string for comic
	resp, err := http.Get(url)
	//checks for error
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() //closes http response body when the function is done

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("recieved status %d", resp.StatusCode)
	}

	var comic Comic //declaring variable to store parsed JSON data

	err = json.NewDecoder(resp.Body).Decode(&comic) //reads JSON data from http, fills comic struct with that data, directly updating it through pointer
	if err != nil {
		return nil, err
	}
	return &comic, nil
} //end get Comic

func downloadImage(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()          //closes http connection to prevent memory leak
	file, err := os.Create(filename) //creates new file
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body) //copies all data from resp.Body into file
	if err != nil {
		return err
	}
	return nil
}

func getLastComicNum() (int, error) {
	resp, err := http.Get("https://xkcd.com/info.0.json")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close() //clean up network connection

	var latest Comic

	err = json.NewDecoder(resp.Body).Decode(&latest) // parse json into struct
	if err != nil {
		return 0, err
	}

	return latest.Num, nil //send comic number back
}

func sanitizeTitle(title string) string {
	safeTitle := strings.ReplaceAll(title, " ", "_")
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, c := range invalidChars {
		safeTitle = strings.ReplaceAll(safeTitle, c, "")
	}
	return safeTitle
}
