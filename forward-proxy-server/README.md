# Forward Proxy Server

A small HTTP forward proxy written in Go. It forwards ordinary HTTP requests, opens TCP tunnels for HTTPS `CONNECT` requests, blocks configured hosts, and scans ordinary HTTP response bodies for configured words.

## Features

- Forwards HTTP requests to their destination.
- Supports HTTPS connections through `CONNECT` tunneling.
- Blocks hosts listed in `forbidden-hosts.txt`.
- Checks response bodies for words listed in `banned-words.txt` and replaces a matching response with `502 Bad Gateway` and `response not safe`.
- Writes request and error logs to `forwardProxyServer.log`.

## Requirements

- Go installed and available on your `PATH`.
- The two configuration files described below in the process's working directory.

The code uses Go's standard library and does not require third-party packages. If you put this project in a Go module, initialize it with `go mod init <module-name>` and run it with `go run .`.

## Configuration

Create `forbidden-hosts.txt` and `banned-words.txt` beside the directory from which you run the server. Put one entry per line. Blank lines and lines beginning with `//` are ignored.

Example `forbidden-hosts.txt`:

```text
example.org
ads.example.net
```

Example `banned-words.txt`:

```text
sample phrase
another word
```

Host matching and banned-word matching are based on the code's string checks; entries are not treated as wildcard patterns. Banned words are checked without regard to case.

## Run

Start the server from the project directory:

```bash
go run .
```

The proxy listens on `:8080` (all network interfaces). To use it with `curl`:

```bash
curl --proxy http://localhost:8080 http://example.com/
curl --proxy http://localhost:8080 https://example.com/
```

The second command uses a `CONNECT` tunnel. Configure your browser or application to use `http://localhost:8080` as its HTTP proxy to test it there as well.

## How it works

1. For an HTTP request, the server checks the requested host, builds an outbound request, and forwards it.
2. It reads the complete response body and checks it for banned words.
3. If a banned word is found, it returns an error response; otherwise it sends the upstream status, headers, and body to the client.
4. For `CONNECT`, it dials the requested host and relays bytes between the client and destination until one copy operation ends.

## Current scope and limitations

- HTTPS traffic is tunneled as encrypted bytes. The proxy does not decrypt or inspect HTTPS responses, so banned-word scanning applies only to ordinary HTTP responses.
- The response body is fully read into memory before it is checked and returned. Large responses can use substantial memory.
- Host filtering is not a complete URL, DNS, or IP-address security boundary. This is a learning project; do not expose it to untrusted networks without adding authentication, careful destination/IP validation, request limits, and robust connection handling.
- The server binds to all interfaces on port `8080`; restrict access with firewall rules or change the listen address if it should only be reachable locally.
- `CONNECT` tunnels are opaque, so host checks are the only filtering applied to their contents.

## Logs

The server creates or appends to `forwardProxyServer.log` in its working directory. It records informational and error messages there.

## License

No license is specified yet. Add a `LICENSE` file before distributing the project if you want to grant others permission to use, modify, or redistribute it.
