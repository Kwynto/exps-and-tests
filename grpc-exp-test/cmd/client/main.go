package main

import (
	"context"
	"fmt"
	ms "grpc-exp-test/internal/proto/message_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var port = ":8080"

func AboutToSayIt(ctx context.Context, m ms.MessageServiceClient, text string) (*ms.Response, error) {
	request := &ms.Request{
		Text:    text,
		Subtext: "New Message!",
	}
	r, err := m.SayIt(ctx, request)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func main() {
	// // Профилирование процессора go tool pprof -http=localhost:8081 ./data/cpuProfile.out
	// cpuFile, err := os.Create("./data/cpuProfile.out")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer cpuFile.Close()
	// pprof.StartCPUProfile(cpuFile)
	// defer pprof.StopCPUProfile()

	// for i := 0; i < 100; i++ {

	// conn, err := grpc.Dial(port, grpc.WithInsecure())
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Dial:", err)
		return
	}
	client := ms.NewMessageServiceClient(conn)
	r, err := AboutToSayIt(context.Background(), client, "My Message!")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Response Text:", r.Text)
	fmt.Println("Response SubText:", r.Subtext)

	// }

	// // Профилирование памяти go tool pprof -http=localhost:8081 ./data/memoryProfile.out
	// memory, err := os.Create("./data/memoryProfile.out")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// // Запись профилирования памяти
	// err = pprof.WriteHeapProfile(memory)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// memory.Close()

}
