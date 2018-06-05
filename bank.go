package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	// Setup
	viper.SetConfigName("default")
	viper.AddConfigPath(".")

	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	// Environment variable setup
	viper.SetEnvPrefix("bank")
	viper.BindEnv("broker")
	viper.BindEnv("port")
	viper.BindEnv("user")
	viper.BindEnv("pass")
	viper.BindEnv("virthost")
	viper.BindEnv("name")
	viper.BindEnv("requestTopic")
	viper.BindEnv("responseTopic")
	viper.BindEnv("autorespond")
	viper.BindEnv("minTerm")
	viper.BindEnv("maxTerm")
	viper.BindEnv("minAmount")
	viper.BindEnv("maxAmount")
	viper.BindEnv("minConsumerRate")
	viper.BindEnv("maxConsumerRate")

	// Create new connection
	aCon := amqpConnection{}
	aCon.connectToBroker()

	// Connect to channel
	aCon.connectToChannel()

	// Declare a queue
	queue := viper.Get("requestTopic").(string)
	aCon.declareQueue(queue)

	// Consume messages from queue
	aCon.consumeFromQueue()
}
