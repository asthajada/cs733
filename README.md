# Raft

Raft is a consensus algorithm like Paxos that is easy to understand. This package implements Raft in Golang. It is somewhat different from original implementation as read requests are not replicated.

#Execution

go build
./assignment4 <server-id>

#Code Structure

###Client Handler

Contains HandleClient class that listens to client requests and gives them response

###Raft Node
Replicates the client requests, maintains effective communication between servers

###Raft State Machine
Actual implementation of raft mechanisms. Handles vote request/respone,append entries request/response, timeouts etc for follower,candidate and leader

