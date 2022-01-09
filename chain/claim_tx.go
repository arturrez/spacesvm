// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package chain

import (
	"bytes"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/quarkvm/parser"
)

var _ UnsignedTransaction = &ClaimTx{}

type ClaimTx struct {
	*BaseTx `serialize:"true" json:"baseTx"`
}

func (c *ClaimTx) Execute(db database.Database, blockTime uint64) error {
	// Restrict address prefix to be owned by pk
	// [33]byte prefix is reserved for pubkey
	if len(c.Prefix) == crypto.SECP256K1RPKLen && !bytes.Equal(c.Sender[:], c.Prefix) {
		return ErrPublicKeyMismatch
	}

	// Prefix keys only exist if they are still valid
	exists, err := HasPrefix(db, c.Prefix)
	if err != nil {
		return err
	}
	if exists {
		return ErrPrefixNotExpired
	}

	// Anything previously at the prefix was previously removed...
	newInfo := &PrefixInfo{
		Owner:       c.Sender,
		Created:     blockTime,
		LastUpdated: blockTime,
		Expiry:      blockTime + ExpiryTime,
		Units:       1,
	}
	if err := PutPrefixInfo(db, c.Prefix, newInfo, 0); err != nil {
		return err
	}
	return nil
}

// [prefixUnits] requires the caller to produce more work to get prefixes of
// a shorter length because they are more desirable. This creates a "lottery"
// mechanism where the people that spend the most mining power will win the
// prefix.
//
// [prefixUnits] should only be called on a prefix that is valid
func prefixUnits(p []byte) uint64 {
	desirability := parser.MaxKeySize - len(p)
	if len(p) > ClaimTier2Size {
		return uint64(desirability * ClaimTier3Multiplier)
	}
	if len(p) > ClaimTier1Size {
		return uint64(desirability * ClaimTier2Multiplier)
	}
	return uint64(desirability * ClaimTier1Multiplier)
}

func (c *ClaimTx) FeeUnits() uint64 {
	return c.LoadUnits() + prefixUnits(c.Prefix)
}

func (c *ClaimTx) LoadUnits() uint64 {
	return c.BaseTx.LoadUnits() * ClaimFeeMultiplier
}
