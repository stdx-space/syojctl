# syojctl

A command-line tool for interacting with the SYOJ (Sing Yin Online Judge) platform, built with Go and Cobra.

## Features

- Login to SYOJ via API using the `syojctl login` command
- Show problem details with the `syojctl show-problem` command
- Submit solution code for problems with the `syojctl submit` command
- Retrieve and store authentication credentials (Token and TokenId) in XDG standard directories
- Accept username and password via command line flags or environment variables
- Command-line interface built with Cobra
- Modular architecture with separate packages
- Colorful and structured logging with charmbracelet/log
- Beautiful terminal rendering with charmbracelet/glamour

## Usage

### Running directly with Go
```bash
# Show help
go run main.go --help

# Login to SYOJ with flags
go run main.go login -u your-email@example.com -p your-password

# Login to SYOJ with environment variables
export SYOJ_USERNAME=your-email@example.com
export SYOJ_PASSWORD=your-password
go run main.go login

# Show problem details (requires login first)
go run main.go show-problem I002

# Submit solution code from a file
go run main.go submit I001 -i solution.cpp

# Submit solution code from standard input
cat solution.cpp | go run main.go submit I001
```

### Building and running the executable
```bash
# Build the executable
go build -o syojctl

# Show help
./syojctl --help

# Login to SYOJ with flags
./syojctl login -u your-email@example.com -p your-password

# Login to SYOJ with environment variables
export SYOJ_USERNAME=your-email@example.com
export SYOJ_PASSWORD=your-password
./syojctl login

# Show problem details (requires login first)
./syojctl show-problem I002

# Submit solution code from a file
./syojctl submit I001 -i solution.cpp

# Submit solution code from standard input
cat solution.cpp | ./syojctl submit I001
```

The login command will:
1. Login to SYOJ using provided credentials or environment variables
2. Save the authentication tokens to the XDG config directory (`~/.config/syojctl/credentials.json` on most Unix-like systems)

The show-problem command will:
1. Load the saved authentication credentials
2. Fetch the problem details from the SYOJ API
3. Render the problem in a beautifully formatted way in your terminal

The submit command will:
1. Load the saved authentication credentials
2. Read source code from a file or standard input
3. Submit the code to SYOJ for evaluation

## Architecture

- `main.go` - Main application entry point
- `cmd/` - Command implementations
  - `root.go` - Root command definition
  - `login.go` - Login command implementation
  - `show_problem.go` - Show problem command implementation
  - `submit.go` - Submit command implementation
- `api/` - API client implementation
  - `client.go` - SYOJ API client
- `credentials/` - Credentials management
  - `manager.go` - Credentials saving/loading using XDG standard
- `go.mod` - Go module definition
- `syojctl` - Built executable
- `README.md` - This file

## Dependencies

- `github.com/spf13/cobra` - Command-line interface library
- `github.com/charmbracelet/log` - Colorful and structured logging library
- `github.com/charmbracelet/glamour` - Markdown rendering for terminals
- `github.com/adrg/xdg` - XDG Base Directory Specification implementation

## API Endpoints Used

- `POST https://syoj.org/api/login` - Authentication endpoint
- `GET https://syoj.org/api/problems/{problem-id}` - Problem details endpoint
- `POST https://syoj.org/api/submit` - Code submission endpoint

## Commands

- `syojctl login` - Login to the SYOJ platform
  - `-u, --username string` - Username for SYOJ login (defaults to SYOJ_USERNAME environment variable)
  - `-p, --password string` - Password for SYOJ login (defaults to SYOJ_PASSWORD environment variable)
- `syojctl show-problem [problem-id]` - Show details of a specific problem
- `syojctl submit [problem-id]` - Submit solution code for a problem
  - `-i, --input string` - Input file containing source code (reads from stdin if not specified)
  - `-l, --language string` - Programming language for submission
- `syojctl help` - Help about any command
- `syojctl completion` - Generate autocompletion scripts