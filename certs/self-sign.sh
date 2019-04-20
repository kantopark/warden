#!/usr/bin/env bash
# Generates a self-signed certificate for local testing purposes!
# Do not use in production!

openssl genrsa -out domain.key 2048
openssl rsa -in domain.key -out domain.key
openssl req -sha256 -new -key domain.key -out server.csr -subj '/CN=localhost'
openssl x509 -req -sha256 -days 365 -in server.csr -signkey domain.key -out domain.crt

rm ./server.csr

