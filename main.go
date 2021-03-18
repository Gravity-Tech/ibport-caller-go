package main

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/user/caller/token"
	"github.com/user/caller/ibport"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mr-tron/base58"
)


func main() {

	nodeUrl := "https://bsc-dataseed.binance.org"

	senderPrivKeyStr := "0x2d523a2573c8f40eaba9e0cf1d7ed1d972872522ef1f57c47be23e0da00c376b"

	tokenAddress := "0xc4b6F32B84657E9f6a73fE119f0967bE5bA8CF05"
	ibportAddress := "0x8c0e11a6E692d02f71598AB5050083ED691Eb760"

	receiverAddress := "3PG4DtkmcZjPGRpXujbD44ZKydo1D9Y6r2N"
	amount := int64(1)


	amountBig := big.NewInt(amount)

	senderPrivKey, err := hexutil.Decode(senderPrivKeyStr)
	if err != nil {
		return
	}

	sysCtx := context.Background()

	ethClient, err := ethclient.DialContext(sysCtx, nodeUrl)

	ethPrivKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: secp256k1.S256(),
		},
		D: new(big.Int),
	}
	ethPrivKey.D.SetBytes(senderPrivKey)
	ethPrivKey.PublicKey.X, ethPrivKey.PublicKey.Y = ethPrivKey.PublicKey.Curve.ScalarBaseMult(senderPrivKey)

	transactOpt := bind.NewKeyedTransactor(ethPrivKey)

	tokenAddressBytes, err := hexutil.Decode(tokenAddress)
	if err != nil {
		return
	}

	token, err := token.NewToken(common.BytesToAddress(tokenAddressBytes), ethClient)
	if err != nil {
		return
	}

	_, err = token.Approve(transactOpt, common.BytesToAddress(tokenAddressBytes), amountBig)
	if err != nil {
		return
	}

	ibportAddressBytes, err := hexutil.Decode(ibportAddress)
	if err != nil {
		return
	}

	ibport, err := ibport.NewIBPort(common.BytesToAddress(ibportAddressBytes), ethClient)
	if err != nil {
		return
	}

	receiverAddressBytes, err := base58.Decode(receiverAddress)
	if err != nil {
		return
	}

	var receiverAddressHex [32]byte
    copy(receiverAddressHex[:], receiverAddressBytes)

	_, err = ibport.CreateTransferUnwrapRequest(transactOpt, amountBig, receiverAddressHex)
	if err != nil {
		return
	}

	return

}
