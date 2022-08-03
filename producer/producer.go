package main

import (
	_ "dtimer/docs"
	"dtimer/model"
	"encoding/json"
	"fmt"
	"github.com/go-stomp/stomp/v3"
	"github.com/google/uuid"
	"log"
	"net"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Swagger Producer API
// @version 1.0
// @description Timer Producer API.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:7070

const queueName = "/queue/timer"
const activeMQURL = "localhost:61613"
const basePath = "/v1/"

var v = validator.New()
var activeMQ *stomp.Conn
var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
	stomp.ConnOpt.Login("admin", "admin"),
}

func main() {
	netConn, err := net.Dial("tcp", activeMQURL)
	if err != nil {
		log.Fatalln("activemq dial fail", err)
	}
	activeMQ, err = stomp.Connect(netConn, options...)
	if err != nil {
		log.Fatalln("activemq dial fail", err)
	}

	e := echo.New()
	e.GET(basePath+"swagger/*", echoSwagger.WrapHandler)
	e.POST(basePath+"timer/start", startTimer)

	data, _ := json.MarshalIndent(e.Routes(), "", "  ")
	log.Println("Echo routes loaded: ", string(data))
	e.Logger.Fatal(e.Start(":7070"))
}

// Start Timer godoc
// @Summary Start a timer.
// @Description  Start a timer of duration in secs and get a callback once timeout on the mention callback url
// @Tags Timer
// @Accept  json
// @Produce  json
// @Param startTimer body StartTimerRequest true "request body"
// @Success 200 {object} Response
// @Router /v1/timer/start [post]
func startTimer(c echo.Context) error {

	req := new(model.StartTimerRequest)
	err := bindAndValidate(c, req)
	if err != nil {
		return returnJson(c, http.StatusBadRequest, "validation failed")
	}

	qdata := generateQueueData(req)

	data, err := json.Marshal(qdata)
	if err != nil {
		return returnJson(c, http.StatusInternalServerError, "marshal json fail")
	}

	// Produce message to active MQ
	err = activeMQ.Send(queueName, "text/plain",
		data, stomp.SendOpt.Header("AMQ_SCHEDULED_DELAY", fmt.Sprintf("%d", req.Duration*1000)))
	if err != nil {
		log.Println("ERROR: failed to write to queue", err)
		return returnJson(c, http.StatusInternalServerError, "fail to send data to queue")
	}

	return c.JSON(http.StatusOK, &model.Response{
		TimerID: qdata.TimerID,
		Message: "success",
	})
}

func generateQueueData(req *model.StartTimerRequest) *model.QueueData {
	q := &model.QueueData{}
	id, _ := uuid.NewUUID()
	q.TimerID = id.String()
	q.Duration = req.Duration
	q.CallbackMethod = req.CallbackMethod
	q.CallbackURL = req.CallbackURL
	return q
}

func returnJson(c echo.Context, status int, message string) error {
	return c.JSON(status, &model.Response{
		Message: message,
	})
}

func bindAndValidate(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := v.Struct(req); err != nil {
		return err
	}
	return nil
}
