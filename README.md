# net-vigil

Network Traffic Monitoring and Analysis based Local Threat Intelligence Center

## How to run

To compile this project, you need to install Node.js and Go in advance.

Firstly, build frontend resources

```bash
$ cd frontend
$ npm run build
```

Then, build go project

```bash
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
```

```go
type SockTabEntry struct {
  ino        string
  LocalAddr  *SockAddr
  RemoteAddr *SockAddr
  State      SkState
  UID        uint32
  Process    *Process
}
```