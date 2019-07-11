package main

import (
	"context"
	"fmt"
	"grpc-weather/weatherpb"
	"log"
	"os"

	"google.golang.org/grpc"
)

func main() {

	elements := os.Args[1:]
	var city string
	for _, name := range elements {
		city += name + " "
	}

	fmt.Println("Weather client: ", city)

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := weatherpb.NewWeatherServiceClient(cc)

	res, err := c.WeatherDetails(context.Background(), &weatherpb.WeatherRequest{
		Location: city,
	})

	if err != nil {
		log.Fatalf("error when creating blog: %v", err)
	}

	fmt.Printf("Temperature description is: %v\n", res.Weather.GetDescription())
}
