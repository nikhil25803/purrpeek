<pre align="center">

   _                                          _    
  | |                                        | |   
 / __)  _ __  _   _ _ __ _ __ _ __   ___  ___| | __
 \__ \ | '_ \| | | | '__| '__| '_ \ / _ \/ _ \ |/ /
 (   / | |_) | |_| | |  | |  | |_) |  __/  __/   < 
  |_|  | .__/ \__,_|_|  |_|  | .__/ \___|\___|_|\_\
       | |                   | |                   
       |_|                   |_|                   

</pre>

<p align="center">
  <a href="https://github.com/nikhil25803/purrpeek/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/nikhil25803/purrpeek/actions/workflows/ci.yml/badge.svg?branch=main"></a>
  <a href="https://github.com/nikhil25803/purrpeek/blob/main/go.mod"><img alt="Go 1.26.5" src="https://img.shields.io/badge/Go-1.26.5-00ADD8?logo=go&logoColor=white"></a>
  <a href="https://github.com/nikhil25803/purrpeek/blob/main/LICENSE"><img alt="MIT License" src="https://img.shields.io/badge/License-MIT-yellow.svg"></a>
  <a href="https://github.com/nikhil25803/homebrew-tap/blob/main/Formula/purrpeek.rb"><img alt="Homebrew available" src="https://img.shields.io/badge/Homebrew-available-FBB040?logo=homebrew&logoColor=black"></a>
  <a href="https://github.com/nikhil25803/purrpeek/blob/main/flake.nix"><img alt="Nix flake available" src="https://img.shields.io/badge/Nix-flake_available-5277C3?logo=nixos&logoColor=white"></a>
  <a href="https://github.com/nikhil25803/scoop-bucket/blob/main/bucket/purrpeek.json"><img alt="Scoop available" src="https://img.shields.io/badge/Scoop-available-53B7F5?logo=windows&logoColor=white"></a>
  <a href="#linux-installer"><img alt="Linux installer available" src="https://img.shields.io/badge/Linux-installer_available-FCC624?logo=linux&logoColor=black"></a>
  <a href="#installation"><img alt="Debian package coming soon" src="https://img.shields.io/badge/Debian-coming_soon-A81D33?logo=debian&logoColor=white"></a>
</p>

Purrpeek is a cross-platform, cat-approved CLI for quickly inspecting your operating system, hardware, storage, network, power, shell, and terminal. It supports macOS, Linux, and Windows.

Collection is best-effort: if a system detail is unavailable, Purrpeek still displays everything it collected successfully.

|                                                   Preview 1                                                    |                                                   Preview 2                                                   |
| :------------------------------------------------------------------------------------------------------------: | :-----------------------------------------------------------------------------------------------------------: |
|            **Ghostty**<br>![Purrpeek running in Ghostty with Mongo artwork](asset/ss_ghostty_1.png)            |            **Ghostty**<br>![Purrpeek running in Ghostty with Snow artwork](asset/ss_ghostty_2.png)            |
| **macOS Terminal**<br>![Purrpeek Braille output in macOS Terminal with Mongo artwork](asset/ss_mac_term_1.png) | **macOS Terminal**<br>![Purrpeek Braille output in macOS Terminal with Snow artwork](asset/ss_mac_term_2.png) |

> **Note:** Images are rendered in terminals with supported image protocols; other terminals use the Braille fallback.

## Table of contents

1. [Installation](#installation)
   1. [Linux installer](#linux-installer)
   1. [Release archives](#release-archives)
2. [Local setup](#local-setup)
   1. [Prerequisites](#prerequisites)
3. [CLI usage](#cli-usage)
4. [Development commands](#development-commands)
5. [Package builds](#package-builds)
6. [System information](#system-information)
7. [Configuration](#configuration)
   1. [Purrpeek configuration](#purrpeek-configuration)
   2. [Greeting localization](#greeting-localization)
8. [License](#license)

## Installation

| Package manager | Commands                                                                                                                                          |
| --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| Homebrew        | `brew install nikhil25803/tap/purrpeek`                                                                                                           |
| Nix             | Install: `nix profile install github:nikhil25803/purrpeek`<br><br>Run: `nix run github:nikhil25803/purrpeek`                                      |
| Scoop           | Setup: `scoop bucket add nikhil25803 https://github.com/nikhil25803/scoop-bucket`<br><br>Install: `scoop install purrpeek`<br><br>Run: `purrpeek` |
| AUR             | Coming soon                                                                                                                                       |
| WinGet          | Coming soon                                                                                                                                       |
| Debian          | Coming soon                                                                                                                                       |
| Fedora          | Coming soon                                                                                                                                       |
| Snap            | Coming soon                                                                                                                                       |

### Linux installer

Install the latest release on Linux (amd64 or arm64):

```sh
curl -fsSL https://raw.githubusercontent.com/nikhil25803/purrpeek/main/scripts/install.sh | sh
```

The script verifies the release checksum and installs `purrpeek` to `/usr/local/bin`, using `sudo` when required. Review [the installer](scripts/install.sh) before running it if you prefer not to pipe a remote script directly into a shell.

### Release archives

Download the archive for your operating system and architecture from [GitHub Releases](https://github.com/nikhil25803/purrpeek/releases). On Linux or macOS, extract it and install the binary:

```sh
tar -xzf purrpeek_<version>_<os>_<arch>.tar.gz
sudo install -m 755 purrpeek /usr/local/bin/purrpeek
```

Windows releases are provided as `.zip` archives. Extract `purrpeek.exe` and place it in a directory on `PATH`.

## Local setup

### Prerequisites

- [Git](https://git-scm.com/)
- [Go 1.26.5](https://go.dev/dl/) or a compatible version

Clone the repository:

```sh
git clone https://github.com/nikhil25803/purrpeek.git
cd purrpeek
```

Set up, test, and build the project without Make:

```sh
# macOS and Linux
./scripts/setup.sh

# Windows Command Prompt
scripts\setup.bat
```

The scripts download dependencies, run all tests, and build the executable under `bin/`. To perform the same steps manually:

```sh
go mod download
go test ./...
# macOS and Linux
go build -o bin/purrpeek ./cmd/purrpeek

# Windows
go build -o bin/purrpeek.exe ./cmd/purrpeek
```

Run the built executable with `./bin/purrpeek` on macOS/Linux or `bin\purrpeek.exe` on Windows.

## CLI usage

| Command                      | Description                                             |
| ---------------------------- | ------------------------------------------------------- |
| `purrpeek`                   | Render artwork and system details.                      |
| `purrpeek --json`            | Print the complete system report as JSON.               |
| `purrpeek --verbose`         | Render normally and show collection warnings on stderr. |
| `purrpeek --json --verbose`  | Print JSON and show collection warnings on stderr.      |
| `purrpeek --help`            | Show all CLI flags.                                     |

Warnings are quiet by default and written to standard error only with `--verbose`, keeping JSON output usable by other tools.

From the repository, use `go run ./cmd/purrpeek` in place of `purrpeek` or run `make run-purrpeek`.

## Development commands

| Command                 | Description                                |
| ----------------------- | ------------------------------------------ |
| `./scripts/setup.sh`    | Set up, test, and build on macOS or Linux. |
| `scripts\setup.bat`     | Set up, test, and build on Windows.        |
| `make test`             | Run all Go tests.                          |
| `make run-purrpeek`     | Run Purrpeek from source.                  |
| `make build-purrpeek`   | Build `bin/purrpeek`.                      |
| `make build-windows`    | Build `bin/purrpeek.exe` for Windows.      |

## Package builds

Nix commands require Nix; AUR and Debian commands require Docker.

| Command          | Description                                                  |
| ---------------- | ------------------------------------------------------------ |
| `make nix-check` | Validate the flake on all supported systems.                 |
| `make nix-build` | Build the package with Nix.                                  |
| `make aur-check` | Build the x86_64 AUR package interactively.                  |
| `make aur-build` | Build the x86_64 AUR package.                                |
| `make deb-build` | Build the amd64 Debian package and place artifacts in `dist/`. |

## System information

| Category         | Information provided                                                                |
| ---------------- | ----------------------------------------------------------------------------------- |
| Operating system | Username, hostname, OS name and version, kernel version, and architecture           |
| Uptime           | Human-readable uptime and boot time                                                 |
| Time             | Local time, timezone, and UTC offset                                                |
| CPU              | Model, physical and logical core counts, usage, and optional frequency              |
| GPU              | Detected graphics processor models                                                  |
| Memory           | Total, used, and available memory with usage percentage                             |
| Disk             | Home usage, mounted volumes, filesystems, mount points, capacity, and usage         |
| Network          | Primary interface, IPv4 and IPv6 addresses, MAC address, MTU, and interface details |
| Batteries        | Battery names and charge percentages when present                                   |
| Shell            | Name, version, and executable path                                                  |
| Terminal         | Name, version, `TERM`, `COLORTERM`, width, and height                               |

## Configuration

Purrpeek reads optional user files from the platform configuration directory:

| Platform | Purrpeek configuration                                      | Greeting localization                                   |
| -------- | ----------------------------------------------------------- | ------------------------------------------------------- |
| macOS    | `~/Library/Application Support/purrpeek/purrpeek-conf.yaml` | `~/Library/Application Support/purrpeek/greetings.json` |
| Linux    | `$XDG_CONFIG_HOME/purrpeek/purrpeek-conf.yaml`              | `$XDG_CONFIG_HOME/purrpeek/greetings.json`              |
| Windows  | `%AppData%\purrpeek\purrpeek-conf.yaml`                     | `%AppData%\purrpeek\greetings.json`                     |

On Linux, `$XDG_CONFIG_HOME` is normally `~/.config`. User configuration is read on every invocation, so these files do not require rebuilding Purrpeek.

### Purrpeek configuration

The YAML configuration controls which bundled images may be selected and which system fields appear in normal terminal output:

```yaml
images:
  - mongo_no_bg.png
  - snow_no_bg.png

render:
  cpu:
    model:
      name: CPU
      description: Processor model
      enabled: true
```

For each render field:

- `name` is the label displayed in the terminal.
- `description` documents what the field contains.
- `enabled` decides whether the field is rendered.

The user file is merged over the [embedded defaults](internal/conf/purrpeek-conf.yaml), so it only needs to contain values you want to change. Purrpeek randomly chooses from valid configured image names. Malformed YAML, unknown fields, and unreadable configuration files stop the command with a concise error.

Repository contributors can change the embedded YAML or add images under `internal/asset/images/`. Because these files are embedded in the executable, rebuild afterward:

```sh
make build-purrpeek
./bin/purrpeek
```

### Greeting localization

Purrpeek chooses a random available language and greeting on each run. The local time selects one of four periods:

| Period      | Local time  |
| ----------- | ----------- |
| `morning`   | 05:00–11:59 |
| `afternoon` | 12:00–16:59 |
| `evening`   | 17:00–20:59 |
| `night`     | 21:00–04:59 |

Create an optional `greetings.json` in the user configuration directory to add a language or replace phrases for an existing language and period:

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

User entries merge with the bundled catalog by language and period. Phrase arrays for matching periods replace the bundled arrays, while all unspecified translations remain available. Blank and duplicate phrases are removed. Malformed JSON, unsupported period names, control characters, and unreadable files stop the command with an error.

Repository contributors can edit the [bundled greeting catalog](internal/asset/localisation/greetings.json). It is embedded at build time, so run `make build-purrpeek` before testing those changes in `bin/purrpeek`.

## License

Purrpeek is available under the [MIT License](LICENSE). Copyright © 2026 Nikhil Raj.
