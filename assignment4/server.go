package main

import (
	"bufio"
	"fmt"
	"github.com/cs733-iitb/cs733/assignment1/fs"
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
	RequestMap map[int]chan *filesystem.Msg
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

func New(index int,config raftnode.Config) (ch *HandleClient){
	Register()

	ch = &HandleClient{
		//raftnode : raftnode.New(config.Cluster[index].Id,config),
		raftnode : raftnode.New(config),
		Port: config.Ports[index],
		RequestMap: make( map[int]chan *filesystem.Msg),
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

	go func () {
		for {
			connection,_ :=accept.AcceptTCP()
			go ch.serve(connection)
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

func (ch *HandleClient) reply(conn *net.TCPConn, msg *filessystem.Msg) bool {
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
		msg, msgerr, fatalerr := fs.GetMsg(reader)
		if fatalerr != nil || msgerr != nil {
			ch.reply(conn, &filessystem.Msg{Kind: 'M'})
			conn.Close()
			break
		}

		if msgerr != nil {
			if (!ch.reply(conn, &filessystem.Msg{Kind: 'M'})) {
				conn.Close()
				break
			}
		}

		if msg.Kind == 'r' {
            // Do not replicate the read request, directly serve
            response := filesystem.ProcessMsg(msg)
            if !ch.reply(conn, response) {
                conn.Close()
                return
            }
            continue    // Continue to serve another requests
        }

        requestId, waitchannel:= ch.RegRequest()

        request := Request{ServerId:ch.raftnode.ID(), Id:requestId, Req:*msg}
        ch.raftnode.Append(request)

        // Replication
        select {
        case response := <-waitchannel:
            ch.DeregRequest(request)                //Removing the request from the map

            if !ch.reply(conn, &response) {    // Response reply to the client
                return
            }
        case  <- time.After(100*time.Second) : //connection timeout
            ch.DeregRequest(request)

            ch.reply(conn, &filesystem.Msg{Kind:'I'})  // Reply with internal error
            conn.Close()
            return
        }

	}
}


func (ch *HandleClient) RegRequest() (reqId int, waitChan chan filessystem.Msg) {
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
