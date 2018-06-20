package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"time"

	e "github.com/nombiezinja/chstub/entities"
	u "github.com/nombiezinja/chstub/utils"
	"github.com/streadway/amqp"
)

func parseNum() int {
	total, err := strconv.Atoi(os.Args[1])
	u.FailOnError(err, "Failed to parse commandline arg for num")
	return total
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	u.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	u.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"stuffs", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	u.FailOnError(err, "Failed to declare a queue")

	num := parseNum()

	for i := 0; i < num; i++ {
		p := genPayload(i)
		pJson, err := json.Marshal(p)
		u.FailOnError(err, "Marshal json failed")

		body := string(pJson)

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})

		u.FailOnError(err, "Failed to publish a message")
	}
}

var colours = [5]string{"yellow", "pink", "aquamarine", "blue", "purple"}

func genPayload(i int) e.Payload {

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	p := e.Payload{
		ID:     i,
		Colour: colours[r.Intn(len(colours))],
	}

	return p
}
