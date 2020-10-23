#!/bin/bash

set -euo pipefail

setupConfigFile() {
	mkdir config
	cat <<-EOF > config/user.dots.toml
	[[files]]
	file = "$1"
	EOF
}

setup() {
	read -rs -d\| input
	echo "$input"

	[[ $1 =~ folder ]] && {
		mkdir "$1.rel" && cd "$1.rel"
		setupConfigFile "folder-dot"
		eval "$input"
		cd ..

		mkdir "$1.subdir" && cd "$1.subdir"
		setupConfigFile "subdir/folder-dot"
		eval "$input"
		cd ..

		mkdir "$1.root" && cd "$1.root"
		setupConfigFile "/subdir/folder-dot"
		eval "$input"
		cd ..

		return 0
	}

	mkdir "$1.rel" && cd "$1.rel"
	setupConfigFile "a-dot-file"
	cd ..

	mkdir "$1.subdir" && cd "$1.subdir"
	setupConfigFile "subdir/a-dot-file"
	cd ..

	mkdir "$1.root" && cd "$1.root"
	setupConfigFile "/subdir/a-dot-file"
	cd ..

	return 0
}

skel() {
	mkdir dest
	mkdir "$1"
}

[[ ${1:-""} == teardown ]] && {
	rm -r misc.*
	rm -r file.*
	rm -r folder.*

	exit 0
}

[[ ${1:-""} =~ h|help ]] && {
	cat > >&1 <<-EOF
		:
	EOF

	exit 0
}

declare -r fileName="a-dot-file"
declare -r folderName="folder-dot"

cd "$PWD"
cd ../testdata

setup file.1 <<-EOF
	skel user
	touch user/$fileName|
EOF

setup file.2 <<-EOF
	skel user
	touch dest/$fileName|
EOF


setup file.3 <<-EOF
	skel user
	touch user/$fileName
	touch dest/$fileName|
EOF


setup file.4 <<-EOF
	skel user
	cat > user/$fileName <<< "same content"
	cat > dest/$fileName <<< "same content"|
EOF

# setup file.5 <<-EOF
# 	skel user
# 	cat > user/$fileName <<< "user file"
# 	cat > dest/$fileName <<< "dest file"
# EOF

setup folder.1 <<-EOF
	skel user
	mkdir user/$folderName|
EOF

setup folder.2 <<-EOF
	skel user
	mkdir dest/$folderName|
EOF

setup folder.3 <<-EOF
	skel user
	mkdir user/$folderName
	mkdir dest/$folderName
	cat > dest/$folderName/$fileName <<< "some content"|
EOF

setup folder.4 <<-EOF
	skel user
	mkdir user/$folderName
	cat > user/$folderName/$fileName <<< "some content"
	mkdir dest/$folderName|
EOF

setup folder.5 <<-EOF
	skel user
	mkdir user/$folderName
	cat > user/$folderName/$fileName <<< "user content"
	mkdir dest/$folderName
	cat > dest/$folderName/$fileName <<< "dest content"|
EOF
