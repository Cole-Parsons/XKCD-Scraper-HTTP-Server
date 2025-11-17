# *XKCD Downloader & Server*

A Go-based XKCD Scraper that evolved from a simple downloader into a full-featured tool, featuring a command-line interface, multithreaded downloads, and an HTTP server with a REST API. Built incrementally, it demonstrates staged development, refactoring, and feature extension.

---

## Server Precompiled Quick Run (Windows)  
1. [Download build v4](https://github.com/Cole-Parsons/XKCD-Scraper-HTTP-Server/blob/main/xkcd-server4.exe) (Other builds available at bottom of readme)  
2. cd to where the build is located  
3. Run `.\xkcd-server4.exe`  

---

# How to Build/Run From Source  

### Clone Repository  
```bash
git clone https://github.com/Cole-Parsons/XKCD-Scraper-HTTP-Server.git
cd XKCD-Downloader
```
## Run and Build the XKCD Server  
```bash
go build -o xkcd_server.exe
.\xkcd_server.exe
```

---

## Version 1–4 | Go Project  
A multi-version Go project that evolves from a simple comic downloader into a full REST API server for XKCD comics.  
This project demonstrates progressive software development in Go — including CLI design, unit testing, concurrency, and HTTP server implementation.  

---

### XKCD Server (v4)  
This version turns the project into an HTTP-based REST API  

*REST Endpoints*
```text
GET    /comic/{id}        Returns JSON about whether a comic is downloaded  
POST   /comic/{id}        Requests the server to download that comic  
GET    /download/{id}     Returns the comic file if it exists, else 404  
```

---

### Multithreading (v3)  
Uses Goroutines to download multiple comics concurrently.  

*CLI FLag*  
`--threads=3`          Number of comics to download at once (default: 3) 

---

### CLI Interface (v2)  
Adds a full command-line interface (CLI) for better user control.  

**CLI Flags:**  
```text
--version            Prints program version  
--parser=regex/html  Choose how to parse XKCD HTML pages (Default uses JSON)  
--download-all       Continue downloading all comics even if some are already downloaded  
(default behavior)     Uses JSON parsing, stops when an already downloaded comic is found
```
---

### Initial Downloader (v1)  
A Go program that downloads all comics from [XKCD](https://xkcd.com/).  
When rerun, it automatically skips comics that have already been downloaded — ensuring efficient updates.  

---

## Server Builds  
[Version 4](https://github.com/Cole-Parsons/XKCD-Scraper-HTTP-Server/blob/main/xkcd-server4.exe)  
[Version 3](https://github.com/Cole-Parsons/XKCD-Scraper-HTTP-Server/blob/main/xkcd-server3.exe)  
[Version 2](https://github.com/Cole-Parsons/XKCD-Scraper-HTTP-Server/blob/main/xkcd-server2.exe)  
[Version 1](https://github.com/Cole-Parsons/XKCD-Scraper-HTTP-Server/blob/main/xkcd-server1.exe)  

## Server Source
[Version 4](https://github.com/Cole-Parsons/XKCD-Scraper-HTTP-Server/blob/main/xkcd-server4.exe)  
[Version 3](https://github.com/Cole-Parsons/XKCD-Scraper-HTTP-Server/blob/54ddce8277f81a37309409cde733578b4c5b9094/Project1.go)
[Version 2](https://github.com/Cole-Parsons/XKCD-Scraper-HTTP-Server/blob/5ca8b19af337ff2714c8d6424c76857a2cb06869/Project1.go)
[Version 1](https://github.com/Cole-Parsons/XKCD-Scraper-HTTP-Server/blob/3ab881d18802de36a7789b4f456c0d5825c49ef4/Project1.go)

## Cross-Platform Support    
Test on Windows, MacOS, and Linux using Virtual Box  

## Related Projects  
[XKCD-Client](https://github.com/Cole-Parsons/XKCD-Client.git)  
[XKCD-Frontend](https://github.com/Cole-Parsons/XKCD-frontend.git)  

## To run the entirety of the Scraper see:  
[Docker-compose](https://github.com/Cole-Parsons/XKCD-docker-compose.git)



