package main

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/Terry-Bui/go-microservice-example/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
)

// Repository is the interface for vessel-service
type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

// VesselRepository - datastore for vessel
type VesselRepository struct {
	vessels []*pb.Vessel
}

// FindAvailable - checks a specification against a map of vessels.
// A vessel is matched and return if the capacity and max weight of
// the vessel is less than the specification.
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New("No vessel found according to specification")
}

// grpc service handler
type service struct {
	repo Repository
}

func (s *service) FindAvailable (ctx context.Context, req *pb.Specification, res *pb.Response) error {
	// Find next available vessel
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return err
	}
	// Set vessel as part of response message type
	res.Vessel = vessel
	return nil
}

func main() {
	// Create a array of vessels
	vessels := []*pb.Vessel{
		&pb.Vessel{Id: "vessel001", Name: "Boaty McBoatface". MaxWeight: 200000, Capacity: 500},
	}
	repo := &VesselRepository{vessels}

	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.version("latest"),
	)
	srv.Init()

	pb.RegisterVesselServiceHandler(srv.Server(), &service{repo})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
