CS3031 Lab 2: Securing the Cloud
================================

This project was developed by Kyrill Zorin.
For this project I developed a secure cloud storage application.  
The project consists of a client and a server, as well as some helpful setup scripts.  
The client and server are written in [Go](https://golang.org).  
The server uses the [RethinkDB Database](http://rethinkdb.com) to store files, keys and users.  
The client provides a CLI interface to connect to the server and carry out actions.  

## Installation and Compilation

To compile the project you will need Go which can be installed by following the [official installation instructions](https://golang.org/doc/install).  
You can then install the additional libraries and dependencies by running the *getdependencies.sh* script in this repo.  
You can then compile the client, server and initDB programs by running the *compile.sh* script in their respective folders.  
To install the RethinkDB database you can follow the [official installation instructions](http://rethinkdb.com/docs/install).  
To initialize the database and create the required tables and indices you will need to run the initDB program once.  
This is necessary before first running the server.  
In order to run the server or initDB programs please make sure that RethinkDB is running.  

## Usage and Configuration
The client, server and initDB programs use a config.toml file for handling configuration.  
The config file is in [TOML](https://github.com/toml-lang/toml) format.  

For the client, valid config paramaters are:  

  * ClientUser (The client user, default = "test")  
  * Server (The cloud server, default = "127.0.0.1:8080")  

For the server, valid config paramaters are:  

  * DBHost (The RethinkDB host, default = "127.0.0.1")  
  * Port = (The port to run the surver on, default = "8080")  

For the initDB program, valid config paramater is:  

  * DBHost (The RethinkDB host, default = "127.0.0.1")  

Th config is loaded on application startup.  
The server and initDB programs can be loaded by simply running them in a Terminal without any arguments.  
The client program needs to be run with arguments otherwise it will simply print usage instructions.  

Below are the valid client commands:  

  client register \<username>  
  client upload \<filepath> \<filename>  
  client download \<user> \<filename> \<outputpath>  
  client share \<filename> \<user>...  
  client revoke \<filename> \<user>...  
  client -h | --help  

\<foo> indicates a variable.  
\... means one or more variables, in this case users.  
The share and revoke commands can be used to act on one or multiple users simultaneously.  
The help screen shows the application name and usage instructions.  

## Implementation and Protocol

The client, server and initDB programs are written in [Go](https://golang.org).  
I primarily relied on Go's standard libraries for functionality but also used several open source libraries.  
The getdependencies.sh and compile.sh scripts are written in [Bash](https://www.gnu.org/software/bash).  
The server uses the [RethinkDB Database](http://rethinkdb.com) to store files, keys and users.  
The initDB program is a simple utility which will create the necessary DB, tables and indices.  
I used the [GoRethink library](https://github.com/dancannon/gorethink) to interface with the database.  
For requests and data transfer I used [JSON encoding](http://www.json.org/).  
The server uses the [HttpRouter library](https://github.com/julienschmidt/httprouter) for multiplexing requests.  
For easily rendering JSON server responses I used the [Render library](https://github.com/unrolled/render).  
The client provides a CLI interface to connect to the server and carry out actions.  
To create the CLI interface I used the [docopt library](https://github.com/docopt/docopt.go) which parses ClI arguments from a usage message.  
For configuration in the programs I used the [Viper library](https://github.com/spf13/viper) which loads configuration from a file.  
The config file is in [TOML](https://github.com/toml-lang/toml) format.  
For cryptography I used Go's standard crypto libraries.  
For encrypting shared secret keys and signing messages I used RSA encryption.  
The private key is generated and stored on the client in an X.509 encoded PEM file.  
It is loaded on startup and used by the client to sign file upload, file sharing and file revocation requests.  
It is also used to decrypt shared secrets before accessing a file.  
The server stores user public keys which are used for verifying signed requests as well as used by clients to encorypt shared secrets.  
For encrypting a file using a shared secret I use AES256 encryption.  
A secure key is randomly generated for the file before encrypting and uploading it.  

Before being able to access other commands a user must first register on the server with their username and public key.  
The client makes a JSON request to the server's */register* HTTP endpoint and receives a response with the status.  
If the username has already been taken by someone else the registration will fail with an error message and  
the user will need to pick a new username to successfully register.  
Once a user is successfully register with the server they will be able to access other commands.  
For getting a user they can make a request to the */users/\<user>* endpoint.  
The server responds with either the requested user or an error message.  
For getting a file the client can make a request to the */users/\<owner>/\<file>* endpoint.  
\<owner> is the file's owner.
The server responds with either the requested file or an error message.  
For getting a list of file users the client can make a request to the */users/\<owner>/\<file>/users* endpoint.  
File users are those who have access to the file by owning an encoded shared secret for it.  
The server responds with either the list of file users or an error message.  
For getting a file key the client can make a request to the */users/\<owner>/\<file>/key/\<user>* endpoint.  
\<user> is the user who is trying to access the file.  
The server responds with either the requested file key or an error message.  

To upload a file the client can make a request to the */uploadfile* endpoint.  
The file is first encrypted on the client using a generated shared secret.  
The upload request is then signed with the client's private key and sent to the server.  
The server verifies the request, saves the file in the database and responds with the status.  
The user then uploads the shared key encrypted with their public key so that they can safely retrieve it at any time.  

To share a file with a user the client encodes the shared secret key using that user's public key.  
A signed request with the key is then made to the */sharefile* endpoint.  
The server verifies the request, saves the file key in the database and responds with the status.  
My implementation provides file access management on a per-file basis to provide greater control to the file owner.  
When downloading a file the user gets the file and the filekey and decrypts the file key using their private key.  
They can then use the decrypted shared secret to decrypt the file itself.  
If a user does not have access to the file the client will simply exit with an error message.  
If a web browser is used to access the file it will show the encrypted data.  

To revoke file access the client sends a signed request to the */revokefile* endpoint.  
The server verifies the request, removes the file key from the database and responds with the status.  
The client then creates a new shared secret, re-encrypts the file using it and re-uploads the file to the server.  
The new shared secret is then shared with the remaining file users so that they can still access the file.  

The code is commented and provides some further imformation regarding the implementation.
