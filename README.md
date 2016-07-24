# Go Irc

This is a simple Internet relay chat application written in go as a coding excercise.
It consists of a utility package and two executables, the server and the client.

## Instructions

Run bin/server command. Then join one or more clients with unique IP address + port
combinations. You need to give the server's IP address as a parameter e.g.
`bin/client 172.17.41.218`, the client always expects port 50500 for the server.

In order to successfully delvier messages to each client, you need to run the
`/WHOIS` command, which will provide you with a list of all connected client Id's.
Then you write a message, press enter and it should be delivered to all the other
clients.

## Server

When you start the server, it checks out your IP address and starts listening
port 50500, so that port must be free. Server listens for four command, JOIN,
WHOAMI, PEOPLE and MESSAGE. They take parameters seoarated by colon.

Examples:
```
JOIN:172.12.43.18:50501 // Adds the user to the server AddressBook
WHOAMI:172.12.43.18:50501 // Gets the user id from the AddressBook
PEOPLE:5 // Looks up the AddressBook for all other user ids, except 5
MESSAGE:1,3,4:Hello World! // Sends message "Hello World!" for users 1,3 and 4
```

## Client

When you start the client, you will have to provide the IP address of the server,
as an argument, then the client connects you to the server. All servers listen to
the port 50500, so you don't have to care about the correct port. The client
automatically detects your IP address, so that you don't have to care about
sending the JOIN request (as specified above section).

After successfully connecting to a server, the Client will also automatically
query, who is online. You can later do the same manually by typing "/WHOIS".
This will update "otherUsers"-list, which is used by default as recipient list,
when sending a public message.

So there are three custom commands /WHOAMI /WHOIS and /PRIVATE:
```
/WHOAMI // Gets the user id from the AddressBook
/WHOIS // Looks up the AddressBook for all other user ids, except 5
/PRIVATE 1,3,4 Hello World! // Sends message "Hello World!" for users 1,3 and 4
```

## TODOS / Versions

Here is a list of versions and TODOs:

#### v-0.1.0

- Added initial working client (without private messaging)
- Added initial working server, with tests
- Added minor utils package, with tests

#### v-0.1.1

- Fix some of the data types to match the specifications
- Fix the bug, which causes panic upon sending less than six character messages

#### v-0.2.0

- Refactor more utilities
- Refactor the brainfarts away
- Learn more about channels and consider refactorings
- Get more familiar with Go naming conventions and best practices, implement

#### v-0.3.0 (TODO)

- Write more tests: improve server tests and write tests for the client
- Consider fault tolerance for the server (prevent runtime errors, if any)
- Consider adding tests for the client
- Consider performance issues more in detail
- Consider breaking the router to smaller parts if it feels like a good idea