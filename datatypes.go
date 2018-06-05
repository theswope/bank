package main

import (
	"log"

	"github.com/spf13/viper"
)

// {"requestDeadline": 0, "amount": 10000, "term": 5, "requestId": "28a4f6b7-a0d4-4cb6-84dd-592c0b50a61d", "consumerRate": 6885}

type request struct {
	RequestDeadline int    `json:"requestDeadline,omitempty"`
	Amount          int    `json:"amount,omitempty"`
	Term            int    `json:"term,omitempty"`
	RequestID       string `json:"requestId,omitempty"`
	ConsumerRate    string `json:"consumerRate,omitempty"`
}

type response struct {
	BankName  string  `json:"bankName,omitempty"`
	RequestID string  `json:"requestId,omitempty"`
	QuoteRate float64 `json:"quoteRate,omitempty"`
}

func (r *request) processRequest() *response {
	if !r.isValidRequest() {
		log.Printf("Invalid request: %v", r)

		return nil
	}

	quoteRate := getRandomQuoteRate()

	log.Printf("New quote rate: %f\n", quoteRate)

	return &response{
		BankName:  viper.Get("name").(string),
		RequestID: r.RequestID,
		QuoteRate: quoteRate,
	}
}

func (r *request) isValidRequest() bool {
	// Compare all features to the parameters of the bank
	validTerm := isValidTerm(r.Term)
	validAmount := isValidAmount(r.Amount)
	validConsumerRate := isValidConsumerRate(r.ConsumerRate)

	return validTerm && validAmount && validConsumerRate
}

func isValidTerm(term int) bool {
	minCorrect := false
	maxCorrect := false

	minTerm := viper.Get("minTerm")
	if minTerm != 0 {
		minCorrect = term >= minTerm.(int)
	} else {
		minCorrect = true
	}

	maxTerm := viper.Get("maxTerm")
	if maxTerm != 0 {
		maxCorrect = term <= maxTerm.(int)
	} else {
		maxCorrect = true
	}

	return minCorrect && maxCorrect
}

func isValidAmount(amount int) bool {
	minCorrect := false
	maxCorrect := false

	minAmount := viper.Get("minAmount").(int)
	if minAmount != 0 {
		minCorrect = amount >= minAmount
	} else {
		minCorrect = true
	}

	maxAmount := viper.Get("maxAmount").(int)
	if maxAmount != 0 {
		maxCorrect = amount <= maxAmount
	} else {
		maxCorrect = true
	}

	return minCorrect && maxCorrect
}

func isValidConsumerRate(consumerRate string) bool {
	minCorrect := false
	maxCorrect := false

	minConsumerRate := viper.Get("minConsumerRate").(string)
	if minConsumerRate != "" {
		minCorrect = consumerRate >= minConsumerRate
	} else {
		minCorrect = true
	}

	maxConsumerRate := viper.Get("maxConsumerRate").(string)
	if maxConsumerRate != "" {
		maxCorrect = consumerRate <= maxConsumerRate
	} else {
		maxCorrect = true
	}

	return minCorrect && maxCorrect
}
