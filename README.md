# link-lint

A small Go CLI that crawls a website starting from a URL you provide, follows linked pages, and surfaces **dead** or **problematic** links by testing each page over HTTP.

## What it does

- Starts at the URL you pass to `scrape`.
- Parses HTML, collects `<a href>` links (absolute `http`/`https` URLs and same-site paths resolved against the starting URL).
- Fetches each discovered URL concurrently (goroutines) with a **5 second** timeout per request.
- **Tests each page** and classifies the response:
  - **Dead page** — HTTP status **greater than 400** (e.g. 404, 500). Printed in red in the form: `https://example.com/path is dead page`.
  - **Redirect** — status **between 300 and 399**. Printed as: `https://example.com/path is a redirected page`.
- Network or read errors are printed in red with the URL and error details.
- When the run finishes, it prints how long scraping took.

## Requirements

- [Go](https://go.dev/dl/) **1.25** or compatible toolchain (see `go.mod`).

## Install / build

From the repository root:

```bash
go build -o link-lint .
```

Or run without installing a binary:

```bash
go run . scrape https://example.com
```

## Usage

```text
link-lint scrape <url>
```

**Example:**

```bash
./link-lint scrape https://example.com
```

`scrape` expects exactly one argument: the starting URL (include `https://` or `http://` as appropriate).

## Dependencies

- [github.com/spf13/cobra](https://github.com/spf13/cobra) — CLI structure
- [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html) — HTML parsing

## Project layout

```text
.
├── main.go              # entrypoint
├── cmd/link_lint/       # Cobra commands (root + scrape)
└── pkg/link_lint/       # crawl, parse, and dead-page checks
```

## Notes

- Output uses ANSI colors in the terminal for dead-link messages.
- This tool is intended for **personal or light crawling**; respect `robots.txt`, rate limits, and terms of service for any site you scan.
