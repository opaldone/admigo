# Admigo
This is a Web Server with an example of the implementation of the admin panel

## It depends on these great packages
![@gorilla](https://avatars.githubusercontent.com/u/489566?s=15&v=4) https://github.com/gorilla/csrf \
![@julienschmidt](https://avatars.githubusercontent.com/u/944947?s=15&v=4) https://github.com/julienschmidt/httprouter \
![@letsencrypt](https://avatars.githubusercontent.com/u/9289019?s=15&v=4) https://pkg.go.dev/golang.org/x/crypto/acme/autocert

## How to install and compile
##### Clonning
```bash
git clone https://github.com/opaldone/admigo.git
```
##### Go to the root "admigo" directory
```bash
cd admigo
```
##### Set the GOPATH variable to the current directory "admigo" to avoid cluttering the global GOPATH directory
```bash
export GOPATH=$(pwd)
```
##### Go to the source folder
```bash
cd src/admigo
```
##### Installing the required Golang packages
```bash
go mod init
```
```bash
go mod tidy
```
##### Return to the "admigo" root directory, You can see the "admigo/pkg" folder that contains the required Golang packages
```bash
cd ../..
```
##### Compiling by the "r" bash script
> r - means "run", b - means "build"
```bash
./r b
```
##### Creating the required folders structure and copying the frontend part by the "u" bash script
> The "u" script is a watching script then for stopping press Ctrl+C \
> u - means "update"
```bash
./u
```
##### You can check the "admigo/bin" folder. It should contain the necessary structure of folders and files
```bash
ls -lash --group-directories-first bin
```
##### Start the server
```bash
./r
```
## About config
The config file is located here __admigo/bin/config.json__
```JavaScript
{
  "appname": "Admigo",  // Just a name of application
  "address": "0.0.0.0", // IP address of the server, zeros mean current host
  "port": 8443, // Port, don't forget to open for firewall
  "static": "static", // The folder that stores the frontend part of the site
  "acme": false, // Set to true if You need to use acme/autocert
  "acmehost": "opaldone.click", // The domain name, set acme: true
  "dirCache": "./certs", // The folder where acme/autocert will store the keys, set acme: true
  "crt": "./certs_local/server.crt", // The path to your HTTPS cert key, set acme: false
  "key": "./certs_local/server.key", // The path to your HTTPS key, set acme: false
  "lang": "en", // Language of the site
  "db": { // Settings of db connection
    "host": "localhost",
    "driver": "postgres",
    "user": "user_login",
    "password": "user_password",
    "dbname": "admigo_db",
    "sslmode": "disable"
  },
  "mail": { // Settings of mail handler
    "from": "some_user@some_site.com",
    "host": "smtp.site.com",
    "password": "email_password",
    "port": 465,
    "username": "user_name@some_site.com",
    "gotourl": "https://opaldone.click:8443" // URL of the site
  }
}
```
