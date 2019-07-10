package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"grpc-weather/weatherpb"
)

type server struct{}

func (server *server) WeatherDetails(context context.Context, req *weatherpb.WeatherRequest) (*weatherpb.WeatherResponse, error) {
	data := &weatherpb.Weather{
		Description:    "Very cold",
		Found:          true,
		Temperature:    11.0,
		TemperatureMax: 12.0,
		TemperatureMin: 5.0,
	}

	return &weatherpb.WeatherResponse{
		Weather: data,
	}, nil
}

func main() {

	// if we crash the go code, we get the file name and the line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logar("WeatherService start")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	weatherpb.RegisterWeatherServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		logar("Starting server....")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
	}()

	// Wait for control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch

	logar("Stoping the server")
	s.Stop()
	logar("Closing the listener")
	lis.Close()

}

func logar(formato string, valores ...interface{}) {
	log.Printf(fmt.Sprintf("%s\n", formato), valores...)
}
