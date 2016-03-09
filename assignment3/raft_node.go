package main 

import(
	"fmt"

	cluster "github.com/cs733-iitb/cluster"
)

type Config struct {
	cluster []NetConfig 
	Id int
	LogDir string
	ElectionTimeout int
	HeartbeatTimeout int
}

type NetConfig struct { 
	Id int
	Host string	
	Port int
}

type Node interface{

	// Client's message to Raft node
	Append([]byte)
       
    // A channel for client to listen on. What goes into Append must come out of here at some point.
	//CommitChannel() <- chan CommitInfo

	// Last known committed index in the log. This could be -1 until the system stabilizes.
	CommittedIndex() int

	// Returns the data at a log index, or an error.
   // Get(index int) (err, []byte)

	// Node's id
	Id() int

	// Id of leader. -1 if unknown
    LeaderId() int

    // Signal to shut down all goroutines, stop sockets, flush log and close it, cancel timers.
	Shutdown()
	

}

type Event interface{

}

type Time interface{

}

type RaftNode struct { // implements Node interface
	server RaftStateMachine
    eventCh chan Event
    timeoutCh <-chan Time
    serverOfCluster cluster.Server

}

func (rn *RaftNode) ID() int {
	return rn.server.myID
	
}

func (rn *RaftNode) LeaderId() int {
	return rn.server.leaderID
	
}

func getLeader(rn []RaftNode) RaftNode{
	leaderID:=rn[0].server.leaderID

	var ldr RaftNode
	for i := 0; i < len(rn); i++ {
		if rn[i].server.myID==leaderID {
			ldr=rn[i]
		}else{
			ldr=rn[0] 			//this will never happen , there will always be a leader
		}


	}
	return ldr
	
}

func makeRafts() []RaftNode {
	
	myNetConfig := make([]NetConfig,0)
	myNetConfig=append(myNetConfig,NetConfig{Id:100,Host:"localhost" ,Port:7070} )
	


	myConfig:=Config{cluster:myNetConfig,Id:100}

	mynode:=New(myConfig)
	myraftnodearray:=make([]RaftNode,0)
	myraftnodearray=append(myraftnodearray,mynode)

	myNetConfig = make([]NetConfig,0)

	myNetConfig=append(myNetConfig,NetConfig{Id:200,Host:"localhost" ,Port:7080} )


	myConfig=Config{cluster:myNetConfig,Id:100}

	mynode=New(myConfig)
	
	myraftnodearray=append(myraftnodearray,mynode)

	return myraftnodearray


}

func New (config Config) RaftNode {

	var raft RaftNode
	
	fmt.Println(config)

	s:=RaftStateMachine{term:1,votedFor:-1 ,state: "follower", myID:config.Id}
	raft.server=s


	configCluster:=cluster.Config{ InboxSize: 100,
  		 Peers: []cluster.PeerConfig{
   		  {Id:100,Address:"localhost:7070"},
   		  {Id:200,Address:"localhost:7080"}, 		
   		}}

  clusterServer,_:=cluster.New(config.Id,configCluster)
  raft.serverOfCluster=clusterServer

  return raft;
	
}


