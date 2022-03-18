package mailingservice

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	cachingservice "github.com/ericklima-ca/bago/services/caching_service"
	amqp "github.com/rabbitmq/amqp091-go"
)

func sendMessageToQueue(msg []byte) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	checkError(err, "connection failed")
	defer conn.Close()

	ch, err := conn.Channel()
	checkError(err, "channel communication failed")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"mail", // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
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
			Body:         msg,
		},
	)
}

func checkError(err error, msg string) {
	if err != nil {
		log.Fatalf("%v: %v", err.Error(), msg)
	}
}

func SendConfirmationEmail(id uint, name, email, url string) {
	strId := strconv.Itoa(int(id))
	token := cachingservice.GetToken("signup", id)
	msg, _ := json.Marshal(map[string]interface{}{
		"to":      email,
		"subject": "Email confirmation!",
		"body": `<p>` + name + `Please confirm your email by clicking the link below:</p>
				<p>https://` + url + `/api/auth/verify/signup/` + strId + `/` + token + `</p>`,
	})
	sendMessageToQueue(msg)
}

func SendRecoveryEmail(id uint, name, email, url string) {
	strId := strconv.Itoa(int(id))
	token := cachingservice.GetToken("recovery", id)
	msg, _ := json.Marshal(map[string]interface{}{
		"to":      email,
		"subject": "Password recovery confirmation!",
		"body": `<p>` + name + `Please confirm your new password by clicking the link below:</p>
				<p>https://` + url + `/api/auth/verify/recovery/` + strId + `/` + token + `</p>`,
	})
	sendMessageToQueue(msg)
}
