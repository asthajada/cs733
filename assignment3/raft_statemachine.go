package main  

import (
	"strconv"
	"sort"
	"time"
	"errors"
	
)
	
type RaftStateMachine struct {			


  Term int				//local Term of the server
  VotedFor int  		// id of the server To whom vote is given,-1 if not given To anyone
  Log []Log
  CommitIndex int
  State string
  LeaderID int			//ID of the leader
  NextIndex []int
  MatchIndex []int 
  MyID int
  PeerID []int 
  VoteReceived []int      // status of votes of peers

  }	

type Log struct{
	Term int
	Data []byte
}

const noOfServers int=5

//type Event interface{}

const serverID int=123;

type  VoteRequest  struct{
	Term int				//Term for which it is candidate for election
	CandidateID int         // ID of the server requesting a vote
	LastLogIndex int 		//Index of the last Log entry
	LastLogTerm int      //last Term stored in its Log
}

type VoteResponse struct{
	From int           //To know From where the vote is received
	Term int
	VoteGranted int     //-1 if not voted, 0 if no and 1 if yes
}

type AppendEntriesRequest struct{
	Term int
	LeaderID int
	PrevLogIndex int
	PrevLogTerm int
	Log []Log
	LeaderCommit int

}

type AppendEntriesResponse struct{
	Term int
	Success bool
	From int 		//To know From where the AppendEntriesResponse is received
	Count int 		//To keep track of how many entries were sent
	LastLogIndex int
}
type Append struct{
	Data []byte

}
type Timeout struct{

}

type Action interface{

}

type Send struct{
	To int
	Event interface{ }
	
}

type Commit struct{
	Index int
	Data []byte
	Err error

}

type Alarm struct{
	Time time.Time 
}

type LogStore struct{
	Index int
	Data []byte

}

type  StateStore struct{
	Term int
	VotedFor int
}



//if equal then check votedfor== candidate iD then grant


func onVoteRequest(obj VoteRequest,s *RaftStateMachine) []Action{




	action:=make([]Action,0)

	if (s.State=="follower"){

		if(s.Term>obj.Term){
			//follower gets the request for lower order Term
			//reject the vote

			action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID, Term: s.Term,VoteGranted: 0}})


		} else if (s.Term==obj.Term){
			//Term of follower and candidate are equal

			if(s.Log[len(s.Log)-1].Term>obj.LastLogTerm){
				//reject the vote
				action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID, Term: s.Term,VoteGranted: 0}})

			} else if (s.Log[len(s.Log)-1].Term==obj.LastLogTerm && len(s.Log)-1>obj.LastLogIndex) {
				//follower is more knowledgeable than candidate
				//reject the vote
				
				action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID, Term: s.Term,VoteGranted: 0}})
				
			} else{

				if (s.VotedFor==-1 || s.VotedFor==obj.CandidateID){
				s.Term=obj.Term
				s.VotedFor=obj.CandidateID



				for i := 0; i < len(s.VoteReceived); i++ {

					s.VoteReceived[i]=-1

				}

				action=append(action,StateStore{Term:s.Term,VotedFor:obj.CandidateID})

				//grant the vote
				
				action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID,Term: s.Term,VoteGranted: 1}})
				

				}else{

				//follower has already given vote To another candidate so reject the vote
				action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID,Term: s.Term,VoteGranted: 0}})
				}
			}

		} else if (s.Term < obj.Term) {

			if(s.Log[len(s.Log)-1].Term>obj.LastLogTerm){

				
				s.Term=obj.Term
				s.VotedFor=-1

				for i := 0; i < len(s.VoteReceived); i++ {

					s.VoteReceived[i]=-1
				}


				action=append(action,StateStore{Term:s.Term,VotedFor:s.VotedFor})
				//reject the vote

				action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID, Term: s.Term,VoteGranted: 0}})

			}else if (s.Log[len(s.Log)-1].Term==obj.LastLogTerm && len(s.Log)-1>obj.LastLogIndex) {
				//follower has more Log entries than candidate
				
				s.Term=obj.Term
				s.VotedFor=-1

				for i := 0; i < len(s.VoteReceived); i++ {
					s.VoteReceived[i]=-1

				}


				action=append(action,StateStore{Term:s.Term,VotedFor:s.VotedFor})

				//reject the vote

				action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID, Term: s.Term,VoteGranted: 0}})
				
			} else{

				if (s.VotedFor==-1 || s.VotedFor==obj.CandidateID){
				s.Term=obj.Term
				s.VotedFor=obj.CandidateID

				for i := 0; i < len(s.VoteReceived); i++ {

					s.VoteReceived[i]=-1

				}


				action=append(action,StateStore{Term:s.Term,VotedFor:obj.CandidateID})

				//grant the vote
				
				action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID,Term: s.Term,VoteGranted: 1}})
				

				}else{
				//follower has already given vote To another candidate so reject the vote
				s.Term=obj.Term
				s.VotedFor=-1


				for i := 0; i < len(s.VoteReceived); i++ {

					s.VoteReceived[i]=-1
				}

				action=append(action,StateStore{Term:s.Term,VotedFor:s.VotedFor})

				//reject the vote

				action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID,Term: s.Term,VoteGranted: 0}})
				}
			}

			
		}
			
	}else if (s.State=="candidate" || s.State=="leader") {
		if (s.Term<obj.Term) {

			s.State="follower"


			if(s.Log[len(s.Log)-1].Term>obj.LastLogTerm){

				
				s.Term=obj.Term
				s.VotedFor=-1

				for i := 0; i < len(s.VoteReceived); i++ {

					s.VoteReceived[i]=-1
				}


				action=append(action,StateStore{Term:s.Term,VotedFor:s.VotedFor})
				//reject the vote

				action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID, Term: s.Term,VoteGranted: 0}})

			}else if (s.Log[len(s.Log)-1].Term==obj.LastLogTerm && len(s.Log)-1>obj.LastLogIndex) {
				//follower has more Log entries than candidate

				
				s.Term=obj.Term
				s.VotedFor=-1

				for i := 0; i < len(s.VoteReceived); i++ {

					s.VoteReceived[i]=-1
				}

				action=append(action,StateStore{Term:s.Term,VotedFor:s.VotedFor})

				//reject the vote

				action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID, Term: s.Term,VoteGranted: 0}})
				
			} else{	


				if (s.VotedFor==-1 || s.VotedFor==obj.CandidateID){
				s.Term=obj.Term
				s.VotedFor=obj.CandidateID

				for i := 0; i < len(s.VoteReceived); i++ {

					s.VoteReceived[i]=-1

				}

				action=append(action,StateStore{Term:s.Term,VotedFor:obj.CandidateID})

				//grant the vote
				
				action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID,Term: s.Term,VoteGranted: 1}})
				

				}else{

				//follower has already given vote To another candidate so reject the vote
				s.Term=obj.Term
				s.VotedFor=-1

				for i := 0; i < len(s.VoteReceived); i++ {

					s.VoteReceived[i]=-1

				}

				action=append(action,StateStore{Term:s.Term,VotedFor:s.VotedFor})

				//reject the vote

				action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID,Term: s.Term,VoteGranted: 0}})
				}
			}


		} else if (s.Term>=obj.Term) {
			//reject the vote
			
			action=append(action,Send{To: obj.CandidateID,Event:VoteResponse{From:s.MyID,Term: s.Term,VoteGranted: 0}})
			
		}

	}
return action
}

func onVoteResponse(obj VoteResponse,s *RaftStateMachine) []Action{
	

	var majority,Count0,Count1 int
	majority= (len(s.PeerID)+1)/2
	Count0=0
	Count1=1
	action:=make([]Action,0)
	myLog := make([]Log,0)

	if (s.State=="follower") {

		//follower is not supposed To handle VoteResponse 

		//update Term if any higher order Term

		if(s.Term<obj.Term){
			s.Term=obj.Term


			for i := 0; i < len(s.VoteReceived); i++ {

				s.VoteReceived[i]=-1
			}
			action=append(action,StateStore{Term:s.Term,VotedFor:-1})

		}

	} else if (s.State=="leader") {
		
		//if leader gets a VoteResponse it will just ignore it and it has already changed itself From candidate To leader, it cannot get VoteResponse From any higher order Term

	}else if (s.State=="candidate") {

		if(obj.Term==s.Term){
				
			s.VoteReceived[obj.From-1]=obj.VoteGranted

			for i := 0; i < len(s.VoteReceived); i++ {
				if(s.VoteReceived[i]==1){
					
					Count1++
					
				}else if (s.VoteReceived[i]==0) {
					
					Count0++
				}

				if(Count1==majority){

					s.LeaderID=s.MyID
					s.State="leader"

				

					//send empty append entries request
					myLog=make([]Log,0)

					for i := 0; i < len(s.PeerID); i++ {
					

						action=append(action, Send{To:s.PeerID[i], Event:AppendEntriesRequest{Log:myLog,Term:s.Term,LeaderID:s.MyID}})
						

					}
			
				//	action=append(action,Alarm{Time: time.Now().Add(time.Duration(100) * time.Millisecond) })
					action=append(action,Alarm{Time: time.Now().Add(time.Duration(10) * time.Millisecond) })
					
					break

				} else if (Count0==majority) {
					s.State="follower"
					action=append(action,Alarm{Time: time.Now().Add(time.Duration(10) * time.Millisecond) })
					break
				}
				
			}

		}

	}

	return action
}

func onAppendEntriesRequest(obj AppendEntriesRequest,s *RaftStateMachine)[]Action {

	action:=make([]Action,0)

	if(s.State=="follower"){

		if len(obj.Log)==0 {
			s.LeaderID=obj.LeaderID


			if obj.LeaderCommit > s.CommitIndex {
					
						s.CommitIndex = minimum(obj.LeaderCommit,obj.PrevLogIndex + 1)
					
						action = append(action, Commit{Index: s.CommitIndex, Data: s.Log[s.CommitIndex].Data})
					}

			return action
		}

		if(s.Term>obj.Term){


			

			action=append(action,Send{To: obj.LeaderID,Event:AppendEntriesResponse{Term: s.Term,Success: false,From:s.MyID}})

		}else{
			
			



			if(len(s.Log)-1!=obj.PrevLogIndex){

				

				//different last Log Index then do not append

				action=append(action,Send{To: obj.LeaderID,Event:AppendEntriesResponse{Term: s.Term,Success: false,From:s.MyID}})

			} else {
				
				if( obj.PrevLogIndex != -1 && len(s.Log) != 0 && s.Log[len(s.Log)-1].Term != obj.PrevLogTerm){

					
					// last Log Index matches but different Terms at last Index
					action=append(action,Send{To: obj.LeaderID,Event:AppendEntriesResponse{Term: s.Term,Success: false,From:s.MyID}})

				}else{

					//delete entries From LastLogIndex and then append

					Index := obj.PrevLogIndex + 1
					

					if len(obj.Log) > 0 {
				
						s.Log=s.Log[:Index]
						s.Log = append(s.Log, obj.Log...)
					}

					s.LeaderID = obj.LeaderID
					//update Term and voted for if Term in append entry rpc is greater than follower current Term
					if obj.Term > s.Term {
						s.Term = obj.Term
						s.VotedFor = -1
						
						action = append(action, StateStore{Term: s.Term, VotedFor: s.VotedFor})
					}
					
					if obj.LeaderCommit > s.CommitIndex {
					
						s.CommitIndex = minimum(obj.LeaderCommit, Index-1)

						action = append(action, Commit{Index: s.CommitIndex, Data: s.Log[s.CommitIndex].Data})
					}

					action = append(action, Send{To: obj.LeaderID, Event: AppendEntriesResponse{Term: s.Term, Success: true,From:s.MyID,LastLogIndex:len(s.Log)-1}})


				}

			}
				 
		}

	} else if (s.State=="candidate") {

		
		if (obj.Term<s.Term) {

			action=append(action,Send{To: obj.LeaderID,Event:AppendEntriesResponse{Term: s.Term,Success: false,From:s.MyID}})
			
		}else if(obj.Term == s.Term) {

			if(len(s.Log)-1!=obj.PrevLogIndex){

				//different last Log Index then do not append

				action=append(action,Send{To: obj.LeaderID,Event:AppendEntriesResponse{Term: s.Term,Success: false,From:s.MyID}})

			} else {
				
				if( obj.PrevLogIndex != -1 && len(s.Log) != 0 && s.Log[len(s.Log)-1].Term != obj.PrevLogTerm){
					
					// last Log Index matches but different Terms at last Index
					action=append(action,Send{To: obj.LeaderID,Event:AppendEntriesResponse{Term: s.Term,Success: false,From:s.MyID}})

				}else{

					//delete entries From LastLogIndex and then append

					Index := obj.PrevLogIndex + 1
					

					if len(obj.Log) > 0 {
				
						action = append(action, LogStore{Index: Index, Data: obj.Log[obj.PrevLogIndex+1].Data})
						s.Log = append(s.Log, obj.Log...)
					}

					s.LeaderID = obj.LeaderID
					//update Term and voted for if Term in append entry rpc is greater than follower current Term
					if obj.Term > s.Term {
						s.Term = obj.Term
						s.VotedFor = -1
						
						action = append(action, StateStore{Term: s.Term, VotedFor: s.VotedFor})
					}
					
					if obj.LeaderCommit > s.CommitIndex {
					
						s.CommitIndex = minimum(obj.LeaderCommit, Index-1)

						
						action = append(action, Commit{Index: s.CommitIndex, Data: s.Log[s.CommitIndex].Data})
					}


					action = append(action, Send{To: obj.LeaderID, Event: AppendEntriesResponse{Term: s.Term, Success: true,From:s.MyID}})
					s.State="follower"


				}

			}


		}else{
	
		//when obj.Term > s.Term
				s.Term=obj.Term
				s.VotedFor=-1

				for i := 0; i < len(s.VoteReceived); i++ {

					s.VoteReceived[i]=-1

				}

				s.State="follower"
				action = append(action, StateStore{Term: s.Term, VotedFor: s.VotedFor})

			
			if(len(s.Log)-1!=obj.PrevLogIndex){

				//different last Log Index then do not append

				action=append(action,Send{To: obj.LeaderID,Event:AppendEntriesResponse{Term: s.Term,Success: false,From:s.MyID}})

			} else {
				
				if( obj.PrevLogIndex != -1 && len(s.Log) != 0 && s.Log[len(s.Log)-1].Term != obj.PrevLogTerm){
					
					// last Log Index matches but different Terms at last Index
					action=append(action,Send{To: obj.LeaderID,Event:AppendEntriesResponse{Term: s.Term,Success: false,From:s.MyID}})

				}else{

					
					//delete entries From LastLogIndex and then append

					Index := obj.PrevLogIndex + 1
					

					if len(obj.Log) > 0 {
				
						action = append(action, LogStore{Index: Index, Data: obj.Log[obj.PrevLogIndex+1].Data})
						s.Log = append(s.Log, obj.Log...)
					}

					s.LeaderID = obj.LeaderID
					
					if obj.LeaderCommit > s.CommitIndex {
					
						s.CommitIndex = minimum(obj.LeaderCommit, Index-1)

						
						action = append(action, Commit{Index: s.CommitIndex, Data: s.Log[s.CommitIndex].Data})
					}


					action = append(action, Send{To: obj.LeaderID, Event: AppendEntriesResponse{Term: s.Term, Success: true,From:s.MyID}})
					


				}

			}

		}
	}else{
		//s.State=leader

		//getting request for higher order Term
		if(s.Term<obj.Term){
			s.State="follower"
			s.Term=obj.Term
			s.VotedFor=-1

			for i := 0; i < len(s.VoteReceived); i++ {

				s.VoteReceived[i]=-1
			}

			action = append(action, StateStore{Term: s.Term, VotedFor: s.VotedFor})

		}

		//leader cannot get request From same order or higher order Term

	}


	
	return action
}

func onAppendEntriesResponse(obj AppendEntriesResponse, s *RaftStateMachine)[]Action {
	
	action:=make([]Action,0)
	if(s.State=="follower"){
		if (s.Term < obj.Term) {
			s.Term=obj.Term
			s.VotedFor=-1

			for i := 0; i < len(s.VoteReceived); i++ {

				s.VoteReceived[i]=-1

			}

			action = append(action, StateStore{Term: s.Term, VotedFor: s.VotedFor})
		}
	}else if (s.State=="candidate") {

		if (s.Term < obj.Term) {
			s.Term=obj.Term
			s.VotedFor=-1

			for i := 0; i < len(s.VoteReceived); i++ {

				s.VoteReceived[i]=-1

			}

			action = append(action, StateStore{Term: s.Term, VotedFor: s.VotedFor})
		}
		
	}else{
		//s.State=leader

		if (s.Term < obj.Term) {
			s.Term=obj.Term
			s.VotedFor=-1


			for i := 0; i < len(s.VoteReceived); i++ {

				s.VoteReceived[i]=-1
			}
		
			action = append(action, StateStore{Term: s.Term, VotedFor: s.VotedFor})
		}else{

					//checking if appendrequest was Successfull
			if(obj.Success==false){
				//decrease NextIndex 

				if(s.NextIndex[obj.From-1]>1){

					s.NextIndex[obj.From-1]=s.NextIndex[obj.From-1]-1
				}


				//updating last Log Index and Term
				newLastIndex:=obj.LastLogIndex-1
				newlastTerm:=0
				if(newLastIndex>=0){
					newlastTerm=s.Log[newLastIndex].Term
				}else{
					newlastTerm=0
				}

				//obtaining Data From newLastIndex To the length of the Log
				Log:=s.Log[newLastIndex+1:len(s.Log)]

				//resending AppendEntriesRequest
		
				action = append(action, Send{To: obj.From, Event: AppendEntriesRequest{ LeaderID: s.LeaderID, Term: s.Term,PrevLogIndex: newLastIndex, PrevLogTerm: newlastTerm, Log: Log, LeaderCommit: s.CommitIndex}})
			
			}else {
				// AppendEntryRequest Successfull
					
				//updating MatchIndex
				if((obj.LastLogIndex+obj.Count)+1>s.MatchIndex[obj.From]){
					s.MatchIndex[obj.From]=obj.LastLogIndex+obj.Count

				}

				// sorting MatchIndex and storing in new array To find whether the entries can be commited
				newMatchIndex := make([]int, len(s.PeerID)-1)
				copy(newMatchIndex, s.MatchIndex)
				sort.IntSlice(newMatchIndex).Sort()


				//commit is found by looking the MatchIndexes and finding which is present on the majority
			

				newCommit:=newMatchIndex[(len(s.PeerID)-1)/2]     // for 5 servers looking the last second entry

				//the entry is commited only if it is present on majority and for the same Term
				if(newCommit>s.CommitIndex && s.Log[newCommit].Term==s.Term){
					s.CommitIndex=newCommit

					action = append(action,Commit{Index:s.CommitIndex,Data:s.Log[obj.LastLogIndex].Data})

				}
				
		
			}
		}
		
	}

return action

}


func onTimeout(obj Timeout,s *RaftStateMachine) []Action{
	
	action:=make([]Action,0)

	if (s.State=="follower" || s.State=="candidate") {
		s.State="candidate"

		//increment the Term
		s.Term=s.Term+1

		//vote for self
		s.VotedFor=s.MyID


		for i := 0; i < len(s.VoteReceived); i++ {

			s.VoteReceived[i]=-1

		}

		action=append(action,StateStore{Term:s.Term,VotedFor:s.MyID})

		//reset election timer
		
		action=append(action,Alarm{Time: time.Now().Add(time.Duration(100) * time.Millisecond) })

        //send VoteRequests To everyone
		for i := 0; i < len(s.PeerID); i++ {
		
			action=append(action, Send{To:s.PeerID[i], Event:VoteRequest{Term:s.Term , CandidateID:s.MyID, LastLogIndex:len(s.Log)-1, LastLogTerm:s.Log[len(s.Log)-1].Term}})
			
		}

	
	} else{
		//s.State=leader

	
		myLog := make([]Log,0)

		action=append(action,Alarm{Time: time.Now().Add(time.Duration(100) * time.Millisecond) })

		//send empty append entries request
		for i := 0; i < len(s.PeerID); i++ {

			action=append(action, Send{To:s.PeerID[i], Event:AppendEntriesRequest{Log:myLog,Term:s.Term,LeaderID:s.LeaderID,LeaderCommit:s.CommitIndex}})
			
		}
		
	}


return action

}

func onAppend(obj Append,s *RaftStateMachine)[]Action{

	
action:=make([]Action,0)

	if(s.State=="follower" || s.State=="candidate"){


		action=append(action,Commit{Data:obj.Data,Err:errors.New("leader:"+strconv.Itoa(s.LeaderID)) })  //as follower or candidate cannot append To Log send LeaderID To the client
		

	}else{
		
			length:=len(s.Log)

			newLog:= make([]Log,0)


			s.Log=append(s.Log,Log{s.Term,obj.Data})
			newLog=append(newLog,Log{s.Term,obj.Data})
		
			action = append(action, LogStore{Index: length, Data: s.Log[length].Data})

		

			//send empty append entries request
			for i := 0; i < len(s.PeerID); i++ {
		
				action=append(action, Send{To:s.PeerID[i], Event:AppendEntriesRequest{Log:newLog,Term:s.Term,PrevLogIndex:len(s.Log)-2,PrevLogTerm:s.Log[len(s.Log)-2].Term,LeaderID:s.LeaderID}})
			
		}

	}


return action
}


func (s *RaftStateMachine) ProcessEvent(getEvent interface{}) []Action{

	action:=make([]Action,0)

	

		switch getEvent.(type){
		case VoteRequest:
			objVoteRequest:=getEvent.(VoteRequest)		//creating object of type VoteRequest From the Event
			action=onVoteRequest(objVoteRequest,s)
		case AppendEntriesRequest:
			objAppendEntriesRequest:=getEvent.(AppendEntriesRequest)		//creating object of type AppendEntriesRequest From the Event
			action=onAppendEntriesRequest(objAppendEntriesRequest,s)
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

