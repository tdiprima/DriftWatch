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
./driftwatch scan
./driftwatch baseline
./driftwatch diff
```

### Example output

```
$ ./driftwatch scan
DriftWatch v0.1.0

Host:   vulcan
OS:     Ubuntu 24.04.2 LTS
Kernel: 6.8.0-49-generic

Disk Usage:
  /                    63%
  /home                41%

Listening Ports:
  0.0.0.0:22 (tcp)
  0.0.0.0:443 (tcp)
  0.0.0.0:9100 (tcp)

Local Users:
  root (uid 0)
  alex (uid 1000)

Running Services:
  sshd.service
  nginx.service
  chronyd.service
```

### Example diff output

```
$ ./driftwatch diff
DriftWatch — Changes detected

+ NEW Listening Port
  0.0.0.0:8080 (tcp)

+ NEW User
  backupsvc

- Service REMOVED
  nginx.service

~ Disk CHANGE
  /: 63% -> 74%

4 change(s) detected.
```

When nothing changed:

```
$ ./driftwatch diff
DriftWatch — No drift detected.

System matches baseline.
```

### What gets scanned

| Check     | Source                       |
|-----------|------------------------------|
| Host info | `hostname`, `/etc/os-release`, `uname -r` |
| Disk      | `df`                         |
| Ports     | `ss -tuln`                   |
| Users     | `/etc/passwd` (UID >= 1000 + root) |
| Services  | `systemctl list-units`       |

## Commands

| Command    | Description                          |
|------------|--------------------------------------|
| `baseline` | Capture the trusted baseline snapshot |
| `scan`     | Display current system state          |
| `diff`     | Compare current state to baseline     |
| `version`  | Print version information             |

## Project Structure

```
driftwatch/
├── cmd/driftwatch/
│   └── main.go              CLI entry point
├── internal/
│   ├── scanner/
│   │   ├── host.go           Hostname, OS, kernel
│   │   ├── disk.go           Filesystem usage
│   │   ├── ports.go          Listening ports
│   │   ├── users.go          Local user accounts
│   │   └── services.go       Running systemd services
│   ├── snapshot/
│   │   └── snapshot.go       Snapshot struct, Capture(), Save(), LoadBaseline()
│   └── diff/
│       └── diff.go           Compare two snapshots, produce list of changes
├── go.mod
└── README.md
```

## Data Storage

Baseline saved as JSON at `~/.driftwatch/baseline.json`. Running `baseline` again overwrites it.

## Install (optional)

Copy the binary somewhere on your `PATH`:

```bash
sudo cp driftwatch /usr/local/bin/
```
