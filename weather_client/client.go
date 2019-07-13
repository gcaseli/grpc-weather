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

	if len(os.Args) == 0 {
		log.Fatalf("Missing location parameter")
	}

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

	if res.Weather.Found {
		fmt.Printf("Description is: %s\n", res.Weather.GetDescription())
		fmt.Printf("Temperature is: %.1fºC\n", res.Weather.GetTemperature())
		fmt.Printf("Temperature Max is: %.1fºC\n", res.Weather.GetTemperatureMax())
		fmt.Printf("Temperature Min is: %.1fºC\n", res.Weather.GetTemperatureMin())
		fmt.Printf("Country is: %s\n", res.Weather.GetCountry())
	} else {
		fmt.Printf("Could not find information for location %s\n", city)
	}
}
