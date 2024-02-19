package main

import (
	"broker/event"
	"broker/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayLoad  `json:"log,omitempty"`
	Mail   MailPayLoad `json:"mail,omitempty"`
}

type MailPayLoad struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayLoad struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItemViaRPC(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)

	default:
		app.errorJson(w, errors.New("unknown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayLoad) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceUrl := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJson(w, http.StatusAccepted, payload)

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	//create json and send to auth micro service
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	//make sure we get back correct status code.
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("invalid credentials"))
		return
	}
	if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling auth service"))
		return
	}

	//read response from auth service
	var jsonFromServive jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromServive)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if jsonFromServive.Error {
		app.errorJson(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromServive.Data

	app.writeJson(w, http.StatusAccepted, payload)

}

func (app *Config) sendMail(w http.ResponseWriter, mail MailPayLoad) {
	jsonData, _ := json.MarshalIndent(mail, "", "\t")

	//call the mail service
	mailServiceURL := "http://mailer-service/send"

	//post to mail service

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	//make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling mail service : response CODE != 202"))
	}

	//send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to" + mail.To

	app.writeJson(w, http.StatusAccepted, payload)

}

//emmit an event to rabbit MQ
func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayLoad) {

	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJson(w, err)
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via rabbitMQ"

	app.writeJson(w, http.StatusAccepted, payload)

}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmiiter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayLoad{
		Name: name,
		Data: msg,
	}

	json, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(json), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

//log via rpc
func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayLoad) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		log.Println(err)
		app.errorJson(w, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCserver.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	//connect to grpc server

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer conn.Close()

	//client
	c := logs.NewLogServieClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})

	if err != nil {
		app.errorJson(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJson(w, http.StatusAccepted, payload)

}
