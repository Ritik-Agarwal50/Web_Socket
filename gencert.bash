#!/bin/bash
echo "this script is used to generate self-signed certificate"
echo "creating server.key"
openssl genrsa -out server.key 2024
openssl ecparam -genkey -name secp384r1 -out server.key
echo "creating server.crt"
openssl req -new -x509 -sha256 -key server.key -out server.crt -batch -days 3650