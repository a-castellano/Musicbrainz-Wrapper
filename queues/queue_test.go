// +build integration_tests

package queues

import (
	"bytes"
	commontypes "github.com/a-castellano/music-manager-common-types/types"
	config "github.com/a-castellano/music-manager-config-reader/config_reader"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

type RoundTripperMock struct {
	Response *http.Response
	RespErr  error
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (rtm *RoundTripperMock) RoundTrip(*http.Request) (*http.Response, error) {
	return rtm.Response, rtm.RespErr
}

func TestSendDie(t *testing.T) {

	var queueConfig config.Config

	queueConfig.Server.Host = "rabbitmq"
	queueConfig.Server.Port = 5672
	queueConfig.Server.User = "guest"
	queueConfig.Server.Password = "guest"

	queueConfig.Incoming.Name = "incoming"
	queueConfig.Outgoing.Name = "outgoing"

	var job commontypes.Job

	job.ID = 0
	job.Status = true
	job.Finished = false
	job.Type = commontypes.Die

	encodedJob, _ := commontypes.EncodeJob(job)

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{
	"error": "",
	"iTotalRecords": 0,
	"iTotalDisplayRecords": 0,
	"sEcho": 0,
	"aaData": [
		]
}
	`))}}}

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel in TestSendDie")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"incoming", // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue in TestSendDie")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         encodedJob,
		})

	jobManagementError := StartJobManagement(queueConfig, client)
	if jobManagementError != nil {
		t.Errorf("StartJobManagement should return no errors when die is processed.")
	}

}

func TestSendNoArtistsFound(t *testing.T) {

	var queueConfig config.Config

	queueConfig.Server.Host = "rabbitmq"
	queueConfig.Server.Port = 5672
	queueConfig.Server.User = "guest"
	queueConfig.Server.Password = "guest"

	queueConfig.Incoming.Name = "incoming"
	queueConfig.Outgoing.Name = "outgoing"

	var infoRetrieval commontypes.InfoRetrieval
	var job commontypes.Job

	infoRetrieval.Type = commontypes.ArtistName
	infoRetrieval.Artist = "AnyArtist"

	retrievalData, _ := commontypes.EncodeInfoRetrieval(infoRetrieval)

	job.Data = retrievalData
	job.ID = 0
	job.Status = true
	job.Finished = false
	job.Type = 1

	encodedJob, _ := commontypes.EncodeJob(job)

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"created":"2021-01-04T18:29:36.962Z","count":0,"offset":0,"artists":[]}
	`))}}}

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ in TestSendNoArtistsFound.")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel in TestSendNoArtistsFound.")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"incoming", // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare incoming queue in TestSendNoArtistsFound.")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         encodedJob,
		})

	failOnError(err, "Failed to send first job in TestSendNoArtistsFound.")
	var dieJob commontypes.Job

	dieJob.ID = 0
	dieJob.Status = true
	dieJob.Finished = false
	dieJob.Type = commontypes.Die

	encodedDieJob, _ := commontypes.EncodeJob(dieJob)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         encodedDieJob,
		})

	failOnError(err, "Failed to send die job in TestSendNoArtistsFound.")

	jobManagementError := StartJobManagement(queueConfig, client)

	if jobManagementError != nil {
		t.Errorf("StartJobManagement should return no errors when die is processed.")
	}

	outgoingCh, err := conn.Channel()
	failOnError(err, "Failed to open outgoing channel in TestSendNoArtistsFound.")
	defer outgoingCh.Close()

	msgs, err := outgoingCh.Consume(
		"outgoing", // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	forever := make(chan bool)
	var receivedData []byte

	go func() {
		for d := range msgs {

			receivedData = d.Body
			d.Ack(false)
			forever <- false
		}
	}()

	<-forever
	decodedJob, decodedJobErr := commontypes.DecodeJob(receivedData)

	if decodedJob.Type != commontypes.ArtistInfoRetrieval {
		t.Errorf("Decoded job type should be ArtistInfoRetrieval in TestSendNoArtistsFound. It's %d.", decodedJob.Type)
	}
	if decodedJobErr != nil {
		t.Errorf("DecodeJob should return no errors.")
	}

	if decodedJob.Error != "Artist retrieval failed: No artist was found." {
		t.Errorf("DecodeJob error should be 'Artist retrieval failed: No artist was found.', not '%s'.", decodedJob.Error)
	}
}

func TestFailedConfig(t *testing.T) {

	var queueConfig config.Config

	queueConfig.Server.Host = "127.0.0.1"
	queueConfig.Server.Port = 5672
	queueConfig.Server.User = "guest"
	queueConfig.Server.Password = "nopassword"

	queueConfig.Incoming.Name = "incoming"
	queueConfig.Outgoing.Name = "outgoing"

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"created":"2021-01-04T18:29:36.962Z","count":0,"offset":0,"artists":[]}
	`))}}}

	jobManagementError := StartJobManagement(queueConfig, client)
	if jobManagementError == nil {
		t.Errorf("StartJobManagement should return an error when credentials are invalid.")
	}

}

func TestSendArtistsFound(t *testing.T) {

	var queueConfig config.Config

	queueConfig.Server.Host = "rabbitmq"
	queueConfig.Server.Port = 5672
	queueConfig.Server.User = "guest"
	queueConfig.Server.Password = "guest"

	queueConfig.Incoming.Name = "incoming"
	queueConfig.Outgoing.Name = "outgoing"

	var infoRetrieval commontypes.InfoRetrieval
	var job commontypes.Job

	infoRetrieval.Type = commontypes.ArtistName
	infoRetrieval.Artist = "Manowar"

	retrievalData, _ := commontypes.EncodeInfoRetrieval(infoRetrieval)

	job.Data = retrievalData
	job.ID = 0
	job.Status = true
	job.Finished = false
	job.Type = 1

	encodedJob, _ := commontypes.EncodeJob(job)

	client := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`
{"created":"2021-01-04T18:59:46.493Z","count":6,"offset":0,"artists":[{"id":"00eeed6b-5897-4359-8347-b8cd28375331","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":100,"name":"Manowar","sort-name":"Manowar","country":"US","area":{"id":"489ce91b-6658-3307-9877-795b68554c98","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"United States","sort-name":"United States","life-span":{"ended":null}},"begin-area":{"id":"27b29a82-d6e5-490b-80a6-965c9f77c741","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Auburn","sort-name":"Auburn","life-span":{"ended":null}},"isnis":["0000000115122690"],"life-span":{"begin":"1980","ended":null},"aliases":[{"sort-name":"Man o' War","name":"Man o' War","locale":null,"type":null,"primary":null,"begin-date":null,"end-date":null}],"tags":[{"count":1,"name":"honor"},{"count":1,"name":"war"},{"count":1,"name":"norse mythology"},{"count":1,"name":"hard rock"},{"count":7,"name":"heavy metal"},{"count":1,"name":"glory"},{"count":1,"name":"symphonic metal"},{"count":1,"name":"power metal"}]},{"id":"2f3d8c8b-cd65-49e2-9a02-e89764411f88","score":39,"name":"Womanowar","sort-name":"Womanowar","life-span":{"ended":null}},{"id":"44f82b18-151f-4e98-a6f1-33b18ad6a46f","type":"Person","type-id":"b6e035f4-3ce9-331c-97df-83397230b0df","score":39,"gender-id":"36d3d30a-839d-3eda-8cb3-29be4384e4a9","name":"Josef Manowarda","sort-name":"Manowarda, Josef","gender":"male","country":"DE","area":{"id":"85752fda-13c4-31a3-bee5-0e5cb1f51dad","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"Germany","sort-name":"Germany","life-span":{"ended":null}},"life-span":{"begin":"1890-07-03","end":"1942-12-24","ended":true}},{"id":"4c96a93c-85ab-4aaa-83f2-bd1ad759e63b","type":"Person","type-id":"b6e035f4-3ce9-331c-97df-83397230b0df","score":38,"gender-id":"36d3d30a-839d-3eda-8cb3-29be4384e4a9","name":"Eric Adams","sort-name":"Adams, Eric","gender":"male","country":"US","area":{"id":"489ce91b-6658-3307-9877-795b68554c98","type":"Country","type-id":"06dd0ae4-8c74-30bb-b43d-95dcedf961de","name":"United States","sort-name":"United States","life-span":{"ended":null}},"begin-area":{"id":"27b29a82-d6e5-490b-80a6-965c9f77c741","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Auburn","sort-name":"Auburn","life-span":{"ended":null}},"disambiguation":"Louis Marullo, Manowar","life-span":{"begin":"1954-07-12","ended":null},"aliases":[{"sort-name":"Marullo, Louis","type-id":"d4dcd0c0-b341-3612-a332-c0ce797b25cf","name":"Louis Marullo","locale":null,"type":"Legal name","primary":null,"begin-date":null,"end-date":null}]},{"id":"963474e4-c3af-4d58-acbd-36a4b13f001c","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":36,"name":"Womenowar","sort-name":"Womenowar","area":{"id":"226c4dca-ef2a-4d4b-ba25-4118d116557a","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Birmingham","sort-name":"Birmingham","life-span":{"ended":null}},"begin-area":{"id":"226c4dca-ef2a-4d4b-ba25-4118d116557a","type":"City","type-id":"6fd8f29a-3d0a-32fc-980d-ea697b69da78","name":"Birmingham","sort-name":"Birmingham","life-span":{"ended":null}},"disambiguation":"UK female Manowar tribute","life-span":{"begin":"2018","ended":null}},{"id":"fb1125a1-47b9-4c38-9944-4374ae785ed1","type":"Group","type-id":"e431f5f6-b5d2-343d-8b36-72607fffb74b","score":36,"name":"Hanowar","sort-name":"Hanowar","disambiguation":"Bad UK Manowar tribute","life-span":{"ended":null}}]}
	`))}}}

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ in TestSendNoArtistsFound.")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel in TestSendNoArtistsFound.")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"incoming", // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare incoming queue in TestSendNoArtistsFound.")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         encodedJob,
		})

	failOnError(err, "Failed to send first job in TestSendNoArtistsFound.")
	var dieJob commontypes.Job

	dieJob.ID = 0
	dieJob.Status = true
	dieJob.Finished = false
	dieJob.Type = commontypes.Die

	encodedDieJob, _ := commontypes.EncodeJob(dieJob)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         encodedDieJob,
		})

	failOnError(err, "Failed to send die job in TestSendNoArtistsFound.")

	jobManagementError := StartJobManagement(queueConfig, client)

	if jobManagementError != nil {
		t.Errorf("StartJobManagement should return no errors when die is processed.")
	}

	outgoingCh, err := conn.Channel()
	failOnError(err, "Failed to open outgoing channel in TestSendNoArtistsFound.")
	defer outgoingCh.Close()

	msgs, err := outgoingCh.Consume(
		"outgoing", // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	forever := make(chan bool)
	var receivedData []byte

	go func() {
		for d := range msgs {

			receivedData = d.Body
			d.Ack(false)
			forever <- false
		}
	}()

	<-forever
	decodedJob, decodedJobErr := commontypes.DecodeJob(receivedData)

	if decodedJob.Type != commontypes.ArtistInfoRetrieval {
		t.Errorf("Decoded job type should be ArtistInfoRetrieval in TestSendNoArtistsFound. It's %d.", decodedJob.Type)
	}
	if decodedJobErr != nil {
		t.Errorf("DecodeJob should return no errors.")
	}

	if decodedJob.Error != "" {
		t.Errorf("DecodeJob error should be empty, found '%s'.", decodedJob.Error)
	}

	if len(decodedJob.Result) == 0 {
		t.Errorf("DecodeJob  should have result.")
	}

	retrievedData, _ := commontypes.DecodeArtistInfo(decodedJob.Result)

	if retrievedData.Data.Name != "Manowar" {
		t.Errorf("Retrieved data name should be Manowar, not %s.", retrievedData.Data.Name)
	}

	if retrievedData.Data.Country != "US" {
		t.Errorf("Retrieved data country should be Sweden, not %s.", retrievedData.Data.Country)
	}

	if len(retrievedData.ExtraData) != 0 {
		t.Errorf("Retrieved extradata should contain 1 entry, not %d.", len(retrievedData.ExtraData))
	}

}
