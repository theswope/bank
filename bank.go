package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	viper.BindEnv("rulebaseurl")
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

	// Responde to rulebase with current settings
	if !updateRulebase() {
		log.Printf("Error: Couldn't update rulebase")
		return
	}

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

func updateRulebase() bool {
	url := viper.Get("rulebaseurl").(string)

	jsonObj := &rulebase{
		BankID:          viper.Get("name").(string),
		Topic:           viper.Get("requestTopic").(string),
		MinTerm:         viper.Get("minTerm").(int),
		MaxTerm:         viper.Get("maxTerm").(int),
		MinAmount:       viper.Get("minAmount").(int),
		MaxAmount:       viper.Get("maxAmount").(int),
		MinConsumerRate: viper.Get("minConsumerRate").(int),
		MaxConsumerRate: viper.Get("maxConsumerRate").(int),
	}

	jsonStr, err := json.Marshal(jsonObj)
	if err != nil {
		log.Printf("Error during marshal of rulebase object")
		return false
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Printf("response Status: %s", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	if resp.Status == "200 OK" {
		return true
	}

	return false
}

type rulebase struct {
	BankID          string `json:"bankId,omitempty"`
	Topic           string `json:"topic,omitempty"`
	MinTerm         int    `json:"minTerm,omitempty"`
	MaxTerm         int    `json:"maxTerm,omitempty"`
	MinAmount       int    `json:"minAmount,omitempty"`
	MaxAmount       int    `json:"maxAmount,omitempty"`
	MinConsumerRate int    `json:"minConsumerRate,omitempty"`
	MaxConsumerRate int    `json:"maxConsumerRate,omitempty"`
}
