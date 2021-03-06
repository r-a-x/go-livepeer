package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Segment struct {
	StreamID              string
	SegmentSequenceNumber *big.Int
	DataHash              string
}

func (s *Segment) Hash() common.Hash {
	return crypto.Keccak256Hash([]byte(s.StreamID), common.LeftPadBytes(s.SegmentSequenceNumber.Bytes(), 32), []byte(s.DataHash))
}
