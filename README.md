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
##### Go to the folder with source code
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
> The "u" script reads sub file "watch_files" \
> C_FOLDERS - the array of folders to simple copy \
> W_FILES - the array of files whose changes are tracked
```bash
./watch_files
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
  // Just a name of application
  "appname": "Admigo",
  // IP address of the server, zeros mean current host
  "address": "0.0.0.0",
  // Port, don't forget to open it in the firewall
  "port": 8443,
  // The folder that stores the frontend part of the site
  "static": "static",
  // Set "acme": true if You need to use acme/autocert
  // false - if You use self-signed certificates
  "acme": false,
  // The array of domain names, set "acme": true
  "acmehost": [
    "opaldone.click",
    "206.189.101.23",
    "www.opaldone.click"
  ],
  // The folder where acme/autocert will store the keys, set "acme": true
  "dirCache": "./certs",
  // The paths to your self-signed HTTPS keys, set "acme": false
  "crt": "./certs_local/server.crt",
  "key": "./certs_local/server.key",
  // Language of the site
  "lang": "en",
  // Settings of db connection
  "db": {
    "host": "localhost",
    "driver": "postgres",
    "user": "user_login",
    "password": "user_password",
    "dbname": "admigo_db",
    "sslmode": "disable"
  },
  // Settings of the mail handler
  "mail": {
    "from": "some_user@some_site.com",
    "host": "smtp.site.com",
    "password": "email_password",
    "port": 465,
    "username": "user_name@some_site.com",
    // URL of the site
    "gotourl": "https://opaldone.click:8443"
  }
  // Settings of the map handler
  "map": {
    // URL for the web socket to exchange location positions
    "ws": "wss://admigo.so:8088",
    // Start point of map focus
    "startpoint": "57.989287,56.213889",
    // URL of the service to make routers
    "routeurl": "https://api.openrouteservice.org/v2/directions/driving-car/geojson",
    // URL key for api.openrouteservice.org
    "routekey": "xxx",
    // The city to search addresses
    "city": "London",
    // The language used in addresses
    "lang": "en"
  }
}
```
