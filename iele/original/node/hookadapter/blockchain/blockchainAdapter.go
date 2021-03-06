package blockchainadapter

import (
	"errors"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
	m "github.com/ElrondNetwork/elrond-vm/iele/original/node/iele-testing-kompiled/ieletestingmodel"
)

// Blockchain is an adapter between K and the outside world
// This class is specific to only 1 generated interpreter
type Blockchain struct {
	// Upstream is the world state callback, which is common to all VMs
	Upstream vmi.BlockchainHook

	// AddressLength is the expected length of an address, in bytes
	AddressLength int
}

// ConvertKIntToAddress takes a K Int and converts it to an address with the correct number of bytes,
// will pad left with zeroes, based on the configured address length
func (b *Blockchain) ConvertKIntToAddress(addrAsK m.KReference, ms *m.ModelState) ([]byte, bool) {
	addrInt, isInt := ms.GetBigInt(addrAsK)
	if !isInt {
		return []byte{}, false
	}
	addrBytes := addrInt.Bytes()
	if len(addrBytes) > b.AddressLength {
		return []byte{}, false
	}
	result := make([]byte, b.AddressLength)

	i := len(addrBytes) - 1
	j := b.AddressLength - 1
	for i >= 0 {
		result[j] = addrBytes[i]
		i--
		j--
	}

	return result, true
}

// GetBalance adapts between K model and elrond function
func (b *Blockchain) GetBalance(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, ms *m.ModelState) (m.KReference, error) {
	acctAddr, isAddr := b.ConvertKIntToAddress(c, ms)
	if !isAddr {
		return m.NoResult, errors.New("invalid account address provided to blockchain hook GetBalance")
	}
	result, err := b.Upstream.GetBalance(acctAddr)
	if err != nil {
		return m.NoResult, err
	}
	return ms.FromBigInt(result), nil
}

// GetNonce adapts between K model and elrond function
func (b *Blockchain) GetNonce(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, ms *m.ModelState) (m.KReference, error) {
	acctAddr, isAddr := b.ConvertKIntToAddress(c, ms)
	if !isAddr {
		return m.NoResult, errors.New("invalid account address provided to blockchain hook GetNonce")
	}
	result, err := b.Upstream.GetNonce(acctAddr)
	if err != nil {
		return m.NoResult, err
	}
	return ms.FromUint64(result), nil
}

// IsCodeEmpty adapts between K model and elrond function
func (b *Blockchain) IsCodeEmpty(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, ms *m.ModelState) (m.KReference, error) {
	acctAddr, isAddr := b.ConvertKIntToAddress(c, ms)
	if !isAddr {
		return m.NoResult, errors.New("invalid account address provided to blockchain hook IsCodeEmpty")
	}
	result, err := b.Upstream.IsCodeEmpty(acctAddr)
	if err != nil {
		return m.NoResult, err
	}
	return m.ToKBool(result), nil
}

// AccountExists adapts between K model and elrond function
func (b *Blockchain) AccountExists(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, ms *m.ModelState) (m.KReference, error) {
	acctAddr, isAddr := b.ConvertKIntToAddress(c, ms)
	if !isAddr {
		return m.NoResult, errors.New("invalid account address provided to blockchain hook AccountExists")
	}
	result, err := b.Upstream.AccountExists(acctAddr)
	if err != nil {
		return m.NoResult, err
	}
	return m.ToKBool(result), nil
}

// GetStorageData adapts between K model and elrond function
func (b *Blockchain) GetStorageData(kaddr m.KReference, kindex m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, ms *m.ModelState) (m.KReference, error) {
	acctAddr, isAddr := b.ConvertKIntToAddress(kaddr, ms)
	if !isAddr {
		return m.NoResult, errors.New("invalid account address provided to blockchain hook GetStorageData")
	}
	index, isInt2 := ms.GetBigInt(kindex)
	if !isInt2 {
		return m.NoResult, errors.New("invalid argument(s) provided to blockchain hook")
	}
	result, err := b.Upstream.GetStorageData(acctAddr, index.Bytes())
	if err != nil {
		return m.NoResult, err
	}
	return ms.IntFromBytes(result), nil
}

// GetCode adapts between K model and elrond function
func (b *Blockchain) GetCode(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, ms *m.ModelState) (m.KReference, error) {
	acctAddr, isAddr := b.ConvertKIntToAddress(c, ms)
	if !isAddr {
		return m.NoResult, errors.New("invalid account address provided to blockchain hook GetCode")
	}
	result, err := b.Upstream.GetCode(acctAddr)
	if err != nil {
		return m.NoResult, err
	}
	return ms.NewString(string(result)), nil
}

// GetBlockhash adapts between K model and elrond function
func (b *Blockchain) GetBlockhash(c m.KReference, lbl m.KLabel, sort m.Sort, config m.KReference, ms *m.ModelState) (m.KReference, error) {
	nonce, isInt := ms.GetUint64(c)
	if !isInt {
		return m.NoResult, errors.New("invalid argument(s) provided to blockchain hook")
	}
	result, err := b.Upstream.GetBlockhash(nonce)
	if err != nil {
		return m.NoResult, err
	}
	return ms.IntFromBytes(result), nil
}
