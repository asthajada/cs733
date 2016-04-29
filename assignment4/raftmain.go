package main

import (
	//"fmt"
	raftnode "github.com/asthajada/cs733/assignment4/raftnode"
	"os"
	"strconv"
	"os/exec"
	"log"
	"time"
)

func main() {

	/*expectArgs := func (n int) {
		if len(os.Args) < n {
			fmt.Println("\nCommandline Arg error ! \n Usage : ./assignment4 <id>\n")
			os.Exit(1)
		}
	}

	expectArgs(1)*/

	

	myNetConfig := make([]raftnode.NetConfig,0)
	myNetConfig=append(myNetConfig,raftnode.NetConfig{Id:1,Host:"localhost" ,Port:2000} )
	myNetConfig=append(myNetConfig,raftnode.NetConfig{Id:2,Host:"localhost" ,Port:3000} )
	myNetConfig=append(myNetConfig,raftnode.NetConfig{Id:3,Host:"localhost" ,Port:4000} )
	myNetConfig=append(myNetConfig,raftnode.NetConfig{Id:4,Host:"localhost" ,Port:5000} )
	myNetConfig=append(myNetConfig,raftnode.NetConfig{Id:5,Host:"localhost" ,Port:6000} )

	myclientports:= make([]int,0)
	myclientports=append(myclientports,9001)
	myclientports=append(myclientports,9002)
	myclientports=append(myclientports,9003)
	myclientports=append(myclientports,9004)	
	myclientports=append(myclientports,9005)


		myConfig:=raftnode.Config{Cluster:myNetConfig,Id:1,InboxSize:100,OutboxSize:100,ElectionTimeout:1000,HeartbeatTimeout:200,LogDir:"logdirectory",Ports:myclientports}
	//config.MCluster = mcluster


	time.Sleep(3* time.Second)


	if len(os.Args) == 1 {
		cmds := make([]*exec.Cmd, 5)
		for i := 0; i < 5; i++ {
			// os.Args[0] works reliably to re-use the currently running executable
			// whether you built and ran it or used `go run`
			cmds[i] = exec.Command(os.Args[0], strconv.Itoa(i))
			err := cmds[i].Start()
			if err != nil {
				log.Println(err)
			}
		}
	

	} else {
		index,_:=strconv.Atoi(os.Args[1])
		myConfig.Id = index + 1
		clientHandler := NewClient(index, myConfig)
		clientHandler.start()
	}
}