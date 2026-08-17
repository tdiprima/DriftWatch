# DriftWatch

System drift detection for Linux. Take a baseline snapshot, then find out what changed.

## Requirements

- Go 1.21 or later
- Linux (Ubuntu, RHEL, Rocky Linux)

## Build

```bash
go build -o driftwatch ./cmd/driftwatch/
```

## Run

```bash
./driftwatch version
```

## Commands

| Command    | Description                          |
|------------|--------------------------------------|
| `baseline` | Capture the trusted baseline snapshot |
| `scan`     | Display current system state          |
| `diff`     | Compare current state to baseline     |
| `version`  | Print version information             |

## Install (optional)

Copy the binary somewhere on your `PATH`:

```bash
sudo cp driftwatch /usr/local/bin/
```
