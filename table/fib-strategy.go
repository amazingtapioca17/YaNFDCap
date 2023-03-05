/* YaNFD - Yet another NDN Forwarding Daemon
 *
 * Copyright (C) 2020-2021 Eric Newberry.
 *
 * This file is licensed under the terms of the MIT License, as found in LICENSE.md.
 */

package table

import (
	"github.com/named-data/YaNFD/ndn"
	enc "github.com/zjkmxy/go-ndn/pkg/encoding"
)

// FibStrategyEntry represents an entry in the FIB-Strategy table.
type FibStrategyEntry interface {
	Name() *ndn.Name
	GetStrategy() *ndn.Name
	GetNextHops() []*FibNextHopEntry
}

// baseFibStrategyEntry represents information that all
// FibStrategyEntry implementations should include.
type baseFibStrategyEntry struct {
	component   ndn.NameComponent
	ppcomponent enc.Component
	name        *ndn.Name
	ppname      *enc.Name
	nexthops    []*FibNextHopEntry
	strategy    *ndn.Name
	ppstrategy  *enc.Name
}

// FibNextHopEntry represents a nexthop in a FIB entry.
type FibNextHopEntry struct {
	Nexthop uint64
	Cost    uint64
}

// FibStrategy represents the functionality that a FIB-strategy table should implement.
type FibStrategy interface {
	FindNextHops(name *ndn.Name) []*FibNextHopEntry
	FindNextHops1(name *enc.Name) []*FibNextHopEntry
	FindStrategy(name *ndn.Name) *ndn.Name
	FindStrategy1(name *enc.Name) *enc.Name
	InsertNextHop(name *ndn.Name, nextHop uint64, cost uint64)
	InsertNextHop1(name *enc.Name, nextHop uint64, cost uint64)
	ClearNextHops(name *ndn.Name)
	ClearNextHops1(name *enc.Name)
	RemoveNextHop(name *ndn.Name, nextHop uint64)
	RemoveNextHop1(name *enc.Name, nextHop uint64)

	GetAllFIBEntries() []FibStrategyEntry

	SetStrategy(name *ndn.Name, strategy *ndn.Name)
	UnsetStrategy(name *ndn.Name)
	SetStrategy1(name *enc.Name, strategy *enc.Name)
	UnsetStrategy1(name *enc.Name)
	GetAllForwardingStrategies() []FibStrategyEntry
}

// FibStrategy is a table containing FIB and Strategy entries for given prefixes.
var FibStrategyTable FibStrategy

// Name returns the name associated with the baseFibStrategyEntry.
func (e *baseFibStrategyEntry) Name() *ndn.Name {
	return e.name
}

// GetStrategy returns the strategy associated with the baseFibStrategyEntry.
func (e *baseFibStrategyEntry) GetStrategy() *ndn.Name {
	return e.strategy
}

// GetNexthops gets the nexthops of the specified entry.
func (e *baseFibStrategyEntry) GetNextHops() []*FibNextHopEntry {
	return e.nexthops
}
