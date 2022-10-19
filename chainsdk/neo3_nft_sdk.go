package chainsdk

import (
	"math/big"
)

func (s *Neo3Sdk) GetAndCheckNFTUrl(queryAddr, asset, owner, tokenId string) (string, error) {
	tokenOwner, err := s.Nep11OwnerOf(asset, tokenId)
	if tokenOwner != owner {
		return "", err
	}
	tokenUrl, err := s.Nep11TokenUrl(asset, tokenId)
	return tokenUrl, nil
}

func (s *Neo3Sdk) GetNFTTokenUri(asset, tokenId string) (string, error) {
	tokenUrl, err := s.Nep11TokenUrl(asset, tokenId)
	if err != nil {
		return "", err
	}
	return tokenUrl, nil
}

func (s *Neo3Sdk) GetNFTBalance(asset, owner string) (*big.Int, error) {
	balance, err := s.Nep11BalanceOf(asset, owner)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (s *Neo3Sdk) GetOwnerNFTsByIndex(queryAddr, asset, owner string, start, length int) (map[string]string, error) {
	m := make(map[string]string)
	tokenIds, _ := s.Nep11TokensOf(asset, owner)
	for _, id := range tokenIds {
		url, _ := s.Nep11TokenUrl(asset, id)
		m[id] = url
	}
	return m, nil
}
