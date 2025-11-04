# *XKCD Downloader & Server*

## Version 1–4 | Go Project  
A multi-version Go project that evolves from a simple comic downloader into a full REST API server for XKCD comics.  
This project demonstrates progressive software development in Go — including CLI design, unit testing, concurrency, and HTTP server implementation.  

---

### Initial Downloader (v1)  
A Go program that downloads all comics from [XKCD](https://xkcd.com/).  
When rerun, it automatically skips comics that have already been downloaded — ensuring efficient updates.  

---

### CLI Interface (v2)  
Adds a full command-line interface (CLI) for better user control.  

**CLI Flags:**  
```text
`--version `           Prints program version  
`--parser=regex/html`  Choose how to parse XKCD HTML pages (Default uses JSON)  
`--download-all`       Continue downloading all comics even if some are already downloaded  
(default behavior)     Uses JSON parsing, stops when an already downloaded comic is found
```
---

### Multithreading (v3)  
Uses Goroutines to download multiple comics concurrently.  

**New CLI FLag*  
`--threads=3`          Number of comics to download at once (default: 3)  

---

### XKCD Server (v4)  
Renamed to XKCD Server, this version turns the project into an HTTP-based REST API  

*REST Endpoints*
```text
`GET`    `/comic/{id}`      Returns JSON about whether a comic is downloaded  
`POST`   `/comic/{id}`      Requests the server to download that comic  
`GET`    `/download/{id}`   Returns the comic file if it exists, else 404  
```

---

# How to run

### Clone Repository
```bash
git clone https://github.com/Cole-Parsons/XKCD-Scraper-HTTP-Server.git
cd XKCD-Downloader
```

## Run the Downloader (v1-v3)
### Default run  
`go run Project1.go`

### Using CLI Flags
`go run Project1.go --version`
`go run Project1.go --parser=regex/html`
`go run Project1.go --download-all`
`go run Project1.go --threads=5`

## Run the XKCD Server (v4)
```bash
go build -o xkcd_server Project1.go
./xkcd_server -server
```

### Example server use
```bash
curl -X GET http://localhost:8080/comic/614
curl -X POST http://localhost:8080/comic/614
curl -O http://localhost:8080/download/614
```

## Cross-Platform Support    
Test on Windows, MacOS, and Linux using Virtual Box  


