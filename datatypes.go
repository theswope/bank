package main

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
)

// {"requestDeadline": 0, "amount": 10000, "term": 5, "requestId": "28a4f6b7-a0d4-4cb6-84dd-592c0b50a61d", "consumerRate": 6885}

type request struct {
	RequestDeadline int    `json:"requestDeadline,omitempty"`
	Amount          int    `json:"amount,omitempty"`
	Term            int    `json:"term,omitempty"`
	RequestID       string `json:"requestId,omitempty"`
	ConsumerRate    int    `json:"consumerRate,omitempty"`
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

	if !validTerm {
		log.Printf("Invalid term")
	}
	if !validAmount {
		log.Printf("Invalid amount")
	}
	if !validConsumerRate {
		log.Printf("Invalid consumer rate")
	}

	return validTerm && validAmount && validConsumerRate
}

func isValidTerm(term int) bool {
	minCorrect := false
	maxCorrect := false

	minTerm := viper.Get("minTerm").(string)
	minT, err := strconv.Atoi(minTerm)
	if err != nil {
		return false
	}

	if minT != 0 {
		minCorrect = term >= minT
	} else {
		minCorrect = true
	}

	maxTerm := viper.Get("maxTerm").(string)
	maxT, err := strconv.Atoi(maxTerm)
	if err != nil {
		return false
	}

	if maxT != 0 {
		maxCorrect = term <= maxT
	} else {
		maxCorrect = true
	}

	return minCorrect && maxCorrect
}

func isValidAmount(amount int) bool {
	minCorrect := false
	maxCorrect := false

	minAmount := viper.Get("minAmount").(string)
	minA, err := strconv.Atoi(minAmount)
	if err != nil {
		return false
	}

	if minA != 0 {
		minCorrect = amount >= minA
	} else {
		minCorrect = true
	}

	maxAmount := viper.Get("maxAmount").(string)
	maxA, err := strconv.Atoi(maxAmount)
	if err != nil {
		return false
	}

	if maxA != 0 {
		maxCorrect = amount <= maxA
	} else {
		maxCorrect = true
	}

	return minCorrect && maxCorrect
}

func isValidConsumerRate(consumerRate int) bool {
	minCorrect := false
	maxCorrect := false

	minConsumerRate := viper.Get("minConsumerRate").(string)
	minCR, err := strconv.Atoi(minConsumerRate)
	if err != nil {
		return false
	}

	if minCR != 0 {
		minCorrect = consumerRate >= minCR
	} else {
		minCorrect = true
	}

	maxConsumerRate := viper.Get("maxConsumerRate").(string)
	maxCR, err := strconv.Atoi(maxConsumerRate)
	if err != nil {
		return false
	}

	if maxCR != 0 {
		maxCorrect = consumerRate <= maxCR
	} else {
		maxCorrect = true
	}

	return minCorrect && maxCorrect
}
