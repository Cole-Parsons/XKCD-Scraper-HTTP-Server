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
`GET`    `/comic/{id}`      Returns JSON about whether a comic is downloaded  
`POST`   `/comic/{id}`      Requests the server to download that comic  
`GET`    `/download/{id}`   Returns the comic file if it exists, else 404  

---

## Cross-Platform Support    
Test on Windows, MacOS, and Linux using Virtual Box  
