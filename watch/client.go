package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

func doClientWork(client *rpc.Client) {
	// khởi chạy một Goroutine riêng biệt để giám sát khóa thay đổi
	go func() {
		var keyChanged string
		// lời gọi `watch` synchronous sẽ block cho đến khi
		// có khóa thay đổi hoặc timeout
		err := client.Call("KVStoreService.Watch", 30, &keyChanged)
		if err != nil {
			log.Fatal("Call Watch: ", err)
		}
		fmt.Println("watch:", keyChanged)
	}()

	err := client.Call(
		//  giá trị KV được thay đổi bằng phương thức `Set`
		"KVStoreService.Set", [2]string{"abc", "abc-value"},
		new(struct{}),
	)
	//  set lại lần nữa để giá trị value của key 'abc' thay đổi
	err = client.Call(
		"KVStoreService.Set", [2]string{"abc", "another-value"},
		new(struct{}),
	)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 3)
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	doClientWork(client)
}
