# grpc-weather

Simple example of using gRPC and call a HTTP service that provide weather information

### GO DEPENDENCIES ### 

```
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
```

### Build and using ### 

You'll need do have an API KEY from `https://openweathermap.org/api`

``` make build-server ```

And the run this:

```OPEN_WEATHER_KEY="yourApiKey" ./weather_server/server```

If everything goes well:

``` 2019/07/13 09:28:54 server.go:102: Server started and listen port 50051.... ```

Similarly, you can build the client:

```make build-client ```

And then run it, providing a location:

``` ./weather_client/client Sao Paulo ```

An example output:

``` Weather client:  Sao Paulo 
Description is: clear sky
Temperature is: 17.5ºC
Temperature Max is: 19.5ºC
Temperature Min is: 15.5ºC
Country is: BR ``` 