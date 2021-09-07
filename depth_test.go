package marketdata

import (
	"reflect"
	"testing"
)

func getPrice(px string) Price {
	v, _ := ParsePrice(px)
	return v
}

func getQuantity(qty string) Quantity {
	v, _ := ParseQuantity(qty)
	return v
}

func createEmptyPriceLevel() PriceLevelEntry {
	var px Price
	PriceInit(&px)

	var qty Quantity
	QuantityInit(&qty)

	return PriceLevelEntry{Price: px, Quantity: qty, NumberOfOrders: Int32NullValue}
}

func TestLevelEntryArray_PushBack(t *testing.T) {
	type fields struct {
		buffer []PriceLevelEntry
		size   int
		cmp    Comparator
	}
	type args struct {
		Price          Price
		Quantity       Quantity
		NumberOfOrders int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    LevelResult
		wantErr bool
	}{
		{
			name: "Push Back to an empty container",
			fields: fields{buffer: []PriceLevelEntry{
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
			},
				size: 0,
				cmp:  compareAsc,
			},
			args:    args{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			want:    LevelResult{Action: NEW, LevelIndex: 1},
			wantErr: false,
		},
		{
			name: "Push Back to container with an entry",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
			},
				size: 1,
				cmp:  compareAsc,
			},
			args:    args{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			want:    LevelResult{Action: NEW, LevelIndex: 2},
			wantErr: false,
		},
		{
			name: "Push Back to full container",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("11"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("12"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("13"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
				size: 5,
				cmp:  compareAsc,
			},
			args:    args{Price: getPrice("15"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			want:    LevelResult{Action: NONE, LevelIndex: 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LevelEntryArray{
				levels: tt.fields.buffer,
				size:   tt.fields.size,
				cmp:    tt.fields.cmp,
			}
			got, err := l.PushBack(tt.args.Price, tt.args.Quantity, tt.args.NumberOfOrders)
			if (err != nil) != tt.wantErr {
				t.Errorf("LevelEntryArray.PushBack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LevelEntryArray.PushBack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelEntryArray_Delete(t *testing.T) {
	type fields struct {
		buffer []PriceLevelEntry
		size   int
		cmp    Comparator
	}
	type args struct {
		Price Price
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       LevelResult
		wantErr    bool
		wantLevels []PriceLevelEntry
	}{
		{
			name: "Offers: Delete front",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("11"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("14"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("15"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
				size: 5,
				cmp:  compareAsc,
			},
			args:    args{Price: getPrice("9")},
			want:    LevelResult{Action: DELETE, LevelIndex: 1},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("11"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("14"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("15"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
		},
		{
			name: "Offers: Delete middle",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("11"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("14"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("15"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
				size: 5,
				cmp:  compareAsc,
			},
			args:    args{Price: getPrice("11")},
			want:    LevelResult{Action: DELETE, LevelIndex: 3},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("14"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("15"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
		},
		{
			name: "Offers: Delete back",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("11"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("14"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("15"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
				size: 5,
				cmp:  compareAsc,
			},
			args:    args{Price: getPrice("15")},
			want:    LevelResult{Action: DELETE, LevelIndex: 5},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("11"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("14"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
		},
		{
			name: "Offers: Delete price does not exist",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("11"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("14"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("15"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
				size: 5,
				cmp:  compareAsc,
			},
			args:    args{Price: getPrice("13")},
			want:    LevelResult{Action: NONE, LevelIndex: 0},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("11"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("14"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("15"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
		},
		{
			name: "Bids: Delete front",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("8"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("6"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("5"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("2"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
				size: 5,
				cmp:  compareDesc,
			},
			args:    args{Price: getPrice("9")},
			want:    LevelResult{Action: DELETE, LevelIndex: 1},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("8"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("6"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("5"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("2"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
		},
		{
			name: "Bids: Delete middle",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("8"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("6"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("5"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("2"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
				size: 5,
				cmp:  compareDesc,
			},
			args:    args{Price: getPrice("5")},
			want:    LevelResult{Action: DELETE, LevelIndex: 4},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("8"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("6"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("2"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
		},
		{
			name: "Bids: Delete back",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("8"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("6"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("5"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("2"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
				size: 5,
				cmp:  compareDesc,
			},
			args:    args{Price: getPrice("2")},
			want:    LevelResult{Action: DELETE, LevelIndex: 5},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("8"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("6"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("5"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
		},
		{
			name: "Bids: Delete price does not exist",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("8"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("6"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("5"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("2"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
				size: 5,
				cmp:  compareDesc,
			},
			args:    args{Price: getPrice("3")},
			want:    LevelResult{Action: NONE, LevelIndex: 0},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("8"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("6"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("5"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("2"), Quantity: getQuantity("1"), NumberOfOrders: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LevelEntryArray{
				levels: tt.fields.buffer,
				size:   tt.fields.size,
				cmp:    tt.fields.cmp,
			}
			got, err := l.Delete(tt.args.Price)
			if (err != nil) != tt.wantErr {
				t.Errorf("LevelEntryArray.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LevelEntryArray.Delete() = %v, want %v", got, tt.want)
			}
			gotLevels := l.levels
			if !reflect.DeepEqual(gotLevels[:l.size], tt.wantLevels) {
				t.Errorf("buffer = %v, want %v", got, tt.wantLevels)
			}

			if l.size != len(tt.wantLevels) {
				t.Errorf("l.size = %v, want %v", l.size, len(tt.wantLevels))
			}
		})
	}
}

func TestLevelEntryArray_InsertOrUpdate(t *testing.T) {
	type fields struct {
		buffer []PriceLevelEntry
		size   int
		cmp    Comparator
	}
	type args struct {
		Price          Price
		Quantity       Quantity
		NumberOfOrders int32
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       LevelResult
		wantErr    bool
		wantLevels []PriceLevelEntry
	}{
		// Offers
		{
			name: "Offers: Insert front",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("12"), Quantity: getQuantity("2"), NumberOfOrders: 2},
				// Empty levels
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
			},
				size: 2,
				cmp:  compareAsc,
			},
			args:    args{Price: getPrice("9"), Quantity: getQuantity("3"), NumberOfOrders: 1},
			want:    LevelResult{Action: NEW, LevelIndex: 1},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("9"), Quantity: getQuantity("3"), NumberOfOrders: 1},
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("12"), Quantity: getQuantity("2"), NumberOfOrders: 2},
			},
		},
		{
			name: "Offers: Insert middle",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("12"), Quantity: getQuantity("2"), NumberOfOrders: 2},
				// Empty levels
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
			},
				size: 2,
				cmp:  compareAsc,
			},
			args:    args{Price: getPrice("11"), Quantity: getQuantity("3"), NumberOfOrders: 1},
			want:    LevelResult{Action: NEW, LevelIndex: 2},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("11"), Quantity: getQuantity("3"), NumberOfOrders: 1},
				{Price: getPrice("12"), Quantity: getQuantity("2"), NumberOfOrders: 2},
			},
		},
		{
			name: "Offers: Insert back",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("12"), Quantity: getQuantity("2"), NumberOfOrders: 2},
				// Empty levels
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
			},
				size: 2,
				cmp:  compareAsc,
			},
			args:    args{Price: getPrice("15"), Quantity: getQuantity("3"), NumberOfOrders: 1},
			want:    LevelResult{Action: NEW, LevelIndex: 3},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("12"), Quantity: getQuantity("2"), NumberOfOrders: 2},
				{Price: getPrice("15"), Quantity: getQuantity("3"), NumberOfOrders: 1},
			},
		},
		{
			name: "Offers: Update level",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("12"), Quantity: getQuantity("2"), NumberOfOrders: 3},
				// Empty levels
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
			},
				size: 2,
				cmp:  compareAsc,
			},
			args:    args{Price: getPrice("12"), Quantity: getQuantity("5"), NumberOfOrders: 3},
			want:    LevelResult{Action: CHANGE, LevelIndex: 2},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("12"), Quantity: getQuantity("5"), NumberOfOrders: 3},
			},
		},
		// Bids
		{
			name: "Bids: Insert front",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("9"), Quantity: getQuantity("2"), NumberOfOrders: 2},
				// Empty levels
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
			},
				size: 2,
				cmp:  compareDesc,
			},
			args:    args{Price: getPrice("12"), Quantity: getQuantity("3"), NumberOfOrders: 1},
			want:    LevelResult{Action: NEW, LevelIndex: 1},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("12"), Quantity: getQuantity("3"), NumberOfOrders: 1},
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("9"), Quantity: getQuantity("2"), NumberOfOrders: 2},
			},
		},
		{
			name: "Bids: Insert middle",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("7"), Quantity: getQuantity("2"), NumberOfOrders: 2},
				// Empty levels
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
			},
				size: 2,
				cmp:  compareDesc,
			},
			args:    args{Price: getPrice("8"), Quantity: getQuantity("3"), NumberOfOrders: 1},
			want:    LevelResult{Action: NEW, LevelIndex: 2},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("8"), Quantity: getQuantity("3"), NumberOfOrders: 1},
				{Price: getPrice("7"), Quantity: getQuantity("2"), NumberOfOrders: 2},
			},
		},
		{
			name: "Bids: Insert back",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("7"), Quantity: getQuantity("2"), NumberOfOrders: 2},
				// Empty levels
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
			},
				size: 2,
				cmp:  compareDesc,
			},
			args:    args{Price: getPrice("6"), Quantity: getQuantity("3"), NumberOfOrders: 1},
			want:    LevelResult{Action: NEW, LevelIndex: 3},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("7"), Quantity: getQuantity("2"), NumberOfOrders: 2},
				{Price: getPrice("6"), Quantity: getQuantity("3"), NumberOfOrders: 1},
			},
		},
		{
			name: "Bids: Update level",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("7"), Quantity: getQuantity("2"), NumberOfOrders: 2},
				// Empty levels
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
			},
				size: 2,
				cmp:  compareDesc,
			},
			args:    args{Price: getPrice("7"), Quantity: getQuantity("4"), NumberOfOrders: 4},
			want:    LevelResult{Action: CHANGE, LevelIndex: 2},
			wantErr: false,
			wantLevels: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("7"), Quantity: getQuantity("4"), NumberOfOrders: 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LevelEntryArray{
				levels: tt.fields.buffer,
				size:   tt.fields.size,
				cmp:    tt.fields.cmp,
			}
			got, err := l.InsertOrUpdate(tt.args.Price, tt.args.Quantity, tt.args.NumberOfOrders)
			if (err != nil) != tt.wantErr {
				t.Errorf("LevelEntryArray.InsertOrUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LevelEntryArray.InsertOrUpdate() = %v, want %v", got, tt.want)
			}
			gotLevels := l.levels
			if !reflect.DeepEqual(gotLevels[:l.size], tt.wantLevels) {
				t.Errorf("buffer = %v, want %v", got, tt.wantLevels)
			}

			if l.size != len(tt.wantLevels) {
				t.Errorf("l.size = %v, want %v", l.size, len(tt.wantLevels))
			}
		})
	}
}

func TestLevelEntryArray_Clear(t *testing.T) {
	type fields struct {
		buffer []PriceLevelEntry
		size   int
		cmp    Comparator
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Clear",
			fields: fields{buffer: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("7"), Quantity: getQuantity("2"), NumberOfOrders: 2},
			},
				size: 2,
				cmp:  compareDesc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LevelEntryArray{
				levels: tt.fields.buffer,
				size:   tt.fields.size,
				cmp:    tt.fields.cmp,
			}
			l.Clear()
			if l.size != 0 {
				t.Errorf("l.size = %v, want %v", l.size, 0)
			}
		})
	}
}

func Test_compareAsc(t *testing.T) {
	type args struct {
		a Price
		b Price
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "compare(9, 10)",
			args: args{a: getPrice("9"), b: getPrice("10")},
			want: -1,
		},
		{
			name: "compare(10, 10)",
			args: args{a: getPrice("10"), b: getPrice("10")},
			want: 0,
		},
		{
			name: "compare(11, 10)",
			args: args{a: getPrice("11"), b: getPrice("10")},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareAsc(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("compareAsc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compareDesc(t *testing.T) {
	type args struct {
		a Price
		b Price
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "compare(9, 10)",
			args: args{a: getPrice("9"), b: getPrice("10")},
			want: 1,
		},
		{
			name: "compare(10, 10)",
			args: args{a: getPrice("10"), b: getPrice("10")},
			want: 0,
		},
		{
			name: "compare(11, 10)",
			args: args{a: getPrice("11"), b: getPrice("10")},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareDesc(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("compareDesc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testNewLevelEntryArray(t *testing.T) {
	type args struct {
		capacity int
		cmp      Comparator
	}
	tests := []struct {
		name string
		args args
		want *LevelEntryArray
	}{
		{
			name: "NewLevelEntryArray init",
			args: args{capacity: 2, cmp: compareAsc},
			want: &LevelEntryArray{levels: []PriceLevelEntry{
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
			}, size: 0, cmp: compareAsc},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLevelEntryArray(tt.args.capacity, tt.args.cmp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLevelEntryArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevelEntryArray_Size(t *testing.T) {
	type fields struct {
		levels []PriceLevelEntry
		size   int
		cmp    Comparator
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Init empty size",
			fields: fields{levels: []PriceLevelEntry{
				// Empty levels
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
			},
				size: 0,
				cmp:  compareAsc,
			},
			want: 0,
		},
		{
			name: "Init size",
			fields: fields{levels: []PriceLevelEntry{
				{Price: getPrice("10"), Quantity: getQuantity("1"), NumberOfOrders: 1},
				{Price: getPrice("12"), Quantity: getQuantity("2"), NumberOfOrders: 3},
				// Empty levels
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
				createEmptyPriceLevel(),
			},
				size: 2,
				cmp:  compareAsc,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LevelEntryArray{
				levels: tt.fields.levels,
				size:   tt.fields.size,
				cmp:    tt.fields.cmp,
			}
			if got := l.Size(); got != tt.want {
				t.Errorf("LevelEntryArray.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
