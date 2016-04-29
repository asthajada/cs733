package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"encoding/gob"
	"time"

	raftnode "github.com/asthajada/cs733/assignment4/raftnode"
	filesystem "github.com/asthajada/cs733/assignment4/fs"

)

var crlf = []byte{'\r', '\n'}

type HandleClient struct {
	raftnode *raftnode.RaftNode
	Port      int
	RequestMap map[int]chan filesystem.Msg
	FileSystem *filesystem.FS
	Counter int //Request Counter
	Address string //HostAddress
	ExitWaitGr *sync.WaitGroup
	MapLock *sync.RWMutex


}

type Request struct{
	Id int
	ServerId int
	Req filesystem.Msg
}

func NewClient(index int,config raftnode.Config) (ch *HandleClient){
	Register()

	ch = &HandleClient{
		//raftnode : raftnode.New(config.Cluster[index].Id,config),
		raftnode : raftnode.New(config),
		Port: config.Ports[index],
		RequestMap: make( map[int]chan filesystem.Msg),
		FileSystem: &filesystem.FS{Dir: make(map[string]*filesystem.FileInfo, 1000), Gversion:0},
		Counter:0,
		Address: config.Cluster[index].Host,
		MapLock: &sync.RWMutex{},
		ExitWaitGr : &sync.WaitGroup{},


	}

	return ch


}

func (ch *HandleClient) run(){

	ch.ExitWaitGr.Add(1)

	ch.raftnode.StartRaftNodes()

//	go ch.ListenActions()

	address,_ := net.ResolveTCPAddr("tcp",ch.Address + ":" +strconv.Itoa(ch.Port))

	accept,_:=net.ListenTCP("tcp",address)

	go ch.handleCommit() 	
	/*go func () {
        CommitLoop:
        for {

            select {
            case commitAction, ok := <-ch.raftnode.CommitChannel():
                if ok {
                	fmt.Println("got commit action")
                    ch.handleCommit(commitAction)
                } else {
                    break CommitLoop
                }
            }
        }
        ch.ExitWaitGr.Done()
    }()

*/
    go func () {
      
        for {
                tcp_conn, _ := accept.AcceptTCP()
                
                ch.ExitWaitGr.Add(1)     
                go ch.serve(tcp_conn)    // Start serve thread
            }
        
        
        ch.ExitWaitGr.Done()
        
    }()


}

func (ch *HandleClient) start() {
	ch.run()
	ch.ExitWaitGr.Wait()
}



func check(obj interface{}) {
	if obj != nil {
		fmt.Println(obj)
		os.Exit(1)
	}
}

func  reply(conn *net.TCPConn, msg *filesystem.Msg) bool {
	var err error
	write := func(data []byte) {
		if err != nil {
			return
		}
		_, err = conn.Write(data)
	}
	var resp string
	switch msg.Kind {
	case 'C': // read response
		resp = fmt.Sprintf("CONTENTS %d %d %d", msg.Version, msg.Numbytes, msg.Exptime)
	case 'O':
		resp = "OK "
		if msg.Version > 0 {
			resp += strconv.Itoa(msg.Version)
		}
	case 'F':
		resp = "ERR_FILE_NOT_FOUND"
	case 'V':
		resp = "ERR_VERSION " + strconv.Itoa(msg.Version)
	case 'M':
		resp = "ERR_CMD_ERR"
	case 'I':
		resp = "ERR_INTERNAL"
	default:
		fmt.Printf("Unknown response kind '%c'", msg.Kind)
		return false
	}
	resp += "\r\n"
	write([]byte(resp))
	if msg.Kind == 'C' {
		write(msg.Contents)
		write(crlf)
	}
	return err == nil
}

func (ch *HandleClient) serve(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		msg, msgerr, fatalerr := filesystem.GetMsg(reader)
		if fatalerr != nil || msgerr != nil {
			reply(conn, &filesystem.Msg{Kind: 'M'})
			conn.Close()
			break
		}

		if msgerr != nil {
			if (!reply(conn, &filesystem.Msg{Kind: 'M'})) {
				conn.Close()
				break
			}
		}

		if msg.Kind == 'r' {
            // Do not replicate the read request, directly serve
            response := ch.FileSystem.ProcessMsg(msg)
            if !reply(conn, response) {
                conn.Close()
                return
            }
            continue    // Continue to serve another requests
        } else {

	        requestId, waitchannel:= ch.RegRequest()

	        request := Request{ServerId:ch.raftnode.ID(), Id:requestId, Req:*msg}
	        ch.raftnode.Append(request)

	        // Replication
	        select {
	        case response := <-waitchannel:
	            ch.DeregRequest(requestId)                //Removing the request from the map

	            if !reply(conn, &response) {    // Response reply to the client
	                return
	            }
	        case  <- time.After(100*time.Second) : //connection timeout
	            ch.DeregRequest(requestId)

	            reply(conn, &filesystem.Msg{Kind:'I'})  // Reply with internal error
	            conn.Close()
	            return
	        }
    }

	}
}

func (ch *HandleClient) handleCommit () {
  /*  var response *filesystem.Msg

    fmt.Printf("\n Commitaction obj : %+v", commitAction)
    request := commitAction.Data.(Request)
    fmt.Printf("\n Req obj : %+v", request)


 
    //if commitAction.Err == nil {                        
        response = ch.FileSystem.ProcessMsg(&request.Req)      // process the request
     /*else {
        switch commitAction.Err.(type) {
        case raftnode.Error_NotLeader:                       // Not a leader, redirect error
            errorNotLeader := commitAction.Err.(rsm.Error_NotLeader)
            response = &filesystem.Msg{
                Kind            : 'R',
                RedirectAddr    : ch.Raft.ServerList[ errorNotLeader.LeaderId ] }
       
        }
    }

   // chd.Raft.UpdateLastApplied(commitAction.Index)      // Update last applied

    fmt.Println("Response  obj:", response)
    // Reply only if the client has requested this server
    if request.Id == ch.raftnode.ID() {
        ch.SendToWaitCh(request.Id, response)      // Send response to corresponding serve thread
    }
    */

    for {
		select {
		case commit, ok :=  <-ch.raftnode.CommitChannel() :
			if ok {

				var response *filesystem.Msg

				reqObj := commit.Data.(Request)

				//if commit.Err == nil {
					response = ch.FileSystem.ProcessMsg(&reqObj.Req)
				/*} else {
					id, _ := strconv.Atoi(commit.Err.Error())
					//If no leader elected yet, try again with same raftnode
					if id == -1 {
						id = ch.rn.Id()
					}
					response = &filesystem.Msg{Kind:'R', RedirectURL:GetRedirectURLFromId(id)}
				}*/

				// Reply only if the client has requested this server
				if reqObj.ServerId == ch.raftnode.ID() {
					//waitCh := ch.(reqObj.Id)
					//waitCh <- response
					ch.SendToWaitCh(reqObj.Id, response)
				}
			} else {
				ch.ExitWaitGr.Done()
				return
			}
		default:

		}
	}
}


func (ch *HandleClient) SendToWaitCh (reqId int, msg *filesystem.Msg) {
    ch.MapLock.RLock()
    conn, ok := ch.RequestMap[reqId]    // Extract wait channel from map
    if ok {                             // If request was not de-registered due to timeout
        conn <- *msg
    }
    ch.MapLock.RUnlock()
}


func (ch *HandleClient) RegRequest() (reqId int, waitChan chan filesystem.Msg) {
    waitChan = make(chan filesystem.Msg)
    ch.MapLock.Lock()
    ch.Counter++
    reqId = ch.Counter
    ch.RequestMap [ reqId ] = waitChan
    ch.MapLock.Unlock()
    return reqId, waitChan
}


func (ch *HandleClient) DeregRequest(reqId int) {
    ch.MapLock.Lock()
    close(ch.RequestMap[reqId])
    delete(ch.RequestMap, reqId)
    ch.MapLock.Unlock()
}


/*func serverMain() {
	tcpaddr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	check(err)
	tcp_acceptor, err := net.ListenTCP("tcp", tcpaddr)
	check(err)

	for {
		tcp_conn, err := tcp_acceptor.AcceptTCP()
		check(err)
		go serve(tcp_conn)
	}
}

func main() {
	serverMain()
}*/

func Register() {
	gob.Register(Request{})
	gob.Register(filesystem.Msg{})
}
