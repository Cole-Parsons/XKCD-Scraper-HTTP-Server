package main

import (
	"encoding/json" //parse JSON data
	"flag"          // for cli inputs
	"fmt"           //print messages to console
	"io"            // reading and writing streams of data
	"net/http"      //make http requests
	"os"            //to interact with file system
	"path/filepath" //builds safe file paths across operating systems
	"regexp"
	"strconv" //converts between strings and numbers
	"strings" //for string manipulation
	"sync"    //Routine coordination

	"golang.org/x/net/html" // html parser
)

// Global maps and lock for server state
var (
	downloading = make(map[int]bool)
	downloaded  = make(map[int]bool)
	mu          sync.Mutex
)

//REST Handlers

func handleGetComic(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/comic/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid comic number", http.StatusBadRequest)
		return
	}
	mu.Lock()
	defer mu.Unlock()

	status := map[string]bool{
		"downloaded":    downloaded[id],
		"isDownloading": downloading[id],
	}
	json.NewEncoder(w).Encode(status)
}

func handlePostComic(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/comic/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid comic number", http.StatusBadRequest)
		return
	}

	mu.Lock()
	if downloading[id] {
		mu.Unlock()
		http.Error(w, "comic is already downloading", http.StatusConflict)
		return
	}
	downloading[id] = true
	mu.Unlock()

	go func() {
		defer func() {
			mu.Lock()
			downloading[id] = false
			downloaded[id] = true
			mu.Unlock()
		}()

		comic, err := fetchComic(id, "json")
		if err != nil {
			fmt.Println("Error downloading comic:", err)
			return
		}

		safeTitle := sanitizeTitle(comic.Title)
		filename := filepath.Join("comics", fmt.Sprintf("%d-%s.png", comic.Num, safeTitle))

		if err := downloadImage(comic.Img, filename); err != nil {
			fmt.Println("Download failed:", err)
		} else {
			fmt.Println("Saved:", filename)
		}
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Download Started"))
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/download/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid comic number", http.StatusBadRequest)
		return
	}

	files, _ := os.ReadDir("comics")
	for _, f := range files {
		if strings.HasPrefix(f.Name(), fmt.Sprintf("%d-", id)) {
			http.ServeFile(w, r, filepath.Join("comics", f.Name()))
			return
		}
	}
	http.NotFound(w, r)
}

// Main Program
func main() {
	versionFlag := flag.Bool("version", false, "Print program version")
	parserFlag := flag.String("parser", "json", "Choose parsing method, html or regex")
	downloadAllFlag := flag.Bool("download-all", false, "download all the comics including ones already downloaded")
	threadsFlag := flag.Int("threads", 3, "choose how many goroutines for the run")
	serverFlag := flag.Bool("server", false, "Run as HTTP server")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Comic downloader v3.0")
		return
	}

	fmt.Println("Parser Method:", *parserFlag)

	if *downloadAllFlag {
		fmt.Println("Downloading all comics even if they exist")
	} else {
		fmt.Println("Stopping when comic is already downloaded")
	}

	folder := "comics"
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		fmt.Println("Error making folder")
		return
	}

	// Start server mode
	if *serverFlag {
		http.HandleFunc("/comic/", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handleGetComic(w, r)
			case http.MethodPost:
				handlePostComic(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})
		http.HandleFunc("/download/", handleDownload)
		fmt.Println("Server running on http://localhost:8080")
		http.ListenAndServe(":8080", nil)
		return
	}

	lastNum, err := getLastComicNum()
	if err != nil {
		fmt.Println("Error getting latest comic")
		return
	}

	var wg sync.WaitGroup
	comicChan := make(chan int)

	for t := 0; t < *threadsFlag; t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range comicChan {
				comic, err := fetchComic(i, *parserFlag)
				if err != nil {
					fmt.Println("Skipping comic", i, ":", err)
					continue
				}
				safeTitle := sanitizeTitle(comic.Title)
				filename := fmt.Sprintf("%s/%d-%s.png", folder, comic.Num, safeTitle)

				if _, err := os.Stat(filename); err == nil {
					if !*downloadAllFlag {
						fmt.Println("File already Exists:", filename)
						fmt.Printf("Ending downloader because comic %s is already downloaded\n", filename)
						return
					}
				}

				err = downloadImage(comic.Img, filename)
				if err != nil {
					fmt.Println("Error downloading Comic: ", err)
					continue
				}
				fmt.Println("Saved: ", filename)
			}
		}()
	}

	for i := 1; i <= lastNum; i++ {
		comicChan <- i
	}
	close(comicChan)

	wg.Wait()
	fmt.Println("Program finished running")
}

//Supporting Functions

type Comic struct {
	Num   int    `json:"num"`
	Title string `json:"title"`
	Img   string `json:"img"`
	Alt   string `json:"alt"`
}

func getComic(num int) (*Comic, error) {
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", num)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received status %d", resp.StatusCode)
	}

	var comic Comic
	err = json.NewDecoder(resp.Body).Decode(&comic)
	if err != nil {
		return nil, err
	}
	return &comic, nil
}

func downloadImage(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

func getLastComicNum() (int, error) {
	resp, err := http.Get("https://xkcd.com/info.0.json")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var latest Comic
	err = json.NewDecoder(resp.Body).Decode(&latest)
	if err != nil {
		return 0, err
	}
	return latest.Num, nil
}

func sanitizeTitle(title string) string {
	safeTitle := strings.ReplaceAll(title, " ", "_")
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, c := range invalidChars {
		safeTitle = strings.ReplaceAll(safeTitle, c, "")
	}
	return safeTitle
}

//HTML + Regex Parsing

func getComicHTML(num int) (*Comic, error) {
	url := fmt.Sprintf("https://xkcd.com/%d/", num)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var comicImg, altText, titleText string
	var traverse func(n *html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == "comic" {
					for c := n.FirstChild; c != nil; c = c.NextSibling {
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
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	if comicImg == "" {
		return nil, fmt.Errorf("comic not found")
	}
	return &Comic{Num: num, Title: titleText, Img: comicImg, Alt: altText}, nil
}

func getComicRegex(num int) (*Comic, error) {
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

	re := regexp.MustCompile(`(?s)<div id="comic">.*?(?:<a [^>]*>)?<img[^>]*src="(//[^"]+)"[^>]*title="(.*?)"[^>]*alt="(.*?)"`)
	matches := re.FindStringSubmatch(htmlContent)

	if len(matches) < 4 {
		return nil, fmt.Errorf("comic not found: %v", matches)
	}
	return &Comic{
		Num:   num,
		Title: matches[3],
		Img:   "https:" + matches[1],
		Alt:   matches[2],
	}, nil
}

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
