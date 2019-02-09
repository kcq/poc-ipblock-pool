package pool

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/segmentio/ksuid"
)

const (
	poolInfoKey          = "poc/pool/info"
	poolLockKey          = "poc/pool/.lock"
	poolBlocksKeyPrefix  = "poc/pool/blocks"
	defaultBaseSubnet    = "169.254.0.0/16"
	defaultStartRange    = "169.254.51.0"
	defaultEndRange      = "169.254.255.244"
	defaultPoolBlockSize = 4
)

// Pool errors
var (
	//
	ErrBlockNotFound = errors.New("Block not found")
)

// StoreConfig contains the Pool Store configurations
type StoreConfig struct {
	Address    string
	Scheme     string
	Datacenter string
}

// Config contains the Pool (Manager) configurations
type Config struct {
	StartRange    string
	EndRange      string
	PoolBlockSize int64
	Store         *StoreConfig
}

// Info contains the Pool metadata persisted in the Pool Store
type Info struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Next  string `json:"next"`
}

// NewPoolInfo creates a new Pool Info object
func NewPoolInfo(start, end, next string) *Info {
	info := Info{
		Start: start,
		End:   end,
		Next:  next,
	}

	return &info
}

// BlockInfo contains the IP Block metadata persisted in the Pool Store
type BlockInfo struct {
	ID    string `json:"id"`
	Start string `json:"start"`
	Key   string `json:"key"`
}

// NewBlockInfo creates a new IP Block Info object
func NewBlockInfo(start, key string) *BlockInfo {
	id, err := ksuid.NewRandom()
	if err != nil {
		panic(err)
	}

	info := BlockInfo{
		ID:    id.String(),
		Start: start,
		Key:   key,
	}

	return &info
}

// Manager is responsible for managing the IP Block Pool
type Manager struct {
	store         *Store
	info          *Info
	startIP       net.IP
	endIP         net.IP
	nextBlock     net.IP
	poolBlockSize int64
	startRange    string
	endRange      string
}

// New creates a new Pool Manager object
func New(configInfo *Config, store *Store) *Manager {
	if store == nil && configInfo != nil && configInfo.Store != nil {
		store = NewStoreWithConfig(configInfo.Store)
	}

	if store == nil {
		fmt.Println("pool.New: using the default Store...")
		store = NewStore(nil)
	}
	pool := Manager{
		store:         store,
		poolBlockSize: defaultPoolBlockSize,
		startRange:    defaultStartRange,
		endRange:      defaultEndRange,
	}

	if configInfo != nil {
		if configInfo.StartRange != "" {
			pool.startRange = configInfo.StartRange
		}

		if configInfo.EndRange == "" {
			pool.endRange = configInfo.EndRange
		}

		if configInfo.PoolBlockSize == 0 {
			pool.poolBlockSize = configInfo.PoolBlockSize
		}
	}

	fmt.Printf("pool.New: manager => %+v\n", pool)

	pool.init()
	return &pool
}

func (pool *Manager) init() {
	fmt.Println("Pool.init: trying to get the pool lock...")

	lock := pool.store.GetLock()
	lockCh, err := lock.Lock(nil)
	if err != nil {
		panic(err)
	}
	if lockCh == nil {
		panic("did not lock")
	}
	defer lock.Unlock()
	fmt.Println("Pool.init: got the pool lock...")

	pool.info = pool.store.GetPool()

	if pool.info == nil {
		fmt.Println("Pool Info - not initialized yet...")

		pool.info = NewPoolInfo(pool.startRange, pool.endRange, pool.startRange)
		pool.store.SavePool(pool.info)

		pool.startIP = net.ParseIP(pool.startRange)
		pool.endIP = net.ParseIP(pool.endRange)
		pool.nextBlock = pool.startIP

	} else {
		fmt.Printf("Pool Info - restored => %#v\n", pool.info)
		pool.startIP = net.ParseIP(pool.info.Start)
		pool.endIP = net.ParseIP(pool.info.End)
		pool.nextBlock = net.ParseIP(pool.info.Next)
	}
}

func (pool *Manager) nextBlockFromRange() string {
	fmt.Printf("nextBlockFromRange - pool.nextBlock => %#v\n", pool.nextBlock)
	//NOTE: nextBlock needs to be fresh when nextBlockFromRange is called
	allocated := pool.nextBlock.String()

	if ipVal := pool.nextBlock.To4(); ipVal != nil {
		ipNum := big.NewInt(0).SetBytes(ipVal)
		pool.nextBlock = net.IP(ipNum.Add(ipNum, big.NewInt(pool.poolBlockSize)).Bytes())

		//NOTE: info needs to be fresh when nextBlockFromRange is called
		pool.info.Next = pool.nextBlock.String()
		pool.store.SavePool(pool.info)

		fmt.Println("nextBlockFromRange - updated current pool info (nextBlock)...")

	} else {
		panic("nextBlockFromRange - unexpected IP value")
	}

	return allocated
}

// Lookup returns the IP Block metadata by the IP Block start address 
// or the Block Key or nil if the IP Block is not allocated yet
func (pool *Manager) Lookup(ipBlock, blockKey string) *BlockInfo {
	if ipBlock != "" {
		return pool.store.GetBlock(ipBlock)
	} else if blockKey != "" {
		return pool.store.FindBlock(blockKey)
	}

	return nil
}

// Allocate returns the newly allocated IP Block or an existing IP Block
// if the provided Block Key matches an existing IP Block allocation
func (pool *Manager) Allocate(blockKey string, delayUnlock bool) *BlockInfo {
	fmt.Println("Pool.Allocate - Trying to get the pool lock...")

	lock := pool.store.GetLock()
	lockCh, err := lock.Lock(nil)
	if err != nil {
		panic(err)
	}
	if lockCh == nil {
		panic("did not lock")
	}

	if !delayUnlock {
		defer lock.Unlock()
	}

	fmt.Println("Pool.Allocate - Got the pool lock...")

	if blockKey != "" {
		if blockInfo := pool.store.FindBlock(blockKey); blockInfo != nil {
			fmt.Println("Pool.Allocate - Already allocated... Returning existing record")
			return blockInfo
		}
	}

	//TODO - refresh pool info here
	blockStart := pool.nextBlockFromRange()
	fmt.Println("Pool.Allocate - Allocated IP block =>", blockStart)

	blockInfo := NewBlockInfo(blockStart, blockKey)
	pool.store.SaveBlock(blockInfo)

	if delayUnlock {
		go func() {
			fmt.Println("Pool.Allocate - Keeping the lock for 15 seconds to demo concurrent IP block allocation...")
			time.Sleep(15 * time.Second)
			fmt.Println("Pool.Allocate - delayed lock release...")
			lock.Unlock()
		}()
	}

	return blockInfo
}

// Free releases the selected IP Block allocation 
// based on the provided IP Block starting address or its Block Key
func (pool *Manager) Free(ipBlock, blockKey string) error {
	fmt.Println("Pool.Free - Trying to get the pool lock...")
	lock := pool.store.GetLock()
	lockCh, err := lock.Lock(nil)
	if err != nil {
		panic(err)
	}
	if lockCh == nil {
		panic("did not lock")
	}
	defer lock.Unlock()

	fmt.Println("Pool.Free - Got the pool lock...")

	if ipBlock != "" {
		if blockInfo := pool.store.GetBlock(ipBlock); blockInfo != nil {
			fmt.Println("Pool.Free - Found record by IP")
			pool.store.RemoveBlock(ipBlock)
			return nil
		}
	} else if blockKey != "" {
		if blockInfo := pool.store.FindBlock(blockKey); blockInfo != nil {
			fmt.Println("Pool.Free - Found record by Key")
			pool.store.RemoveBlock(blockInfo.Start)
			return nil
		}
	}

	return ErrBlockNotFound
}

// Store represents the Pool data store (Consul is used as the store backend)
type Store struct {
	consul *api.Client
	kvAPI  *api.KV
}

// GetLock returns the data store lock object
// You must explicitly acquired/lock the lock object to ensure exclusive access to the Store
func (s *Store) GetLock() *api.Lock {
	lock, err := s.consul.LockKey(poolLockKey)
	if err != nil {
		panic(err)
	}

	return lock
}

// GetRecord returns the selected record
func (s *Store) GetRecord(key string) []byte {
	pair, _, err := s.kvAPI.Get(key, nil)
	if err != nil {
		panic(err)
	}

	if pair == nil || pair.Value == nil {
		return nil
	}

	return pair.Value
}

// SaveRecord saves the provided record in the Store
func (s *Store) SaveRecord(key string, data []byte) {
	pair := &api.KVPair{Key: key, Value: data}
	if _, err := s.kvAPI.Put(pair, nil); err != nil {
		panic(err)
	}
}

// RemoveRecord removes the selected record from the Store
func (s *Store) RemoveRecord(key string) {
	if _, err := s.kvAPI.Delete(key, nil); err != nil {
		panic(err)
	}
}

// FindBlock returns the BlockInfo object selected by IP Block Key
func (s *Store) FindBlock(key string) *BlockInfo {
	//NOTE: this is a hacky way to find the record by key (good enough for a PoC :-))
	if pairs, _, err := s.kvAPI.List(poolBlocksKeyPrefix, nil); err != nil {
		panic(err)
	} else if len(pairs) == 0 {
		return nil
	} else {
		for _, p := range pairs {
			var record map[string]string
			if err := json.Unmarshal(p.Value, &record); err != nil {
				panic(err)
			}

			var block BlockInfo
			if err := json.Unmarshal(p.Value, &block); err != nil {
				panic(err)
			}

			if block.Key == key {
				return &block
			}
		}
	}

	return nil
}

// GetBlock returns the BlockInfo object selected by the IP Block starting address
func (s *Store) GetBlock(blockStart string) *BlockInfo {
	key := fmt.Sprintf("%s/%s", poolBlocksKeyPrefix, blockStart)
	raw := s.GetRecord(key)
	if raw == nil {
		return nil
	}

	var block BlockInfo
	if err := json.Unmarshal(raw, &block); err != nil {
		panic(err)
	}

	return &block
}

// SaveBlock saves the provided BlockInfo object
func (s *Store) SaveBlock(block *BlockInfo) {
	key := fmt.Sprintf("%s/%s", poolBlocksKeyPrefix, block.Start)

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(block); err != nil {
		panic(err)
	}

	s.SaveRecord(key, buf.Bytes())
}

// RemoveBlock removes the BlockInfo object selected by the IP Block starting address
func (s *Store) RemoveBlock(blockStart string) {
	key := fmt.Sprintf("%s/%s", poolBlocksKeyPrefix, blockStart)
	s.RemoveRecord(key)
}

// SavePool saves the Pool metadata
func (s *Store) SavePool(pool *Info) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(pool); err != nil {
		panic(err)
	}

	s.SaveRecord(poolInfoKey, buf.Bytes())
}

// GetPool restores the Pool metadata from the Store backend
func (s *Store) GetPool() *Info {
	raw := s.GetRecord(poolInfoKey)
	if raw == nil {
		return nil
	}

	var pool Info
	if err := json.Unmarshal(raw, &pool); err != nil {
		panic(err)
	}

	return &pool
}

// NewStore creates a new Store object based on the provided backend config
func NewStore(config *api.Config) *Store {
	fmt.Println("pool.NewStore...")
	if config == nil {
		fmt.Println("pool.NewStore: using the default Consul config...")
		config = api.DefaultConfig()
	}

	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	store := Store{
		consul: client,
		kvAPI:  client.KV(),
	}

	return &store
}

// NewStoreWithConfig creates a new Store object based on the provided Store config
func NewStoreWithConfig(configInfo *StoreConfig) *Store {
	config := api.DefaultConfig()

	if configInfo != nil {
		if configInfo.Address != "" {
			config.Address = configInfo.Address
		}

		if configInfo.Scheme != "" {
			config.Scheme = configInfo.Scheme
		}

		if configInfo.Datacenter != "" {
			config.Datacenter = configInfo.Datacenter
		}
	}

	return NewStore(config)
}
