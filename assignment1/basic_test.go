package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
	"sync"
	
)


func TestTCPSimple(t *testing.T){
	var wg sync.WaitGroup
 
 	go serverMain()
	wg.Add(3)

	for i:=0;i<3;i++ {
		
		go TCPSimple(&wg,t)
		
	}

	wg.Wait()
}

// Simple serial check of getting and setting
func TCPSimple(wg *sync.WaitGroup,t *testing.T) {
	defer wg.Done()


    time.Sleep(1 * time.Second) // one second is enough time for the server to start
	name := "hi.txt"
	contents := "bye"
	exptime := 300000
	
	conn, err := net.Dial("tcp", "localhost:8080")
	
	if err != nil {
		t.Error(err.Error())
	}
	
	scanner := bufio.NewScanner(conn)
	
	
	
	// Write to a file with expiry time
	
    fmt.Fprintf(conn, "write %v %v %v\r\n%v\r\n", name, len(contents),exptime, contents)

	scanner.Scan() // read first line
	resp := scanner.Text() // extract the text from the buffer
	arr := strings.Split(resp, " ") // split into OK and <version>
	expect(t, arr[0], "OK")
	version, err := strconv.ParseInt(arr[1], 10, 64) // parse version as number

	if err != nil {
		t.Error("Non-numeric version found")
	}
	
	fmt.Println(scanner.Text())


	//reading from a file

	fmt.Fprintf(conn, "read %v\r\n", name) // try a read now
	scanner.Scan()

	fmt.Println(scanner.Text())
	arr = strings.Split(scanner.Text(), " ")
	expect(t, arr[0], "CONTENTS")
	//expect(t, arr[1], fmt.Sprintf("%v", version)) // expect only accepts strings, convert int version to string
	expect(t, arr[2], fmt.Sprintf("%v", len(contents)))	
	scanner.Scan()
	expect(t, contents, scanner.Text())

	fmt.Println(scanner.Text())
	

	// Write to a file without providing expirytime
    fmt.Fprintf(conn, "write %v %v\r\n%v\r\n", name, len(contents), contents)

	scanner.Scan() // read first line
	resp = scanner.Text() // extract the text from the buffer
	arr = strings.Split(resp, " ") // split into OK and <version>
	expect(t, arr[0], "OK")
	version, err = strconv.ParseInt(arr[1], 10, 64) // parse version as number

	if err != nil {
		t.Error("Non-numeric version found")
	}
	
	fmt.Println(scanner.Text())
	
    //compare and swap with right version number
	fmt.Fprintf(conn, "cas %v %v %v\r\n%v\r\n", name,version, len(contents), contents)
	scanner.Scan()
	fmt.Println(scanner.Text())

	
	fmt.Fprintf(conn, "read %v\r\n", name) // try a read now
	scanner.Scan()

	fmt.Println(scanner.Text())
	arr = strings.Split(scanner.Text(), " ")
	expect(t, arr[0], "CONTENTS")
	expect(t, arr[2], fmt.Sprintf("%v", len(contents)))	
	scanner.Scan()
	expect(t, contents, scanner.Text())

	fmt.Println(scanner.Text())
	
	
	//compare and swap with wrong version number
    contents="hello"
	version=123
	fmt.Fprintf(conn, "cas %v %v %v\r\n%v\r\n", name,version, len(contents), contents)
	scanner.Scan()
	fmt.Println(scanner.Text())

	//deleting a file

	fmt.Fprintf(conn, "delete %v\r\n", name) // try a read now
	scanner.Scan()
	
	fmt.Println(scanner.Text())
	
	
}

// Useful testing function
func expect(t *testing.T, a string, b string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}