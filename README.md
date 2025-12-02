# DicewarePW

DicewarePW is a easy-to-use Diceware passphrase generator written in Go. It generates strong, memorable passphrases using the [German Diceware wordlist](https://github.com/bjoernalbers/diceware-wordlist-german).

## Features

- **Cryptographically Secure**: Uses Go's `crypto/rand` for secure random number generation.
- **Adjustable Length**: Generate passphrases with 5, 6, 7, or 8 words.
- **Entropy Calculation**: Displays the entropy of the generated passphrase.
- **Spice It Up**: Optionally add a special character to your passphrase for extra security.
- **One-Click Copy**: Easily copy the passphrase to your clipboard.
- **Self-Contained**: Single binary with embedded UI assets (HTML, CSS, JS).

## Installation

### Prerequisites

- Go 1.16 or higher (for `embed` support).

### Build from Source

1.  Clone the repository:
    ```bash
    git clone https://github.com/yourusername/dicewarepw.git
    cd dicewarepw
    ```

2.  Build the binary:
    ```bash
    go build -o dicewarepw ./cmd/web
    ```

## Usage

Run the application:

```bash
./dicewarepw
```

The server will start on `http://0.0.0.0:4000` by default.

### Configuration

You can configure the port using the `PORT` environment variable:

```bash
PORT=8080 ./dicewarepw
```

## Credits

- **Wordlist**: [German Diceware Wordlist](https://github.com/bjoernalbers/diceware-wordlist-german) by Björn Albers.
