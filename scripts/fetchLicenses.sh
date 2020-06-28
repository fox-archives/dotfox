#!/bin/sh

licenses="$(curl -o- --silent https://github.com/github/choosealicense.com/tree/gh-pages/_licenses \
	| tac | tac | grep -Pio "(?<=gh-pages/_licenses\/).*(?=.txt\">)")"

mkdir licenses >/dev/null 2>&1
year="2020"
name="Edwin Kofler"
for license in $licenses; do
	printf "\033[0;94m%s\033[0m\n" "working on $license"
	licenseText="$(curl -o- --silent https://raw.githubusercontent.com/github/choosealicense.com/gh-pages/_licenses/$license.txt)"
	licenseText="$(echo "$licenseText" | tr '\n' '~' | sed -E 's/(---~.*---~~)(.*)/\2/g' | tr '~' '\n' | sed "s/\[year\]/$year/" | sed "s/\[fullname\]/$name/")"
	echo "$licenseText" > "licenses/${license}.txt"
done

