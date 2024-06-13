package connection

import (
	"backend/server/env"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQ struct{}

var RabbitMqConnection *amqp.Connection
var RabbitMqChannel *amqp.Channel // save session connection

type RabbitMqStructure struct {
	Exchange      string
	Queue         string
	ExchangeRetry string
	QueueRetry    string
}

func (ref RabbitMQ) Connect() {
	var err error

	rabbit_url := env.GetRabbitUrl()
	RabbitMqConnection, err = amqp.Dial(rabbit_url)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", err)
	}

	RabbitMqChannel, err = RabbitMqConnection.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel", err)
	}

	log.Println("✅ RabbitMQ Connected")
}

func (ref RabbitMQ) Disconnect() {
	RabbitMqConnection.Close()
	RabbitMqChannel.Close()

	log.Println("✅ RabbitMQ Disconnect")
}

func (ref RabbitMQ) CreateModel(topic string) RabbitMqStructure {
	return RabbitMqStructure{
		Exchange:      topic + "-exchange",
		Queue:         topic + "-queue",
		ExchangeRetry: topic + "-exchange-retry",
		QueueRetry:    topic + "-queue-retry",
	}
}

func (ref RabbitMQ) CreateConsumer(topic string) (<-chan amqp.Delivery, *amqp.Connection, *amqp.Channel) {
	var err error

	value := ref.CreateModel(topic)

	rabbit_url := env.GetRabbitUrl()
	Connection, err := amqp.Dial(rabbit_url)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ", err)
	}

	Channel, err := Connection.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel", err)
	}
	Channel.Qos(1, 0, false)

	// Deklarasi exchange example-exchange
	err = Channel.ExchangeDeclare(
		value.Exchange, // Nama exchange
		"direct",       // Tipe exchange
		true,           // Durable
		false,          // Delete when unused
		false,          // Internal
		false,          // No-wait
		nil,            // Arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// Deklarasi queue example-queue
	_, err = Channel.QueueDeclare(
		value.Queue, // Nama queue
		true,        // Durable
		false,       // Delete when unused
		false,       // Exclusive
		false,       // No-wait
		amqp.Table{
			"x-dead-letter-exchange":    value.ExchangeRetry, // Menunjukkan ke exchange utama
			"x-dead-letter-routing-key": "/",                 // Routing key untuk pesan yang gagal
		}, // Arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Binding queue example-queue ke exchange example-exchange
	err = Channel.QueueBind(
		value.Queue,    // Nama queue
		"/",            // Routing key
		value.Exchange, // Nama exchange
		false,          // No-wait
		nil,            // Arguments
	)
	failOnError(err, "Failed to bind a queue")

	// --------------------------------------------------------

	// Deklarasi exchange retry
	err = Channel.ExchangeDeclare(
		value.ExchangeRetry, // Nama exchange retry
		"direct",            // Tipe exchange
		true,                // Durable
		false,               // Delete when unused
		false,               // Internal
		false,               // No-wait
		nil,                 // Arguments
	)
	failOnError(err, "Failed to declare a retry exchange")

	// Deklarasi queue retry
	_, err = Channel.QueueDeclare(
		value.QueueRetry, // Nama queue retry
		true,             // Durable
		false,            // Delete when unused
		false,            // Exclusive
		false,            // No-wait
		amqp.Table{
			"x-dead-letter-exchange":    value.Exchange, // Menunjukkan ke exchange utama
			"x-message-ttl":             1000 * 30,      // TTL (dalam milidetik)
			"x-dead-letter-routing-key": "/",            // Routing key untuk pesan yang gagal
		}, // Arguments
	)
	failOnError(err, "Failed to declare a retry queue")

	// Binding queue retry ke exchange retry
	err = Channel.QueueBind(
		value.QueueRetry,    // Nama queue retry
		"/",                 // Routing key
		value.ExchangeRetry, // Nama exchange retry
		false,               // No-wait
		nil,                 // Arguments
	)
	failOnError(err, "Failed to bind a retry queue")

	// --------------------------------------------------------
	log.Printf("✅ Consumer \"%s\" is Running\n", topic)
	msgs, err := Channel.Consume(
		value.Queue,
		"/",
		false, // set auto-acknowledge to false
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to create a consumer")

	return msgs, Connection, Channel
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
