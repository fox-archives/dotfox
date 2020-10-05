#!/bin/bash

set -euo pipefail

cd testdata

cmd="go run ../"

dirs=("basic" "dual-file-exist")

for dir in "${dirs[@]}"; do
	$cmd \
		--dot-dir "$PWD/$dir" \
		--dest-dir "$PWD/$dir/dest" \
		user apply

	tree "$PWD/$dir"
done
