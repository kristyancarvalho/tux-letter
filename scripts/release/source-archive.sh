#!/usr/bin/env sh
set -eu

version="${1:?version is required}"
ref="${2:-HEAD}"
out_dir="${3:-dist}"
name="tux-letter-${version}"
tmp_index="$(mktemp)"

cleanup() {
  rm -f "$tmp_index"
}
trap cleanup EXIT

mkdir -p "$out_dir"
GIT_INDEX_FILE="$tmp_index" git read-tree "$ref"
GIT_INDEX_FILE="$tmp_index" git rm -r --cached --ignore-unmatch packaging/aur specs >/dev/null
tree="$(GIT_INDEX_FILE="$tmp_index" git write-tree)"
git archive --format=tar --mtime="1970-01-01T00:00:00Z" --prefix="${name}/" "$tree" | gzip -n > "${out_dir}/${name}-src.tar.gz"
