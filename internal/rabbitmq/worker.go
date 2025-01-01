package rabbitmq

import (
	"encoding/json"
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"golang-api-starter/internal/interfaces"
	"log"
	"time"
)

type UploadRequest struct {
	File string `json:"file"`
}

func HandleTestQueue(rbmqWorker interfaces.IRbmqWorker) {
	queueName := *cfg.RabbitMqConf.Queues.TestQueue
	rabbitMQ, err := NewRabbitMQ(GetUrl(), queueName)
	if err != nil {
		logger.Fatalf(err.Error())
	}
	defer rabbitMQ.Close()

	msgs, err := rabbitMQ.channel.Consume(
		rabbitMQ.queue.Name, // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		log.Fatal(err)
	}

	logger.Infof("Test worker started. Waiting for messages...")

	for msg := range msgs {
		var uploadRequest UploadRequest
		err := json.Unmarshal(msg.Body, &uploadRequest)
		if err != nil {
			logger.Errorf(err.Error())
			continue
		}

		// Process the file upload here
		logger.Infof("%s: Received file upload: %d bytes, content: %v\n", queueName, len(uploadRequest.File), uploadRequest)

		// Simulate processing time
		time.Sleep(2 * time.Second)

		logger.Infof("%s:File upload processed successfully", queueName)
	}
}

func HandleLogQueue(rbmqWorker interfaces.IRbmqWorker) {
	queueName := *cfg.RabbitMqConf.Queues.LogQueue
	rabbitMQ, err := NewRabbitMQ(GetUrl(), queueName)
	if err != nil {
		logger.Fatalf(err.Error())
	}
	defer rabbitMQ.Close()

	msgs, err := rabbitMQ.channel.Consume(
		rabbitMQ.queue.Name, // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		log.Fatal(err)
	}

	logger.Infof("Log worker started. Waiting for logs...")

	for msg := range msgs {
		rbmqWorker.HandleLogFromQueue(msg.Body)
		time.Sleep(2 * time.Second)

		logger.Infof("%s: log processed successfully", queueName)
	}
}

func HandleEmailQueue(rbmqWorker interfaces.IRbmqWorker) {
	queueName := *cfg.RabbitMqConf.Queues.EmailQueue
	rabbitMQ, err := NewRabbitMQ(GetUrl(), queueName)
	if err != nil {
		logger.Fatalf(err.Error())
	}
	defer rabbitMQ.Close()

	msgs, err := rabbitMQ.channel.Consume(
		rabbitMQ.queue.Name, // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		log.Fatal(err)
	}

	logger.Infof("Email worker started. Waiting for emails...")

	for msg := range msgs {
		rbmqWorker.HandleEmailFromQueue(msg.Body)
		time.Sleep(2 * time.Second)

		logger.Infof("%s: email sent successfully", queueName)
	}
}

func RunWorker(rbmqWorker interfaces.IRbmqWorker) {
	// Open a dummy channel to hold this RunWorker without exit
	forever := make(chan bool)

	var rbChanns = []func(interfaces.IRbmqWorker){HandleTestQueue, HandleLogQueue, HandleEmailQueue}

	for _, handler := range rbChanns {
		go handler(rbmqWorker)
	}

	<-forever
}
