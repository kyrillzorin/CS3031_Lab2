CS3031 Lab 2: Securing the Cloud
================================

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
