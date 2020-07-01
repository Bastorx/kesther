#!/bin/sh

# Notes:
# - Since 'curl' is not available, we use 'wget' (through busybox),
# - In order to be authorized, the hostname MUST be 'localhost' or '127.0.0.1'.

busybox wget -qO - http://localhost:${PORT}/reset
echo
