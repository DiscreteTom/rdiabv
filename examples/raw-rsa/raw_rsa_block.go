package main

import (
	"math/big"

	"github.com/DiscreteTom/rawrsa"
	"github.com/DiscreteTom/rdiabv"
)

// RawRsaBlock use big.Int to store data & tag.
type RawRsaBlock struct {
	Data *big.Int
	Tag  *big.Int
	Key  *rawrsa.RawRsa
}

// NewRawRsaBlock will generate a new RawRsaBlock with data 1 and tag 1.
func NewRawRsaBlock(key *rawrsa.RawRsa) (block *RawRsaBlock) {
	block = &RawRsaBlock{}
	block.Data = new(big.Int).SetInt64(1)
	block.Tag = new(big.Int).SetInt64(1)
	block.Key = key
	return
}

// Copy will create a copy of the current block.
func (block *RawRsaBlock) Copy() rdiabv.Block {
	copied := &RawRsaBlock{}
	copied.Data = new(big.Int).Set(block.Data)
	copied.Tag = new(big.Int).Set(block.Tag)
	copied.Key = block.Key
	return copied
}

// Validate will check whether the encrypted data equals to the tag.
func (block *RawRsaBlock) Validate() (ret bool) {
	return block.Data.Cmp(block.Key.RawDecrypt(block.Tag)) == 0
}

// Merge will add fields of x to the fields of y.
func (block *RawRsaBlock) Merge(x rdiabv.Block, y rdiabv.Block) rdiabv.Block {
	// type assertion
	blockX, _ := x.(*RawRsaBlock)
	blockY, _ := y.(*RawRsaBlock)
	// m1^e * m2^e mod N == (m1*m2)^e mod N, which means (data1 * data2) mod N matches (tag1 * tag2) mod N
	block.Data.Mul(blockX.Data, blockY.Data)
	block.Data.Mod(block.Data, block.Key.N)
	block.Tag.Mul(blockX.Tag, blockY.Tag)
	block.Tag.Mod(block.Tag, block.Key.N)
	return block
}
