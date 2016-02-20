package main 

import (
	"strconv"
	
)
	
type node struct {			


  term int				//local term of the server
  votedFor int  		// id of the server to whom vote is given,-1 if not given to anyone
  log []Log
  commitIndex int
  state string
  leaderID int			//ID of the leader
 // nextIndex []int
 // matchIndex []int 
  //eventCh chan event
  myID int
  peerID []int 
  voteReceived []int      // status of votes of peers

  }	

type Log struct{
	term int
	data []byte
}

const noOfServers int=5

type event interface{}

const serverID int=123;

type  VoteRequest  struct{
	term int				//term for which it is candidate for election
	candidateID int         // ID of the server requesting a vote
	lastLogIndex int 		//index of the last log entry
	lastLogTerm int      //last term stored in its log
}

type VoteResponse struct{
	from int           //
	term int
	voteGranted int     //-1 if not voted, 0 if no and 1 if yes
}

type AppendEntriesRequest struct{
	term int
	leaderID int
	prevLogIndex int
	prevLogTerm int
	log []Log
	leaderCommit int

}

type AppendEntriesResponse struct{
	term int
	success bool
	from int

}
type Append struct{
	data []byte

}
type Timeout struct{

}

type Action interface{

}

type Send struct{
	to int
	event interface{ }
	
}

type Commit struct{
	index int
	data []byte
	err string

}

type Alarm struct{
	time int 
}

type LogStore struct{
	index int
	data []byte

}

type  StateStore struct{
	term int
	votedFor int
}



//if equal then check votedfor== candidate iD then grant


func onVoteRequest(obj VoteRequest,s *node) []Action{

	action:=make([]Action,0)

	if (s.state=="follower"){

		if(s.term>obj.term){
			//follower gets the request for lower order term
			//reject the vote

			action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID, term: s.term,voteGranted: 0}})


		} else if (s.term==obj.term){
			//term of follower and candidate are equal

			if(s.log[len(s.log)-1].term>obj.lastLogTerm){
				//reject the vote
				action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID, term: s.term,voteGranted: 0}})

			} else if (s.log[len(s.log)-1].term==obj.lastLogTerm && len(s.log)-1>obj.lastLogIndex) {
				//follower is more knowledgeable than candidate
				//reject the vote
				
				action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID, term: s.term,voteGranted: 0}})
				
			} else{

				if (s.votedFor==-1 || s.votedFor==obj.candidateID){
				s.term=obj.term
				s.votedFor=obj.candidateID


				action=append(action,StateStore{term:s.term,votedFor:obj.candidateID})

				//grant the vote
				
				action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID,term: s.term,voteGranted: 1}})
				

				}else{

				//follower has already given vote to another candidate so reject the vote
				action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID,term: s.term,voteGranted: 0}})
				}
			}

		} else if (s.term < obj.term) {

			if(s.log[len(s.log)-1].term>obj.lastLogTerm){

				
				s.term=obj.term
				s.votedFor=-1

				action=append(action,StateStore{term:s.term,votedFor:s.votedFor})
				//reject the vote

				action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID, term: s.term,voteGranted: 0}})

			}else if (s.log[len(s.log)-1].term==obj.lastLogTerm && len(s.log)-1>obj.lastLogIndex) {
				//follower has more log entries than candidate
				
				s.term=obj.term
				s.votedFor=-1

				action=append(action,StateStore{term:s.term,votedFor:s.votedFor})

				//reject the vote

				action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID, term: s.term,voteGranted: 0}})
				
			} else{

				if (s.votedFor==-1 || s.votedFor==obj.candidateID){
				s.term=obj.term
				s.votedFor=obj.candidateID

				action=append(action,StateStore{term:s.term,votedFor:obj.candidateID})

				//grant the vote
				
				action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID,term: s.term,voteGranted: 1}})
				

				}else{
				//follower has already given vote to another candidate so reject the vote
				s.term=obj.term
				s.votedFor=-1

				action=append(action,StateStore{term:s.term,votedFor:s.votedFor})

				//reject the vote

				action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID,term: s.term,voteGranted: 0}})
				}
			}

			
		}
			
	}else if (s.state=="candidate" || s.state=="leader") {
		if (s.term<obj.term) {

			s.state="follower"


			if(s.log[len(s.log)-1].term>obj.lastLogTerm){

				
				s.term=obj.term
				s.votedFor=-1

				action=append(action,StateStore{term:s.term,votedFor:s.votedFor})
				//reject the vote

				action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID, term: s.term,voteGranted: 0}})

			}else if (s.log[len(s.log)-1].term==obj.lastLogTerm && len(s.log)-1>obj.lastLogIndex) {
				//follower has more log entries than candidate

				
				s.term=obj.term
				s.votedFor=-1

				action=append(action,StateStore{term:s.term,votedFor:s.votedFor})

				//reject the vote

				action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID, term: s.term,voteGranted: 0}})
				
			} else{


				if (s.votedFor==-1 || s.votedFor==obj.candidateID){
				s.term=obj.term
				s.votedFor=obj.candidateID

				action=append(action,StateStore{term:s.term,votedFor:obj.candidateID})

				//grant the vote
				
				action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID,term: s.term,voteGranted: 1}})
				

				}else{

				//follower has already given vote to another candidate so reject the vote
				s.term=obj.term
				s.votedFor=-1

				action=append(action,StateStore{term:s.term,votedFor:s.votedFor})

				//reject the vote

				action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID,term: s.term,voteGranted: 0}})
				}
			}


		} else if (s.term>=obj.term) {
			//reject the vote
			
			action=append(action,Send{to: obj.candidateID,event:VoteResponse{from:s.myID,term: s.term,voteGranted: 0}})
			
		}

	}
return action
}

func onVoteResponse(obj VoteResponse,s *node) []Action{

	var majority,count0,count1 int
	majority= (len(s.peerID)+1)/2
	count0=0
	count1=1
	action:=make([]Action,0)
	myLog := make([]Log,0)

	if (s.state=="follower") {

		//follower is not supposed to handle VoteResponse 

		//update term if any higher order term

		if(s.term<obj.term){
			s.term=obj.term
			action=append(action,StateStore{term:s.term,votedFor:-1})

		}

	} else if (s.state=="leader") {
		
		//if leader gets a VoteResponse it will just ignore it and it has already changed itself from candidate to leader, it cannot get VoteResponse from any higher order term

	}else if (s.state=="candidate") {

		if(obj.term==s.term){
			s.voteReceived[obj.from-1]=obj.voteGranted

			for i := 1; i < len(s.voteReceived); i++ {
				if(s.voteReceived[i]==1){
					
					count1++
					
				}else if (s.voteReceived[i]==0) {
					
					count0++
				}

				if(count1==majority){
				
					s.leaderID=s.myID
					s.state="leader"
					//send empty append entries request

					
					for i := 0; i < len(s.peerID); i++ {
		
						action=append(action, Send{to:s.peerID[i], event:AppendEntriesRequest{log:myLog,term:s.term,leaderID:s.myID}})
						

					}
			
					action=append(action,Alarm{25})
					break

				} else if (count0==majority) {
					s.state="follower"
					action=append(action,Alarm{25})
					break
				}
				
			}

		}

	}

	return action
}

func onApendEntriesRequest(obj AppendEntriesRequest,s *node)[]Action {

	action:=make([]Action,0)

	if(s.state=="follower"){

		if(s.term>obj.term){
			

			action=append(action,Send{to: obj.leaderID,event:AppendEntriesResponse{term: s.term,success: false,from:s.myID}})

		}else{
			
			

			if(len(s.log)-1!=obj.prevLogIndex){

				//different last log index then do not append

				action=append(action,Send{to: obj.leaderID,event:AppendEntriesResponse{term: s.term,success: false,from:s.myID}})

			} else {
				
				if( obj.prevLogIndex != -1 && len(s.log) != 0 && s.log[len(s.log)-1].term != obj.prevLogTerm){
					
					// last log index matches but different terms at last Index
					action=append(action,Send{to: obj.leaderID,event:AppendEntriesResponse{term: s.term,success: false,from:s.myID}})

				}else{

					
					//delete entries from lastLogIndex and then append

					index := obj.prevLogIndex + 1
					

					if len(obj.log) > 0 {
				
						action = append(action, LogStore{index: index, data: obj.log[obj.prevLogIndex+1].data})
						s.log = append(s.log, obj.log...)
					}

					s.leaderID = obj.leaderID
					//update term and voted for if term in append entry rpc is greater than follower current term
					if obj.term > s.term {
						s.term = obj.term
						s.votedFor = -1
						
						action = append(action, StateStore{term: s.term, votedFor: s.votedFor})
					}
					
					if obj.leaderCommit > s.commitIndex {
					
						s.commitIndex = minimum(obj.leaderCommit, index-1)

						
						action = append(action, Commit{index: s.commitIndex, data: s.log[s.commitIndex].data})
					}

					action = append(action, Send{to: obj.leaderID, event: AppendEntriesResponse{term: s.term, success: true,from:s.myID}})


				}

			}
				 
		}

	} else if (s.state=="candidate") {

		
		if (obj.term<s.term) {

			action=append(action,Send{to: obj.leaderID,event:AppendEntriesResponse{term: s.term,success: false,from:s.myID}})
			
		}else if(obj.term == s.term) {

			if(len(s.log)-1!=obj.prevLogIndex){

				//different last log index then do not append

				action=append(action,Send{to: obj.leaderID,event:AppendEntriesResponse{term: s.term,success: false,from:s.myID}})

			} else {
				
				if( obj.prevLogIndex != -1 && len(s.log) != 0 && s.log[len(s.log)-1].term != obj.prevLogTerm){
					
					// last log index matches but different terms at last Index
					action=append(action,Send{to: obj.leaderID,event:AppendEntriesResponse{term: s.term,success: false,from:s.myID}})

				}else{

					//delete entries from lastLogIndex and then append

					index := obj.prevLogIndex + 1
					

					if len(obj.log) > 0 {
				
						action = append(action, LogStore{index: index, data: obj.log[obj.prevLogIndex+1].data})
						s.log = append(s.log, obj.log...)
					}

					s.leaderID = obj.leaderID
					//update term and voted for if term in append entry rpc is greater than follower current term
					if obj.term > s.term {
						s.term = obj.term
						s.votedFor = -1
						
						action = append(action, StateStore{term: s.term, votedFor: s.votedFor})
					}
					
					if obj.leaderCommit > s.commitIndex {
					
						s.commitIndex = minimum(obj.leaderCommit, index-1)

						
						action = append(action, Commit{index: s.commitIndex, data: s.log[s.commitIndex].data})
					}


					action = append(action, Send{to: obj.leaderID, event: AppendEntriesResponse{term: s.term, success: true,from:s.myID}})
					s.state="follower"


				}

			}


		}else{
	
		//when obj.term > s.term
				s.term=obj.term
				s.votedFor=-1
				s.state="follower"
				action = append(action, StateStore{term: s.term, votedFor: s.votedFor})

			
			if(len(s.log)-1!=obj.prevLogIndex){

				//different last log index then do not append

				action=append(action,Send{to: obj.leaderID,event:AppendEntriesResponse{term: s.term,success: false,from:s.myID}})

			} else {
				
				if( obj.prevLogIndex != -1 && len(s.log) != 0 && s.log[len(s.log)-1].term != obj.prevLogTerm){
					
					// last log index matches but different terms at last Index
					action=append(action,Send{to: obj.leaderID,event:AppendEntriesResponse{term: s.term,success: false,from:s.myID}})

				}else{

					
					//delete entries from lastLogIndex and then append

					index := obj.prevLogIndex + 1
					

					if len(obj.log) > 0 {
				
						action = append(action, LogStore{index: index, data: obj.log[obj.prevLogIndex+1].data})
						s.log = append(s.log, obj.log...)
					}

					s.leaderID = obj.leaderID
					
					if obj.leaderCommit > s.commitIndex {
					
						s.commitIndex = minimum(obj.leaderCommit, index-1)

						
						action = append(action, Commit{index: s.commitIndex, data: s.log[s.commitIndex].data})
					}


					action = append(action, Send{to: obj.leaderID, event: AppendEntriesResponse{term: s.term, success: true,from:s.myID}})
					


				}

			}

		}
	}else{
		//s.state=leader

		//getting request for higher order term
		if(s.term<obj.term){
			s.state="follower"
			s.term=obj.term
			s.votedFor=-1

			action = append(action, StateStore{term: s.term, votedFor: s.votedFor})

		}

		//leader cannot get request from same order or higher order term

	}


	
	return action
}

func onAppendEntriesResponse(obj AppendEntriesResponse, s *node)[]Action {
	
	action:=make([]Action,0)
	if(s.state=="follower"){
		if (s.term < obj.term) {
			s.term=obj.term
			s.votedFor=-1

			action = append(action, StateStore{term: s.term, votedFor: s.votedFor})
		}
	}else if (s.state=="candidate") {

		if (s.term < obj.term) {
			s.term=obj.term
			s.votedFor=-1

			action = append(action, StateStore{term: s.term, votedFor: s.votedFor})
		}
		
	}

return action

}


func onTimeout(obj Timeout,s *node) []Action{
	

	action:=make([]Action,0)

	if (s.state=="follower" || s.state=="candidate") {
		s.state="candidate"

		//increment the term
		s.term=s.term+1

		//vote for self
		s.votedFor=s.myID
		action=append(action,StateStore{term:s.term,votedFor:s.myID})

		//reset election timer
		
		action=append(action,Alarm{time: 25 })

        //send VoteRequests to everyone
		for i := 0; i < len(s.peerID); i++ {
		
			action=append(action, Send{to:s.peerID[i], event:VoteRequest{term:s.term , candidateID:s.myID, lastLogIndex:len(s.log)-1, lastLogTerm:s.log[len(s.log)-1].term}})
			
		}

	
	} else{
		//s.state=leader

	

		myLog := make([]Log,0)

		action=append(action,Alarm{time: 25 })

		//send empty append entries request
		for i := 0; i < len(s.peerID); i++ {
		
			action=append(action, Send{to:s.peerID[i], event:AppendEntriesRequest{log:myLog,term:s.term}})
			
		}
		
	}


return action

}

func onAppend(obj Append,s *node)[]Action{

action:=make([]Action,0)

	if(s.state=="follower" || s.state=="candidate"){


		action=append(action,Commit{data:obj.data,err:"leader:"+strconv.Itoa(s.leaderID) })  //as follower or candidate cannot append to log send leaderID to the client
		

	}else{
		
			length:=len(s.log)

			s.log=append(s.log,Log{s.term,obj.data})
		
			action = append(action, LogStore{index: length, data: s.log[length].data})

	}


return action
}


func (s *node) ProcessEvent(getEvent interface{}) []Action{

	action:=make([]Action,0)

	

		switch getEvent.(type){
		case VoteRequest:
			objVoteRequest:=getEvent.(VoteRequest)		//creating object of type VoteRequest from the event
			action=onVoteRequest(objVoteRequest,s)
		case AppendEntriesRequest:
			objAppendEntriesRequest:=getEvent.(AppendEntriesRequest)		//creating object of type AppendEntriesRequest from the event
			action=onApendEntriesRequest(objAppendEntriesRequest,s)
		case Timeout:
			objTimeout:=getEvent.(Timeout)
			action=onTimeout(objTimeout,s)
		case VoteResponse:
			objVoteResponse:=getEvent.(VoteResponse)
			action=onVoteResponse(objVoteResponse,s)
		case AppendEntriesResponse:
			objAppendEntriesResponse:=getEvent.(AppendEntriesResponse)
			action=onAppendEntriesResponse(objAppendEntriesResponse,s)
		case Append:
			objAppend:=getEvent.(Append)
			action=onAppend(objAppend,s)
		

		}


return action

}

func  minimum(a int,b int)(int) {
	if (a<b){
		return a
	} else{
		return b
	}

}

