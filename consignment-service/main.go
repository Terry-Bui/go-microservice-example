// consignment-service/main.go
package main

import (
	"fmt"
	// Import generated protobuf code
	pb "github.com/Terry-Bui/go-microservice-example/consignment-service/proto/consignment"
	vesselProto "github.com/Terry-Bui/go-microservice-example/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

// Repository is the interface for ShippingService Server API
type Repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// ConsignmentRepository is the datastore
type ConsignmentRepository struct {
	consignments []*pb.Consignment
}

// Create a consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

// GetAll returns all consignments
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo         Repository
	vesselClient vesselProto.VesselServiceClient
}

// A create method for ShippingService, takes a context and a request as an argument,
// this is then handled by the gRPC server
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity: int32(len(req.Containers)),
	}
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}

	// Set VesselId as vessel returned from vessel service
	req.VesselId = vesselResponse.Vessel.Id

	// Save consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	res.Consignments = s.repo.GetAll()
	return nil
}

func main() {
	repo := &Repository{}

	srv := micro.NewService(
		// Name must match package name in protobuf definition
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	// Init will parse command line flags
	srv.Init()
	// Register ShippingService handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}

}
