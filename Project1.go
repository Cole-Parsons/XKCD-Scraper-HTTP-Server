package main

import (
	"encoding/json" //parse JSON data
	"flag"          // for cli inputs
	"fmt"           //print messages to console
	"io"            // reading and writing streams of data
	"net/http"      //make http requests
	"net/url"
	"os" //to interact with file system
	"regexp"
	"strings" //for string manipulation

	"golang.org/x/net/html" // html parse
)

func main() {

	versionFlag := flag.Bool("version" ,false, "Print program version")
	parserFlag := flag.String("parser" ,json, " Choose parsing method, html or regex")
	downloadAllFlag := flag.Bool("download-all", false, "download all the comics including ones already downloaded")
	flag.Parse()

	if *versionFlag{
		fmt.Println("Comic downloader v2.0")
		return
	}
	
	fmt.Println("Parser Method: ", *parserFlag)

	if *downloadAllFlag {
		fmt.Println("Downloading all comics even if they exist")
	}else {
		fmt.Println("Stopping when comic is already downloaded")
	}


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

		Comic, err := fetchComic(i, *parserFlag)
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
			if !*downloadAllFlag {
				fmt.Println("Stopping because download-all not set")
				break
			}
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

//recursive html
func getComicHTML (num int) (*Comic, error) {
	url := fmt.Sprintf("https://xkcd.com/%d/", num)
	resp, err := http.Get(url)
	if err != nil{
		return nil, err
	}

	defer resp.Body.Close()
	
	//Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil{
		return nil, err
	}
	var comicImg, altText, titleText string
	//recursive function to traverse HTML nodes
	//looking for div#comic
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, attr := n.Attr{
				if attr.Key == "id" && attr.Val == "comic" {
					//found div#comic, now looking for img tag
					for c := n.FirstChild; c != nil; c = c.NextSibling{
						if c.Type == html.ElementNode && c.Data == "img" {
							for _, imgAttr := range c.Attr {
								if imgAttr.Key == "src" {
									comicImg = "https:" + imgAttr.Val
								}
								if imgAttr.Key == "title" {
									altText = imgAttr.Val
								}
								if imgAttr.Key == "alt" {
									titleText = imgAttr.Val
								}
								for c := n.FirstChild; c != nil; c = c.NextSibling {
									traverse(c)
								}
								if comicImg == "" {
									return nil, fmt.Errorf("comic not found")
								}
							}
						}
					}
				}
			}
		}

		//recurse into child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc) //start recussion from root

	if comicImg == "" {
		return nil, fmt.Errorf("comic not found")
	}

	return &Comic{
		Num: num,
		Title: titleText,
		Img: comicImg,
		Alt: altText,
	}, nil
}

//regex html
func getComicRegex (num int) (*Comic, error) {
	resp, err := http.Get(fmt.Sprintf("https://xkcd.com/%d/", num))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	htmlContent := string(bodyBytes)

	//applying regex
	re := regexp.MustCompile(`<div id="comic">.*?<img src="(.*?)".*?title="(.*?)".*?alt="(.*?)".*?>`)
	matches := re.FindStringSubmatch(htmlContent)

	if len(matches) < 4 {
		return nil, fmt.Errorf("comic not found")
	}

	return &Comic {
		Num: num, 
		Title: matches[3], //alt
		Img: "https:" + matches[1],
		Alt: matches[2], //title
	}, nil
}

//decides what method to use
func fetchComic(num int, parser string) (*Comic, error) {
	switch parser {
    case "", "json":
        return getComic(num)
    case "html":
        return getComicHTML(num)
    case "regex":
        return getComicRegex(num)
    default:
        return nil, fmt.Errorf("invalid parser type: %s", parser)
	}
}