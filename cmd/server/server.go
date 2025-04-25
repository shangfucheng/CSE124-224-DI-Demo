package main

import (
	utils "chat_server/internal"
	"chat_server/internal/protos/pb"
	"context"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	clients map[string]chan *pb.Message
	mutex   sync.Mutex
}

func newServer() *ChatServer {
	return &ChatServer{
		clients: make(map[string]chan *pb.Message),
	}
}

func (s *ChatServer) Join(req *pb.JoinRequest, stream pb.ChatService_JoinServer) error {
	msgCh := make(chan *pb.Message, 100)
	s.mutex.Lock()
	s.clients[req.Username] = msgCh
	s.mutex.Unlock()
	defer func() {
		s.mutex.Lock()
		delete(s.clients, req.Username)
		s.mutex.Unlock()
	}()

	go func() {
		<-stream.Context().Done()
		log.Panicln("user disconnect", req.Username)
		s.mutex.Lock()
		delete(s.clients, req.Username)
		s.mutex.Unlock()
	}()

	for msg := range msgCh {
		if err := stream.Send(msg); err != nil {
			log.Panicln("err send msg", err)
			return err
		}
	}
	return nil
}

func (s *ChatServer) SendMessage(ctx context.Context, msg *pb.Message) (*pb.Ack, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	msg.Timestamp = time.Now().Unix()
	for user, ch := range s.clients {
		ch <- msg
		log.Println("send data from: %v to %v", msg.Username, user)
	}

	return &pb.Ack{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", utils.ServerAddr)
	if err != nil {
		log.Fatalln("error start listener")
	}
	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, newServer())
	if err := s.Serve(lis); err != nil {
		log.Fatalln("serve fail")
	}
}
