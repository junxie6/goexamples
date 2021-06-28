sudo -u root sh -c 'echo "127.0.0.1       example.com" >> /etc/hosts'

# >>>>>>>>>>>>>>>>>> 根证书 <<<<<<<<<<<<<<<<<<<<<<
# 生成根证书私钥: ca.key
openssl genrsa -out ca.key 2048

# 生成自签名根证书: ca.crt
openssl req -new -key ca.key -x509 -days 3650 -out ca.crt -subj /C=CN/ST=GuangDong/O="Localhost Ltd"/CN="Localhost Root"

# >>>>>>>>>>>>>>>>>> 服务器证书 <<<<<<<<<<<<<<<<<<<<<<
# 生成服务器证书私钥: ca.key
openssl genrsa -out server.key 2048

# 生成服务器证书请求: server.csr
openssl req -new -nodes -key server.key -out server.csr -subj "/C=CN/ST=GuangDong/L=Guangzhou/O=Localhost Server/CN=*.example.com"

# 签名服务器证书: server.crt
#openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt
openssl x509 -req -extfile <(printf "subjectAltName=DNS:example.com,DNS:www.example.com") -days 365 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt

# >>>>>>>>>>>>>>>>>> 客户端证书 <<<<<<<<<<<<<<<<<<<<<<
# 生成客户端证书私钥: ca.key
openssl genrsa -out client.key 2048

# 生成客户端证书请求: client.csr
openssl req -new -nodes -key client.key -out client.csr -subj "/C=CN/ST=GuangDong/L=Guangzhou/O=Localhost Client/CN=client.example.com"

# 签名客户端证书: client.crt
#openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt
openssl x509 -req -extfile <(printf "subjectAltName=DNS:client.example.com") -days 365 -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt