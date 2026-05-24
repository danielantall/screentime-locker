# Screentime-locker

## Description
Set limits you can't circumvent. Sets a macOS Screen Time passcode **you can't remember**.

## Motivation
Our phones can be addictive. Setting up screen time is great to set limits on our phone usage, but constantly pressing "Remind me in 15 minutes" is unproductive. Screentime-locker makes it an extremely long process to reset your passcode, making it easier for you to lock-in.

Generates a random 4-digit code, then walks you through typing it into System Settings one digit at a time, with fake digits and deletions mixed in so you never see the full passcode. Two rounds (enter + confirm), then it's gone from memory. Nothing is stored or logged.

![demo](./demo.gif)

## Quick Start

```bash
brew install go
go run .
```

Open **System Settings → Screen Time → Lock Screen Time Settings** before starting.

## Usage

### Prerequisites

- **macOS** (Screen Time is a macOS feature)
- **Go 1.21+** installed (`brew install go`)

### Run from source

```bash
go run .
```

### Build a binary

```bash
go build -o screentime-locker .
./screentime-locker
```

### Walkthrough

1. **Welcome screen** — press `Enter` to start.
2. **Setup** — open **System Settings → Screen Time → Lock Screen Time Settings** so the passcode input is visible.
3. **First entry** — the app shows one digit (or a delete instruction) at a time. Type or delete exactly what it tells you in the System Settings passcode field. Fake digits and deletions are mixed in so you never see the full code.
4. **Confirmation** — macOS asks you to re-enter the passcode. The app generates a new distraction pattern for the same code. Follow the steps again.
5. **Done** — the passcode is set. It was never displayed in full and is not stored anywhere.

**Controls:**

| Key | Action |
|---|---|
| `Enter` / `Space` | Advance to the next step |
| `Ctrl+C` | Quit at any time |
| `q` | Quit (welcome & done screens only) |

> **⚠️ Warning:** There is no way to recover the passcode. To remove it, go to **Apple ID → Screen Time → Change Passcode** (requires your Apple ID password) or perform a full device reset.

## Contributing

Contributions are welcome! Here's how to get started:

1. **clone** this repo.
2. **Create a branch** for your change:
   ```bash
   git checkout -b feature/my-change
   ```
3. **Make your changes.** The codebase is organized as:
   - `main.go` — entry point, initializes the Bubble Tea program.
   - `model.go` — TUI model, update logic, and all views/styles.
   - `sequence.go` — passcode generation and distraction-sequence algorithm.
4. **Test locally:**
   ```bash
   go build ./...
   go vet ./...
   ```
5. **Commit** with a clear message and **open a pull request** against `main`.

### Guidelines

- Keep the dependency footprint small (currently just [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss)).
- Follow standard Go formatting (`gofmt`).
- The passcode must **never** be logged, stored, or displayed in full — this is a core design constraint.
