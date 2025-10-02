package main

import (
	"encoding/json" //parse JSON data
	"fmt"           //print messages to console
	"io"            // reading and writing streams of data
	"net/http"      //make http requests
	"os"            //to interact with file system
)

func main() {
	x := 0 //first comic
	y := 0 //last comic

	fmt.Println("Enter the range of comics you want. Ex. comic 10 to 25 inclusive")
	fmt.Print("First Comic you want: ")
	fmt.Scan(&x)
	fmt.Println()
	fmt.Print("Enter the last comic you want: ")
	fmt.Scan(&y)
	fmt.Println()

	for i := x; i <= y; i++ {
		Comic, err := getComic(i) //calls get comic ; gets comic JSON
		if err != nil {
			fmt.Println("Error fetching comic:", err)
			return
		}

		//prints all info about the comic
		fmt.Println("Comic #: ", Comic.Num)
		fmt.Println("Comic Title: ", Comic.Title)
		fmt.Println("Comic url: ", Comic.Img)
		fmt.Println("Alt text: ", Comic.Alt)

		//dynamically creates file name for individual comic
		filename := fmt.Sprintf("%d-%s.png", Comic.Num, Comic.Title)

		//downloads comic to disk using dynamic file name
		err = downloadImage(Comic.Img, filename)

		//confirms if image saved correctly
		if err != nil {
			fmt.Println("Error downloading Comic: ", err)
			return
		}
		fmt.Println("Comic saved Successfully")
	}
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
