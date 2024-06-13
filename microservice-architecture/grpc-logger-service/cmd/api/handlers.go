package main

import (
	"context"

	"github.com/PhilippWegner/hsbo-ma-inf-wegner-philipp/microservice-architecture/grpc-logger-service/database"
	"github.com/PhilippWegner/hsbo-ma-inf-wegner-philipp/microservice-architecture/grpc-logger-service/model"
)

type LogServiceServer struct {
	model.UnimplementedLogServiceServer
}

var mongo = database.ConnectMongodb()

func (api *LogServiceServer) WriteLog(ctx context.Context, in *model.LogRequest) (*model.LogResponse, error) {
	log := database.LogEntry{
		Name: in.LogEntry.GetName(),
		Data: in.LogEntry.GetData(),
	}

	err := mongo.InsertLog(log)
	if err != nil {
		return nil, err
	}

	return &model.LogResponse{
		Result: "Log written successfully!",
	}, nil
}
