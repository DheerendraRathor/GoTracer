package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"

	"net/http"

	"math"

	"github.com/DheerendraRathor/GoTracer/models"
	"github.com/DheerendraRathor/GoTracer/net/constants"
	"github.com/DheerendraRathor/GoTracer/utils"
	"github.com/gorilla/websocket"
	"gopkg.in/cheggaaa/pb.v1"
)

const (
	pongTimeout = 60 * time.Second
)

var renderSpecFile string
var agentsFile string

func init() {
	flag.StringVar(&renderSpecFile, "spec", "sample_world.json", "Name of JSON file containing rendering spec")
	flag.StringVar(&agentsFile, "agents", "agents.json", "Name of JSON file containing list of agents")
}

type Agent struct {
	URL           string
	Cores         int
	Conn          *websocket.Conn
	Env           models.Specification
	ResultChannel chan<- models.Pixel
	Completed     chan<- *Agent
	WorkDone      bool
}

func (a *Agent) Initialize() {
	workDone := false

	defer func() {
		a.WorkDone = workDone
		a.Completed <- a
		a.Conn.Close()
	}()

	a.Conn.SetPingHandler(
		func(appData string) error {
			err := a.Conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(pongTimeout))
			if err == websocket.ErrCloseSent {
				return nil
			} else if e, ok := err.(net.Error); ok && e.Temporary() {
				return nil
			}
			return err
		},
	)

	renderRequest := messages.RenderRequestMessage{
		Type: messages.RenderRequest,
		Data: a.Env,
	}

	err := a.Conn.WriteJSON(renderRequest)
	if err != nil {
		log.Printf("Unable to send JSON message to agent. %s\n", err)
		return
	}

	var message messages.WebSocketMessage
	var pixelMessage messages.PixelResultMessage

	for {
		_, rawMsg, err := a.Conn.ReadMessage()
		if err != nil {
			fmt.Println("Error in reading message:", err)
			return
		}

		json.Unmarshal(rawMsg, &message)

		messageType := message.Type
		switch messageType {
		case messages.PixelResult:
			json.Unmarshal(rawMsg, &pixelMessage)
			pixel := pixelMessage.Data
			a.ResultChannel <- pixel
		case messages.RenderingCompleted:
			workDone = true
			return
		}
	}
}

func main() {

	flag.Parse()

	specFile, e := ioutil.ReadFile(renderSpecFile)
	if e != nil {
		panic(fmt.Sprintf("Render Spec File error: %v\n", e))
	}

	agentsFile, e := ioutil.ReadFile(agentsFile)
	if e != nil {
		panic(fmt.Sprintf("Agents list file error: %v\n", e))
	}

	agents := make([]string, 0)
	json.Unmarshal(agentsFile, &agents)

	var env models.Specification
	json.Unmarshal(specFile, &env)

	pngImage := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{env.Image.Width, env.Image.Height},
	})

	connectedAgents := []*Agent{}

	myClient := &http.Client{
		Timeout: time.Second * 60,
	}

	var agentStatus messages.AgentStatus

	var renderedChannel = make(chan models.Pixel, 100)
	var workDoneChannel = make(chan *Agent, len(agents))
	var availableCores = 0

	for _, agent := range agents {
		wsURL := fmt.Sprintf("ws://%s", agent)
		httpUrl := fmt.Sprintf("http://%s/status", agent)

		response, err := myClient.Get(httpUrl)
		if err != nil {
			log.Printf("Agent '%s' is unavailable. \n\nError: %s\n", agent, err.Error())
			continue
		}

		json.NewDecoder(response.Body).Decode(&agentStatus)
		response.Body.Close()

		if !agentStatus.Available {
			log.Printf("Agent '%s' is busy\n", agent)
			continue
		}

		var dialer *websocket.Dialer

		conn, _, err := dialer.Dial(wsURL, nil)
		if err != nil {
			fmt.Printf("Unable to dial ws connection to agent %s with error %s\n", agent, err)
			continue
		}

		connectedAgent := Agent{
			Cores:         agentStatus.Cores,
			Conn:          conn,
			URL:           agent,
			ResultChannel: renderedChannel,
			Completed:     workDoneChannel,
		}

		availableCores += agentStatus.Cores

		connectedAgents = append(connectedAgents, &connectedAgent)
	}

	if len(connectedAgents) == 0 {
		log.Fatalln("Unable to connect to any agent. Please try again later")
	}

	linesPerCore := int(math.Ceil(float64(env.Image.Width) / float64(availableCores)))

	minWidth, maxWidth := 0, 0

	for _, connectedAgent := range connectedAgents {
		cores := connectedAgent.Cores
		linesToAgent := linesPerCore * cores
		agentEnv := models.Specification(env)

		maxWidth = minWidth + linesToAgent + 1
		if maxWidth > env.Image.Width {
			maxWidth = env.Image.Width
		}

		agentEnv.Image.Patch[2] = minWidth
		agentEnv.Image.Patch[3] = maxWidth

		minWidth = maxWidth
		connectedAgent.Env = agentEnv

		go connectedAgent.Initialize()
	}

	var wg sync.WaitGroup
	wg.Add(1)

	total := env.Image.Width * env.Image.Height
	progressBar := pb.StartNew(total)
	progressBar.ShowFinalTime = true
	progressBar.ShowTimeLeft = false

	go func() {
		defer wg.Done()
		var pixel models.Pixel
		agentsToWaitFor := len(connectedAgents)
		for {
			select {
			case pixel = <-renderedChannel:
				progressBar.Increment()
				rgbaColor := color.RGBA{pixel.Color[0], pixel.Color[1], pixel.Color[2], 255}
				pngImage.Set(pixel.I, pixel.J, rgbaColor)
			case agent := <-workDoneChannel:
				agentsToWaitFor -= 1
				log.Printf("Agent '%s' finished. Status: %t", agent.URL, agent.WorkDone)
				if agentsToWaitFor == 0 {
					return
				}
			}
		}
	}()

	wg.Wait()

	pngFile := utils.CreateNestedFile(env.Image.OutputFile)
	defer pngFile.Close()

	png.Encode(pngFile, pngImage)
}
