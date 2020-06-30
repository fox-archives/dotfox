#!/bin/sh

requires gorun

# https://blog.cloudflare.com/using-go-as-a-scripting-language-in-linux/
# https://www.kernel.org/doc/html/latest/admin-guide/binfmt-misc.html
# https://wiki.ubuntu.com/gorun

mount | grep binfmt_misc

echo ':golang:E::go::/usr/local/bin/gorun:OC' | sudo tee /proc/sys/fs/binfmt_misc/register
:golang:E::go::/usr/local/bin/gorun:OC
