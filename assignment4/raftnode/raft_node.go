package raftnode

import(

	"encoding/gob"
	"time"
	"math/rand"
	"os"
	"encoding/json"
	"strconv"
	"fmt"
	"io/ioutil"
	"sync"

	cluster "github.com/cs733-iitb/cluster"
	log "github.com/cs733-iitb/log"
)

type Config struct {
	Cluster []NetConfig
	Id int
	LogDir string
	ElectionTimeout int
	HeartbeatTimeout int
	InboxSize  int
	OutboxSize int
	Ports []int

}

type NetConfig struct {
	Id int
	Host string
	Port int

}

type JsonStateStore struct {
	Term int `json:"Term"`
	VotedFor int `json:"VotedFor"`
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
    shutdownCh chan Event
    serverOfCluster cluster.Server
    timer *time.Timer
    commitCh chan *Commit
    log *log.Log
    LogDir string
    wg *sync.WaitGroup


}

func (rn *RaftNode) ID() int {
	return rn.server.MyID

}

func (rn *RaftNode) LeaderId() int {
	return rn.server.LeaderID

}

func (rn *RaftNode) Shutdown() {
	//add waitgroup
	rn.wg.Add(1)
	close(rn.shutdownCh)  //close the shudown channel
	rn.wg.Wait()
}

func (rn *RaftNode) CommittedIndex() int {
	return rn.server.CommitIndex

}

func (rn *RaftNode) StartRaftNodes(){
	go rn.ProcessEvents()
}
func (rn *RaftNode) Append(data interface{}) {
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

	// Start timer
	rn.timer = time.NewTimer(time.Duration(rn.server.ElectionTimeout + rand.Intn(rn.server.ElectionTimeout)) * time.Millisecond)

	for {

		var ev Event

		select{
			//reading from eventCh
			case ev = <- rn.eventCh:

				actions := rn.server.ProcessEvent(ev)
				rn.doActions(actions)


			case ev = <- rn.timer.C:

				actions := rn.server.ProcessEvent(Timeout{})

				rn.doActions(actions)

			case ev := <- rn.serverOfCluster.Inbox():

				actions := rn.server.ProcessEvent(ev.Msg)
			


				rn.doActions(actions)

			case _, err := <- rn.shutdownCh:
				if !err {


					rn.timer.Stop()
					close(rn.commitCh)

					rn.server.State="follower"


					close(rn.eventCh)
					rn.serverOfCluster.Close()

					rn.wg.Done()

					return
				}





		}

	}


}

func (rn *RaftNode) doActions(actions []Action) {

	for _,action := range actions {
	switch action.(type) {

			case Send:
			
				actionname := action.(Send)
				
				rn.serverOfCluster.Outbox() <- &cluster.Envelope{Pid:actionname.To, Msg:actionname.Event}
			case Alarm:
				//resetting the timer
				rn.timer.Stop()
				alarmAction := action.(Alarm)
				rn.timer.Reset(alarmAction.Time)

			case Commit:



				//output commit obtained from statemachine into the Commit Channel
				newaction:=action.(Commit)
				
				rn.CommitChannel() <- &newaction


			case LogStore:

				logstore:=action.(LogStore)

				rn.log.Append(logstore.Data)


			case StateStore:

				//creating files for persistent storage of State


				statestore:=action.(StateStore)

				//writing statestore into the persistent storage
				err:=writeFile("statestore"+ strconv.Itoa(rn.server.MyID),statestore )

				if err!=nil {

					//reading from the persistent storage
					_,err:=readFile("statestore"+strconv.Itoa(rn.server.MyID))

					if err!=nil {

					}
				}


		}


	}

}


func New (config Config) *RaftNode {

	Register()

	var raft RaftNode

	directories:=config.LogDir + strconv.Itoa(config.Id)

	lg, _ := log.Open(directories)

	lastindex := lg.GetLastIndex()

	intlastindex:= int(lastindex)


	mypeers:=make([]int,noOfServers)
	myLog := make([]Log,0)
	myLog=append(myLog,Log{0,[]byte("hello")})
	mynextIndex:=make([]int,noOfServers)
	mymatchIndex:=make([]int,noOfServers)
	myVoteReceived:=make([]int,noOfServers)


	statestore:="statestore"+strconv.Itoa(config.Id)

	file, e := ioutil.ReadFile("./"+statestore)
	if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }
 

    var jsontype JsonStateStore
    json.Unmarshal(file, &jsontype)
  
	var wait sync.WaitGroup
	raft.wg=&wait


	configCluster := cluster.Config{
        Peers: []cluster.PeerConfig{},
        InboxSize: config.InboxSize,
        OutboxSize: config.OutboxSize,
    }


    for i := 0; i < len(config.Cluster); i++ {
  
    	configCluster.Peers=append(configCluster.Peers,cluster.PeerConfig{Id:config.Cluster[i].Id,Address:fmt.Sprint(config.Cluster[i].Host,":",config.Cluster[i].Port)})
    	mypeers[i]=config.Cluster[i].Id

    }

    s:=RaftStateMachine{
						Term:jsontype.Term,
						VotedFor:jsontype.VotedFor ,
						State: "follower",
						MyID:config.Id,
						PeerID:mypeers,
						Log:myLog,
						CommitIndex:0,
						LeaderID:-1,
						NextIndex:mynextIndex,
 						MatchIndex:mymatchIndex,
 						VoteReceived:myVoteReceived,
						HeartbeatTimeout:config.HeartbeatTimeout,
	    				ElectionTimeout:config.ElectionTimeout,
					}
	raft.server=s

	for i := 0; i <=intlastindex ; i++ {
		
		newLog,_:=lg.Get(int64(i))

		raft.server.Log=append(raft.server.Log,Log{Term:0,Data:newLog})




	}




	eventChannel:=make(chan Event,1000)
	commitChannel := make(chan *Commit,1000)
	shutdownChannel:=make(chan Event,1000)
	raft.eventCh=eventChannel
	raft.commitCh = commitChannel

	raft.shutdownCh= shutdownChannel
  	clusterServer,_:=cluster.New(config.Id,configCluster)
  	raft.serverOfCluster=clusterServer

  	raft.log=lg

  	raft.LogDir=config.LogDir


  	return &raft;

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
