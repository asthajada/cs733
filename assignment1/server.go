//File server created by Astha Jada on 21 Janaury 2016

package main

import (
	"fmt"
	"net"
	"strings"
	"log"
	"strconv"
	"math/rand"
	"time"
	"sync"
)


type fileContent struct {			//to store the details of the file

  numbytes,exptime,version int64
  content []byte
  creationtime time.Time

  }
 
 //mutex to handle multiple clients concurrently 
 var mu = &sync.RWMutex{}
 
var(

 mymap=make(map[string]fileContent)
)
const PORT = ":8080"

var expirytime time.Time


//function to handle the read command
func handleRead(cmd []string, conn net.Conn){
		var time bool
		time=false
	
		mu.RLock()
		if val,ok := mymap[cmd[1]]; ok {	//checking if the file is present in the map
		mu.RUnlock()
		
		
		if val.exptime!=0 {
		time=expirytime.After(val.creationtime)
						
		}else{
		time=true
		}
		
		
			
			if time{				//checking if the file has expired
		    		_ ,err := conn.Write([]byte("CONTENTS "+strconv.FormatInt(val.version,10)+" "+strconv.FormatInt(val.numbytes,10)+" "+strconv.FormatInt(val.exptime,10)+"\r\n"+string(val.content)+"\r\n"))
					if err != nil {
					errMsg := "ERR_INTERNAL\r\n"
					conn.Write([]byte(errMsg))
					return
					}
			} else{
					mu.Lock()
					delete(mymap,cmd[1])
					mu.Unlock()
					errMsg := "ERR_FILE_NOT_FOUND\r\n"
					conn.Write([]byte(errMsg))
					return
				  }
				  
		}else{
		
			errMsg := "ERR_FILE_NOT_FOUND\r\n"
			conn.Write([]byte(errMsg))
			return
	
		}
		
	
		
}

// to handle the write command
func handleWrite(command []string,data []byte)(int64,error){

		key:=command[1]
	
		numBytes,err:=strconv.ParseInt(command[2],0,64)
	
		if err!=nil{
			return 0,err
		}

		
  		var expiry int64
  		
		if len(command)==4{				//checking if the expirytime is provided by the user
		
			expiry,err=strconv.ParseInt(command[3],0,64)
			if err!=nil{
				 return 0,err
			}
	
			
		}else{			
			expiry=0			//default value if expiry is not provided
			
			
		}
	
		
		version:=rand.Int63()			//generation of random number for version
		

		curr:=time.Now()
		expirytime=curr.Add(time.Duration(expiry)*time.Second)
			
		mu.Lock()
		mymap[key] = fileContent{numBytes,expiry,version,data,curr}
		mu.Unlock()
	

		return version,err

}


func serverMain() {
 	
 	fmt.Println("Launching server...")   

	ln,err := net.Listen("tcp", PORT)   

	if err != nil {
		log.Printf("Error in listener: %v", err.Error())		
	}


	for {     
		// accept connection on port   
		conn, _ := ln.Accept()  

		//handle multiple clients
		go processClient(conn)
		
	
	}
}


func processClient(conn net.Conn) {

for {
	dataBuffer:= make([]byte, 1024)
	
    
	size,err:=conn.Read(dataBuffer[0:])
	if err != nil {
		log.Printf("Error in reading from client")
	}
	dataBuffer = dataBuffer[:size]
	
	data1 := strings.Split(string(dataBuffer),"\r\n")[0]

	
	command:=strings.Split(data1," ")
	
	
	if len(command) == 0 {
		errMsg := "ERR_CMD_ERR\r\n"
		conn.Write([]byte(errMsg))
	}

	switch command[0] {

		case "read":
			
			if len(command) == 2 && len(command[1])<=250{
			handleRead(command, conn)
			
			} else {
				errMsg := "ERR_CMD_ERR\r\n"
				conn.Write([]byte(errMsg))
				return
			}	
			
		case "write":
		
			if len(command)>4 || len(command)<3 || len(command[1])>250{
				errMsg := "ERR_CMD_ERR\r\n"
				conn.Write([]byte(errMsg))
			return;
				
			} else {
			
			data2 := strings.Split(string(dataBuffer),"\r\n")[1]
		
			i, _ := strconv.Atoi(command[2])
			
			if i!=len(data2){
		
				errMsg := "ERR_CMD_ERR\r\n"
				conn.Write([]byte(errMsg))
				return;
				
				}
				
			version,writeError:=handleWrite(command,[]byte(data2))
			if writeError!=nil{
				 errMsg :="ERR_CMD_ERR\r\n"
				 conn.Write([]byte(errMsg))
			 	return;
				}
			
			conn.Write([]byte("OK "+ strconv.FormatInt(version,10)+"\r\n"))
			
	
			}
	
		case "cas":
		
			
		
			var time bool
			time=false
			
			if len(command)<4 || len(command)>5 || len(command[1])>250{
				errMsg := "ERR_CMD_ERR\r\n"
				conn.Write([]byte(errMsg))
				return;
				
			} else {
		
				data2 := strings.Split(string(dataBuffer),"\r\n")[1]
				
				mu.Lock()
				if val,ok := mymap[command[1]]; ok {	
				mu.Unlock()
				
				
				if val.exptime!=0{
				time=expirytime.After(val.creationtime)
				
				}else{
				time=true
				}
				
					if time{     //same as done in write
			
						if command[2]==strconv.FormatInt(val.version,10) {
			
							val.content=[]byte(data2)
				
							val.version=rand.Int63()
			
							conn.Write([]byte("OK "+strconv.FormatInt(val.version,10)+"\r\n"))
			
						} else {
				
							conn.Write([]byte("ERR_VERSION "+ strconv.FormatInt(val.version,10) +"\r\n"))
			
						}	
				
			 	 }	else {
			  		mu.Lock()
		     		delete(mymap,command[1])
		     		mu.Unlock()
		     		
		        	errMsg := "ERR_FILE_NOT_FOUND\r\n"
					conn.Write([]byte(errMsg))
					return
			  
			  		}
			  
			 } else {
			
				errMsg := "ERR_FILE_NOT_FOUND\r\n"
				conn.Write([]byte(errMsg))
				return
				
			}
			
	
		}
			
		case "delete":
				
			var time bool
			time=false
			
	    	if len(command) == 2 && len(command[1])<=250{
			
				if val,ok := mymap[command[1]]; ok {		
				
					if val.exptime!=0 {
					time=expirytime.After(val.creationtime)
					
					}else{
					time=true
					}
				
					if time{
	
						mu.Lock()
						delete(mymap,command[1])
						mu.Unlock()
				
						conn.Write([]byte("OK\r\n"))
			
					} else {
						errMsg := "ERR_FILE_NOT_FOUND\r\n"
						conn.Write([]byte(errMsg))
						return
					
					}
			
				
				}else{
					errMsg := "ERR_FILE_NOT_FOUND\r\n"
					conn.Write([]byte(errMsg))
					return
				}	
		
			}else {
				errMsg := "ERR_CMD_ERR\r\n"
				conn.Write([]byte(errMsg))
				return	
			}
		
}		
			
}

  	
}


func main() {
       serverMain()
}
