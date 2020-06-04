package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/naspinall/Hive/pkg/config"
	"github.com/naspinall/Hive/pkg/models"
	"google.golang.org/grpc"
)

const port = 50000

func listen() {

	cfg := config.LoadConfig()
	dbCfg := cfg.Database

	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(true),
		models.WithMeasurements(),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("Cannot listen on port %d", port)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	WithMeasurements(grpcServer, services.Measurement)

	// Listen
	grpcServer.Serve(lis)

}
