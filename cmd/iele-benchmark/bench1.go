package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"path/filepath"
	"testing"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"
	cryptohook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-crypto"
	compiler "github.com/ElrondNetwork/elrond-vm/iele/compiler"
	eiele "github.com/ElrondNetwork/elrond-vm/iele/elrond/node/endpoint"
)

func benchmarkManyErc20SimpleTransfers(b *testing.B, nrTransfers int) {

	contractPathFilePath := filepath.Join(elrondTestRoot, "iele-examples/erc20_elrond.iele")

	compiledBytes := compiler.AssembleIeleCode(contractPathFilePath)
	decoded, err := hex.DecodeString(string(compiledBytes))
	if err != nil {
		panic(err)
	}

	world := worldhook.NewMock()

	contractAddrHex := "c0879ac700000000000000000000000000000000000000000000000000000000"
	account1AddrHex := "acc1000000000000000000000000000000000000000000000000000000000000"
	account2AddrHex := "acc2000000000000000000000000000000000000000000000000000000000000"

	contractAddr, _ := hex.DecodeString(contractAddrHex)
	account1Addr, _ := hex.DecodeString(account1AddrHex)
	account2Addr, _ := hex.DecodeString(account2AddrHex)

	constractStorage := make(map[string][]byte)
	constractStorage[storageKey("01"+account1AddrHex)] = big.NewInt(2000000000).Bytes()
	constractStorage[storageKey("00")] = big.NewInt(2000000000).Bytes() // total supply

	world.AcctMap.PutAccount(&worldhook.Account{
		Exists:  true,
		Address: contractAddr,
		Nonce:   0,
		Balance: big.NewInt(0),
		Storage: constractStorage,
		Code:    decoded,
	})

	world.AcctMap.PutAccount(&worldhook.Account{
		Exists:  true,
		Address: account1Addr,
		Nonce:   0,
		Balance: hexToBigInt("e8d4a51000"),
		Storage: make(map[string][]byte),
		Code:    []byte{},
	})

	world.AcctMap.PutAccount(&worldhook.Account{
		Exists:  true,
		Address: account2Addr,
		Nonce:   0,
		Balance: hexToBigInt("e8d4a51000"),
		Storage: make(map[string][]byte),
		Code:    []byte{},
	})

	account2AsArg := hexToArgument(account2AddrHex)

	// create the VM and allocate some memory
	vm := eiele.NewElrondIeleVM(
		eiele.TestVMType, eiele.ElrondDefault,
		world, cryptohook.KryptoHookMockInstance)

	if b != nil { // nil when debugging
		b.ResetTimer()
	}

	for benchMarkRepeat := 0; benchMarkRepeat < 1; benchMarkRepeat++ {
		for txi := 0; txi < nrTransfers; txi++ {
			input := &vmi.ContractCallInput{
				RecipientAddr: contractAddr,
				Function:      "transfer",
				VMInput: vmi.VMInput{
					CallerAddr: account1Addr,
					Arguments: [][]byte{
						account2AsArg,
						[]byte{1},
					},
					CallValue:   big.NewInt(0),
					GasPrice:    1,
					GasProvided: 100000,
				},
			}

			output, err := vm.RunSmartContractCall(input)
			if err != nil {
				panic(err)
			}

			if output.ReturnCode != vmi.Ok {
				panic(fmt.Sprintf("returned non-zero code: %d", output.ReturnCode))
			}

			lastReturnCode = output.ReturnCode
		}
	}
}
