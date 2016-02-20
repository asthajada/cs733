package main
	
import (
	"fmt"
	"testing"
	"reflect"
)

	
func TestFollowerVoteRequest (t *testing.T){

	//follower's term is higher than that of candidate

	myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})

	s:=&node{term:3,votedFor:-1,log: myLog ,state: "follower", myID:2}

	actionarray:=s.ProcessEvent(VoteRequest{candidateID:1 ,term:2,lastLogIndex:2,lastLogTerm:2})
	actionname,event,_,resp,term,voteRecvFrom,_:=obtainDetails(actionarray,0)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,3)
 	expectInt(t,voteRecvFrom,2)
 	
 	//term of follower is same as that of candidate but lastLogTerm of follower is higher than that of candidate

 	myLog = make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})
	myLog=append(myLog,Log{3,[]byte("3")})
	myLog=append(myLog,Log{4,[]byte("4")})


	s=&node{term:6,votedFor:-1,log: myLog ,state: "follower",myID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:3,lastLogTerm:3})

 	
 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,0)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)


 	//term of follower is same as that of candidate and lastLogTerm of follower is also same as that of candidate but log index of follower is greater

	s=&node{term:6,votedFor:-1,log: myLog ,state: "follower",myID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:2,lastLogTerm:4})

 	
 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,0)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)

 	//term of follower is same as that of candidate and lastLogTerm of follower is also same as that of candidate and log index of follower is same as that of candidate

 
	s=&node{term:6,votedFor:-1,log: myLog ,state: "follower",myID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:4,lastLogTerm:4})

 	actionname,_,_,_,term,_,votedFor:=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,6)
	expectInt(t,votedFor,1)


 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,1)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,1)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)

 	//term of follower is same as that of candidate and lastLogTerm of follower is also same as that of candidate and log index of follower is same as that of candidate but already votedFor another candidate

	s=&node{term:6,votedFor:3,log: myLog ,state: "follower",myID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:4,lastLogTerm:4})

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,0)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	
 	//term of follower is lesser than that of candidate but lastLogterm of follower is greater than that of candidate

	s=&node{term:5,votedFor:3,log: myLog ,state: "follower",myID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:4,lastLogTerm:3})

 	actionname,_,_,_,term,_,votedFor=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,6)
	expectInt(t,votedFor,-1)

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,1)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	
 	//term of follower is lesser than that of candidate but lastLogterm of follower is same as that of candidate but lastlogindex of follower is greater than that of candidate

 	s=&node{term:5,votedFor:3,log: myLog ,state: "follower",myID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:2,lastLogTerm:4})

 	actionname,_,_,_,term,_,votedFor=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,6)
	expectInt(t,votedFor,-1)

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,1)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)

//term of follower is lesser as that of candidate and lastLogTerm of follower is also same as that of candidate and log index of follower is same as that of candidate and follower has not voted for any candidate

	s=&node{term:5,votedFor:-1,log: myLog ,state: "follower",myID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:4,lastLogTerm:4})

 	actionname,_,_,_,term,_,votedFor=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,6)
	expectInt(t,votedFor,1)


 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,1)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,1)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)

 	//term of follower is lesser as that of candidate and lastLogTerm of follower is also same as that of candidate and log index of follower is same as that of candidate but already votedFor another candidate


	s=&node{term:5,votedFor:3,log: myLog ,state: "follower",myID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:4,lastLogTerm:4})

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,0)

 	actionname,_,_,_,term,_,votedFor=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,6)
	expectInt(t,votedFor,-1)

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,1)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 

}

func TestCandidateVoteRequest(t *testing.T) {
	//the candidate has lower term than candidate requesting vote and last logterm of this candidate is greater than lastlogterm of candidate requesting vote
	myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})
	myLog=append(myLog,Log{3,[]byte("3")})
	myLog=append(myLog,Log{4,[]byte("4")})

	s:=&node{term:5,votedFor:3,log: myLog ,state: "candidate",myID:2}

 	actionarray:=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:4,lastLogTerm:3})

 	actionname,event,_,resp,term,voteRecvFrom,_:=obtainDetails(actionarray,0)

 	actionname,_,_,_,term,_,votedFor:=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,6)
	expectInt(t,votedFor,-1)

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,1)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	expectString(t,s.state,"follower")

 	//the candidate has lower term than candidate requesting vote and last logterm of this candidate is same as that of candidate requesting vote but last log index of this candidate is greater than that of candidate requesting vote

	s=&node{term:5,votedFor:3,log: myLog ,state: "candidate",myID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:2,lastLogTerm:4})

 	actionname,_,_,_,term,_,votedFor=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,6)
	expectInt(t,votedFor,-1)

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,1)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	expectString(t,s.state,"follower")


 	//term of candidate is lesser as that of candidate requesting vote and lastLogTerm of candidate is also same as that of candidate and log index of candidate is same as that of candidate and candidate has not voted for any candidate

	s=&node{term:5,votedFor:-1,log: myLog ,state: "candidate",myID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:3,lastLogTerm:4})

 	actionname,_,_,_,term,_,votedFor=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,6)
	expectInt(t,votedFor,1)


 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,1)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,1)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	expectString(t,s.state,"follower")

 	//term of candidate is lesser as that of candidate requesting vote and lastLogTerm of candidate is also same as that of candidate and log index of candidate is same as that of candidate and candidate has already voted for any candidate


	s=&node{term:5,votedFor:3,log: myLog ,state: "candidate",myID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:4,lastLogTerm:4})

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,0)

 	actionname,_,_,_,term,_,votedFor=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,6)
	expectInt(t,votedFor,-1)

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,1)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	expectString(t,s.state,"follower")

 	//term of candidate is same as that of candidate requesting vote
	s=&node{term:6,votedFor:3,log: myLog ,state: "candidate",myID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:2,lastLogTerm:4})

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,0)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	expectString(t,s.state,"candidate")

	//term of candidate is greater as that of candidate requesting vote
	s=&node{term:6,votedFor:3,log: myLog ,state: "candidate",myID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:5,lastLogIndex:2,lastLogTerm:4})

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,0)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	expectString(t,s.state,"candidate")
 	 	

}


func TestLeaderVoteRequest(t *testing.T){
	//the leader has lower term than candidate requesting vote and last logterm of leader is greater than lastlogterm of candidate requesting vote
	myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})
	myLog=append(myLog,Log{3,[]byte("3")})
	myLog=append(myLog,Log{4,[]byte("4")})

	s:=&node{term:5,votedFor:3,log: myLog ,state: "leader",myID:2,leaderID:2}

 	actionarray:=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:4,lastLogTerm:3})

 	actionname,event,_,resp,term,voteRecvFrom,_:=obtainDetails(actionarray,0)

 	actionname,_,_,_,term,_,votedFor:=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,6)
	expectInt(t,votedFor,-1)

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,1)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	expectString(t,s.state,"follower")

 	//the leader has lower term than candidate requesting vote and last logterm of leader is same as that of candidate requesting vote but last log index of leader is greater than that of candidate requesting vote

	s=&node{term:5,votedFor:3,log: myLog ,state: "leader",myID:2,leaderID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:2,lastLogTerm:4})

 	actionname,_,_,_,term,_,votedFor=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,6)
	expectInt(t,votedFor,-1)

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,1)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	expectString(t,s.state,"follower")


 	//term of leader is lesser as that of candidate requesting vote and lastLogTerm of leader is also same as that of candidate and log index of leader is same as that of candidate and leader has not voted for any candidate

	s=&node{term:5,votedFor:-1,log: myLog ,state: "leader",myID:2,leaderID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:3,lastLogTerm:4})

 	actionname,_,_,_,term,_,votedFor=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,6)
	expectInt(t,votedFor,1)


 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,1)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,1)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	expectString(t,s.state,"follower")

 	//term of leader is lesser as that of candidate requesting vote and lastLogTerm of leader is also same as that of candidate and log index of leader is same as that of candidate and leader has already voted for any candidate


	s=&node{term:5,votedFor:3,log: myLog ,state: "leader",myID:2,leaderID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:4,lastLogTerm:4})

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,0)

 	actionname,_,_,_,term,_,votedFor=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,6)
	expectInt(t,votedFor,-1)

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,1)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	expectString(t,s.state,"follower")

 	//term of leader is same as that of candidate requesting vote
	s=&node{term:6,votedFor:3,log: myLog ,state: "leader",myID:2,leaderID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:6,lastLogIndex:2,lastLogTerm:4})

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,0)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	expectString(t,s.state,"leader")

	//term of leader is greater as that of candidate requesting vote
	s=&node{term:6,votedFor:3,log: myLog ,state: "leader",myID:2,leaderID:2}

 	actionarray=s.ProcessEvent(VoteRequest{candidateID:1 ,term:5,lastLogIndex:2,lastLogTerm:4})

 	actionname,event,_,resp,term,voteRecvFrom,_=obtainDetails(actionarray,0)

 	expectString(t,actionname,"Send")
 	expectString(t,event,"VoteResponse")
 	expectInt(t,resp,0)
 	expectInt(t,term,6)
 	expectInt(t,voteRecvFrom,2)
 	expectString(t,s.state,"leader")
 	 	


}

func TestFollowerVoteResponse(t *testing.T) {

	myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})

	//getting voteResponse from higher order term
	s:=&node{term:3,votedFor:-1,log: myLog ,state: "follower",leaderID: 1, myID:2}

	actionarray:=s.ProcessEvent(VoteResponse{from:2 ,term:4,voteGranted:0})

	actionname,_,_,_,term,_,_:=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")

	expectInt(t,term,4)

	//getting voteResponse form lower order term

	s=&node{term:3,votedFor:-1,log: myLog ,state: "follower",leaderID: 1, myID:2}

	_=s.ProcessEvent(VoteResponse{from:2 ,term:1,voteGranted:0})

	expectInt(t,s.term,3)
	


	
}

func TestLeaderVoteResponse(t *testing.T) {

	myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})

	s:=&node{term:3,votedFor:-1,log: myLog ,state: "leader",leaderID: 1, myID:1}

	actionarray:=s.ProcessEvent(VoteResponse{from:2 ,term:2,voteGranted:0})
	expectInt(t,len(actionarray),0)

	
}

func TestCandidateVoteResponse(t *testing.T) {
	var i int
	myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})


	mypeerID:=make([]int,0)
	mypeerID=append(mypeerID,1,2,3,4,5)

	myvoteReceived:=make([]int,5)
	myvoteReceived[0]=1

	for i := 1; i < len(myvoteReceived); i++ {
		myvoteReceived[i]=-1
	}

	

	s:=&node{term:2,votedFor:-1,log: myLog ,state: "candidate",leaderID: -1, myID:1,peerID:mypeerID,voteReceived:myvoteReceived}

	//checking if on getting the majority client converts into leader
	_=s.ProcessEvent(VoteResponse{from:2 ,term:2,voteGranted:1})
	expectString(t,s.state,"candidate")
	actionarray:=s.ProcessEvent(VoteResponse{from:3 ,term:2,voteGranted:1})
	expectString(t,s.state,"leader")
	expectInt(t,s.leaderID,1)

	for i = 0; i < len(s.peerID); i++ {
				
		actionname,event,_,_,term,_,_:=obtainDetails(actionarray,i)

		expectString(t,actionname,"Send")
		expectString(t,event,"AppendEntriesRequest")
		expectInt(t,term,2)

	}

	actionname,_,_,_,_,_,_:=obtainDetails(actionarray,i)
	expectString(t,actionname,"Alarm")




	// checking that multiple VoteResponses from same server does not add up to majority 

	myvoteReceived=make([]int,5)
	myvoteReceived[0]=1

	for i := 1; i < len(myvoteReceived); i++ {
		myvoteReceived[i]=-1
	}

	s=&node{term:2,votedFor:-1,log: myLog ,state: "candidate",myID:1,peerID:mypeerID,voteReceived:myvoteReceived}

	_=s.ProcessEvent(VoteResponse{from:2 ,term:2,voteGranted:1})
	expectString(t,s.state,"candidate")
	_=s.ProcessEvent(VoteResponse{from:2 ,term:2,voteGranted:1})
	expectString(t,s.state,"candidate")	

	myvoteReceived=make([]int,5)
	myvoteReceived[0]=1

	for i := 1; i < len(myvoteReceived); i++ {
		myvoteReceived[i]=-1
	}

	s=&node{term:2,votedFor:-1,log: myLog ,state: "candidate",leaderID: -1, myID:1,peerID:mypeerID,voteReceived:myvoteReceived}

	_=s.ProcessEvent(VoteResponse{from:2 ,term:2,voteGranted:1})
	expectString(t,s.state,"candidate")
	_=s.ProcessEvent(VoteResponse{from:3 ,term:2,voteGranted:0})
	expectString(t,s.state,"candidate")	
	_=s.ProcessEvent(VoteResponse{from:4 ,term:2,voteGranted:0})
	expectString(t,s.state,"candidate")
	_=s.ProcessEvent(VoteResponse{from:5 ,term:2,voteGranted:1})
	expectString(t,s.state,"leader")
	expectInt(t,s.leaderID,1)

	for i = 0; i < len(s.peerID); i++ {
		
		
		actionname,event,_,_,term,_,_:=obtainDetails(actionarray,i)

		expectString(t,actionname,"Send")
		expectString(t,event,"AppendEntriesRequest")
		expectInt(t,term,2)

	}

	actionname,_,_,_,_,_,_=obtainDetails(actionarray,i)
	expectString(t,actionname,"Alarm")


	//converting to follower if majority of vote are denied

	myvoteReceived=make([]int,5)
	myvoteReceived[0]=1

	for i := 1; i < len(myvoteReceived); i++ {
		myvoteReceived[i]=-1
	}

	s=&node{term:2,votedFor:-1,log: myLog ,state: "candidate", myID:1,peerID:mypeerID,voteReceived:myvoteReceived}
	_=s.ProcessEvent(VoteResponse{from:2 ,term:2,voteGranted:0})
	expectString(t,s.state,"candidate")
	_=s.ProcessEvent(VoteResponse{from:3 ,term:2,voteGranted:0})
	expectString(t,s.state,"candidate")	
	actionarray=s.ProcessEvent(VoteResponse{from:4 ,term:2,voteGranted:0})
	expectString(t,s.state,"follower")

	actionname,_,_,_,_,_,_=obtainDetails(actionarray,0)
	expectString(t,actionname,"Alarm")

	
}	

func TestFollowerAppendEntriesRequest(t *testing.T){

	myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})

	//term of follower is higher than that of leader

	s:=&node{term:4,votedFor:-1,log: myLog ,state: "follower",leaderID: 1, myID:2}

	newlog:=make([]Log,0)
	newlog=append(newlog,Log{1,[]byte("1")})
	newlog=append(newlog,Log{2,[]byte("2")})
	

	actionarray:=s.ProcessEvent(AppendEntriesRequest{term:3,leaderID:1,prevLogIndex:2,prevLogTerm:2,log:newlog})
	
	actionname,event,resp1,_,_,_,_:=obtainDetails(actionarray,0)
	expectString(t,actionname,"Send")
	expectString(t,event,"AppendEntriesResponse")
	expectBool(t,resp1,false)

	//same terms but different log lengths

	s=&node{term:4,votedFor:-1,log: myLog ,state: "follower",leaderID: 1, myID:2}
	newlog=append(newlog,Log{3,[]byte("3")})
	actionarray=s.ProcessEvent(AppendEntriesRequest{term:4,leaderID:1,prevLogIndex:2,prevLogTerm:3,log:newlog})

	actionname,event,resp1,_,_,_,_=obtainDetails(actionarray,0)
	expectString(t,actionname,"Send")
	expectString(t,event,"AppendEntriesResponse")
	expectBool(t,resp1,false)

	//same log lengths but different terms at that index

	s=&node{term:4,votedFor:-1,log: myLog ,state: "follower",leaderID: 1, myID:2}

	newlog=make([]Log,0)
	newlog=append(newlog,Log{1,[]byte("1")})
	newlog=append(newlog,Log{3,[]byte("2")})
	actionarray=s.ProcessEvent(AppendEntriesRequest{term:4,leaderID:1,prevLogIndex:1,prevLogTerm:3,log:newlog})

	actionname,event,resp1,_,_,_,_=obtainDetails(actionarray,0)
	expectString(t,actionname,"Send")
	expectString(t,event,"AppendEntriesResponse")
	expectBool(t,resp1,false)

	//same term, lastlogindex and lastlogterm

	s=&node{term:4,votedFor:-1,log: myLog ,state: "follower",leaderID: 1, myID:2,commitIndex:1}

	newlog=make([]Log,0)
	newlog=append(newlog,Log{3,[]byte("1")})
	newlog=append(newlog,Log{4,[]byte("2")})
	newlog=append(newlog,Log{4,[]byte("2")})

	actionarray=s.ProcessEvent(AppendEntriesRequest{term:4,leaderID:1,prevLogIndex:1,prevLogTerm:2,log:newlog,leaderCommit:2})
	

	actionname,event,resp1,_,_,_,_=obtainDetails(actionarray,0)
	expectString(t,actionname,"LogStore")
	
	actionname,_,_,_,index,_,_:=obtainDetails(actionarray,1)
	expectString(t,actionname,"Commit")
	expectInt(t,index,1)

	actionname,event,resp1,_,term,_,_:=obtainDetails(actionarray,2)
	expectString(t,actionname,"Send")
	expectString(t,event,"AppendEntriesResponse")
	expectBool(t,resp1,true)
	expectInt(t,term,4)

	//if term of leader is higher

	s=&node{term:4,votedFor:-1,log: myLog ,state: "follower",leaderID: 1, myID:2,commitIndex:1}

	newlog=make([]Log,0)
	newlog=append(newlog,Log{3,[]byte("1")})
	newlog=append(newlog,Log{4,[]byte("2")})
	newlog=append(newlog,Log{4,[]byte("2")})

	actionarray=s.ProcessEvent(AppendEntriesRequest{term:5,leaderID:1,prevLogIndex:1,prevLogTerm:2,log:newlog,leaderCommit:2})

	actionname,event,resp1,_,_,_,_=obtainDetails(actionarray,0)
	expectString(t,actionname,"LogStore")
	
	actionname,_,_,_,term,_,votedFor:=obtainDetails(actionarray,1)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,5)
	expectInt(t,votedFor,-1)

	actionname,_,_,_,index,_,_=obtainDetails(actionarray,2)
	expectString(t,actionname,"Commit")
	expectInt(t,index,1)

	actionname,event,resp1,_,term,_,_=obtainDetails(actionarray,3)
	expectString(t,actionname,"Send")
	expectString(t,event,"AppendEntriesResponse")
	expectBool(t,resp1,true)
	expectInt(t,term,5)


}

func TestCandidateAppendEntriesRequest(t *testing.T){

myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})

	//term of candidate is higher than that of leader

	s:=&node{term:4,votedFor:-1,log: myLog ,state: "candidate",leaderID: 1, myID:2}

	newlog:=make([]Log,0)
	newlog=append(newlog,Log{1,[]byte("1")})
	newlog=append(newlog,Log{2,[]byte("2")})
	

	actionarray:=s.ProcessEvent(AppendEntriesRequest{term:3,leaderID:1,prevLogIndex:2,prevLogTerm:2,log:newlog})
	
	actionname,event,resp1,_,_,_,_:=obtainDetails(actionarray,0)
	expectString(t,actionname,"Send")
	expectString(t,event,"AppendEntriesResponse")
	expectBool(t,resp1,false)

	//same terms but different log lengths

	s=&node{term:4,votedFor:-1,log: myLog ,state: "candidate",leaderID: 1, myID:2}
	newlog=append(newlog,Log{3,[]byte("3")})
	actionarray=s.ProcessEvent(AppendEntriesRequest{term:4,leaderID:1,prevLogIndex:2,prevLogTerm:3,log:newlog})

	actionname,event,resp1,_,_,_,_=obtainDetails(actionarray,0)
	expectString(t,actionname,"Send")
	expectString(t,event,"AppendEntriesResponse")
	expectBool(t,resp1,false)

	//same log lengths but different terms at that index

	s=&node{term:4,votedFor:-1,log: myLog ,state: "candidate",leaderID: 1, myID:2}

	newlog=make([]Log,0)
	newlog=append(newlog,Log{1,[]byte("1")})
	newlog=append(newlog,Log{3,[]byte("2")})
	actionarray=s.ProcessEvent(AppendEntriesRequest{term:4,leaderID:1,prevLogIndex:1,prevLogTerm:3,log:newlog})

	actionname,event,resp1,_,_,_,_=obtainDetails(actionarray,0)
	expectString(t,actionname,"Send")
	expectString(t,event,"AppendEntriesResponse")
	expectBool(t,resp1,false)

	//same term, lastlogindex and lastlogterm

	s=&node{term:4,votedFor:-1,log: myLog ,state: "candidate",leaderID: 1, myID:2,commitIndex:1}

	newlog=make([]Log,0)
	newlog=append(newlog,Log{3,[]byte("1")})
	newlog=append(newlog,Log{4,[]byte("2")})
	newlog=append(newlog,Log{4,[]byte("2")})

	actionarray=s.ProcessEvent(AppendEntriesRequest{term:4,leaderID:1,prevLogIndex:1,prevLogTerm:2,log:newlog,leaderCommit:2})
	

	actionname,event,resp1,_,_,_,_=obtainDetails(actionarray,0)
	expectString(t,actionname,"LogStore")
	
	actionname,_,_,_,index,_,_:=obtainDetails(actionarray,1)
	expectString(t,actionname,"Commit")
	expectInt(t,index,1)

	actionname,event,resp1,_,term,_,_:=obtainDetails(actionarray,2)
	expectString(t,actionname,"Send")
	expectString(t,event,"AppendEntriesResponse")
	expectBool(t,resp1,true)
	expectInt(t,term,4)
	expectString(t,s.state,"follower")

   //if term of leader is higher

	s=&node{term:4,votedFor:-1,log: myLog ,state: "candidate",leaderID: 1, myID:2,commitIndex:1}

	newlog=make([]Log,0)
	newlog=append(newlog,Log{3,[]byte("1")})
	newlog=append(newlog,Log{4,[]byte("2")})
	newlog=append(newlog,Log{4,[]byte("2")})

	actionarray=s.ProcessEvent(AppendEntriesRequest{term:5,leaderID:1,prevLogIndex:1,prevLogTerm:2,log:newlog,leaderCommit:2})


	actionname,_,_,_,term,_,votedFor:=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,5)
	expectInt(t,votedFor,-1)
	expectString(t,s.state,"follower")

	actionname,event,resp1,_,_,_,_=obtainDetails(actionarray,1)
	expectString(t,actionname,"LogStore")
	
	actionname,_,_,_,index,_,_=obtainDetails(actionarray,2)
	expectString(t,actionname,"Commit")
	expectInt(t,index,1)

	actionname,event,resp1,_,term,_,_=obtainDetails(actionarray,3)
	expectString(t,actionname,"Send")
	expectString(t,event,"AppendEntriesResponse")
	expectBool(t,resp1,true)
	expectInt(t,term,5)
	expectString(t,s.state,"follower")


}

func TestLeaderAppendEntriesRequest(t *testing.T){

	myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})

	//getting request from higher order term

	s:=&node{term:4,votedFor:-1,log: myLog ,state: "leader",leaderID: 1, myID:2}

	newlog:=make([]Log,0)
	newlog=append(newlog,Log{1,[]byte("1")})
	newlog=append(newlog,Log{2,[]byte("2")})
	

	actionarray:=s.ProcessEvent(AppendEntriesRequest{term:5,leaderID:1,prevLogIndex:2,prevLogTerm:2,log:newlog})
	
	actionname,_,_,_,term,_,votedFor:=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,5)
	expectInt(t,votedFor,-1)
	expectString(t,s.state,"follower")


}

func TestFollowerAppendEntriesResponse(t *testing.T) {
	
	myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})

	//getting response from higher order term

     s:=&node{term:4,votedFor:-1,log: myLog ,state: "follower",leaderID: 1, myID:2}

	actionarray:=s.ProcessEvent(AppendEntriesResponse{term:5,success:false,from:3})

	actionname,_,_,_,term,_,votedFor:=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,5)
	expectInt(t,votedFor,-1)

	
}

func TestCandidateAppendEntriesResponse(t *testing.T) {

	myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})

	//getting response from higher order term

     s:=&node{term:4,votedFor:-1,log: myLog ,state: "candidate",leaderID: 1, myID:2}

	actionarray:=s.ProcessEvent(AppendEntriesResponse{term:5,success:false,from:3})

	actionname,_,_,_,term,_,votedFor:=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,5)
	expectInt(t,votedFor,-1)

	
}

func TestFollowerTimeout(t *testing.T){

	mypeerID:=make([]int,0)
	mypeerID=append(mypeerID,1,2,3,4,5)

	myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})

	s:=&node{myID:2,term:2,votedFor:-1,state: "follower", peerID:mypeerID,log:myLog}

	actionarray:=s.ProcessEvent(Timeout{})
	
	actionname,_,_,_,term,_,votedFor:=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,3)
	expectInt(t,votedFor,2)


	actionname,_,_,_,_,_,_=obtainDetails(actionarray,1)
	expectString(t,actionname,"Alarm")
	expectString(t,s.state,"candidate")

	for i := 0; i < len(s.peerID); i++ {
		
		actionname,event,_,_,term,_,_:=obtainDetails(actionarray,i+2)

		expectString(t,actionname,"Send")
		expectString(t,event,"VoteRequest")
		expectInt(t,term,3)

	}
	

}

func TestCandidateTimeout(t *testing.T) {

	mypeerID:=make([]int,0)
	mypeerID=append(mypeerID,1,2,3,4,5)

	myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})

	s:=&node{myID:2,term:2,votedFor:-1,state: "candidate", peerID:mypeerID,log:myLog}

	actionarray:=s.ProcessEvent(Timeout{})
	
	actionname,_,_,_,term,_,votedFor:=obtainDetails(actionarray,0)
	expectString(t,actionname,"StateStore")
	expectInt(t,term,3)
	expectInt(t,votedFor,2)
	expectString(t,s.state,"candidate")


	actionname,_,_,_,_,_,_=obtainDetails(actionarray,1)
	expectString(t,actionname,"Alarm")

	for i := 0; i < len(s.peerID); i++ {
		
		actionname,event,_,_,term,_,_:=obtainDetails(actionarray,i+2)

		expectString(t,actionname,"Send")
		expectString(t,event,"VoteRequest")
		expectInt(t,term,3)

	}
	
}

func TestLeaderTimeout(t *testing.T){
	mypeerID:=make([]int,0)
	mypeerID=append(mypeerID,1,2,3,4,5)

	s:=&node{myID:2,term:2,votedFor:-1,state: "leader", peerID:mypeerID}
	actionarray:=s.ProcessEvent(Timeout{})

	actionname,_,_,_,_,_,_:=obtainDetails(actionarray,0)
	expectString(t,actionname,"Alarm")

	for i := 0; i < len(s.peerID); i++ {
		
		actionname,event,_,resp,term,_,_:=obtainDetails(actionarray,i+1)

		expectString(t,actionname,"Send")
		expectString(t,event,"AppendEntriesRequest")
		expectInt(t,term,2)
		expectInt(t,resp,0)

	}

}

func TestFollowerAppend(t *testing.T) {

	s:=&node{myID:2,term:2,votedFor:-1,state: "follower",leaderID:3}
	actionarray:=s.ProcessEvent(Append{[]byte("1234")})

	actionname,_,_,_,_,_,_:=obtainDetails(actionarray,0)
	expectString(t,actionname,"Commit")
	actiontype := actionarray[0]
	Object:=actiontype.(Commit)
	expectString(t,Object.err,"leader:3")
	
}

func TestCandidateAppend(t *testing.T) {

	s:=&node{myID:2,term:2,votedFor:-1,state: "candidate",leaderID:3}
	actionarray:=s.ProcessEvent(Append{[]byte("1234")})

	actionname,_,_,_,_,_,_:=obtainDetails(actionarray,0)
	expectString(t,actionname,"Commit")
	actiontype := actionarray[0]
	Object:=actiontype.(Commit)
	expectString(t,Object.err,"leader:3")
	
}

func TestLeaderAppend(t *testing.T) {
    myLog := make([]Log,0)
	myLog=append(myLog,Log{1,[]byte("1")})
	myLog=append(myLog,Log{2,[]byte("2")})

	s:=&node{myID:2,term:2,votedFor:-1,state: "leader",leaderID:2,log:myLog}

	newdata:=make([]byte,0)
	newdata=append(newdata,6)

	actionarray:=s.ProcessEvent(Append{newdata})
	

	actionname,_,_,_,_,_,_:=obtainDetails(actionarray,0)
	expectString(t,actionname,"LogStore")

	
}

func obtainDetails(action []Action,index int) (string,string,bool,int,int,int,int){
	var actionname,event string
	var term,voteRecvFrom,votedFor,resp int
	var resp1 bool
	
	
	actiontype := action[index]
		actionname=reflect.TypeOf(actiontype).Name()
		
		switch actionname{
			case "Send":
						Object:=actiontype.(Send)
						event=reflect.TypeOf(Object.event).Name()
						switch event{
							case "VoteResponse":
								respObject:=Object.event.(VoteResponse)
								voteRecvFrom=respObject.from
								term=respObject.term
								resp=respObject.voteGranted
							case "AppendEntriesResponse":
								respObject:=Object.event.(AppendEntriesResponse)
								term=respObject.term
								resp1=respObject.success
							case "VoteRequest":
								respObject:=Object.event.(VoteRequest)
								term=respObject.term
							case "AppendEntriesRequest":
								respObject:=Object.event.(AppendEntriesRequest)
								term=respObject.term
								resp=len(respObject.log)



						}
			case "LogStore":
				Object:=actiontype.(LogStore)
				term=Object.index
			case "Commit":
				Object:=actiontype.(Commit)
				term=Object.index

			case "Alarm":
			case "StateStore":
				Object:=actiontype.(StateStore)
				term=Object.term
				votedFor=Object.votedFor
				
				
		}
	return actionname,event,resp1,resp,term,voteRecvFrom,votedFor
}


func expectString(t *testing.T, a string, b string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}

func expectInt(t *testing.T, a int, b int) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}

func expectBool(t *testing.T, a bool, b bool) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}


