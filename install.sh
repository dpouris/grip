#!/bin/sh

# Run cleanup
source cleanup.sh

# build binary
go build -o out/

# make bin dir in home and copy binary into it
mkdir ~/bin && cp out/grip ~/bin/grip

# create grip function and add to shell rc
echo 'function grip() { ~/bin/grip $1 $2 $3 }' >> ~/.$1rc

rm -rf out/
