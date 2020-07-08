#!/bin/sh -aux

mount | grep -q binfmt_misc || {
        # see https://www.kernel.org/doc/html/v4.14/admin-guide/binfmt-misc.html
        echo "You don't have binfmt_misc mounted. This configuration is unsupported"
        exit 1
}

[ -e /proc/sys/fs/binfmt_misc/golang ] && {
        echo "You already have an 'interpreter' associated with golang"
        cat /proc/sys/fs/binfmt_misc/golang
        exit
}

# install gorun
go get github.com/erning/gorun

# add kernel support for executing go files with the gorun 'interpreter'
echo ':golang:E::go::gorun:OC' | sudo tee /proc/sys/fs/binfmt_misc/register
