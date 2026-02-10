#!/bin/sh

set -eu

OUT=$(mktemp -d -t nar-hash-XXX)
go mod vendor -o "$OUT"

# go install tailscale.com/cmd/nardump@v1.92.3
SHA=$(nardump --sri "$OUT")

sed -i "s|vendorHash = \".*\";|vendorHash = \"${SHA}\";|g" flake.nix
rm -rf "$OUT"
