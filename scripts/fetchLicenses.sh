#!/usr/bin/env bash
set -euo pipefail

cd ../licenses

licenses="$(curl -o- --silent https://github.com/github/choosealicense.com/tree/gh-pages/_licenses \
	| tac | tac | grep -Pio "(?<=gh-pages/_licenses\/).*(?=.txt\">)")"

function doit() {
	declare -r license="${1:-"default"}"
	test "$license" = "default" && {
		echo "license blank. exiting."
		exit 1
	}

	year="2020"
	name="Edwin Kofler"

	printf "\033[0;94m%s\033[0m\n" "working on $license"
	licenseText="$(curl -o- --silent "https://raw.githubusercontent.com/github/choosealicense.com/gh-pages/_licenses/$license.txt")"
	licenseText="$(echo "$licenseText" | tr '\n' '~' | sed -E 's/(---~.*---~~)(.*)/\2/g' | tr '~' '\n' | sed "s/\[year\]/$year/" | sed "s/\[fullname\]/$name/")"
	echo "$licenseText" >"${license}.txt"
}

# export -f doit

# slow
# for license in $licenses; do
# parallel -P 0 --link --bar --progress --joblog _parallel.log \
# --trim lr \
# doit ::: "$license"
# done

for license in $licenses; do
	doit "$license" &
	while test "$(jobs -p | wc -w)" -ge 10; do sleep 0.1; done
done

# doesn't work
# parallel -P 0 --link --bar --progress --joblog _parallel.log \
# 	--trim lr -a $licenses \
# 	doit ::: {1}

# other
# doit ::: "$(echo "$licenses" | tr " " "\n" | wc -l | xargs seq 1)"
