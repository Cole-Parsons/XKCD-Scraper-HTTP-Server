package main

import (
	"encoding/json" //parse JSON data
	"flag"          // for cli inputs
	"fmt"           //print messages to console
	"io"            // reading and writing streams of data
	"net/http"      //make http requests
	"os"            //to interact with file system
	"path/filepath" //builds safe file paths across operating systems
	"strconv"       //converts between strings and numbers
	"strings"       //for string manipulation
	"sync"          //Routine coordination
)

// Global maps and lock for server state
var (
	downloading = make(map[int]bool)
	downloaded  = make(map[int]bool)
	mu          sync.Mutex
)

//REST Handlers

func handleGetComic(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/comic/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid comic number"}`))
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
	enableCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/comic/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid comic number"}`))
		return
	}

	mu.Lock()
	if downloading[id] {
		mu.Unlock()
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error":"comic is already downloading"}`))
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
	w.Write([]byte(`{"message":"Download Started"}`))
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

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// Main Program
func main() {
	versionFlag := flag.Bool("version", false, "Print program version")
	serverFlag := flag.Bool("server", false, "Run as HTTP server")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Comic downloader v4.0")
		return
	}

	folder := "comics"
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		fmt.Println("Error making folder")
		return
	}

	os.MkdirAll("comics", os.ModePerm)
	initDownloadedMap() //populates download map with comics that exists in the folder

	// Start server mode
	if *serverFlag {
		http.HandleFunc("/comic/", func(w http.ResponseWriter, r *http.Request) {
			enableCORS(w) // MUST be first

			if r.Method == http.MethodOptions { // Preflight
				w.WriteHeader(http.StatusOK)
				return
			}

			idStr := strings.TrimPrefix(r.URL.Path, "/comic/")
			_, err := strconv.Atoi(idStr)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"invalid comic number"}`))
				return
			}

			switch r.Method {
			case http.MethodGet:
				handleGetComic(w, r)
			case http.MethodPost:
				handlePostComic(w, r)
			default:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte(`{"error":"method not allowed"}`))
			}
		})
		http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
			enableCORS(w)
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			handleDownload(w, r)
		})
		fmt.Println("Server running on http://localhost:8080")
		http.ListenAndServe(":8080", nil)
		return
	}
	fmt.Println("Server is online")
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

func fetchComic(num int, parser string) (*Comic, error) {
	switch parser {
	case "", "json":
		return getComic(num)
	default:
		return nil, fmt.Errorf("invalid parser type: %s", parser)
	}
}

func initDownloadedMap() {
	files, _ := os.ReadDir("comics")
	for _, f := range files {
		name := f.Name()

		if f.IsDir() {
			continue
		}

		parts := strings.SplitN(name, "-", 2)
		if len(parts) == 0 {
			continue
		}

		id, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		downloaded[id] = true
		fmt.Println("comic exists: ", id)
	}
}
