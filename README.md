# sample-golang-clientcert

#### CA

```
-- create ca private
$ openssl genrsa -out ca-privatekey.pem 2048

-- create ca csr
$ openssl req -new -key ca-privatekey.pem -out ca-csr.pem

-- create ca cert
$ openssl req -x509 -key ca-privatekey.pem -in ca-csr.pem -out ca-crt.pem -days 3560
```


#### Server Certificate

```
-- create server private
$ openssl genrsa -out server-privatekey.pem

-- create server csr
$ openssl req -new -key server-privatekey.pem -out server-csr.pem

-- create server csr
$ openssl x509 -req -CA ca-crt.pem -CAkey ca-privatekey.pem -CAcreateserial -in server-csr.pem -out server-crt.pem -days 3650
```


#### Client Certitificate

```
-- create client private
$ openssl genrsa -out client-privatekey.pem

-- create server csr
$ openssl req -new -key client-privatekey.pem -out client-csr.pem

-- create server csr
$ openssl x509 -req -CA ca-crt.pem -CAkey ca-privatekey.pem -CAcreateserial -in client-csr.pem -out client-crt.pem -days 3650
```


#### build/start server

```
$ cd cmd/server
$ go build
$ ./server -cert=../../cert/server-crt.pem -key=../../cert/server-privatekey.pem
```


#### build/execute client

```
$ cd cmd/client
$ go build
$ ./client -ca=../../cert/ca-crt.pem -cert=../../cert/client-crt.pem -key=../../cert/client-privatekey.pem
```
