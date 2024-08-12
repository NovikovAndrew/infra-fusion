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
)

type EthClientConfig struct {
	ApiURL         string
	ConnectionType commonType.SubscribeType
	TickerTime     time.Duration
}

type BlockchainClient struct {
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

func NewBlockchainClient(logger *slog.Logger, cfg EthClientConfig) *BlockchainClient {
	return &BlockchainClient{
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

func (bc *BlockchainClient) SubscribeFilterLogs(ctx context.Context, addresses []common.Address, topics [][]common.Hash, outCh chan types.Log) (ethereum.Subscription, error) {
	q := ethereum.FilterQuery{
		Addresses: addresses,
		Topics:    topics,
	}
	return bc.client.SubscribeFilterLogs(ctx, q, outCh)
}

func (bc *BlockchainClient) SubscribeNewHead(ctx context.Context, headerCh chan *types.Header) (ethereum.Subscription, error) {
	return bc.client.SubscribeNewHead(ctx, headerCh)
}

func (bc *BlockchainClient) ContractsLog(ctx context.Context, addresses []common.Address) ([]types.Log, error) {
	q := ethereum.FilterQuery{
		Addresses: addresses,
	}
	return bc.client.FilterLogs(ctx, q)
}

func (bc *BlockchainClient) BlockNumber(ctx context.Context) (uint64, error) {
	return bc.client.BlockNumber(ctx)
}

func (bc *BlockchainClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return bc.client.BlockByNumber(ctx, number)
}

func (bc *BlockchainClient) SubscribeHead(ctx context.Context, headerCh chan *types.Header) (ethereum.Subscription, error) {
	return bc.client.SubscribeNewHead(ctx, headerCh)
}

func (bc *BlockchainClient) Subscribe(ctx context.Context) (chan *types.Block, error) {
	switch bc.cfg.ConnectionType {
	case commonType.WSSubscriptionType:
		err := bc.runWS(ctx)
		return bc.blockCh, err
	case commonType.HTTP:
		bc.runHttp(ctx)
		return bc.blockCh, nil
	}

	return nil, ErrorWrongSubType
}

func (bc *BlockchainClient) GetTxByHash(ctx context.Context, txHash string) (*types.Transaction, bool, error) {
	return bc.client.TransactionByHash(ctx, common.HexToHash(txHash))
}

func (bc *BlockchainClient) GetBlockByHeight(ctx context.Context, blockHeight uint64) (*types.Block, error) {
	return bc.client.BlockByNumber(ctx, big.NewInt(int64(blockHeight)))
}

func (bc *BlockchainClient) Unsubscribe() {
	if bc.sub != nil {
		bc.sub.Unsubscribe()
	}

	bc.client.Close()
	close(bc.blockCh)
}

func (bc *BlockchainClient) runHttp(ctx context.Context) {
	ticker := time.NewTicker(bc.cfg.TickerTime)
	go func() {
		for {
			select {
			case _ = <-ticker.C:
				blockNumber, err := bc.BlockNumber(ctx)
				if err != nil {
					bc.logger.Error("runHttp: failed to get block number", slog.StringValue(err.Error()))
					continue
				}

				if atomic.LoadUint64(&bc.previousBlockNumber) == blockNumber {
					continue
				}
				atomic.StoreUint64(&bc.previousBlockNumber, blockNumber)

				block, err := bc.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
				if err != nil {
					bc.logger.Error("runHttp: failed to get block by number", slog.StringValue(err.Error()))
					continue
				}

				bc.blockCh <- block
			case <-ctx.Done():
				ticker.Stop()
				bc.client.Close()
				close(bc.blockCh)
				return
			}
		}
	}()
}

func (bc *BlockchainClient) runWS(ctx context.Context) error {
	headerCh := make(chan *types.Header)
	sub, err := bc.SubscribeHead(ctx, headerCh)
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
				block, err := bc.BlockByNumber(ctx, header.Number)
				if err != nil {
					bc.logger.Error("runHttp: failed to get block by number", slog.StringValue(err.Error()))
					continue
				}

				bc.blockCh <- block
			case <-ctx.Done():
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
