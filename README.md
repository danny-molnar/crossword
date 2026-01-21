Crossword Engine

A Go-based backend engine for cryptic and quick (“natural”) crosswords, providing a validated puzzle model, solver helpers, and a minimal HTTP API.

This project focuses on correctness, clean domain modelling, and extensibility rather than UI or presentation (for now).

Status: Early development – core engine functional (v0.2.0)

Features

Core domain

Arbitrary grid sizes (tested 5x5 up to 15x15)

Grid to entry detection and numbering (across / down)

Strict puzzle, grid, entry, clue, and enumeration validation

Enumeration parsing (3, 3,5, 4-4, etc.)

Solver helpers

Wordlist loader

Exact anagram helper

Pattern matcher (e.g. TR?C?)

API

Public puzzle view (solutions and answers stripped)

Anonymous solve sessions

Session state persistence (in-memory)

Helper endpoints (anagram / pattern)

Project layout

cmd/
api/ API entrypoint

internal/
api/
handlers/ HTTP handlers
router.go chi router

domain/ core crossword domain model and validation
store/ in-memory stores (puzzles, sessions)
tools/ wordlist, anagram, pattern helpers
util/ shared utilities (ULID IDs)

wordlists/
english.txt sample wordlist

Running locally

Prerequisites

Go 1.22 or newer

Run tests
go test ./...

Start the API
go run ./cmd/api

The server will start on
http://localhost:8080

Example endpoints

Health check
curl http://localhost:8080/v1/health

Fetch a puzzle (public view)
curl http://localhost:8080/v1/puzzles/puz_demo

Create a solve session
curl -X POST http://localhost:8080/v1/puzzles/puz_demo/sessions

Anagram helper
curl "http://localhost:8080/v1/tools/anagram?letters=react&len=5
"

Pattern matcher
curl "http://localhost:8080/v1/tools/pattern?pattern=tr?c?&len=5
"

Design principles

Correctness first: invalid puzzles are rejected early

Explicit domain modelling: grid, entries, clues, and sessions are first-class concepts

No solution leakage: solver APIs never expose answers

Incremental extensibility: storage, auth, and UI can be added later without refactoring the core

Roadmap (rough)

Check and reveal endpoints

Creator-side puzzle creation and editing

Persistent storage (Postgres)

Authentication (optional)

Frontend solver UI

License

TBD (MIT or Apache-2.0 likely)

Author

Danny Molnar
