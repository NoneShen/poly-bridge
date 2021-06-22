package switcheolisten

import (
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego/logs"
	"math/big"
	"poly-bridge/basedef"
	"poly-bridge/chainsdk"
	"poly-bridge/conf"
	"poly-bridge/models"
	"strconv"
	"strings"
)

const (
	_switcheo_crosschainlock    = "make_from_cosmos_proof"
	_switcheo_crosschainunlock  = "verify_to_cosmos_proof"
	_switcheo_lock = "lock"
	_switcheo_unlock = "unlock"
)

type SwitcheoChainListen struct {
	swthCfg *conf.ChainListenConfig
	swthSdk *chainsdk.SwitcheoSdkPro
}

func NewSwitcheoChainListen(cfg *conf.ChainListenConfig) *SwitcheoChainListen{
	swthListen:=&SwitcheoChainListen{}
	swthListen.swthCfg=cfg
	urls:=cfg.GetNodesUrl()
	sdk:=chainsdk.NewSwitcheoSdkPro(urls, cfg.ListenSlot, cfg.ChainId)
	swthListen.swthSdk=sdk
	return swthListen
}

func (this *SwitcheoChainListen) GetLatestHeight() (uint64, error) {
	return this.swthSdk.GetLatestHeight()
}


func (this *SwitcheoChainListen) GetChainListenSlot() uint64 {
	return this.swthCfg.ListenSlot
}

func (this *SwitcheoChainListen) GetChainId() uint64 {
	return this.swthCfg.ChainId
}

func (this *SwitcheoChainListen) GetChainName() string {
	return this.swthCfg.ChainName
}

func (this *SwitcheoChainListen) GetDefer() uint64 {
	return this.swthCfg.Defer
}

func (this *SwitcheoChainListen) HandleNewBlock(height uint64) ([]*models.WrapperTransaction, []*models.SrcTransaction, []*models.PolyTransaction, []*models.DstTransaction, error) {
	block, err := this.swthSdk.GetBlockByHeight(height)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if block == nil {
		return nil, nil, nil, nil, fmt.Errorf("there is no switcheo block!")
	}
	tt :=uint64(block.Block.Time.Unix())
	srcTransactions := make([]*models.SrcTransaction, 0)
	dstTransactions := make([]*models.DstTransaction, 0)

	ccmLockEvent,lockEvents, err := this.getCosmosCCMLockEventByBlockNumber(height)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	ccmUnlockEvent,unlockEvents, err := this.getCosmosCCMUnlockEventByBlockNumber(height)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	for _, lockEvent := range ccmLockEvent {
		if lockEvent.Method == _switcheo_crosschainlock {
			logs.Info("from chain: %d, txhash: %s\n", this.GetChainName(), lockEvent.TxHash)
			srcTransfer := &models.SrcTransfer{}
			for _, v := range lockEvents {
				if v.TxHash == lockEvent.TxHash {
					srcTransfer.ChainId = this.GetChainId()
					srcTransfer.TxHash = lockEvent.TxHash
					srcTransfer.Time = tt
					srcTransfer.From = v.FromAddress
					srcTransfer.To = lockEvent.Contract
					srcTransfer.Asset =v.FromAssetHash
					amount:= new(big.Int).SetUint64(v.Amount)
					srcTransfer.Amount=models.NewBigInt(amount)
					srcTransfer.DstChainId = uint64(v.ToChainId)
					srcTransfer.DstAsset = v.ToAssetHash
					srcTransfer.DstUser = v.ToAddress
					break
				}
			}
			srcTransaction := &models.SrcTransaction{}
			srcTransaction.ChainId = this.GetChainId()
			srcTransaction.Hash = lockEvent.TxHash
			srcTransaction.State = 1
			srcTransaction.Fee = models.NewBigIntFromInt(int64(lockEvent.Fee))
			srcTransaction.Time = tt
			srcTransaction.Height = lockEvent.Height
			srcTransaction.User = lockEvent.User
			srcTransaction.DstChainId = uint64(lockEvent.Tchain)
			srcTransaction.Contract = lockEvent.Contract
			srcTransaction.Key = lockEvent.Txid
			srcTransaction.Param = hex.EncodeToString(lockEvent.Value)
			srcTransaction.SrcTransfer = srcTransfer
			srcTransactions = append(srcTransactions, srcTransaction)
		}
	}
	for _, unLockEvent := range ccmUnlockEvent {
		if unLockEvent.Method == _switcheo_crosschainunlock {
			logs.Info("to chain: %s, txhash: %s\n", this.GetChainName(), unLockEvent.TxHash)
			dstTransfer := &models.DstTransfer{}
			for _, v := range unlockEvents {
				if v.TxHash == unLockEvent.TxHash {
					dstTransfer.ChainId = this.GetChainId()
					dstTransfer.TxHash = unLockEvent.TxHash
					dstTransfer.Time = tt
					dstTransfer.From =unLockEvent.Contract
					dstTransfer.To = v.ToAddress
					dstTransfer.Asset = v.ToAssetHash
					amount := new(big.Int).SetUint64(v.Amount)
					dstTransfer.Amount = models.NewBigInt(amount)
					break
				}
			}
			dstTransaction := &models.DstTransaction{}
			dstTransaction.ChainId = this.GetChainId()
			dstTransaction.Hash = unLockEvent.TxHash
			dstTransaction.State = 1
			dstTransaction.Fee = models.NewBigIntFromInt(int64(unLockEvent.Fee))
			dstTransaction.Time = tt
			dstTransaction.Height = height
			dstTransaction.SrcChainId = uint64(unLockEvent.FChainId)
			dstTransaction.Contract = unLockEvent.Contract
			dstTransaction.PolyHash = unLockEvent.RTxHash
			dstTransaction.DstTransfer = dstTransfer
			dstTransactions = append(dstTransactions, dstTransaction)
		}
	}
	return nil, srcTransactions, nil, dstTransactions, nil
}

func (this *SwitcheoChainListen) getCosmosCCMLockEventByBlockNumber(height uint64) ([]*models.ECCMLockEvent, []*models.LockEvent, error) {
	client := this.swthSdk
	ccmLockEvents := make([]*models.ECCMLockEvent, 0)
	lockEvents := make([]*models.LockEvent, 0)
	query := fmt.Sprintf("tx.height=%d AND make_from_cosmos_proof.status='1'", height)
	res, err := client.TxSearch(height,query, false, 1, 100, "asc")
	if err != nil {
		return ccmLockEvents, lockEvents, err
	}
	if res.TotalCount != 0 {
		pages := ((res.TotalCount - 1) / 100) + 1
		for p := 1; p <= pages; p ++ {
			if p > 1 {
				res, err = client.TxSearch(height,query, false, p, 100, "asc")
				if err != nil {
					return ccmLockEvents, lockEvents, err
				}
			}
			for _, tx := range res.Txs {
				for _, e := range tx.TxResult.Events {
					if e.Type == _switcheo_crosschainlock {
						tchainId, _ := strconv.ParseUint(string(e.Attributes[5].Value), 10, 32)
						value, _ := hex.DecodeString(string(e.Attributes[6].Value))
						ccmLockEvents = append(ccmLockEvents, &models.ECCMLockEvent{
							Method: _switcheo_crosschainlock,
							Txid: string(e.Attributes[1].Value),
							TxHash: strings.ToLower(tx.Hash.String()),
							User: string(e.Attributes[3].Value),
							Tchain: uint32(tchainId),
							Contract: string(e.Attributes[4].Value),
							Height: height,
							Value: value,
							Fee: uint64(tx.TxResult.GasUsed),
						})
					} else if e.Type == _switcheo_lock {
						tchainId, _ := strconv.ParseUint(string(e.Attributes[1].Value), 10, 32)
						amount, _ := strconv.ParseUint(string(e.Attributes[5].Value), 10, 64)
						lockEvents = append(lockEvents, &models.LockEvent{
							Method: _switcheo_lock,
							TxHash: strings.ToLower(tx.Hash.String()),
							FromAddress: string(e.Attributes[3].Value),
							FromAssetHash: string(e.Attributes[0].Value),
							ToChainId: uint32(tchainId),
							ToAssetHash: string(e.Attributes[2].Value),
							ToAddress: string(e.Attributes[4].Value),
							Amount: amount,
						})
					}
				}
			}
		}
	}

	return ccmLockEvents, lockEvents, nil
}

func (this *SwitcheoChainListen) getCosmosCCMUnlockEventByBlockNumber(height uint64) ([]*models.ECCMUnlockEvent, []*models.UnlockEvent, error) {
	client := this.swthSdk
	ccmUnlockEvents := make([]*models.ECCMUnlockEvent, 0)
	unlockEvents := make([]*models.UnlockEvent, 0)
	query := fmt.Sprintf("tx.height=%d", height)
	res, err := client.TxSearch(height,query, false, 1, 100, "asc")
	if err != nil {
		return ccmUnlockEvents, unlockEvents, err
	}
	if res.TotalCount != 0 {
		pages := ((res.TotalCount - 1) / 100) + 1
		for p := 1; p <= pages; p ++ {
			if p > 1 {
				res, err = client.TxSearch(height,query, false, p, 100, "asc")
				if err != nil {
					return ccmUnlockEvents, unlockEvents, err
				}
			}
			for _, tx := range res.Txs {
				for _, e := range tx.TxResult.Events {
					if e.Type == _switcheo_crosschainunlock {
						fchainId, _ := strconv.ParseUint(string(e.Attributes[2].Value), 10, 32)
						ccmUnlockEvents = append(ccmUnlockEvents, &models.ECCMUnlockEvent{
							Method: _switcheo_crosschainunlock,
							TxHash: strings.ToLower(tx.Hash.String()),
							RTxHash: basedef.HexStringReverse(string(e.Attributes[0].Value)),
							FChainId: uint32(fchainId),
							Contract: string(e.Attributes[3].Value),
							Height: height,
							Fee: uint64(tx.TxResult.GasUsed),
						})
					} else if e.Type == _switcheo_unlock {
						amount, _ := strconv.ParseUint(string(e.Attributes[2].Value), 10, 64)
						unlockEvents = append(unlockEvents, &models.UnlockEvent{
							Method: _switcheo_unlock,
							TxHash: strings.ToLower(tx.Hash.String()),
							ToAssetHash: string(e.Attributes[0].Value),
							ToAddress: string(e.Attributes[1].Value),
							Amount: amount,
						})
					}
				}
			}
		}
	}

	return ccmUnlockEvents, unlockEvents, nil
}

func (this *SwitcheoChainListen) GetExtendLatestHeight() (uint64, error) {
	if len(this.swthCfg.ExtendNodes) == 0 {
		return this.GetLatestHeight()
	}
	return this.GetLatestHeight()
}



