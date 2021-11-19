package main

import "C"
import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	simulator2 "github.com/Ethworks/Waffle/simulator"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

type TransactionRequest struct {
	To    *string `json:"to"`
	From  *string `json:"from"`
	Nonce *string `json:"nonce"`

	GasLimit *string `json:"gasLimit"`
	GasPrice *string `json:"gasPrice"`

	Data    *string `json:"data"`
	Value   *string `json:"value"`
	ChainId *uint64 `json:"chainId"`
}

func main() {}

var simulator, _ = simulator2.NewSimulator()

//export getBlockNumber
func getBlockNumber() *C.char {
	bn := simulator.GetLatestBlockNumber()
	return C.CString(bn.String())
}

//export getBalance
func getBalance(account *C.char) *C.char {
	bal, err := simulator.Backend.BalanceAt(context.Background(), common.HexToAddress(C.GoString(account)), nil)
	if err != nil {
		log.Fatal(err)
	}

	return C.CString(bal.String())
}

//export getTransactionCount
func getTransactionCount(account *C.char) C.int {
  count, err :=  simulator.Backend.NonceAt(context.Background(), common.HexToAddress(C.GoString(account)), nil)
  if err != nil {
    log.Fatal(err)
  }

  return C.int(count)
}

//export call
func call(msgJson *C.char) *C.char {
	var msg TransactionRequest

	fmt.Println(C.GoString(msgJson))

	err := json.Unmarshal([]byte(C.GoString(msgJson)), &msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(msg)

	var callMsg ethereum.CallMsg

	if msg.From != nil {
		callMsg.From = common.HexToAddress(*msg.From)
	}
	if msg.To != nil {
		temp := common.HexToAddress(*msg.To)
		callMsg.To = &temp
	}
	if msg.GasLimit != nil {
		value, err := strconv.ParseUint(*msg.GasLimit, 16, 64)
		if err != nil {
			log.Fatal(err)
		}

		callMsg.Gas = value
	}
	if msg.GasPrice != nil {
		callMsg.GasPrice = big.NewInt(0)
		callMsg.GasPrice.SetString(*msg.GasPrice, 16)
	}
	if msg.Data != nil {
		data, err := hex.DecodeString((*msg.Data)[2:])
		if err != nil {
			log.Fatal(err)
		}

		callMsg.Data = data
	}

	res, err := simulator.Backend.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatal(err)
	}

	return C.CString(hex.EncodeToString(res))
}

//export sendTransaction
func sendTransaction(txData *C.char) {

	bytes, err := hex.DecodeString(C.GoString(txData)[2:])
	if err != nil {
		log.Fatal(err)
	}

	tx := &types.Transaction{}
	err = tx.UnmarshalBinary(bytes)
	if err != nil {
		log.Fatal(err)
	}

	err = simulator.Backend.SimulatedBackend.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatal(err)
	}

	simulator.Backend.Commit()

	receipt, err := simulator.Backend.SimulatedBackend.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tx.Hash().String())
	fmt.Println(receipt.ContractAddress.String())
}

//export cgoCurrentMillis
func cgoCurrentMillis() C.long {
	return C.long(time.Now().Unix())
}

//export cgoSeed
func cgoSeed(m C.long) {
	rand.Seed(int64(m))
}

//export cgoRandom
func cgoRandom(m C.int) C.int {
	return C.int(rand.Intn(int(m)))
}