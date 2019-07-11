package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"grpc-weather/weatherpb"
)

const (
	openWeatherMapURL = "http://api.openweathermap.org/data/2.5/weather"
)

type server struct{}

type openWeatherMapResult struct {
	Name string `json:"name"`
	Main struct {
		KelvinTemp    float64 `json:"temp"`
		KelvinTempMin float64 `json:"temp_min"`
		KelvinTempMax float64 `json:"temp_max"`
	} `json:"main"`
	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func (server *server) WeatherDetails(context context.Context, req *weatherpb.WeatherRequest) (*weatherpb.WeatherResponse, error) {

	key := os.Getenv("OPEN_WEATHER_KEY")
	if key == "" {
		log.Fatal("There isn't OPEN_WEATHER_KEY environment")
	}

	queryURL := fmt.Sprintf("%s?q=%s&appid=%s", openWeatherMapURL, url.QueryEscape(req.GetLocation()), key)

	response, err := http.Get(queryURL)
	if err != nil {
		log.Fatalln(err)
	}

	if response.StatusCode != 200 {
		respErr := fmt.Errorf("Unexpected response: %s", response.Status)
		log.Fatalln(fmt.Sprintf("Request failed: %v", respErr))
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var result openWeatherMapResult
	json.Unmarshal(body, &result)

	fmt.Println(string(body))
	fmt.Println(result)

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
