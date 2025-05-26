mkdir -p certs
cd certs

openssl req -new -key server.key -out server.csr -config openssl-san.cnf
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out server.crt -days 365 -sha256 -extensions req_ext -extfile openssl-san.cnf

rm ca.srl server.csr
