# Go KVDB 

Go KVDB is a lightweight and efficient database system designed to store key-value pairs. It provides three main commands: SET, GET, and DELETE, allowing you to store, retrieve, and delete data in a simple and intuitive manner.

## Features

- **SET**: Store a key-value pair in the database, with an optional time-to-live (TTL) value.
    * Example Command : ```SET <key> <value> <ttl>```
- **GET**: Retrieve the value associated with a given key from the database.
    * Example Command: ```GET <key>```
- **DELETE**: Remove a key-value pair from the database.
    * Example Command: ```DEL <key>```

## Planned Features

We have an exciting roadmap for My Awesome Key-Value Database. Here are some planned features that will be implemented in the near future:

- **Basic TCP Server with SET, GET and DEL commands**: Implementing a TCP server that accepts commands like SET and GET similar to Redis
- **Disaster Recovery using Write-Ahead Logging (WAL)**: Implement a write-ahead logging mechanism to ensure data durability and support disaster recovery scenarios
- **Synchronous and Asynchronous Replication**: Enable replication of data across multiple instances of the database, providing fault tolerance and high availability.
- **Leader Election using Raft**: Introduce leader election functionality based on the Raft consensus algorithm to maintain a consistent state across replicated instances.

## Getting Started

To get started with Go KVDB, follow these steps:

1. Clone the repository: `git clone https://github.com/sriramr98/go_kvdb.git`
2. Start database server: `make run`

## Usage

Once the database server is up and running, you can interact with it using any tcp client like `telnet` or `nc`. Here are some examples:

```shell
$ nc localhost 3000
$ SET foo bar 300
$ OK ( Server Response )
$ GET foo
$ bar ( Server Response )
```