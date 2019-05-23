package ielejson

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	oj "github.com/ElrondNetwork/elrond-vm/iele/test-util/orderedjson"
)

// ParseTopLevel ... converts json string to object representation
func ParseTopLevel(jsonString []byte) ([]*Test, error) {

	jobj, err := oj.ParseOrderedJSON(jsonString)
	if err != nil {
		return nil, err
	}

	topMap, isMap := jobj.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled test top level object is not a map")
	}

	var top []*Test
	for _, kvp := range topMap.OrderedKV {
		t, tErr := processTest(kvp.Value)
		if tErr != nil {
			return nil, tErr
		}
		t.TestName = kvp.Key
		top = append(top, t)
	}
	return top, nil
}

func processTest(testObj oj.OJsonObject) (*Test, error) {
	testMap, isTestMap := testObj.(*oj.OJsonMap)
	if !isTestMap {
		return nil, errors.New("unmarshalled test object is not a map")
	}
	test := Test{}

	for _, kvp := range testMap.OrderedKV {

		if kvp.Key == "pre" {
			preMap, isPreMap := kvp.Value.(*oj.OJsonMap)
			if !isPreMap {
				return nil, errors.New("unmarshalled pre object is not a map")
			}
			for _, acctKVP := range preMap.OrderedKV {
				acct, acctErr := processAccount(acctKVP.Value)
				if acctErr != nil {
					return nil, acctErr
				}
				acctAddr, hexErr := processAccountAddress(acctKVP.Key)
				if hexErr != nil {
					return nil, hexErr
				}
				acct.Address = acctAddr
				test.Pre = append(test.Pre, acct)
			}
		}

		if kvp.Key == "blocks" {
			blocksRaw, blocksOk := kvp.Value.(*oj.OJsonList)
			if !blocksOk {
				return nil, errors.New("unmarshalled blocks object is not a list")
			}
			for _, blRaw := range blocksRaw.AsList() {
				bl, blErr := processBlock(blRaw)
				if blErr != nil {
					return nil, blErr
				}
				test.Blocks = append(test.Blocks, bl)
			}
		}

		if kvp.Key == "network" {
			var networkOk bool
			test.Network, networkOk = parseString(kvp.Value)
			if !networkOk {
				return nil, errors.New("test network value not a string")
			}
		}

		if kvp.Key == "blockhashes" {
			var bhsOk bool
			test.BlockHashes, bhsOk = processBigIntList(kvp.Value)
			if !bhsOk {
				return nil, errors.New("unmarshalled blockHashes object is not a list")
			}
		}

		if kvp.Key == "postState" {
			postMap, isPostMap := kvp.Value.(*oj.OJsonMap)
			if !isPostMap {
				return nil, errors.New("unmarshalled postState object is not a map")
			}
			for _, acctKVP := range postMap.OrderedKV {
				acct, acctErr := processAccount(acctKVP.Value)
				if acctErr != nil {
					return nil, acctErr
				}
				acctAddr, hexErr := processAccountAddress(acctKVP.Key)
				if hexErr != nil {
					return nil, hexErr
				}
				acct.Address = acctAddr
				test.PostState = append(test.PostState, acct)
			}
		}
	}

	return &test, nil
}

func processAccount(acctRaw oj.OJsonObject) (*Account, error) {
	acctMap, isMap := acctRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled account object is not a map")
	}

	acct := Account{}
	var nonceOk, balanceOk, codeOk bool

	for _, kvp := range acctMap.OrderedKV {

		if kvp.Key == "nonce" {
			acct.Nonce, nonceOk = parseBigInt(kvp.Value)
			if !nonceOk {
				return nil, errors.New("invalid account nonce")
			}
		}

		if kvp.Key == "balance" {
			acct.Balance, balanceOk = parseBigInt(kvp.Value)
			if !balanceOk {
				return nil, errors.New("invalid account balance")
			}
		}

		if kvp.Key == "storage" {
			storageMap, storageOk := kvp.Value.(*oj.OJsonMap)
			if !storageOk {
				return nil, errors.New("invalid account storage")
			}
			for _, storageKvp := range storageMap.OrderedKV {
				intKey := big.NewInt(0)
				_, keyOk := intKey.SetString(storageKvp.Key, 0)
				if !keyOk {
					return nil, errors.New("invalid account storage key")
				}
				intVal, valOk := parseBigInt(storageKvp.Value)
				if !valOk {
					return nil, errors.New("invalid account storage value")
				}
				stElem := StorageKeyValuePair{Key: intKey, Value: intVal}
				acct.Storage = append(acct.Storage, &stElem)
			}
		}

		if kvp.Key == "code" {
			acct.Code, codeOk = parseString(kvp.Value)
			if !codeOk {
				return nil, errors.New("invalid account code")
			}
		}
	}

	return &acct, nil
}

func processBlock(blockRaw oj.OJsonObject) (*Block, error) {
	blockMap, isMap := blockRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block object is not a map")
	}
	bl := Block{}

	for _, kvp := range blockMap.OrderedKV {

		if kvp.Key == "results" {
			resultsRaw, resultsOk := kvp.Value.(*oj.OJsonList)
			if !resultsOk {
				return nil, errors.New("unmarshalled block results object is not a list")
			}
			for _, resRaw := range resultsRaw.AsList() {
				blr, blrErr := processBlockResult(resRaw)
				if blrErr != nil {
					return nil, blrErr
				}
				bl.Results = append(bl.Results, blr)
			}
		}

		if kvp.Key == "transactions" {
			transactionsRaw, transactionsOk := kvp.Value.(*oj.OJsonList)
			if !transactionsOk {
				return nil, errors.New("unmarshalled block transactions object is not a list")
			}
			for _, trRaw := range transactionsRaw.AsList() {
				tr, trErr := processBlockTransaction(trRaw)
				if trErr != nil {
					return nil, trErr
				}
				bl.Transactions = append(bl.Transactions, tr)
			}
		}

		if kvp.Key == "blockHeader" {
			blh, blhErr := processBlockHeader(kvp.Value)
			if blhErr != nil {
				return nil, blhErr
			}
			bl.BlockHeader = blh
		}
	}

	return &bl, nil
}

func processBlockResult(blrRaw oj.OJsonObject) (*TransactionResult, error) {
	blrMap, isMap := blrRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block result is not a map")
	}

	blr := TransactionResult{}
	var outOk, statusOk, gasOk, logsOk, refundOk bool

	for _, kvp := range blrMap.OrderedKV {

		if kvp.Key == "out" {
			blr.Out, outOk = processBigIntList(kvp.Value)
			if !outOk {
				return nil, errors.New("invalid block result out")
			}
		}

		if kvp.Key == "status" {
			blr.Status, statusOk = parseBigInt(kvp.Value)
			if !statusOk {
				return nil, errors.New("invalid block result status")
			}
		}

		if kvp.Key == "gas" {
			blr.Gas, gasOk = parseBigInt(kvp.Value)
			if !gasOk {
				return nil, errors.New("invalid block result gas")
			}
		}

		if kvp.Key == "logs" {
			blr.Logs, logsOk = parseString(kvp.Value)
			if !logsOk {
				return nil, errors.New("invalid block result logs")
			}
		}

		if kvp.Key == "refund" {
			blr.Refund, refundOk = parseBigInt(kvp.Value)
			if !refundOk {
				return nil, errors.New("invalid block result refund")
			}
		}
	}

	return &blr, nil
}

func processBlockTransaction(blrRaw oj.OJsonObject) (*Transaction, error) {
	bltMap, isMap := blrRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block transaction is not a map")
	}

	blt := Transaction{}
	var nonceOk, functionOk, valueOk, argumentsOk, contractCodeOk, gasPriceOk, gasLimitOk bool

	for _, kvp := range bltMap.OrderedKV {

		if kvp.Key == "nonce" {
			blt.Nonce, nonceOk = parseBigInt(kvp.Value)
			if !nonceOk {
				return nil, errors.New("invalid block transaction nonce")
			}
		}

		if kvp.Key == "function" {
			blt.Function, functionOk = parseString(kvp.Value)
			if !functionOk {
				return nil, errors.New("invalid block transaction function")
			}
		}

		if kvp.Key == "contractCode" {
			blt.ContractCode, contractCodeOk = parseString(kvp.Value)
			if !contractCodeOk {
				return nil, errors.New("invalid block transaction contract code")
			}
		}

		if kvp.Key == "value" {
			blt.Value, valueOk = parseBigInt(kvp.Value)
			if !valueOk {
				return nil, errors.New("invalid block transaction value")
			}
		}

		if kvp.Key == "to" {
			toStr, toOk := parseString(kvp.Value)
			if !toOk {
				return nil, errors.New("invalid block transaction to")
			}
			var toErr error
			blt.To, toErr = processAccountAddress(toStr)
			if toErr != nil {
				return nil, toErr
			}

			// note "to": "0x00" has to yield isCreate=false, even though it parses to 0, just like the 2 cases below
			blt.IsCreate = toStr == "" || toStr == "0x"

		}

		if kvp.Key == "arguments" {
			blt.Arguments, argumentsOk = processBigIntList(kvp.Value)
			if !argumentsOk {
				return nil, errors.New("invalid block transaction arguments")
			}
		}

		if kvp.Key == "contractCode" {
			blt.ContractCode, contractCodeOk = parseString(kvp.Value)
			if !contractCodeOk {
				return nil, errors.New("invalid block transaction contractCode")
			}
		}

		if kvp.Key == "gasPrice" {
			blt.GasPrice, gasPriceOk = parseBigInt(kvp.Value)
			if !gasPriceOk {
				return nil, errors.New("invalid block transaction gasPrice")
			}
		}

		if kvp.Key == "gasLimit" {
			blt.GasLimit, gasLimitOk = parseBigInt(kvp.Value)
			if !gasLimitOk {
				return nil, errors.New("invalid block transaction gasLimit")
			}
		}

		if kvp.Key == "from" {
			fromStr, fromOk := parseString(kvp.Value)
			if !fromOk {
				return nil, errors.New("invalid block transaction from")
			}
			var fromErr error
			blt.From, fromErr = processAccountAddress(fromStr)
			if fromErr != nil {
				return nil, fromErr
			}
		}
	}

	return &blt, nil
}

func processBlockHeader(blhRaw interface{}) (*BlockHeader, error) {
	blhMap, isMap := blhRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block header is not a map")
	}

	blh := BlockHeader{}
	var gasLimitOk, numberOk, difficultyOk, timestampOk, coinbaseOk bool

	for _, kvp := range blhMap.OrderedKV {

		if kvp.Key == "gasLimit" {
			blh.GasLimit, gasLimitOk = parseBigInt(kvp.Value)
			if !gasLimitOk {
				return nil, errors.New("invalid block header gasLimit")
			}
		}

		if kvp.Key == "number" {
			blh.Number, numberOk = parseBigInt(kvp.Value)
			if !numberOk {
				return nil, errors.New("invalid block header number")
			}
		}

		if kvp.Key == "difficulty" {
			blh.Difficulty, difficultyOk = parseBigInt(kvp.Value)
			if !difficultyOk {
				return nil, errors.New("invalid block header difficulty")
			}
		}

		if kvp.Key == "timestamp" {
			blh.UnixTimestamp, timestampOk = parseBigInt(kvp.Value)
			if !timestampOk {
				return nil, errors.New("invalid block header timestamp")
			}
		}

		if kvp.Key == "coinbase" {
			blh.Beneficiary, coinbaseOk = parseBigInt(kvp.Value)
			if !coinbaseOk {
				return nil, errors.New("invalid block header coinbase")
			}
		}
	}

	return &blh, nil
}

func processAccountAddress(addrRaw string) ([]byte, error) {
	if len(addrRaw) == 0 {
		return []byte{}, errors.New("missing account address")
	}
	if !(strings.HasPrefix(addrRaw, "0x") || strings.HasPrefix(addrRaw, "0X")) {
		return []byte{}, errors.New("account address should be hex representation starting with '0x'")
	}
	return hex.DecodeString(addrRaw[2:])
}

func processStringList(obj interface{}) ([]string, bool) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, false
	}
	var result []string
	for _, elemRaw := range listRaw.AsList() {
		str, strOk := elemRaw.(*oj.OJsonString)
		if !strOk {
			return nil, false
		}
		result = append(result, str.Value)
	}
	return result, true
}

func processBigIntList(obj interface{}) ([]*big.Int, bool) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, false
	}
	var result []*big.Int
	for _, elemRaw := range listRaw.AsList() {
		i, iOk := parseBigInt(elemRaw)
		if !iOk {
			return nil, false
		}
		result = append(result, i)
	}
	return result, true
}

func parseBigInt(obj oj.OJsonObject) (*big.Int, bool) {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return nil, false
	}
	result := new(big.Int)
	var parseOk bool
	if len(str.Value) > 0 { // interpret "" as 0
		result, parseOk = result.SetString(str.Value, 0)
		if !parseOk {
			return nil, false
		}
	}

	return result, true
}

func parseString(obj oj.OJsonObject) (string, bool) {
	str, isStr := obj.(*oj.OJsonString)
	if !isStr {
		return "", false
	}
	return str.Value, true
}