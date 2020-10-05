#!/bin/bash

set -euo pipefail

cd testdata

cmd="go run ../"

dirs=("basic")

for dir in "${dirs[@]}"; do
	$cmd \
		--dot-dir "$(pwd)/$dir" \
		--dest-dir "$(pwd)/$dir/dest" \
		user apply
done
