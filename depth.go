package marketdata

import (
	"fmt"
	"math"
)

const Int32NullValue int32 = math.MaxInt32

type Comparator func(a Price, b Price) int

// The Comparator functions are used in the context of comparing prices for the
// bid levels and offer levels. The offer prices should be compared in
// ascending order and the bid prices should be compared in descending order.
// Another way to think of the comparator functions is that compare* function
// should return -1 if `a` is closer to the inside market than `b`

// Compares prices in ascending order
//
// Returns -1 if a < b, 0 if a == b, and 1 if a > b
func compareAsc(a Price, b Price) int {
	if a.Mantissa < b.Mantissa {
		return -1
	} else if a.Mantissa > b.Mantissa {
		return 1
	} else {
		return 0
	}
}

// Compares prices in descending order
//
// Returns -1 if a > b, 0 if a == b, and 1 if a < b
func compareDesc(a Price, b Price) int {
	return compareAsc(b, a)
}

// Entry of a price level
type PriceLevelEntry struct {
	Price          Price
	Quantity       Quantity
	NumberOfOrders int32
}

// Copies the field values
func (p *PriceLevelEntry) Copy(that *PriceLevelEntry) {
	p.Price.Mantissa = that.Price.Mantissa
	p.Quantity.Mantissa = that.Quantity.Mantissa
	p.NumberOfOrders = that.NumberOfOrders
}

type LevelResult struct {
	Action     MDUpdateAction
	LevelIndex int32
}

// Container of price level entries backed by an array (slice to be more specific)
type LevelEntryArray struct {
	levels []PriceLevelEntry
	size   int
	cmp    Comparator
}

// Creates a new level entry array instance
func NewLevelEntryArray(capacity int, cmp Comparator) *LevelEntryArray {
	levels := LevelEntryArray{}
	levels.levels = make([]PriceLevelEntry, 0, capacity)
	for i := 0; i < capacity; i++ {

		var px Price
		PriceInit(&px)

		var qty Quantity
		QuantityInit(&qty)

		entry := PriceLevelEntry{Price: px, Quantity: qty, NumberOfOrders: Int32NullValue}

		levels.levels = append(levels.levels, entry)
	}
	levels.size = 0
	levels.cmp = cmp
	return &levels
}

// Gets the size of the levels
func (l *LevelEntryArray) Size() int {
	return l.size
}

// Gets the price level entry at the specified index
func (l *LevelEntryArray) Get(i int) PriceLevelEntry {
	return l.levels[i]
}

// Pushes back a market data entry
func (l *LevelEntryArray) PushBack(px Price, qty Quantity, numOrders int32) (LevelResult, error) {
	result := LevelResult{Action: NONE, LevelIndex: 0}

	if l.size >= cap(l.levels) {
		return result, fmt.Errorf("Buffer is full")
	}

	pxEntry := &l.levels[l.size]
	pxEntry.Price = px
	pxEntry.Quantity = qty
	pxEntry.NumberOfOrders = numOrders

	// increment size when a new level is added
	l.size++

	// push back is always a new price level
	result.Action = NEW
	result.LevelIndex = int32(l.size)

	return result, nil
}

// Deletes a price level
func (l *LevelEntryArray) Delete(px Price) (LevelResult, error) {
	result := LevelResult{Action: NONE, LevelIndex: 0}

	for i := 0; i < l.size; i++ {
		if l.cmp(px, l.levels[i].Price) == 0 {
			// Found the price
			// Only need to copy elements to shift the array if the index to
			// delete is not the last element
			result.Action = DELETE
			result.LevelIndex = int32(i + 1)
			if i < l.size-1 {
				for j := i; j < l.size-1; j++ {
					l.levels[j].Copy(&l.levels[j+1])
				}
			}
			l.size--
			break
		} else if l.cmp(px, l.levels[i].Price) < 0 {
			break
		}
	}
	return result, nil
}

// Inserts or updates a price level
func (l *LevelEntryArray) InsertOrUpdate(px Price, qty Quantity, numOrders int32) (LevelResult, error) {
	result := LevelResult{Action: NONE, LevelIndex: 0}
	found := false
	for i := 0; i < l.size; i++ {
		pxCompare := l.cmp(px, l.levels[i].Price)
		if pxCompare < 0 { // Insert
			curIndex := cap(l.levels) - 2
			if l.size < cap(l.levels) {
				curIndex = l.size - 1
			}

			// Shift the array to the right to make room for the insert
			for j := curIndex; j >= i; j-- {
				l.levels[j+1].Copy(&l.levels[j])
			}

			// "Insert" the price level
			pxEntry := &l.levels[i]
			pxEntry.Price = px
			pxEntry.Quantity = qty
			pxEntry.NumberOfOrders = numOrders

			// Only increment the size if less than capacity
			if l.size < cap(l.levels) {
				l.size++
			}

			found = true
			result.Action = NEW
			result.LevelIndex = int32(i) + 1
			break
		} else if pxCompare == 0 { // Update
			pxEntry := &l.levels[i]
			pxEntry.Quantity = qty
			pxEntry.NumberOfOrders = numOrders
			found = true
			result.Action = CHANGE
			result.LevelIndex = int32(i) + 1
			break
		}
	}

	if !found && l.size < cap(l.levels) {
		pxEntry := &l.levels[l.size]
		pxEntry.Price = px
		pxEntry.Quantity = qty
		pxEntry.NumberOfOrders = numOrders
		l.size++
		result.Action = NEW
		result.LevelIndex = int32(l.size)
	}
	return result, nil
}

// Clears all entries
//
// Does not actually delete any objects. Just sets size to 0 to avoid
// deallocations.
func (l *LevelEntryArray) Clear() {
	l.size = 0
}
