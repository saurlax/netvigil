# netvigil

Network Traffic Monitoring and Analysis based Local Threat Intelligence Center

## How to run

To compile this project, you need to install Node.js and Go in advance.

Firstly, build frontend resources

```bash
$ cd frontend
$ npm i
$ npm run build
```

Then, build go project

```bash
$ cd ..
$ go build .
```

In order to run the program correctly, you also need to provide a configuration file. You can directly rename `config.example.toml` to `config.toml` to use the default configuration.

## APIs

| Path       | Method | Request                | Response | Description |
| ---------- | ------ | ---------------------- | -------- | ----------- |
| `/records` | GET    | `?sortBy&page`         | Record[] |             |
| `/login`   | POST   | `{username, password}` | Token    |             |

**Types**

```go
type Record struct {
	ID         string
	LocalAddr  string
	RemoteAddr string
	TIX        string
	Location   string
	Reason     string
	Risk       RiskLevel
	Confidence ConfidenceLevel
}

const (
	Unknown RiskLevel = iota
	Safe
	Normal
	Suspicious
	Malicious
)

const (
	Low ConfidenceLevel = iota
	Medium
	High
)
```

```go
type netstat.SockTabEntry struct {
  ino        string
  LocalAddr  *SockAddr
  RemoteAddr *SockAddr
  State      SkState
  UID        uint32
  Process    *Process
}
```

## Problems & Solutions
* `invalid go version '1.21.6': must match format 1.23`

Upgrade your `go` version to at least `1.21.6`

* `Binary was compiled with 'CGO ENABLED=0', go-sqlite3 requires cgo to work. This is a stub`

Add `CGO_ENABLED=1` to your user environment variable. If env is correctly set, you will see `set CGO_ENABLED=1` with the fllowing command
```bash
$ go env
```

* `cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in %PATH%`

Install `gcc` to fix it