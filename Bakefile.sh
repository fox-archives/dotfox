# shellcheck shell=bash
# shellcheck shell=bash

task.build() {
	nimble build dotfox "$@"
}

task.run() {
	nimble run dotfox "$@"
}
