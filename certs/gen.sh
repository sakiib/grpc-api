#!/bin/bash

# remove all the pem files before generating new files
rm -rf *.pem

# --------------------------------------------------------------------------------------------

# 1. Generate CA's private key and self-signed certificate
# addting the -nodes will not require any passphrase, remove it to enter passphrase
# if we add -nodes options then the private key will not be encrypted anymore
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem \
    -subj "/C=BD/ST=Dhaka/L=Uttara/O=AppsCode/OU=SWE/CN=AC/emailAddress=sakibalamin@appscode.com"

echo "Here's the CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text

# --------------------------------------------------------------------------------------------

# 2. Generate web server's private key and certificate signing request (CSR)
# addting the -nodes will not require any passphrase, remove it to enter passphrase
# if we add -nodes options then the private key will not be encrypted anymore
# similar to the step 1. remove the x509, days options & add server info in subject ref
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem \
    -subj "/C=BD/ST=SDhaka/L=SUttara/O=SAppsCode/OU=SSWE/CN=SAC/emailAddress=server@appscode.com"

# --------------------------------------------------------------------------------------------

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
# by default this certificate is valid for 30 days, add -days <day> for changing it: server-req.pem -days 60
# add server-ext.cnf to add the subject alternative names, DNS, Email, IP etc.
openssl x509 -req -in server-req.pem -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem \
    -extfile server-ext.cnf

# --------------------------------------------------------------------------------------------

# printing output in text format
echo "Here's the server's self-signed cert"
openssl x509 -in server-cert.pem -noout -text

# --------------------------------------------------------------------------------------------

# verifying certificate
# pass the ca cert (ca-cert.pem) & the cert (server-cert.pem) we want to verify
# it'll return OK if certificate is verified
openssl verify -CAfile ca-cert.pem  server-cert.pem