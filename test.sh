#!/bin/bash

set -euo pipefail

cd testdata

cmd="go run ../"

dirs=("basic" "dual-folder-exist")

for dir in "${dirs[@]}"; do
	$cmd \
		--dot-dir "$PWD/$dir" \
		user --user-dir "$PWD/$dir/dest" \
		apply

	tree "$PWD/$dir"
done
