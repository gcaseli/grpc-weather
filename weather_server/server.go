package main

// a estrutura 'server' possui um dado provider que é uma interface 'OpenWeatherMapProvider'
// a estrutura é
import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"grpc-weather/weather_server/providers"
	"grpc-weather/weatherpb"
)

type server struct {
	provider providers.OpenWeatherMapProvider
}

var (
	openWeatherMapKey string
)

func init() {

	openWeatherMapKey = os.Getenv("OPEN_WEATHER_KEY")
	if openWeatherMapKey == "" {
		log.Fatal("There isn't OPEN_WEATHER_KEY environment")
	}
}

func (server *server) WeatherDetails(context context.Context, req *weatherpb.WeatherRequest) (*weatherpb.WeatherResponse, error) {

	timeNow := time.Now()
	weatherInfo, err := server.provider.Search(req.GetLocation())
	if err != nil {
		return nil, err
	}
	defer elapsed(timeNow)
	data := &weatherpb.Weather{
		Description:    weatherInfo.Description,
		Found:          weatherInfo.Found,
		Temperature:    weatherInfo.Temperature,
		TemperatureMax: weatherInfo.TemperatureMax,
		TemperatureMin: weatherInfo.TemperatureMin,
		Country:        weatherInfo.Country,
	}

	return &weatherpb.WeatherResponse{
		Weather: data,
	}, nil
}

func main() {

	// if we crash the go code, we get the file name and the line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	server := &server{provider: providers.OpenWeatherMap{WeatherKey: openWeatherMapKey}}

	weatherpb.RegisterWeatherServiceServer(s, server)
	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		logar("Server started and listen port 50051....")
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

func elapsed(start time.Time) {
	log.Printf("[OpenWeatherMap] Request took %02f seconds\n", time.Since(start).Seconds())
}
