#!/bin/bash

set -eo pipefail

echo "Testing local resolver mode..."
echo "example.com" | ./dist/dnsbench -local -concurrency 1 -count 10

echo "Testing external nameserver mode..."
echo "example.com" | ./dist/dnsbench -concurrency 1 -count 10 8.8.8.8
