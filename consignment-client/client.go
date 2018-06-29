// consignment-cli/cli.go
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/Terry-Bui/go-microservice-example/consignment-service/proto/consignment"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"golang.org/x/net/context"
)

const (
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {

	cmd.Init()

	client := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)

	// Contact server and print out its response
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	consignment, err := parseFile(file)

	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	r, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Could not list consignment: %v", err)
	}

	for _, v := range getAll.Consignments {
		log.Println(v)
	}

}
