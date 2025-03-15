# Broadcast Server

A simple WebSocket-based broadcasting server and client built with Go.

## Features

- WebSocket server for real-time messaging
- Supports multiple clients
- CLI for starting a server and connecting clients
- Message broadcasting to all connected clients

## Prerequisites

- Go 1.20+

## Installation

Clone the repository:

```sh
git clone https://github.com/yourusername/broadcast-server.git
cd broadcast-server
```

Build the project:

```sh
go build -o broadcast-server main.go
```

## Usage

### Start the Server

Run the server with a specified port and channel name:

```sh
go run main.go start -p 8080 -n mygroup
```

### Connect a Client

A client can connect to the server and start sending messages:

```sh
go run main.go connect -p 8080 -g mygroup
```

### Sending Messages

Once connected, type a message and press **ENTER** to send.
To disconnect, type:

```sh
exit
```

## Project Structure

```
.
├── cmd                # CLI commands
│   ├── start.go       # Start server command
│   ├── connect.go     # Connect client command
├── server             # WebSocket server implementation
│   ├── server.go      # Server logic
├── config             # General configs
│   ├── constants.go   # constants for app
├── main.go            # Entry point
├── go.mod             # Dependencies
├── go.sum             # Checksum file
└── README.md          # Documentation
```

## Configuration

The default WebSocket server runs on `localhost:8080`. You can change the port using the `-p` flag.

## Contributing

1. Fork the repository
2. Create a new branch: `git checkout -b feature-branch`
3. Commit your changes: `git commit -m "Add new feature"`
4. Push to the branch: `git push origin feature-branch`
5. Open a pull request

## License

MIT License. See [LICENSE](LICENSE) for details.

---

Have fun! 🚀
