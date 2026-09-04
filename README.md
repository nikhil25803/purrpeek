<div align="center">

<pre>
                                       _
 _ __  _   _ _ __ _ __ _ __   ___  ___| | __
| '_ \| | | | '__| '__| '_ \ / _ \/ _ \ |/ /
| |_) | |_| | |  | |  | |_) |  __/  __/   &lt;
| .__/ \__,_|_|  |_|  | .__/ \___|\___|_|\_\
|_|                   |_|
</pre>

<a href="https://github.com/nikhil25803/purrpeek/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/nikhil25803/purrpeek/actions/workflows/ci.yml/badge.svg?branch=main"></a>
<a href="https://github.com/nikhil25803/purrpeek/blob/main/go.mod"><img alt="Go version" src="https://img.shields.io/github/go-mod/go-version/nikhil25803/purrpeek?logo=go"></a>
<a href="https://github.com/nikhil25803/purrpeek/blob/main/LICENSE"><img alt="MIT license" src="https://img.shields.io/github/license/nikhil25803/purrpeek"></a>

</div>

Purrpeek is a cross-platform CLI that gives you a quick, cat-approved view of your operating system, hardware, storage, network, power, shell, and terminal. It supports macOS, Linux, and Windows.

System collection is best-effort: if a detail is unavailable, Purrpeek returns the rest of the report instead of failing the entire command.

## Local setup

### Prerequisites

- [Git](https://git-scm.com/)
- [Go 1.26.5](https://go.dev/dl/) or a compatible version

Clone the repository and download its dependencies:

```sh
git clone https://github.com/nikhil25803/purrpeek.git
cd purrpeek
go mod download
```

Run the tests and start Purrpeek from source:

```sh
go test ./...
go run ./cmd/purrpeek --json
```

Build a local binary:

```sh
# macOS and Linux
go build -o bin/purrpeek ./cmd/purrpeek

# Windows
go build -o bin/purrpeek.exe ./cmd/purrpeek
```

## Useful commands

| Command                                  | Description                         |
| ---------------------------------------- | ----------------------------------- |
| `go run ./cmd/purrpeek --json`           | Print the system report as JSON.    |
| `go run ./cmd/purrpeek --json --verbose` | Print JSON and collection warnings. |
| `go run ./cmd/purrpeek --help`           | Show all CLI flags.                 |
| `make test`                              | Run all Go tests.                   |
| `make run-purrpeek`                      | Run Purrpeek from source.           |
| `make build-purrpeek`                    | Build `bin/purrpeek`.               |

Warnings are written to standard error only when `--verbose` is enabled, so JSON output remains usable by other tools.

## Image configuration

Purrpeek randomly selects one bundled image from `purrpeek-conf.yaml`:

```yaml
images:
  - mongo_no_bg.png
  - snow_no_bg.png
```

Place the file under `purrpeek/` in your operating system's user configuration directory:

| Platform | Configuration file |
| --- | --- |
| macOS | `~/Library/Application Support/purrpeek/purrpeek-conf.yaml` |
| Linux | `$XDG_CONFIG_HOME/purrpeek/purrpeek-conf.yaml` (normally `~/.config/purrpeek/purrpeek-conf.yaml`) |
| Windows | `%AppData%\purrpeek\purrpeek-conf.yaml` |

Available images are `mongo_no_bg.png`, `mongo_purrpeek.png`, `snow_no_bg.png`, and `snow_purrpeek.png`. A missing configuration uses the bundled defaults; malformed or unreadable configuration stops the command with an error.

## Greeting localization

Purrpeek randomly selects a greeting language on each run. To add a language or replace phrases for an existing language, create `greetings.json` beside `purrpeek-conf.yaml` in the platform-specific configuration directory above:

```json
{
  "de": {
    "morning": ["Guten Morgen"],
    "afternoon": ["Guten Tag"],
    "evening": ["Guten Abend"],
    "night": ["Gute Nacht"]
  }
}
```

The supported periods are `morning`, `afternoon`, `evening`, and `night`. User phrases replace the matching language and period while all other bundled translations remain available. The file is optional, but malformed or unreadable greeting files stop the command with an error.

## System information

| Category         | Information provided                                                                      |
| ---------------- | ----------------------------------------------------------------------------------------- |
| Operating system | Username, hostname, OS name and version, kernel version, and architecture                 |
| Uptime           | Human-readable uptime and boot time                                                       |
| Time             | Local time, timezone, and UTC offset                                                      |
| CPU              | Model, physical and logical core counts, usage, and frequency when available              |
| GPU              | Detected graphics processor models                                                        |
| Memory           | Total, used, and available memory with usage percentage                                   |
| Disk             | Home usage, mounted volumes, filesystems, mount points, capacity, and usage               |
| Network          | Primary interface, local IPv4 and IPv6 addresses, MAC address, MTU, and interface details |
| Batteries        | Battery names and charge percentages when present                                         |
| Shell            | Name, version, and executable path                                                        |
| Terminal         | Name, version, `TERM`, `COLORTERM`, width, and height                                     |

## License

Purrpeek is available under the [MIT License](LICENSE). Copyright © 2026 Nikhil Raj.
