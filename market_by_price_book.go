package marketdata

// Very simple and clean data model for message structures to update an order
// book. The idea is that each endpoint (e.g. exchange or liquidity provider)
// has their own idiosyncracies when it comes to publishing market data. The
// structs defined here should contain only the minimal core data structures
// required for maintaining a market by price book.

type MarketByPriceBook struct {
	Bids   *LevelEntryArray
	Offers *LevelEntryArray
}

// Allocates a new market by price book
func NewMarketByPriceBook(depth int) *MarketByPriceBook {
	mbp := MarketByPriceBook{}
	mbp.Bids = NewLevelEntryArray(depth, compareDesc)
	mbp.Offers = NewLevelEntryArray(depth, compareAsc)
	return &mbp
}

// Clears the book
func (m *MarketByPriceBook) Clear() {
	m.Bids.Clear()
	m.Offers.Clear()
}


// Processes market data snapshot full refresh message and updates the book
func (m *MarketByPriceBook) OnSnapshot(snapshot *MDSnapshotFullRefresh) {
	// TODO: update state from snapshot message properties
	m.Clear()
	for i := 0; i < len(snapshot.Entries); i++ {
		e := &snapshot.Entries[i]
		if e.Type == BID {
			r, _ := m.Bids.PushBack(e.Price, e.Size, e.NumberOfOrders)
			e.Action = r.Action
			e.PriceLevelIndex = r.LevelIndex
		} else if e.Type == OFFER {
			r, _ := m.Offers.PushBack(e.Price, e.Size, e.NumberOfOrders)
			e.Action = r.Action
			e.PriceLevelIndex = r.LevelIndex
		}
	}
}


// Processes market data incremental refresh message and updates the book
func (m *MarketByPriceBook) OnIncrementalUpdate(incremental *MDIncrementalRefresh) {
	// TODO: update state from incremental update message properties (e.g. )
	for i := 0; i < len(incremental.Entries); i++ {
		e := &incremental.Entries[i]
		m.Update(e)
	}
}

// Processes a single market data entry and updates the book
func (m *MarketByPriceBook) Update(entry *MDEntry) {
	switch entry.Type {
	case BID:
		if entry.Action == DELETE {
			r, _ := m.Bids.Delete(entry.Price)
			entry.PriceLevelIndex = r.LevelIndex
		} else {
			r, _ := m.Bids.InsertOrUpdate(entry.Price, entry.Size, int32(entry.NumberOfOrders))
			entry.Action = r.Action
			entry.PriceLevelIndex = r.LevelIndex
		}
	case OFFER:
		if entry.Action == DELETE {
			r, _ := m.Offers.Delete(entry.Price)
			entry.PriceLevelIndex = r.LevelIndex
		} else {
			r, _ := m.Offers.InsertOrUpdate(entry.Price, entry.Size, int32(entry.NumberOfOrders))
			entry.Action = r.Action
			entry.PriceLevelIndex = r.LevelIndex
		}
	}
}
