package main

import (
	"encoding/json" //parse JSON data
	"fmt"           //print messages to console
	"io"            // reading and writing streams of data
	"net/http"      //make http requests
	"os"            //to interact with file system
)

func main() {
	Comic, err := getComic(614) //calls get comic
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Comic #: ", Comic.Num)
	fmt.Println("Comic Title: ", Comic.Title)
	fmt.Println("Comic url: ", Comic.Img)
	fmt.Println("Alt text: ", Comic.Alt)

	err = downloadImage(Comic.Img, fmt.Sprintf("%d-%s.png", Comic.Num, Comic.Title))

	if err != nil {
		fmt.Println("Error downloading Comic: ", err)
		return
	}
	fmt.Println("Comic saved Successfully")

}

// collection of related data grouped together
type Comic struct {
	Num   int    `json:"num`
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

	var comic Comic //declaring variable to store parsed JSON data

	err = json.NewDecoder(resp.Body).Decode(&comic) //reads JSON data from http, fills comic struct with that data, directly updating it through pointer
} //end get Comic

func downloadImage(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()                   //closes http connection to prevent memory leak
	file, err := os.Create("ScraperTest.png") //creates new file
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
