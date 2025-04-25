# Go gRPC Chat Server

A simple gRPC-based chat server implementation in Go.

## Prerequisites

- [Go](https://golang.org/dl/) (1.22 or newer)
- [Protocol Buffers Compiler](https://grpc.io/docs/protoc-installation/) (protoc)
- Go plugins for Protocol Buffers

## Installation (Only verified on Mac and Linux)

### 1. Install Go (Assume already done at this point)

### 2. Install Protocol Buffers Compiler
https://protobuf.dev/installation/

#### macOS
```
brew install protobuf
```

#### Linux (Debian/Ubuntu)
```
sudo apt install -y protobuf-compiler
```

#### Windows
```
winget install protobuf
```

### 3. Install Go plugins for Protocol Buffers and gRPC

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 4. Update your PATH

#### For Bash (Linux/macOS)
```
export PATH="$PATH:$(go env GOPATH)/bin”
source ~/.bashrc
```

#### For Zsh (macOS)
```
export PATH="$PATH:$(go env GOPATH)/bin”
source ~/.zshrc
```

#### For Windows
Add `%GOPATH%\bin` to your system PATH.

## Project Setup

### 1. Clone the repository
```
git clone <repository-url>
cd <repository-directory>
```

### 2. Install dependencies
```
go get
```

### 3. Generate Protocol Buffers code
```
protoc --proto_path=internal/protos \
       --go_out=internal/protos/pb \
       --go_opt=paths=source_relative \
       --go-grpc_out=internal/protos/pb \
       --go-grpc_opt=paths=source_relative \
       internal/protos/chat_server.proto
```

## Running the Application

### Start the server (Not included yet, will update after discussion)
```
go run cmd/server/server.go
```

### Start a client
```
go run cmd/client/client.go
```

## Development

### Regenerating Protocol Buffers

If you make changes to the `.proto` files, you'll need to regenerate the Go code:

```
protoc --proto_path=internal/protos \
       --go_out=internal/protos/pb \
       --go_opt=paths=source_relative \
       --go-grpc_out=internal/protos/pb \
       --go-grpc_opt=paths=source_relative \
       internal/protos/chat_server.proto
```

## Project Structure

```
.
├── cmd/
│   ├── client/       # Client application
│   └── server/       # Server application
├── internal/
│   ├── protos/       # Protocol Buffers definitions
│   │   ├── pb/       # Generated Go code
│   │   └── *.proto   # Proto files
|   ├── utils.go      # Has Server Address  
└── README.md
```