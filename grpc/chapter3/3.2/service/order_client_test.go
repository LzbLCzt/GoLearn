package service

import (
	ecommerce "GoLearn/grpc/chapter3/3.2/proto"
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"io"
	"log"
	"testing"
	"time"
)

func TestOrderClient(t *testing.T) {
	conn, err := grpc.Dial("localhost:8089", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := ecommerce.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	//order1 := ecommerce.Order{Id: "101", Items: []string{"iPhone XS", "Mac Book Pro"}, Destination: "San Jose, CA", Price: 2300.00}
	searchStream, _ := client.SearchOrders(ctx, &wrappers.StringValue{Value: "Apple"})
	for {
		searchOrder, err := searchStream.Recv()
		if err == io.EOF {
			log.Printf("EOF")
			break
		}
		log.Print("Search Result : ", searchOrder)
	}
}
