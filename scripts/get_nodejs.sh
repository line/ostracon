#!/usr/bin/env bash


VERSION=v12.15.0
NODE_FULL=node-${VERSION}-linux-x64

mkdir -p ~/.local/bin
mkdir -p ~/.local/node
wget http://nodejs.org/dist/${VERSION}/${NODE_FULL}.tar.gz -O ~/.local/node/${NODE_FULL}.tar.gz
tar -xzf ~/.local/node/${NODE_FULL}.tar.gz -C ~/.local/node/
ln -s ~/.local/node/${NODE_FULL}/bin/node ~/.local/bin/node
ln -s ~/.local/node/${NODE_FULL}/bin/npm ~/.local/bin/npm
export PATH=~/.local/bin:$PATH
npm i -g dredd@13.1.0
ln -s ~/.local/node/${NODE_FULL}/bin/dredd ~/.local/bin/dredd
