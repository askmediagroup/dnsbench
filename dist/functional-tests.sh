#!/bin/bash

set -eo pipefail

echo "Testing local resolver mode..."
echo "example.com" | ./dist/dnsbench run --resolver=local --count=10

echo "Testing external nameserver mode..."
echo "example.com" | ./dist/dnsbench run --resolver=remote --nameserver=8.8.8.8:53 --count=10
