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

## Install (optional)

Copy the binary somewhere on your `PATH`:

```bash
sudo cp driftwatch /usr/local/bin/
```
