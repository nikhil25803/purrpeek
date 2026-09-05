#!/bin/sh

set -eu

repo=https://github.com/nikhil25803/purrpeek

fail() {
	printf 'Error: %s\n' "$1" >&2
	exit 1
}

[ "$(uname -s)" = Linux ] || fail "this installer supports Linux only"

for command in curl tar sha256sum install; do
	command -v "$command" >/dev/null 2>&1 || fail "$command is required"
done

case "$(uname -m)" in
	x86_64 | amd64) arch=amd64 ;;
	aarch64 | arm64) arch=arm64 ;;
	*) fail "unsupported architecture: $(uname -m)" ;;
esac

release_url=$(curl -fsSL -o /dev/null -w '%{url_effective}' "$repo/releases/latest")
tag=${release_url##*/}
case "$tag" in
	v[0-9]*) ;;
	*) fail "could not determine the latest release" ;;
esac

version=${tag#v}
archive="purrpeek_${version}_linux_${arch}.tar.gz"
download_url="$repo/releases/download/$tag"
temp_dir=$(mktemp -d)
trap 'rm -rf "$temp_dir"' EXIT HUP INT TERM

curl -fsSL "$download_url/$archive" -o "$temp_dir/$archive"
curl -fsSL "$download_url/checksums.txt" -o "$temp_dir/checksums.txt"

expected=$(awk -v archive="$archive" '$2 == archive { print $1 }' "$temp_dir/checksums.txt")
[ -n "$expected" ] || fail "checksum not found for $archive"
actual=$(sha256sum "$temp_dir/$archive" | awk '{ print $1 }')
[ "$actual" = "$expected" ] || fail "checksum verification failed"

tar -xzf "$temp_dir/$archive" -C "$temp_dir"
[ -f "$temp_dir/purrpeek" ] || fail "release archive does not contain purrpeek"

if [ -w /usr/local/bin ]; then
	install -m 0755 "$temp_dir/purrpeek" /usr/local/bin/purrpeek
elif command -v sudo >/dev/null 2>&1; then
	sudo install -m 0755 "$temp_dir/purrpeek" /usr/local/bin/purrpeek
else
	fail "write access to /usr/local/bin or sudo is required"
fi

printf 'Installed purrpeek %s to /usr/local/bin/purrpeek\n' "$version"
