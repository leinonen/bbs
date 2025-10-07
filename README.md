# Go BBS System

A modern Bulletin Board System (BBS) written in Go that users can connect to via SSH.

## Features

- SSH-based access (no web browser needed)
- User registration and authentication
- Multiple message boards
- Threaded discussions with replies
- Terminal-based UI with ANSI colors
- SQLite database for persistence
- Admin functionality for board management

## Prerequisites

- Go 1.19 or later
- Git (for downloading dependencies)

## Installation

### Quick Setup (using Makefile)
```bash
# Complete setup in one command
make setup

# Start the BBS
make run
```

### Manual Installation

1. Clone or download this repository:
```bash
cd /path/to/bbs
```

2. Install dependencies:
```bash
go mod download
```

3. Build the BBS:
```bash
go build -o gobbs
```

4. Initialize the database:
```bash
./gobbs -init
```

5. (Optional) Create a configuration file:
```bash
cp config.example.json config.json
# Edit config.json with your preferred settings
```

## Running the BBS

Start the BBS server:
```bash
./gobbs
```

Or with a custom config file:
```bash
./gobbs -config myconfig.json
```

The BBS will start listening on port 2222 by default.

## Connecting to the BBS

Users can connect using any SSH client:

```bash
ssh localhost -p 2222
```

Or if you want to connect with a specific username:
```bash
ssh username@localhost -p 2222
```

## First Time Setup

1. When you first connect, you can:
   - Register a new account
   - Login with existing credentials
   - Continue as a guest (if enabled)

2. The database is pre-populated with three boards:
   - `general` - General discussion
   - `tech` - Technology and programming
   - `random` - Random topics

3. To create an admin user, first register normally, then manually update the database:
```bash
sqlite3 bbs.db
UPDATE users SET is_admin = 1 WHERE username = 'yourusername';
.quit
```

## Configuration

The BBS can be configured using a JSON file. See `config.example.json` for available options:

- `listen_addr`: Address and port to listen on (default: ":2222")
- `database_path`: Path to SQLite database file (default: "bbs.db")
- `server_name`: Name displayed in the BBS (default: "Go BBS System")
- `host_key_path`: Path to SSH host key file (default: "host_key")
- `allow_anonymous`: Allow guest access without login (default: true)
- `max_users`: Maximum concurrent users (default: 100)

## Usage

### Navigation

- Use number keys to select menu options
- Type commands as shown in menus (N for New, V for View, etc.)
- Use Enter to confirm selections

### Posting

1. Browse to a board
2. Press 'N' to create a new post
3. Enter a title
4. Type your message (multiple lines)
5. Type '.' on a new line to finish

### Replying

1. View a post
2. Press 'R' to reply
3. Type your reply
4. Type '.' on a new line to finish

## Security Notes

- The SSH host key is automatically generated on first run
- User passwords are hashed using bcrypt
- Consider disabling anonymous access in production
- Use a firewall to restrict access if needed

## Makefile Commands

The project includes a comprehensive Makefile with the following commands:

- `make build` - Build the BBS binary
- `make run` - Build and run the server
- `make init` - Initialize the database
- `make setup` - Complete setup (deps, build, init)
- `make clean` - Remove build artifacts
- `make test` - Run tests
- `make fmt` - Format code
- `make check` - Run fmt, vet, and tests
- `make help` - Show all available commands

## Development

### Project Structure

```
bbs/
├── main.go           # Entry point
├── config/          # Configuration handling
├── server/          # SSH server implementation
├── bbs/            # Core BBS logic (boards, posts, users)
├── ui/             # Terminal UI
├── database/       # Database layer
├── Makefile         # Build automation
└── Dockerfile       # Container support
```

### Adding Features

The codebase is modular and easy to extend:

1. Add new UI screens in `ui/`
2. Add new data models in `bbs/`
3. Extend the database schema in `database/db.go`

## Troubleshooting

### Port already in use
If port 2222 is already in use, change it in the config file or use a different port.

### SSH connection refused
Make sure the BBS server is running and listening on the correct port.

### Database locked
Ensure only one instance of the BBS is running at a time.

## License

This project is provided as-is for educational purposes.

## Contributing

Contributions are welcome! Feel free to submit issues and pull requests.