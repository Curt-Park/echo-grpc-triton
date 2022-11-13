package main

import (
	"context"
	_ "echo-grpc-triton/docs" // docs is generated by Swag CLI, you have to import it.
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
	"log"
	"net/http" // Package http provides HTTP client and server implementations.
	"time"
)

// Flags contains the information to send requests to Triton inference server.
type Flags struct {
	ModelName    string
	ModelVersion string
	URL          string
	TIMEOUT      int64
}

// Contexts that is necessary to communicate with Triton
type Client struct {
	grpc GRPCInferenceServiceClient
}

// Global variables
var flags = Flags{}
var client = Client{}

// parseFlags parses the arguments and initialize the flags.
func (flags *Flags) parseFlags() {
	// https://github.com/NVIDIA/triton-inference-server/tree/master/docs/examples/model_repository/simple
	flag.StringVar(&flags.ModelName, "m", "simple", "Name of model being served. (Required)")
	flag.StringVar(&flags.ModelVersion, "x", "", "Version of model. Default: Latest Version.")
	flag.StringVar(&flags.URL, "u", "localhost:8001", "Inference Server URL. Default: localhost:8001")
	flag.Int64Var(&flags.TIMEOUT, "t", 10, "Timeout. Default: 10 Sec.")
	flag.Parse()
}

// ConnectToTritonWithGRPC Create GRPC Connection
func (client *Client) connectToTriton(url string) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("couldn't connect to endpoint %s: %v", url, err)
	}
	client.grpc = NewGRPCInferenceServiceClient(conn)
}

// @contact.name  Curt-Park
// @contact.email www.jwpark.co.kr@gmail.com
func main() {
	// Parse the args
	flags.parseFlags()
	log.Println("Flags:", flags)

	// Check the gRPC connection well-established
	client.connectToTriton(flags.URL)

	// Create a server with echo
	e := echo.New()
	// Logger middleware logs the information about each HTTP request
	e.Use(middleware.Logger())

	// APIs
	e.GET("/", getHealthCheck)
	e.GET("/liveness", getServerLiveness)
	e.GET("/readiness", getServerLiveness)
	e.GET("/model-metadata", getModelMetadata)

	// Swagger
	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":8080"))
}

// @Summary     Healthcheck
// @Description It returns true if the api server is alive
// @Accept      json
// @Produce     json
// @Success     200 {object} bool "API server's liveness"
// @Router      / [get]
func getHealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, true)
}

// @Summary     Check Triton's liveness
// @Description It returns true if the triton server is alive
// @Accept      json
// @Produce     json
// @Success     200 {object} bool "Triton server's liveness"
// @Router      /liveness [get]
func getServerLiveness(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(flags.TIMEOUT)*time.Second)
	defer cancel()

	serverLiveRequest := ServerLiveRequest{}
	serverLiveResponse, err := client.grpc.ServerLive(ctx, &serverLiveRequest)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, serverLiveResponse.Live)
}

// @Summary     Check Triton's Readiness
// @Description It returns true if the triton server is ready
// @Accept      json
// @Produce     json
// @Success     200 {object} bool "Triton server's readiness"
// @Router      /readiness [get]
func getServerReadiness(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(flags.TIMEOUT)*time.Second)
	defer cancel()

	serverReadyRequest := ServerReadyRequest{}
	serverReadyResponse, err := client.grpc.ServerReady(ctx, &serverReadyRequest)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, serverReadyResponse.Ready)
}

// @Summary     Get model metadata
// @Description It returns the requested model metadata
// @Accept      json
// @Produce     json
// @Param       model   query    string true "model name"
// @Param       version query    string false "model version"
// @Success     200 {object} ModelMetadataResponse "Triton server's model metadata"
// @Router      /model-metadata [get]
func getModelMetadata(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(flags.TIMEOUT)*time.Second)
	defer cancel()

    // Create status request for a given model
    modelName := c.QueryParam("model")
    modelVersion := c.QueryParam("version")
    log.Println(modelName, modelVersion)
    modelMetadataRequest := ModelMetadataRequest{
        Name: modelName,
        Version: modelVersion,
    }

    // Submit modelMetadata request to server
    modelMetadataResponse, err := client.grpc.ModelMetadata(ctx, &modelMetadataRequest)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, modelMetadataResponse)
}
