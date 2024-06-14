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
	logEntry := database.LogEntry{
		Name: in.LogEntry.GetName(),
		Data: in.LogEntry.GetData(),
	}

	err := mongo.InsertLog(logEntry)
	if err != nil {
		return &model.LogResponse{
			Result: "failed!",
		}, err
	}

	return &model.LogResponse{
		Result: "Log written successfully!",
	}, nil
}
