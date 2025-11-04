# XKCD Scraper HTTP Server

A Golang-based HTTP server that scrapes XKCD comics, serves them via a REST API, and supports multi-threaded downloads using Goroutines.

## Features

- Downloads XKCD comics and stores them locally.
- REST API endpoints:
  - `GET /comic/{id}` – Check if a comic is downloaded or currently downloading.
  - `POST /comic/{id}` – Request a comic to be downloaded.
  - `GET /download/{id}` – Download the comic image if available.
- Multi-threaded downloads using Goroutines.
- Supports CLI flags for version info and server mode.

## Requirements

- Go 1.25+
- Docker (optional, for containerization)

## Setup & Run

### Locally

1. Clone the repository:
   ```bash
   git clone https://github.com/Surfs-Up5/XKCD-Scraper-HTTP-Server.git
   cd Project-1
   go run Project1.go
