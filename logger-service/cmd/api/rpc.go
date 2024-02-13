package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

type RPCserver struct{}

type RPCPlayload struct {
	Name string
	Data string
}

func (r *RPCserver) LogInfo(payload RPCPlayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo [Err]: ", err)
		return err
	}

	*resp = "Processed Payload via RPC :" + payload.Name
	return nil
}
