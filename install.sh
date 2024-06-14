#!/bin/bash

npm install -g yarn cross-env

yarn install
if [ ! -d "dist" ]; then
  yarn build
fi

go mod tidy
