package messages

import (
	"github.com/DheerendraRathor/GoTracer/models"
)

type RenderResponseCode string

const (
	RenderRequestAccepted                 RenderResponseCode = "RenderRequestAccepted"
	RenderRequestTracingAlreadyInProgress                    = "RenderRequestTracingAlreadyInProgress"
)

type WebSocketMessage struct {
	Type        string
	Data        interface{}
	OperationId string
}

type RenderRequestMessage struct {
	Type        string
	Data        models.World
	OperationId string
}

type RenderRequestResponseMessage struct {
	Type        string
	Code        RenderResponseCode
	OperationId string
}

type PixelResultMessage struct {
	Type        string
	Data        models.Pixel
	OperationId string
}

type AgentStatus struct {
	Available bool
	Cores     int
}
