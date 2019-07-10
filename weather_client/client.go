package main

import (
	"context"
	"fmt"
	"grpc-weather/weatherpb"
	"log"

	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Weather client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := weatherpb.NewWeatherServiceClient(cc)

	res, err := c.WeatherDetails(context.Background(), &weatherpb.WeatherRequest{
		Location: "London",
	})

	if err != nil {
		log.Fatalf("error when creating blog: %v", err)
	}

	fmt.Printf("Temperature description is: %v\n", res.Weather.GetDescription())
}
