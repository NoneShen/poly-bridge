/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package crosschainlisten

import (
	"fmt"
	"math"
	"runtime/debug"
	"time"

	"poly-bridge/basedef"
	"poly-bridge/conf"
	"poly-bridge/crosschaindao"
	"poly-bridge/crosschainlisten/ethereumlisten"
	"poly-bridge/crosschainlisten/neo3listen"
	"poly-bridge/crosschainlisten/neolisten"
	"poly-bridge/crosschainlisten/o3listen"
	"poly-bridge/crosschainlisten/ontologylisten"
	"poly-bridge/crosschainlisten/polylisten"
	"poly-bridge/crosschainlisten/switcheolisten"
	"poly-bridge/http/tools"
	"poly-bridge/models"

	"github.com/beego/beego/v2/core/logs"
)

var chainListens [12]*CrossChainListen

func StartCrossChainListen(server string, backup bool, listenCfg []*conf.ChainListenConfig, dbCfg *conf.DBConfig) {
	dao := crosschaindao.NewCrossChainDao(server, backup, dbCfg)
	if dao == nil {
		panic("server is not valid")
	}
	for i, cfg := range listenCfg {
		chainHandle := NewChainHandle(cfg)
		if chainHandle == nil {
			panic(fmt.Sprintf("chain %d handler is invalid", cfg.ChainId))
		}
		chainListen := NewCrossChainListen(chainHandle, dao)
		chainListen.Start()
		chainListens[i] = chainListen
	}
}

func StopCrossChainListen() {
	for _, chainListen := range chainListens {
		if chainListen != nil {
			chainListen.Stop()
		}
	}
}

type ChainHandle interface {
	GetExtendLatestHeight() (uint64, error)
	GetLatestHeight() (uint64, error)
	HandleNewBlock(height uint64) ([]*models.WrapperTransaction, []*models.SrcTransaction, []*models.PolyTransaction, []*models.DstTransaction, int, int, error)
	GetChainListenSlot() uint64
	GetChainId() uint64
	GetChainName() string
	GetDefer() uint64
}

func NewChainHandle(chainListenConfig *conf.ChainListenConfig) ChainHandle {
	if chainListenConfig.ChainId == basedef.ETHEREUM_CROSSCHAIN_ID {
		return ethereumlisten.NewEthereumChainListen(chainListenConfig)
	} else if chainListenConfig.ChainId == basedef.POLY_CROSSCHAIN_ID {
		return polylisten.NewPolyChainListen(chainListenConfig)
	} else if chainListenConfig.ChainId == basedef.NEO_CROSSCHAIN_ID {
		return neolisten.NewNeoChainListen(chainListenConfig)
	} else if chainListenConfig.ChainId == basedef.BSC_CROSSCHAIN_ID {
		return ethereumlisten.NewEthereumChainListen(chainListenConfig)
	} else if chainListenConfig.ChainId == basedef.HECO_CROSSCHAIN_ID {
		return ethereumlisten.NewEthereumChainListen(chainListenConfig)
	} else if chainListenConfig.ChainId == basedef.ONT_CROSSCHAIN_ID {
		return ontologylisten.NewOntologyChainListen(chainListenConfig)
	} else if chainListenConfig.ChainId == basedef.OK_CROSSCHAIN_ID {
		return ethereumlisten.NewEthereumChainListen(chainListenConfig)
	} else if chainListenConfig.ChainId == basedef.O3_CROSSCHAIN_ID {
		return o3listen.NewO3ChainListen(chainListenConfig)
	} else if chainListenConfig.ChainId == basedef.SWITCHEO_CROSSCHAIN_ID {
		return switcheolisten.NewSwitcheoChainListen(chainListenConfig)
	} else if chainListenConfig.ChainId == basedef.NEO3_CROSSCHAIN_ID {
		return neo3listen.NewNeo3ChainListen(chainListenConfig)
	} else if chainListenConfig.ChainId == basedef.MATIC_CROSSCHAIN_ID {
		return ethereumlisten.NewEthereumChainListen(chainListenConfig)
	} else if chainListenConfig.ChainId == basedef.PLT_CROSSCHAIN_ID {
		return ethereumlisten.NewEthereumChainListen(chainListenConfig)
	} else {
		return nil
	}
}

type CrossChainListen struct {
	handle ChainHandle
	db     crosschaindao.CrossChainDao
	exit   chan bool
	height uint64
}

func NewCrossChainListen(handle ChainHandle, db crosschaindao.CrossChainDao) *CrossChainListen {
	crossChainListen := &CrossChainListen{
		handle: handle,
		db:     db,
		exit:   make(chan bool, 0),
	}
	return crossChainListen
}

func (ccl *CrossChainListen) SetHeight(height uint64) {
	ccl.height = height
}

func (ccl *CrossChainListen) Start() {
	logs.Info("start cross chain listen: %s", ccl.handle.GetChainName())
	go ccl.ListenChain()
}

func (ccl *CrossChainListen) Stop() {
	ccl.exit <- true
	logs.Info("stop cross chain listen: %s", ccl.handle.GetChainName())
}

func (ccl *CrossChainListen) ListenChain() {
	for {
		exit := ccl.listenChain()
		if exit {
			close(ccl.exit)
			break
		}
		time.Sleep(time.Second * 5)
	}
}

func (ccl *CrossChainListen) HandleNewBlock(height uint64) (w []*models.WrapperTransaction, s []*models.SrcTransaction, p []*models.PolyTransaction, d []*models.DstTransaction, err error) {
	chain := ccl.handle.GetChainId()
	var locks, unlocks int
	for c := 3; c > 0; c-- {
		w, s, p, d, locks, unlocks, err = ccl.handle.HandleNewBlock(height)
		if err != nil {
			return
		}
		if locks == len(s) && unlocks == len(d) {
			return
		}
		if c > 1 {
			logs.Warn("Possible missing events for chain %v height %v", chain, height)
			time.Sleep(time.Second * 5)
		}
	}
	logs.Error("Possible inconsistent chain %d height %d wrapper %d/%d src %d/%d dst %d/%d", chain, height, len(w), locks, len(s), locks, len(d), unlocks)
	return
}
func (ccl *CrossChainListen) listenChain() (exit bool) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("service start, recover info: %s", string(debug.Stack()))
			exit = false
		}
	}()
	chain, err := ccl.db.GetChain(ccl.handle.GetChainId())
	if err != nil {
		panic(err)
	}
	height, err := ccl.handle.GetLatestHeight()
	if err != nil || height == 0 {
		panic(err)
	}
	if chain.Height == 0 {
		chain.Height = height
	}
	ccl.db.UpdateChain(chain)
	if ccl.height != 0 {
		chain.Height = ccl.height
	}
	logs.Info("cross chain listen, chain: %s, dao: %s......", ccl.handle.GetChainName(), ccl.db.Name())
	ticker := time.NewTicker(time.Second * time.Duration(ccl.handle.GetChainListenSlot()))
	for {
		select {
		case <-ticker.C:
			var height, err = ccl.handle.GetLatestHeight()
			if err != nil || height == 0 || height == math.MaxUint64 {
				logs.Error("listenChain - cannot get chain %s height, err: %s", ccl.handle.GetChainName(), err)
				continue
			}
			extendHeight, err := ccl.handle.GetExtendLatestHeight()
			if err != nil || extendHeight == 0 {
				logs.Error("ListenChain - cannot get chain %s extend height, err: %s", ccl.handle.GetChainName(), err)
			} else if extendHeight >= height+21 {
				logs.Error("ListenChain - chain %s node is too slow, node height: %d, really height: %d", ccl.handle.GetChainName(), height, extendHeight)
			}
			tools.Record(height, "%v.lastest_height", chain.ChainId)
			tools.Record(extendHeight, "%v.watch_height", chain.ChainId)
			tools.Record(chain.Height, "%v.height", chain.ChainId)
			if chain.Height >= height-ccl.handle.GetDefer() {
				continue
			}
			logs.Info("ListenChain - chain %s latest height is %d, listen height: %d", ccl.handle.GetChainName(), height, chain.Height)
			for chain.Height < height-ccl.handle.GetDefer() {
				wrapperTransactions, srcTransactions, polyTransactions, dstTransactions, err := ccl.HandleNewBlock(chain.Height + 1)
				if err != nil {
					logs.Error("HandleNewBlock %d err: %v", chain.Height+1, err)
					break
				}
				chain.Height += 1
				err = ccl.db.UpdateEvents(chain, wrapperTransactions, srcTransactions, polyTransactions, dstTransactions)
				if err != nil {
					logs.Error("UpdateEvents on block %d err: %v", chain.Height+1, err)
					chain.Height -= 1
					break
				}
				tools.Record(len(srcTransactions), "%v.locks", chain.ChainId)
				tools.Record(len(dstTransactions), "%v.unlocks", chain.ChainId)
			}
		case <-ccl.exit:
			logs.Info("cross chain listen exit, chain: %s, dao: %s......", ccl.handle.GetChainName(), ccl.db.Name())
			return true
		}
	}
}
