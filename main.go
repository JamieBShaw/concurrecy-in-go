package main

import (
	"errors"
	"fmt"
	"log"
)

type Owner int

const (
	Bob Owner = iota
	Alice
)

type Wallet struct {
	amount int
}

func NewMultiSigWallet(userOne, userTwo Wallet) (*MultiSigWallet, error) {
	initialFunding, err := calcInitialFunding(userOne, userTwo)
	if err != nil {
		return nil, err
	}

	return &MultiSigWallet{
		bob:     userOne,
		alice:   userTwo,
		funding: initialFunding,
		transactionChan: TransactionChan{
			txs:     []Txs{},
			channel: make(chan Txs),
		},
	}, nil
}

func calcInitialFunding(userOne, userTwo Wallet) (int, error) {
	userOneAmount := userOne.amount
	userTwoAmount := userTwo.amount
	if userOneAmount < 0 || userTwoAmount < 0 {
		return 0, errors.New("invalid amounts")
	}
	return userOne.amount + userTwo.amount, nil
}

type MultiSigWallet struct {
	bob             Wallet
	alice           Wallet
	funding         int
	transactionChan TransactionChan
}

func (msw *MultiSigWallet) aliceSend(amount int) {
	msw.transactionChan.channel <- Txs{
		sender: Alice,
		sendTo: Bob,
		amount: amount,
	}
}

func (msw *MultiSigWallet) bobSend(amount int) {
	msw.transactionChan.channel <- Txs{
		sender: Bob,
		sendTo: Alice,
		amount: amount,
	}
}

type TransactionChan struct {
	channel chan Txs
	txs     []Txs
}

type Txs struct {
	sender Owner
	sendTo Owner
	amount int
}

func (w *Wallet) spend(amount int) {
	w.amount -= amount
}

func (w *Wallet) add(amount int) {
	w.amount += amount
}

func main() {
	strC, err := readCSV("pipeline.csv")
	if err != nil {
		log.Fatalf("Could not read csv %v", err)
	}

	for val := range sanitize(titleCase(strC)) {
		fmt.Printf("finished: %v\n", val)
	}

}
