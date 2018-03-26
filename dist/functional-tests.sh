#!/bin/bash

set -eo pipefail

echo "Testing local resolver mode..."
echo "example.com" | ./dist/dnsbench run local --concurrency 1 --count 10

echo "Testing external nameserver mode..."
echo "example.com" | ./dist/dnsbench run remote 8.8.8.8 --concurrency 1 --count 10
