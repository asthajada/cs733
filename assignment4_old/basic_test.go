package main

import (
	"fmt"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	rafts := makeRafts()
	time.Sleep(4 * time.Second)

	flag := false
	for i := 0; i < len(rafts); i++ {
		if rafts[i].server.State == "leader" {
			flag = true
			break

		}
	}
	expectBool(t, flag, true) //testing if a leader is elected

	ldr := getLeader(rafts)

	ldr.Append([]byte("foo"))
	time.Sleep(1 * time.Second)

	for _, node := range rafts {

		select { // to avoid blocking on channel.

		case ci := <-node.CommitChannel():

			if ci.Err != nil {
				t.Fatal(ci.Err)
			}
			if string(ci.Data) != "foo" {

				t.Fatal("Got different data")
			}

		default:
			t.Fatal("Expected message on all nodes")
		}
	}

	ldr.Append([]byte("bye"))

	time.Sleep(1 * time.Second)

	/*for i := 0; i < len(rafts); i++ {
		fmt.Println("State :",rafts[i].server.State, rafts[i].server.MyID)
		fmt.Println("log contents",rafts[i].server.Log)
	}
	*/

	//fmt.Println("shutting down leader")

	//ldr.Shutdown()

	time.Sleep(1 * time.Second)

	for i := 0; i < len(rafts); i++ {
		fmt.Println("State :", rafts[i].server.State, rafts[i].server.MyID)
		fmt.Println("log contents", rafts[i].server.Log)
	}

}

func expectString(t *testing.T, a string, b string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}

func expectBool(t *testing.T, a bool, b bool) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}
