#!/bin/bash
# Helper script to run tests with XSD validation
# This script sets up the required LD_LIBRARY_PATH for libxml2

# Set library path for linuxbrew libxml2
export LD_LIBRARY_PATH=/home/linuxbrew/.linuxbrew/opt/libxml2/lib:$LD_LIBRARY_PATH

# Run tests with xsdvalidate tag
go test -tags xsdvalidate ./test "$@"
