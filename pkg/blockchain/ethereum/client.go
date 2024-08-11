package ethereum

import (
	"context"
	"log/slog"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	commonType "github.com/NovikovAndrew/infra-fusion/pkg/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

type EthClientConfig struct {
	ApiURL         string
	ConnectionType commonType.SubscribeType
	TickerTime     time.Duration
}

type BlockchainClient struct {
	ctx                 context.Context
	logger              *slog.Logger
	client              *ethclient.Client
	lock                *sync.Mutex
	apiURL              string
	isConnected         bool
	cfg                 EthClientConfig
	previousBlockNumber uint64
	blockCh             chan *types.Block
	sub                 ethereum.Subscription
}

func NewBlockchainClient(ctx context.Context, logger *slog.Logger, cfg EthClientConfig) *BlockchainClient {
	return &BlockchainClient{
		ctx:     ctx,
		logger:  logger.With(commonType.BlockchainAttr, "ethereum"),
		lock:    &sync.Mutex{},
		apiURL:  cfg.ApiURL,
		cfg:     cfg,
		blockCh: make(chan *types.Block),
	}
}

func (bc *BlockchainClient) Connect() error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	if bc.isConnected {
		return nil
	}

	client, err := ethclient.Dial(bc.apiURL)
	if err != nil {
		return err
	}

	bc.isConnected = true
	bc.client = client
	return nil
}

func (bc *BlockchainClient) Disconnect() {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	if !bc.isConnected || bc.client == nil {
		bc.isConnected = false
		return
	}
	bc.client.Close()
	bc.client = nil
	bc.client = nil
}

func (bc *BlockchainClient) IsConnected() bool {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	return bc.isConnected
}

func (bc *BlockchainClient) SubscribeFilterLogs(addresses []common.Address, topics [][]common.Hash, outCh chan types.Log) (ethereum.Subscription, error) {
	q := ethereum.FilterQuery{
		Addresses: addresses,
		Topics:    topics,
	}
	return bc.client.SubscribeFilterLogs(bc.ctx, q, outCh)
}

func (bc *BlockchainClient) SubscribeNewHead(headerCh chan *types.Header) (ethereum.Subscription, error) {
	return bc.client.SubscribeNewHead(bc.ctx, headerCh)
}

func (bc *BlockchainClient) ContractsLog(addresses []common.Address) ([]types.Log, error) {
	q := ethereum.FilterQuery{
		Addresses: addresses,
	}
	return bc.client.FilterLogs(bc.ctx, q)
}

func (bc *BlockchainClient) BlockNumber() (uint64, error) {
	return bc.client.BlockNumber(bc.ctx)
}

func (bc *BlockchainClient) BlockByNumber(number *big.Int) (*types.Block, error) {
	return bc.client.BlockByNumber(bc.ctx, number)
}

func (bc *BlockchainClient) SubscribeHead(headerCh chan *types.Header) (ethereum.Subscription, error) {
	return bc.client.SubscribeNewHead(bc.ctx, headerCh)
}

func (bc *BlockchainClient) Subscribe() (chan *types.Block, error) {
	switch bc.cfg.ConnectionType {
	case commonType.WSSubscriptionType:
		err := bc.runWS()
		return bc.blockCh, err
	case commonType.HTTP:
		bc.runHttp()
		return bc.blockCh, nil
	}

	return nil, ErrorWrongSubType
}

func (bc *BlockchainClient) Unsubscribe() {
	if bc.sub != nil {
		bc.sub.Unsubscribe()
	}

	bc.client.Close()
	close(bc.blockCh)
}

func (bc *BlockchainClient) runHttp() {
	ticker := time.NewTicker(bc.cfg.TickerTime)
	go func() {
		for {
			select {
			case _ = <-ticker.C:
				blockNumber, err := bc.BlockNumber()
				if err != nil {
					log.Error("runHttp: failed to get block number, err - [%s]", err.Error())
					continue
				}

				if atomic.LoadUint64(&bc.previousBlockNumber) == blockNumber {
					continue
				}
				atomic.StoreUint64(&bc.previousBlockNumber, blockNumber)

				block, err := bc.BlockByNumber(big.NewInt(int64(blockNumber)))
				if err != nil {
					log.Error("runHttp: failed to get block, err - [%s]", err.Error())
					continue
				}

				bc.blockCh <- block
			case _ = <-bc.ctx.Done():
				ticker.Stop()
				bc.client.Close()
				close(bc.blockCh)
				return
			}
		}
	}()
}

func (bc *BlockchainClient) runWS() error {
	headerCh := make(chan *types.Header)
	sub, err := bc.SubscribeHead(headerCh)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case _ = <-sub.Err():
				sub.Unsubscribe()
				bc.client.Close()
				close(bc.blockCh)
				return
			case header := <-headerCh:
				block, err := bc.BlockByNumber(header.Number)
				if err != nil {
					log.Error("runWS: failed to get block, err - [%s]", err.Error())
					continue
				}

				bc.blockCh <- block
			case _ = <-bc.ctx.Done():
				sub.Unsubscribe()
				bc.client.Close()
				close(bc.blockCh)
				return
			}
		}
	}()

	bc.sub = sub
	return nil
}

func (bc *BlockchainClient) GetTxByHash(ctx context.Context, txHash string) (*types.Transaction, bool, error) {
	return bc.client.TransactionByHash(ctx, common.HexToHash(txHash))
}

func (bc *BlockchainClient) GetBlockByHeight(ctx context.Context, blockHeight uint64) (*types.Block, error) {
	return bc.client.BlockByNumber(ctx, big.NewInt(int64(blockHeight)))
}
