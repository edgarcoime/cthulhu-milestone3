package handlers

import (
	"log"

	rabbitmq "github.com/wagslane/go-rabbitmq"
)

// EventsHandler handles messages from the cthulhu.events fanout exchange.
func EventsHandler() rabbitmq.Handler {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		log.Printf("[filemanager] consumed message: %s\n", string(d.Body))
		return rabbitmq.Ack
	}
}

// DiagnoseHandler handles messages from the diagnose.events fanout exchange.
// It can publish to other exchanges (e.g., filemanager.requests) based on diagnose results.
func DiagnoseHandler(publisher *rabbitmq.Publisher) rabbitmq.Handler {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		log.Printf("[filemanager][diagnose] consumed message: %s\n", string(d.Body))

		// Example: if diagnose message indicates download needed, trigger it
		// In a real scenario, you'd parse the diagnose message and conditionally publish
		downloadPayload := []byte(`{"triggered_by": "diagnose", "message": "` + string(d.Body) + `"}`)

		err := publisher.Publish(
			downloadPayload,
			[]string{"filemanager.download_request"},
			rabbitmq.WithPublishOptionsExchange("filemanager.requests"),
		)
		if err != nil {
			log.Printf("[filemanager][diagnose] failed to publish download request: %v\n", err)
			return rabbitmq.NackRequeue
		}

		log.Printf("[filemanager][diagnose] published download_request to filemanager.requests\n")
		return rabbitmq.Ack
	}
}

// RequestsHandler handles messages from the filemanager.requests direct exchange.
// It routes based on routing keys to handle different request types.
func RequestsHandler() rabbitmq.Handler {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		// In a direct exchange, d.RoutingKey tells you which action was requested
		switch d.RoutingKey {
		case "filemanager.get_presigned_url":
			log.Printf("[filemanager] handling get_presigned_url request: %s\n", string(d.Body))
			// TODO: Parse JSON, call S3 client, return presigned URL
			return rabbitmq.Ack

		case "filemanager.download_request":
			log.Printf("[filemanager] handling download_request: %s\n", string(d.Body))
			// TODO: Parse JSON, handle download logic
			return rabbitmq.Ack

		default:
			log.Printf("[filemanager] unknown routing key: %s\n", d.RoutingKey)
			return rabbitmq.NackDiscard // Reject unknown routing keys
		}
	}
}
