package main

import(
	"encoding/json" //parse JSON data
	"fmt" //print messages to console
	"io" // reading and writing streams of data
	"net/http" //make http requests
	"os" //to interact with file system
)

//collection of related data grouped together 
type Comic Struct {
	Num int 'json:"num'
	Title string 'json:"title"'
	Img string 'json:"img"'
	Alt string 'json:"alt"'
}

func getComic (num int) (*Comic, error) {
	num := 1
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", num) //dynamically storing the string for comic
	resp, err := https.Get(url)
	
	//checks for error
	if err != nil {
		return nil, err
	}
	
	defer resp.Body.close() //closes http response body when the function is done

	var comic Comic //declaring variable to store parsed JSON data

	err = json.NewDecoder(resp.body).Decode(&Comic) //reads JSON data from http, fills comic struct with that data, directly updating it through pointer




}//end get Comic