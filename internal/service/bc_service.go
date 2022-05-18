package service

import (
	"github.com/alexmolinanasaev/exterr"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type HyperLedService struct {
	contract *gateway.Contract
}

func NewEthService(contract *gateway.Contract) HyperLedApi {
	return &HyperLedService{
		contract: contract,
	}
}

func (e *HyperLedService) AddWallet() exterr.ErrExtender {
	return nil
}

func (e *HyperLedService) BuyTokens(tokenAmount int, buyerAddress string) ([]byte, exterr.ErrExtender) {

	// fromAddress := common.HexToAddress(buyerAddress)
	// nonce, err := e.client.PendingNonceAt(context.Background(), fromAddress)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// value := big.NewInt(0) // in wei (0 eth)
	// gasPrice, err := e.client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// toAddress := common.HexToAddress(buyerAddress)
	// tokenAddress := common.HexToAddress("0x34200e2980E89ab2AAe8A508932Ef9025E2ea150")

	// transferFnSignature := []byte("transfer(address,uint256)")

	// hash := sha3.NewLegacyKeccak256()
	// hash.Write(transferFnSignature)
	// methodID := hash.Sum(nil)[:4]
	// // fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	// paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	// // fmt.Println(hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d

	// amount := new(big.Int)
	// amount.SetString(strconv.Itoa(tokenAmount), 10) // 1 token
	// paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	// // fmt.Println(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	// var data []byte
	// data = append(data, methodID...)
	// data = append(data, paddedAddress...)
	// data = append(data, paddedAmount...)

	// gasLimit, err := e.client.EstimateGas(context.Background(), ethereum.CallMsg{
	// 	To:   &toAddress,
	// 	Data: data,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // fmt.Println(gasLimit) // 23256
	// tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	// res, err := json.Marshal(tx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return nil, nil
}
