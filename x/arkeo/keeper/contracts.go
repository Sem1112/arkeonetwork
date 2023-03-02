package keeper

import (
	"errors"
	"strconv"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	gogotypes "github.com/gogo/protobuf/types"
)

func (k KVStore) setContract(ctx cosmos.Context, key string, record types.Contract) {
	store := ctx.KVStore(k.storeKey)
	buf := k.cdc.MustMarshal(&record)
	if buf == nil {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) getContract(ctx cosmos.Context, key string, record *types.Contract) (bool, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return false, nil
	}

	bz := store.Get([]byte(key))
	if err := k.cdc.Unmarshal(bz, record); err != nil {
		return true, err
	}
	return true, nil
}

// GetContractIterator iterate contract
func (k KVStore) GetContractIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixContract)
}

// GetContract get the entire Contract metadata struct based on given asset
func (k KVStore) GetContract(ctx cosmos.Context, pubkey common.PubKey, chain common.Chain, client common.PubKey) (types.Contract, error) {
	record := types.NewContract(pubkey, chain, client)
	_, err := k.getContract(ctx, k.GetKey(ctx, prefixContract, record.Key()), &record)

	return record, err
}

// SetContract save the entire Contract metadata struct to key value store
func (k KVStore) SetContract(ctx cosmos.Context, record types.Contract) error {
	if record.ProviderPubKey.IsEmpty() || record.Chain.IsEmpty() || record.Client.IsEmpty() {
		return errors.New("cannot save a contract with an empty provider pubkey, chain, or client address")
	}
	k.setContract(ctx, k.GetKey(ctx, prefixContract, record.Key()), record)
	return nil
}

// ContractExists check whether the given contract exist in the data store
func (k KVStore) ContractExists(ctx cosmos.Context, pubkey common.PubKey, chain common.Chain, client common.PubKey) bool {
	record := types.NewContract(pubkey, chain, client)
	return k.has(ctx, k.GetKey(ctx, prefixContract, record.Key()))
}

func (k KVStore) RemoveContract(ctx cosmos.Context, pubkey common.PubKey, chain common.Chain, client common.PubKey) {
	record := types.NewContract(pubkey, chain, client)
	k.del(ctx, k.GetKey(ctx, prefixContract, record.Key()))
}

func (k KVStore) setContractExpirationSet(ctx cosmos.Context, key string, record types.ContractExpirationSet) {
	store := ctx.KVStore(k.storeKey)
	buf := k.cdc.MustMarshal(&record)
	if len(record.Contracts) == 0 {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) getContractExpirationSet(ctx cosmos.Context, key string, record *types.ContractExpirationSet) (bool, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return false, nil
	}

	bz := store.Get([]byte(key))
	if err := k.cdc.Unmarshal(bz, record); err != nil {
		return true, err
	}
	return true, nil
}

func (k KVStore) getContractExpirationSetKey(ctx cosmos.Context, height int64) string {
	return k.GetKey(ctx, prefixContractExpirationSet, strconv.FormatInt(height, 10))
}

// GetContractExpirationSetIterator iterate contract expiration sets
func (k KVStore) GetContractExpirationSetIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixContractExpirationSet)
}

// GetContractExpirationSet get a series of contract expirations
func (k KVStore) GetContractExpirationSet(ctx cosmos.Context, height int64) (types.ContractExpirationSet, error) {
	record := types.ContractExpirationSet{
		Height: height,
	}
	_, err := k.getContractExpirationSet(ctx, k.getContractExpirationSetKey(ctx, height), &record)
	return record, err
}

// SetContractExpirationSet save the series of Contract Expirations
func (k KVStore) SetContractExpirationSet(ctx cosmos.Context, record types.ContractExpirationSet) error {
	if record.Height <= 0 {
		return errors.New("cannot save a contract expiration set with an invalid height (less than or equal to zero)")
	}
	k.setContractExpirationSet(ctx, k.getContractExpirationSetKey(ctx, record.Height), record)
	return nil
}

func (k KVStore) RemoveContractExpirationSet(ctx cosmos.Context, height int64) {
	k.del(ctx, k.GetKey(ctx, prefixContractExpirationSet, strconv.FormatInt(height, 10)))
}

func (kvStore KVStore) GetNextContractId(ctx cosmos.Context) uint64 {
	var contractId uint64
	store := ctx.KVStore(kvStore.storeKey)

	bz := store.Get([]byte(types.KeyNextGlobalContractId))
	if bz == nil {
		// initialize the pool id numbers
		contractId = 1
	} else {
		val := gogotypes.UInt64Value{}
		kvStore.cdc.MustUnmarshal(bz, &val)
		contractId = val.GetValue()
	}

	kvStore.setNextContractId(ctx, contractId+1)
	return contractId
}

func (k KVStore) setNextContractId(ctx cosmos.Context, contractId uint64) {
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: contractId})
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.KeyNextGlobalContractId), bz)
}
