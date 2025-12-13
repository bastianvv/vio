# VIO

VIO is a self-hosted video library server. Think of it as a Navidrome for video media.

It currently focuses on library indexing, metadata modeling, and filesystem correctness, not playback or transcoding.

## Status
- Early development
- Backend scanning is functional
- API is unstable
- UI not implemented yet

## Features
- Incremental filesystem scanning
- Missing media detection
- Series / season / episode modeling
- SQLite backend

## Non-goals (for now)
- Transcoding
- User management
- Streaming UI

## Requirements
- Go 1.22+

## Running
- Clone the repository
- `go build ./...`
- `go build ./cmd/server`
- `./server`

The API will be available at: http://localhost:8080

## Dev Note
This project is currently being developed using a mix between traditional programming and AI-assisted development (ChatGPT).
I work in QA and do not have a formal background in backend or functional programming, so this project is as much about learning as it is about building something useful, with a strong focus on correctness, edge cases, and repeatable testing.

API Testing is done using [Bruno](https://www.usebruno.com/)
Developed using [Zed](https://zed.dev/)
