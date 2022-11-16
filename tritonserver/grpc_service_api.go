package tritonserver

import (
	"context"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http" // Package http provides HTTP client and server implementations.
	"time"
)

// Contexts that is necessary to communicate with Triton
type GRPCInferenceServiceAPIClient interface {
	GetServerLiveness(c echo.Context) error
	GetServerReadiness(c echo.Context) error
	GetModelMetadata(c echo.Context) error
	GetModelInferStats(c echo.Context) error
	LoadModel(c echo.Context) error
	UnloadModel(c echo.Context) error
	Infer(c echo.Context) error
}

// Contexts that is necessary to communicate with Triton
type gRPCInferenceServiceAPIClient struct {
	grpc    GRPCInferenceServiceClient
	url     string
	timeout int64
}

// ConnectToTritonWithGRPC Create GRPC Connection
func NewGRPCInferenceServiceAPIClient(url string, timeout int64) gRPCInferenceServiceAPIClient {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("couldn't connect to endpoint %s: %v", url, err)
	}
	client := gRPCInferenceServiceAPIClient{}
	client.grpc = NewGRPCInferenceServiceClient(conn)
	client.url = url
	client.timeout = timeout
	return client
}

// @Summary     Check Triton's liveness.
// @Description It returns true if the triton server is alive.
// @Accept      json
// @Produce     json
// @Success     200 {object} bool "Triton server's liveness"
// @Router      /liveness [get]
func (client *gRPCInferenceServiceAPIClient) GetServerLiveness(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(client.timeout)*time.Second)
	defer cancel()

	serverLiveRequest := ServerLiveRequest{}
	serverLiveResponse, err := client.grpc.ServerLive(ctx, &serverLiveRequest)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, serverLiveResponse.Live)
}

// @Summary     Check Triton's Readiness.
// @Description It returns true if the triton server is ready.
// @Accept      json
// @Produce     json
// @Success     200 {object} bool "Triton server's readiness"
// @Router      /readiness [get]
func (client *gRPCInferenceServiceAPIClient) GetServerReadiness(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(client.timeout)*time.Second)
	defer cancel()

	serverReadyRequest := ServerReadyRequest{}
	serverReadyResponse, err := client.grpc.ServerReady(ctx, &serverReadyRequest)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, serverReadyResponse.Ready)
}

// @Summary     Get model metadata.
// @Description It returns the requested model metadata
// @Accept      json
// @Produce     json
// @Param       model   query    string                true  "model name"
// @Param       version query    string                false "model version"
// @Success     200     {object} ModelMetadataResponse "Triton server's model metadata"
// @Router      /model-metadata [get]
func (client *gRPCInferenceServiceAPIClient) GetModelMetadata(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(client.timeout)*time.Second)
	defer cancel()

	modelMetadataRequest := ModelMetadataRequest{Name: c.QueryParam("model"), Version: c.QueryParam("version")}
	modelMetadataResponse, err := client.grpc.ModelMetadata(ctx, &modelMetadataRequest)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, modelMetadataResponse)
}

// @Summary     Get model inference statistics.
// @Description It returns the requested model's inference statistics.
// @Accept      json
// @Produce     json
// @Param       model   query    string                  true  "model name"
// @Param       version query    string                  false "model version"
// @Success     200     {object} ModelStatisticsResponse "Triton server's model statistics"
// @Router      /model-stats [get]
func (client *gRPCInferenceServiceAPIClient) GetModelInferStats(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(client.timeout)*time.Second)
	defer cancel()

	modelStatisticsRequest := ModelStatisticsRequest{Name: c.QueryParam("model"), Version: c.QueryParam("version")}
	modelStatisticsResponse, err := client.grpc.ModelStatistics(ctx, &modelStatisticsRequest)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, modelStatisticsResponse)
}

// @Summary     Load a model.
// @Description It requests to load a model. This is only allowed when polling is enabled.
// @Accept      json
// @Produce     json
// @Param       model query    string                      true "model name"
// @Success     200   {object} RepositoryModelLoadResponse "Triton server's model load response"
// @Router      /model-load [post]
func (client *gRPCInferenceServiceAPIClient) LoadModel(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(client.timeout)*time.Second)
	defer cancel()

	modelLoadRequest := RepositoryModelLoadRequest{ModelName: c.QueryParam("model")}
	modelLoadResponse, err := client.grpc.RepositoryModelLoad(ctx, &modelLoadRequest)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, modelLoadResponse)
}

// @Summary     Unload a model.
// @Description It requests to unload a model. This is only allowed when polling is enabled.
// @Accept      json
// @Produce     json
// @Param       model query    string                        true "model name"
// @Success     200   {object} RepositoryModelUnloadResponse "Triton server's model unload response"
// @Router      /model-unload [post]
func (client *gRPCInferenceServiceAPIClient) UnloadModel(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(client.timeout)*time.Second)
	defer cancel()

	modelUnloadRequest := RepositoryModelUnloadRequest{ModelName: c.QueryParam("model")}
	modelUnloadResponse, err := client.grpc.RepositoryModelUnload(ctx, &modelUnloadRequest)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, modelUnloadResponse)
}

// @Summary     Model inference api for the model with bytes a input and a bytes output.
// @Description It outputs a single bytes with a single bytes input.
// @Accept      json
// @Produce     json
// @Param       model   formData string             true  "model name"
// @Param       file    formData file               true  "input"
// @Param       version formData string             false "model version"
// @Success     200     {object} ModelInferResponse "Triton server's inference response"
// @Router      /infer [post]
func (client *gRPCInferenceServiceAPIClient) Infer(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(client.timeout)*time.Second)
	defer cancel()

	// Get the model information
	modelName := c.FormValue("model")
	modelVersion := c.FormValue("version")

	// Get the file
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	fileContent, err := file.Open()
	if err != nil {
		return err
	}
	defer fileContent.Close()
	rawInput, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return err
	}

	// Create request input / output tensors
	size := int64(len(rawInput))
	inferInputs := []*ModelInferRequest_InferInputTensor{{Name: "INPUT", Datatype: "UINT8", Shape: []int64{size}}}
	inferOutputs := []*ModelInferRequest_InferRequestedOutputTensor{{Name: "OUTPUT"}}

	// Create a request
	modelInferRequest := ModelInferRequest{
		ModelName:        modelName,
		ModelVersion:     modelVersion,
		Inputs:           inferInputs,
		Outputs:          inferOutputs,
		RawInputContents: [][]byte{rawInput},
	}

	// Get infer response
	modelInferResponse, err := client.grpc.ModelInfer(ctx, &modelInferRequest)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, modelInferResponse)
}
