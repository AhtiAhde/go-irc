# Go Irc

This is a simple Internet relay chat application written in go as a coding excercise.
It consists of a utility package and two executables, the server and the client.

## Instructions

Run bin/server command. Then join one or more clients to unique IP address + port
combinations. You need to give the server IP address as a parameter e.g.
`bin/client 172.17.41.218`, the client always expects port 50500 for server.

In order to successfully delvier messages to each client, you need to run the
`/WHOIS` command, which will provide you with a list of connected client Id's.
Then you write a message, press enter and it should be delivered to other clients.
The client currently lacks a feature for sending messages to selected clients.

## Server

When you start the server, it checks out your IP address and starts listening
port 50500, so that port must be free. Server listens for three command, JOIN,
PEOPLE and MESSAGE. They take parameters seoarated by colon.

Examples:
```
JOIN:172.12.43.18:50501
PEOPLE:5
MESSAGE:1,3,4:Hello World!
```
As you can see, the JOIN command takes in parameters IP address and PORT. The response
is a welcome message contains the user id assigned to the user.

PEOPLE command takes your own user id as a parameter (as you will want to know which
other id's exist). The response is a comma separated list of user id's.

MESSAGE command sends a message to all the selected clients. It takes two parameters,
a comma separated list of users and the message. The users id is appended to the
message tail. The command builds a message queue of 1000 messages, which are iterated
once per millisecond.

## Client

When you start the client, you will have to provide the IP address of the server,
then the client connects you to the server. All servers listen to the port 50500,
so you don't have to care about correct port. The client automatically detects
your IP address, so that you don't have to care about sending the JOIN request
(as specified above section).

After successfully connecting to a server, the Client will also automatically
query, who is online. You can later do the same manually by typing "/WHOIS".

By default the client sends public messages, where recipients are all other users,
execpt you. In future versions you can also send private messages to one or more users by typing
"/PRIVATE id1,id2,id3...: message goes here".

## TODOS / Versions

Here is a list of versions and TODOs:

#### v-0.1.0

- Added initial working client (without private messaging)
- Added initial working server, with tests
- Added minor utils package, with tests

#### v-0.1.1

- Fix some of the data types to match the specifications
- Fix the bug, which causes panic upon sending less than six character messages

#### v-0.2.0 (TODO)

- Refactor more utilities
- Refactor the brainfarts away
- Learn more about channels and consider refactorings
- Get more familiar with Go naming conventions and best practices, implement
- Write more tests: improve server tests and write tests for the client

#### v-0.3.0 (TODO)

- Consider resource utilization and Event / Subscriber model (how to implement tests)
- Consider proper data types for everything
- Consider usability tweaks
- Consider fault tolerance for the server (prevent runtime errors, if any)

## Refactoring Plan

[This video](https://www.youtube.com/watch?v=f6kdp27TYZs) seems to provide a nice
approach to effective channel patterns, which would work for this problem area.

First I had based the architecture upon the false idea, that Go is a functional
language, which is far from truth. Yeah... should had considered about it more
seriously when I read about structs.

Anyway, the Connections struct and it's sub-constructs should be reformated in
more object oriented ways and some functions should be refactored as methods.

Then I should move those to a package perhaps?

Testability should remain the same; perhaps I don't have to test the server.go,
but only the Connections "object", which perhaps should be renamed as "server".

Struct Server:
-var userIds
-var addressBook (user addresses)
-var messageQueue (struct, with own methods)
-var clientDialer
+method handleRequest
+method handleMessageQueue (rename handleMessageBuffer, later implement with channels)
+method setClientDialer
+method routeRequest (public for testing purposes, I guess)
-method handleClientJoinRequest
-method handlePeopleRequest
-method handleMessageRequest
-method messageQueue.insertNewMessage

The above "class" seems a bit bloated. Perhaps addressBook should deal with join
and people requests and messageQueue should deal with message requests. Even
though the channels will be the ultimate solution, do not care about them now.
