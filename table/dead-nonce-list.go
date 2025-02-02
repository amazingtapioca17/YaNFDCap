/* YaNFD - Yet another NDN Forwarding Daemon
 *
 * Copyright (C) 2020-2021 Eric Newberry.
 *
 * This file is licensed under the terms of the MIT License, as found in LICENSE.md.
 */

package table

import (
	"encoding/binary"
	"time"

	"github.com/cespare/xxhash"
	"github.com/named-data/YaNFD/ndn"
	"github.com/named-data/YaNFD/utils/priority_queue"
)

// DeadNonceList represents the Dead Nonce List for a forwarding thread.
type DeadNonceList struct {
	list            map[uint64]bool
	expirationQueue priority_queue.Queue[uint64, int64]
	Ticker          *time.Ticker
}

// NewDeadNonceList creates a new Dead Nonce List for a forwarding thread.
func NewDeadNonceList() *DeadNonceList {
	d := new(DeadNonceList)
	d.list = make(map[uint64]bool)
	d.Ticker = time.NewTicker(100 * time.Millisecond)
	d.expirationQueue = priority_queue.New[uint64, int64]()
	return d
}

// Find returns whether the specified name and nonce combination are present in the Dead Nonce List.
func (d *DeadNonceList) Find(name *ndn.Name, nonce []byte) bool {
	wire, err := name.Encode().Wire()
	if err != nil {
		return false
	}
	hash := xxhash.Sum64(wire) + uint64(binary.BigEndian.Uint32(nonce))
	_, ok := d.list[hash]
	return ok
}

// Insert inserts an entry in the Dead Nonce List with the specified name and nonce. Returns whether nonce already present.
func (d *DeadNonceList) Insert(name *ndn.Name, nonce []byte) bool {
	wire, err := name.Encode().Wire()
	if err != nil {
		return false
	}
	hash := xxhash.Sum64(wire) + uint64(binary.BigEndian.Uint32(nonce))
	_, exists := d.list[hash]

	if !exists {
		d.list[hash] = true
		d.expirationQueue.Push(hash, time.Now().Add(deadNonceListLifetime).UnixNano())
	}
	return exists
}

// RemoveExpiredEntry removes all expired entries from Dead Nonce List.
func (d *DeadNonceList) RemoveExpiredEntries() {
	evicted := 0
	for d.expirationQueue.Len() > 0 && d.expirationQueue.PeekPriority() < time.Now().UnixNano() {
		hash := d.expirationQueue.Pop()
		delete(d.list, hash)
		evicted += 1

		if evicted >= 100 {
			break
		}
	}
}
