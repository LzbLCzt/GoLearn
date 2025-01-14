package service

import (
	ecommerce "GoLearn/grpc/chapter3/3.2/proto"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
	"testing"
)

var orderMap = make(map[string]ecommerce.Order)

type server struct {
	orderMap map[string]*ecommerce.Order
}

func (s *server) SearchOrders(searchQuery *wrappers.StringValue, stream ecommerce.OrderManagement_SearchOrdersServer) error {
	for key, order := range orderMap {
		log.Print(key, order)
		for _, itemStr := range order.Items {
			log.Print(itemStr)
			if strings.Contains(itemStr, searchQuery.Value) {
				// Send the matching orders in a stream
				err := stream.Send(&order)
				if err != nil {
					return fmt.Errorf("error sending message to stream : %v", err)
				}
				log.Print("Matching Order Found : " + key)
				break
			}
		}
	}
	return nil
}
func TestOrderService(t *testing.T) {
	initSampleData()
	listen, err := net.Listen("tcp", ":8089")
	if err != nil {
		t.Errorf("Error while listening : %v", err)
		return
	}
	s := grpc.NewServer()
	ecommerce.RegisterOrderManagementServer(s, &server{})
	if err = s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initSampleData() {
	orderMap["102"] = ecommerce.Order{Id: "102", Items: []string{"Google Pixel 3A", "Mac Book Pro"}, Destination: "Mountain View, CA", Price: 1800.00}
	orderMap["103"] = ecommerce.Order{Id: "103", Items: []string{"Apple Watch S4"}, Destination: "San Jose, CA", Price: 400.00}
	orderMap["104"] = ecommerce.Order{Id: "104", Items: []string{"Google Home Mini", "Google Nest Hub"}, Destination: "Mountain View, CA", Price: 400.00}
	orderMap["105"] = ecommerce.Order{Id: "105", Items: []string{"Amazon Echo"}, Destination: "San Jose, CA", Price: 30.00}
	orderMap["106"] = ecommerce.Order{Id: "106", Items: []string{"Amazon Echo", "Apple iPhone XS"}, Destination: "Mountain View, CA", Price: 300.00}
}

func (s *server) AddOrder(ctx context.Context, orderReq *ecommerce.Order) (*wrappers.StringValue, error) {
	return nil, nil
}

func (s *server) GetOrder(ctx context.Context, orderId *wrappers.StringValue) (*ecommerce.Order, error) {
	return nil, nil
}

func (s *server) UpdateOrders(stream ecommerce.OrderManagement_UpdateOrdersServer) error {
	return nil
}

func (s *server) ProcessOrders(stream ecommerce.OrderManagement_ProcessOrdersServer) error {
	return nil
}
