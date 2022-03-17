package mailingservice

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	cachingservice "github.com/ericklima-ca/bago/services/caching_service"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SendMessageToQueue(id uint, name, email string) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	checkError(err, "connection failed")
	defer conn.Close()

	ch, err := conn.Channel()
	checkError(err, "channel communication failed")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"mail",       // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	checkError(err, "declaring queue failed")

	ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         formatEmail(id, name, email),
		},
	)
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%v: %v", err.Error(), msg)
	}
}

func formatEmail(id uint, n, e string) []byte {
	strId := strconv.Itoa(int(id))
	token := cachingservice.GetToken("signup", id)
	msg, _ := json.Marshal(map[string]interface{}{
		"to": e,
		"subject": "Email confirmation!",
		"body": `<p>`+ n + `Please confirm your email by clicking the link below:</p>
				<p>https://example.com/api/auth/verify/` + strId + `/` +token + `</p>`,
	})
	return msg
}