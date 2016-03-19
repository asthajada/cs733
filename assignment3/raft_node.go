package main 

import(
	
	"encoding/gob"
	"time"
	"math/rand"
	"os"
	"encoding/json"
	"strconv"

	cluster "github.com/cs733-iitb/cluster"
	mockcluster "github.com/cs733-iitb/cluster/mock"
	log "github.com/cs733-iitb/log"
)

type Config struct {
	cluster []NetConfig 
	Id int
	LogDir string
	ElectionTimeout int
	HeartbeatTimeout int
	InboxSize  int
	OutboxSize int
	
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
	CommitChannel() <- chan Commit

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
    timeoutCh chan Time
    serverOfCluster cluster.Server
    timer *time.Timer
    heartbeattimer *time.Ticker
    commitCh chan *Commit
    log *log.Log
    LogDir string


}

func (rn *RaftNode) ID() int {
	return rn.server.MyID
	
}

func (rn *RaftNode) LeaderId() int {
	return rn.server.LeaderID
	
}

func (rn *RaftNode) CommittedIndex() int {
	return rn.server.CommitIndex
	
}

func (rn *RaftNode) Append(data []byte) {
	rn.eventCh <- Append{Data: data}
}

func (rn *RaftNode) CommitChannel() chan *Commit {

	return rn.commitCh


}


func getLeader(rn []RaftNode) RaftNode{
	LeaderID:=rn[0].server.LeaderID

	var ldr RaftNode
	for i := 0; i < len(rn); i++ {

		if rn[i].server.MyID==LeaderID {

			ldr=rn[i]
			return ldr
		}else{
			ldr=rn[0] 			//this will never happen , there will always be a leader
		}


	}
	return ldr
	
}

func (rn *RaftNode) ProcessEvents() {

	for {

		var ev Event

		select{
			case ev = <- rn.eventCh:
				actions := rn.server.ProcessEvent(ev)
				rn.doActions(actions)
			case ev = <- rn.timeoutCh:

				actions := rn.server.ProcessEvent(ev)
				rn.doActions(actions)
			case ev := <- rn.serverOfCluster.Inbox():
				actions := rn.server.ProcessEvent(ev.Msg)
				rn.doActions(actions)

			case <- rn.heartbeattimer.C:

				if rn.server.State=="leader" {
					
						rn.timeoutCh <- Timeout{}

				}

		}
		
	}

	
}

func (rn *RaftNode) doActions(actions []Action) {

//	rn.timeoutCh=make(chan Time,100)
//	rn.eventCh=make(chan Event,100)
	//rn.commitCh=make(chan Commit,1000)

	for _,action := range actions {
	switch action.(type) {

			case Send:
				rn.timer.Stop()
				rn.timer = time.AfterFunc(time.Duration(1000+rand.Intn(400))*time.Millisecond, func() { rn.timeoutCh <- Timeout{} })
				
				actionname := action.(Send)
				
				rn.serverOfCluster.Outbox() <- &cluster.Envelope{Pid:actionname.To, Msg:actionname.Event}
			case Alarm:
				
				rn.timer.Stop()
				rn.timer = time.AfterFunc(time.Duration(1000+rand.Intn(400))*time.Millisecond, func() { rn.timeoutCh <- Timeout{} })

			case Commit:
			
				newaction:=action.(Commit)
	
				rn.CommitChannel() <- &newaction

			
			case LogStore:

				lg, _ := log.Open(rn.LogDir+ string(rn.server.MyID))

				logstore:=action.(LogStore)

				lg.Append(logstore.Data)
			

			case StateStore:

				statestore:=action.(StateStore)

				err:=writeFile("statestore"+ strconv.Itoa(rn.server.MyID),statestore )

				if err!=nil {
					
					_,err:=readFile("statestore"+strconv.Itoa(rn.server.MyID))

					if err!=nil {
						
					}
				}
				

		}


	}

}

func createMockCluster(config Config)(*mockcluster.MockCluster,error ) {
	
	clconfig := cluster.Config{Peers:[]cluster.PeerConfig{
		{Id:1}, {Id:2}, {Id:3},{Id:4},{Id:5},
	}}
	cluster, err := mockcluster.NewCluster(clconfig)
	if err != nil {return cluster, err}

	return cluster,nil

}


func makeRafts() []RaftNode {
	
	myNetConfig := make([]NetConfig,0)
	myNetConfig=append(myNetConfig,NetConfig{Id:1,Host:"localhost" ,Port:2000} )
	myNetConfig=append(myNetConfig,NetConfig{Id:2,Host:"localhost" ,Port:3000} )
	myNetConfig=append(myNetConfig,NetConfig{Id:3,Host:"localhost" ,Port:4000} )
	myNetConfig=append(myNetConfig,NetConfig{Id:4,Host:"localhost" ,Port:5000} )
	myNetConfig=append(myNetConfig,NetConfig{Id:5,Host:"localhost" ,Port:6000} )
		

	mynodes:=make([]RaftNode,noOfServers)

	myConfig:=Config{cluster:myNetConfig,Id:1,InboxSize:100,OutboxSize:100,ElectionTimeout:1000,HeartbeatTimeout:200,LogDir:"logdirectory"}
	
	mockcluster1,_:=createMockCluster(myConfig)

	for i := 0; i < noOfServers; i++ {
		mynodes[i]=New(myNetConfig[i].Id,myConfig)
		mynodes[i].serverOfCluster=mockcluster1.Servers[mynodes[i].server.MyID]
		go mynodes[i].ProcessEvents()
	}



	time.Sleep(5*time.Millisecond)

	Register()

	return mynodes


}

func New (id int,config Config) RaftNode {

	var raft RaftNode

	directories:=config.LogDir + strconv.Itoa(id)

	lg, _ := log.Open(directories)

	lastindex := lg.GetLastIndex()
	intlastindex:= int(lastindex)

	
	mypeers:=make([]int,noOfServers)
	myLog := make([]Log,0)
	myLog=append(myLog,Log{0,[]byte("hello")})
	mynextIndex:=make([]int,noOfServers+1)
	mymatchIndex:=make([]int,noOfServers+1)
	myVoteReceived:=make([]int,noOfServers+1)



    for i := 0; i < len(config.cluster); i++ {
    	mypeers[i]=config.cluster[i].Id
    	
    }
  
    s:=RaftStateMachine{
						Term:0,
						VotedFor:-1 ,
						State: "follower", 
						MyID:id,
						PeerID:mypeers,
						Log:myLog,
						CommitIndex:0,
						LeaderID:-1,
						NextIndex:mynextIndex,
 						MatchIndex:mymatchIndex,
 						VoteReceived:myVoteReceived,

					}
	raft.server=s

	for i := 0; i <=intlastindex ; i++ {
		newLog,_:=lg.Get(int64(i))

		newLog1:=newLog.([]byte)
		
		raft.server.Log=append(raft.server.Log,Log{Term:0,Data:newLog1})
		
	}


	eventChannel:=make(chan Event,1000)
	timeoutChannel:=make(chan Time,1000)
	commitChannel := make(chan *Commit,1000)
	raft.eventCh=eventChannel
	raft.timeoutCh=timeoutChannel
	raft.commitCh = commitChannel	
  
  	raft.log=lg

  	raft.LogDir=config.LogDir

 
	rand.Seed(time.Now().UnixNano())
	raft.timer = time.AfterFunc(time.Duration(config.ElectionTimeout+rand.Intn(100))*time.Millisecond, func() { raft.timeoutCh <- Timeout{} })
 
 	rand.Seed(time.Now().UnixNano())
  	raft.heartbeattimer = time.NewTicker(time.Duration(config.HeartbeatTimeout+ rand.Intn(100)) * time.Millisecond)

  	return raft;
	
}


func readFile(filename string) (statestore StateStore, err error){
	
		var f *os.File
		if f, err = os.Open(filename); err != nil {
			return statestore, err
		}
		defer f.Close()
		dec := json.NewDecoder(f)
		if err = dec.Decode(&statestore); err != nil {
			return statestore, err
		}
	
	return statestore, nil


}

func writeFile(filename string , statestore StateStore ) (err error){
	
		var f *os.File
		if f, err = os.Create(filename); err != nil {
			return err
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		if err = enc.Encode(statestore); err != nil {
			return err
		}
	
	return nil
}


func Register() {
	gob.Register(VoteRequest{})
	gob.Register(VoteResponse{})
	gob.Register(AppendEntriesRequest{})
	gob.Register(AppendEntriesResponse{})
	gob.Register(Alarm{})
	gob.Register(Send{})
	gob.Register(LogStore{})
	gob.Register(StateStore{})
}
