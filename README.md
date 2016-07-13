# Go Irc

This is a simple Internet relay chat application written in go as a coding excercise.
It consists of a utility package and two executables, the server and the client.

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
message tail.

## Client

When you start the client, you will have to provide the IP address of the server,
then the client connects you to the server. All servers listen to the port 50500,
so you don't have to care about correct port. The client automatically detects
your IP address, so that you don't have to care about sending the JOIN request
(as specified above section).

After successfully connecting to a server, the Client will also automatically
query, who is online. You can later do the same manually by typing "/WHOIS".

By default the client sends public messages, where recipients are all other users,
execpt you. You can also send private messages to one or more users by typing
"/PRIVATE id1,id2,id3...: message goes here".