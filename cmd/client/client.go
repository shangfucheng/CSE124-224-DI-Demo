package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	utils "chat_server/internal"
	pb "chat_server/internal/protos/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// grpc connection
	conn, err := grpc.NewClient(utils.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
		return
	}
	defer conn.Close()

	// create a new client
	client := pb.NewChatServiceClient(conn)

	// read user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	ctx := context.Background()

	// join the chat
	stream, err := client.Join(ctx, &pb.JoinRequest{Username: username})
	if err != nil {
		log.Fatalf("Error joining chat: %v", err)
	}
	fmt.Println("Welcome to the chat! Type your messages below:")

	go func() {
		for {
			// wait for messages from the server
			msg, err := stream.Recv()
			if err != nil {
				log.Fatalf("Error receiving: %v", err)
			}
			fmt.Printf("[%v] %s: %s\n", time.Unix(msg.Timestamp, 0).Format(time.RFC3339), msg.Username, msg.Content)
		}
	}()

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		// clear the last line (input line) just make it look nice
		// this is a hacky way to clear the last line
		// it may not work on all terminals
		fmt.Print("\033[1A\033[2K")

		// send message to server
		ack, err := client.SendMessage(ctx, &pb.Message{
			Username:  username,
			Content:   text,
			Timestamp: time.Now().Unix(),
		})
		if err != nil {
			log.Printf("Send error: %v", err)
		}
		if !ack.Success {
			log.Printf("Send failed")
		}
	}
}
