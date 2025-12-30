# read from here
# https://devopscube.com/create-self-signed-certificates-openssl/

cd /mnt/terik/a_my/admigo_project/admigo/howtocert

openssl req -x509 \
            -sha256 -days 356 \
            -nodes \
            -newkey rsa:2048 \
            -subj "/CN=admigo.so/C=RU/L=Perm" \
            -keyout ./kk/rootCA.key -out ./kk/rootCA.crt

openssl genrsa -out ./kk/server.key 2048

openssl req -new -key ./kk/server.key -out ./kk/server.csr -config csr.conf

openssl x509 -req \
    -in ./kk/server.csr \
    -CA ./kk/rootCA.crt -CAkey ./kk/rootCA.key \
    -CAcreateserial -out ./kk/server.crt \
    -days 365 \
    -sha256 -extfile cert.conf

# set in
# google chrome
chrome://settings/certificates?search=auth
# in firefox
about:preferences#privacy
	# find
	Certificates / View Certificates

# copy here from the kk
/mnt/terik/a_my/admigo_project/admigo/bin/certs_local
server.crt
server.key
