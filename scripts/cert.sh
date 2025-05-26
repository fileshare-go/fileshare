#!/bin/bash

# 创建证书目录
mkdir -p certs && cd certs

# 生成CA私钥
openssl genrsa -out ca.key 2048

# 生成CA证书
openssl req -x509 -new -nodes -key ca.key -sha256 -days 1024 -out ca.crt -subj "/CN=GRPCCA"

# 生成服务器私钥
openssl genrsa -out server.key 2048

# 生成服务器证书签名请求
openssl req -new -key server.key -out server.csr -subj "/CN=localhost"

# 使用CA签署服务器证书
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 500 -sha256

# 生成客户端私钥
openssl genrsa -out client.key 2048

# 生成客户端证书签名请求
openssl req -new -key client.key -out client.csr -subj "/CN=client"

# 使用CA签署客户端证书
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 500 -sha256

# 清理
rm server.csr client.csr ca.srl
