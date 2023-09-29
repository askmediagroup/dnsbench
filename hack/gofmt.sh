#!/usr/bin/env bash

files=$(go fmt ./...)
if [[ $files ]]; then
    echo 'go fmt was not run on the following files:'
    echo "${files}"
    echo 'run \`go fmt ./...\` to resolve.'
    exit 1
fi

