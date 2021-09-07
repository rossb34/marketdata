package marketdata

import (
	"testing"
)

func TestMarketByPriceBook_OnSnapshot(t *testing.T) {
	type fields struct {
		Bids   *LevelEntryArray
		Offers *LevelEntryArray
	}
	type args struct {
		snapshot MDSnapshotFullRefresh
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantSize       int
		wantLevelIndex []int32
	}{
		{
			name:   "Top of book snapshot",
			fields: fields{Bids: NewLevelEntryArray(5, compareDesc), Offers: NewLevelEntryArray(5, compareAsc)},
			args: args{snapshot: MDSnapshotFullRefresh{Entries: []MDEntry{
				{NEW, BID, "FOO", 1, getPrice("9"), getQuantity("2"), 1, 0},
				{NEW, OFFER, "FOO", 1, getPrice("10"), getQuantity("1"), 1, 0},
			}}},
			wantSize:       1,
			wantLevelIndex: []int32{1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MarketByPriceBook{
				Bids:   tt.fields.Bids,
				Offers: tt.fields.Offers,
			}
			m.OnSnapshot(&tt.args.snapshot)

			if m.Offers.Size() != tt.wantSize {
				t.Errorf("m.Offers.Size() = %v, want %v", m.Offers.Size(), tt.wantSize)
			}

			if m.Bids.Size() != tt.wantSize {
				t.Errorf("m.Bids.Size() = %v, want %v", m.Bids.Size(), tt.wantSize)
			}

			for i, v := range tt.args.snapshot.Entries {
				if v.PriceLevelIndex != tt.wantLevelIndex[i] {
					t.Errorf("entry.PriceLevelIndex = %v, want %v", v.PriceLevelIndex, tt.wantLevelIndex[i])
				}
			}

		})
	}
}

func TestMarketByPriceBook_OnIncrementalUpdate(t *testing.T) {
	type fields struct {
		Bids   *LevelEntryArray
		Offers *LevelEntryArray
	}
	type args struct {
		incremental MDIncrementalRefresh
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "Bid delete",
			fields: fields{
				Bids: &LevelEntryArray{levels: []PriceLevelEntry{
					{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
					{Price: getPrice("9"), Quantity: getQuantity("2"), NumberOfOrders: 1},
				}},
				Offers: &LevelEntryArray{levels: []PriceLevelEntry{
					{Price: getPrice("12"), Quantity: getQuantity("1"), NumberOfOrders: 1},
					{Price: getPrice("14"), Quantity: getQuantity("2"), NumberOfOrders: 1},
				}},
			},
			args: args{incremental: MDIncrementalRefresh{Entries: []MDEntry{{Action: DELETE, Type: BID, Symbol: "FOO", RptSequenceNumber: 1, Price: getPrice("10"), Size: getQuantity("0"), NumberOfOrders: 0, PriceLevelIndex: 0}}}},
		},
		{
			name: "Offer delete",
			fields: fields{
				Bids: &LevelEntryArray{levels: []PriceLevelEntry{
					{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
					{Price: getPrice("9"), Quantity: getQuantity("2"), NumberOfOrders: 1},
				}},
				Offers: &LevelEntryArray{levels: []PriceLevelEntry{
					{Price: getPrice("12"), Quantity: getQuantity("1"), NumberOfOrders: 1},
					{Price: getPrice("14"), Quantity: getQuantity("2"), NumberOfOrders: 1},
				}},
			},
			args: args{incremental: MDIncrementalRefresh{Entries: []MDEntry{{Action: DELETE, Type: BID, Symbol: "FOO", RptSequenceNumber: 1, Price: getPrice("12"), Size: getQuantity("0"), NumberOfOrders: 0, PriceLevelIndex: 0}}}},
		},
		{
			name: "Bid insert",
			fields: fields{
				Bids: &LevelEntryArray{levels: []PriceLevelEntry{
					{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
					{Price: getPrice("9"), Quantity: getQuantity("2"), NumberOfOrders: 1},
				}},
				Offers: &LevelEntryArray{levels: []PriceLevelEntry{
					{Price: getPrice("12"), Quantity: getQuantity("1"), NumberOfOrders: 1},
					{Price: getPrice("14"), Quantity: getQuantity("2"), NumberOfOrders: 1},
				}},
			},
			args: args{incremental: MDIncrementalRefresh{Entries: []MDEntry{{Action: NEW, Type: BID, Symbol: "FOO", RptSequenceNumber: 1, Price: getPrice("11"), Size: getQuantity("3"), NumberOfOrders: 1, PriceLevelIndex: 0}}}},
		},
		{
			name: "Offer insert",
			fields: fields{
				Bids: &LevelEntryArray{levels: []PriceLevelEntry{
					{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
					{Price: getPrice("9"), Quantity: getQuantity("2"), NumberOfOrders: 1},
				}},
				Offers: &LevelEntryArray{levels: []PriceLevelEntry{
					{Price: getPrice("12"), Quantity: getQuantity("1"), NumberOfOrders: 1},
					{Price: getPrice("14"), Quantity: getQuantity("2"), NumberOfOrders: 1},
				}},
			},
			args: args{incremental: MDIncrementalRefresh{Entries: []MDEntry{{Action: NEW, Type: BID, Symbol: "FOO", RptSequenceNumber: 1, Price: getPrice("13"), Size: getQuantity("3"), NumberOfOrders: 1, PriceLevelIndex: 0}}}},
		},
		{
			name: "Bid update",
			fields: fields{
				Bids: &LevelEntryArray{levels: []PriceLevelEntry{
					{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
					{Price: getPrice("9"), Quantity: getQuantity("2"), NumberOfOrders: 1},
				}},
				Offers: &LevelEntryArray{levels: []PriceLevelEntry{
					{Price: getPrice("12"), Quantity: getQuantity("1"), NumberOfOrders: 1},
					{Price: getPrice("14"), Quantity: getQuantity("2"), NumberOfOrders: 1},
				}},
			},
			args: args{incremental: MDIncrementalRefresh{Entries: []MDEntry{{Action: NEW, Type: BID, Symbol: "FOO", RptSequenceNumber: 1, Price: getPrice("10"), Size: getQuantity("3"), NumberOfOrders: 2, PriceLevelIndex: 0}}}},
		},
		{
			name: "Offer update",
			fields: fields{
				Bids: &LevelEntryArray{levels: []PriceLevelEntry{
					{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
					{Price: getPrice("9"), Quantity: getQuantity("2"), NumberOfOrders: 1},
				}},
				Offers: &LevelEntryArray{levels: []PriceLevelEntry{
					{Price: getPrice("12"), Quantity: getQuantity("1"), NumberOfOrders: 1},
					{Price: getPrice("14"), Quantity: getQuantity("2"), NumberOfOrders: 1},
				}},
			},
			args: args{incremental: MDIncrementalRefresh{Entries: []MDEntry{{Action: NEW, Type: BID, Symbol: "FOO", RptSequenceNumber: 1, Price: getPrice("14"), Size: getQuantity("3"), NumberOfOrders: 2, PriceLevelIndex: 0}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MarketByPriceBook{
				Bids:   tt.fields.Bids,
				Offers: tt.fields.Offers,
			}
			m.OnIncrementalUpdate(&tt.args.incremental)
		})
	}
}
