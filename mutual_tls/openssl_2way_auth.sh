# Steps to create CA, server and client keys + certificates for SSL 2-way authentication
#
# Reference:
# https://gist.github.com/zapstar/4b51d7cfa74c7e709fcdaace19233443
# https://github.com/zapstar/two-way-ssl-c
# https://github.com/iSECPartners/ssl-conservatory
#
# Move to current directory...
cd .
mkdir keys
cd keys

# Generate a self signed certificate for the CA along with a key.
mkdir -p ca/private
chmod 700 ca/private
# NOTE: I'm using -nodes, this means that once anybody gets
# their hands on this particular key, they can become this CA.
openssl req \
    -x509 \
    -nodes \
    -days 3650 \
    -newkey rsa:4096 \
    -keyout ca/private/ca_key.pem \
    -out ca/ca_cert.pem \
    -subj "/C=US/ST=Acme State/L=Acme City/O=Acme Inc./CN=example.com"

# Create server private key and certificate request
mkdir -p server/private
chmod 700 ca/private
openssl genrsa -out server/private/server_key.pem 4096
openssl req -new \
    -key server/private/server_key.pem \
    -out server/server.csr \
    -subj "/C=US/ST=Acme State/L=Acme City/O=Acme Inc./CN=server.example.com"

# Create client private key and certificate request
mkdir -p client/private
chmod 700 client/private
openssl genrsa -out client/private/client_key.pem 4096
openssl req -new \
    -key client/private/client_key.pem \
    -out client/client.csr \
    -subj "/C=US/ST=Acme State/L=Acme City/O=Acme Inc./CN=client.example.com"

# Generate certificates
openssl x509 -req -days 1460 -in server/server.csr \
    -CA ca/ca_cert.pem -CAkey ca/private/ca_key.pem \
    -CAcreateserial -out server/server_cert.pem
openssl x509 -req -days 1460 -in client/client.csr \
    -CA ca/ca_cert.pem -CAkey ca/private/ca_key.pem \
    -CAcreateserial -out client/client_cert.pem

# Now test both the server and the client
# On one shell, run the following
#openssl s_server -CAfile ca/ca_cert.pem -cert server/server_cert.pem -key server/private/server_key.pem -Verify 1
# On another shell, run the following
#openssl s_client -CAfile ca/ca_cert.pem -cert client/client_cert.pem -key client/private/client_key.pem
# Once the negotiation is complete, any line you type is sent over to the other side.
# By line, I mean some text followed by a keyboard return press.
