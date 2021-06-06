package queues

import (
	"fmt"
	config "github.com/a-castellano/music-manager-config-reader/config_reader"
	"github.com/a-castellano/music-manager-musicbrainz-wrapper/jobs"
	"github.com/streadway/amqp"
	"net/http"
	"strconv"
)

func StartJobManagement(config config.Config, client http.Client) error {

	connection_string := "amqp://" + config.Server.User + ":" + config.Server.Password + "@" + config.Server.Host + ":" + strconv.Itoa(config.Server.Port) + "/"
	conn, err := amqp.Dial(connection_string)

	if err != nil {
		return fmt.Errorf("Failed to stablish connection with RabbitMQ: %w", err)
	}
	defer conn.Close()

	incoming_ch, err := conn.Channel()
	defer incoming_ch.Close()

	if err != nil {
		return fmt.Errorf("Failed to open incoming channel: %w", err)
	}

	outgoing_ch, err := conn.Channel()
	defer outgoing_ch.Close()

	if err != nil {
		return fmt.Errorf("Failed to open outgoing channel: %w", err)
	}

	incoming_q, err := incoming_ch.QueueDeclare(
		config.Incoming.Name,
		true,  // Durable
		false, // DeleteWhenUnused
		false, // Exclusive
		false, // NoWait
		nil,   // arguments
	)

	if err != nil {
		return fmt.Errorf("Failed to declare incoming queue: %w", err)
	}

	outgoing_q, err := outgoing_ch.QueueDeclare(
		config.Outgoing.Name,
		true,  // Durable
		false, // DeleteWhenUnused
		false, // Exclusive
		false, // NoWait
		nil,   // arguments
	)

	if err != nil {
		return fmt.Errorf("Failed to declare outgoing queue: %w", err)
	}

	err = incoming_ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		return fmt.Errorf("Failed to set incoming QoS: %w", err)
	}

	jobsToProcess, err := incoming_ch.Consume(
		incoming_q.Name,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		return fmt.Errorf("Failed to register a consumer: %w", err)
	}

	processJobs := make(chan bool)

	go func() {
		for job := range jobsToProcess {

			die, jobResult, _ := jobs.ProcessJob(job.Body, config.Origin, client)

			if die {
				job.Ack(false)
				processJobs <- false
				return
			}
			err = outgoing_ch.Publish(
				"",              // exchange
				outgoing_q.Name, // routing key
				false,           // mandatory
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         jobResult,
				})
			if err != nil {
				fmt.Println(err)
				return
			}

			job.Ack(false)
		}
		return
	}()

	<-processJobs

	return nil
}
