//go:build devnet
// +build devnet

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

package basedef

const (
	POLY_CROSSCHAIN_ID       = uint64(0)
	ETHEREUM_CROSSCHAIN_ID   = uint64(2)
	ONT_CROSSCHAIN_ID        = uint64(3)
	NEO_CROSSCHAIN_ID        = uint64(4)
	BSC_CROSSCHAIN_ID        = uint64(6)
	HECO_CROSSCHAIN_ID       = uint64(7)
	O3_CROSSCHAIN_ID         = uint64(80)
	NEO3_CROSSCHAIN_ID       = uint64(88)
	OK_CROSSCHAIN_ID         = uint64(90)
	MATIC_CROSSCHAIN_ID      = uint64(13)
	PLT_CROSSCHAIN_ID        = uint64(107)
	ARBITRUM_CROSSCHAIN_ID   = uint64(215)
	XDAI_CROSSCHAIN_ID       = uint64(206)
	ZILLIQA_CROSSCHAIN_ID    = uint64(111)
	FANTOM_CROSSCHAIN_ID     = uint64(208)
	AVAX_CROSSCHAIN_ID       = uint64(209)
	OPTIMISTIC_CROSSCHAIN_ID = uint64(210)
	METIS_CROSSCHAIN_ID      = uint64(300)
	BOBA_CROSSCHAIN_ID       = uint64(400)
	RINKEBY_CROSSCHAIN_ID    = uint64(402)
	OASIS_CROSSCHAIN_ID      = uint64(500)
	HARMONY_CROSSCHAIN_ID    = uint64(800)
	KCC_CROSSCHAIN_ID        = uint64(900)
	HSC_CROSSCHAIN_ID        = uint64(603)
	STARCOIN_CROSSCHAIN_ID   = uint64(318)
	ZKSYNC_CROSSCHAIN_ID     = uint64(940)
	CELO_CROSSCHAIN_ID       = uint64(960)
	CLOVER_CROSSCHAIN_ID     = uint64(970)
	CONFLUX_CROSSCHAIN_ID    = uint64(980)
	RIPPLE_CROSSCHAIN_ID     = uint64(223)
	ASTAR_CROSSCHAIN_ID      = uint64(990)
	APTOS_CROSSCHAIN_ID      = uint64(998)
	BRISE_CROSSCHAIN_ID      = uint64(1010)
	DEXIT_CROSSCHAIN_ID      = uint64(1020)
	CLOUDTX_CROSSCHAIN_ID    = uint64(1030)

	ENV = "devnet"
)

const (
	BSC_NORMAL_GASPRICE   = 5000000000
	ASTAR_NORMAL_GASPRICE = 60000000000
)

var ETH_CHAINS = []uint64{
	ETHEREUM_CROSSCHAIN_ID, BSC_CROSSCHAIN_ID, HECO_CROSSCHAIN_ID, OK_CROSSCHAIN_ID, MATIC_CROSSCHAIN_ID,
	O3_CROSSCHAIN_ID, PLT_CROSSCHAIN_ID, PLT2_CROSSCHAIN_ID, ARBITRUM_CROSSCHAIN_ID, XDAI_CROSSCHAIN_ID, OPTIMISTIC_CROSSCHAIN_ID,
	FANTOM_CROSSCHAIN_ID, AVAX_CROSSCHAIN_ID, METIS_CROSSCHAIN_ID, BOBA_CROSSCHAIN_ID, RINKEBY_CROSSCHAIN_ID, GOERLI_CROSSCHAIN_ID,
	OASIS_CROSSCHAIN_ID, HARMONY_CROSSCHAIN_ID, KCC_CROSSCHAIN_ID, HSC_CROSSCHAIN_ID, ONTEVM_CROSSCHAIN_ID,
	MILKOMEDA_CROSSCHAIN_ID, ZKSYNC_CROSSCHAIN_ID, CELO_CROSSCHAIN_ID, CLOVER_CROSSCHAIN_ID, CONFLUX_CROSSCHAIN_ID, ASTAR_CROSSCHAIN_ID,
	BRISE_CROSSCHAIN_ID, DEXIT_CROSSCHAIN_ID, CLOUDTX_CROSSCHAIN_ID,
}
