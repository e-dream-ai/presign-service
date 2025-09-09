#!/bin/bash

echo "=== API Key Generator ==="
echo

if command -v openssl &> /dev/null; then
    echo "Base64 API Key:"
    api_key=$(openssl rand -base64 32)
    echo "   $api_key"
elif command -v base64 &> /dev/null; then
    echo "Base64 API Key:"
    api_key=$(head -c 32 /dev/urandom | base64)
    echo "   $api_key"
else
    echo "Error: Neither openssl nor base64 command found!"
    exit 1
fi

echo
echo "Usage:"
echo "export API_KEY=\"$api_key\""
echo
echo "Note: Store this key securely and never commit it to version control!"
