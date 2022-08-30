#!/bin/sh

mkdir ~/bin && cp out/grip ~/bin/grip
echo 'function grip() { ~/bin/grip $1 $2 }' >> ~/.zshrc
