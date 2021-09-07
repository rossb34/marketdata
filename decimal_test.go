package marketdata

import (
	"math"
	"reflect"
	"testing"
)

func TestParseDecimal(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want Decimal
	}{
		{
			name: "0.12345678 to decimal",
			args: args{"0.12345678"},
			want: Decimal{Mantissa: 12345678, Exponent: -8},
		},
		{
			name: "-0.12345678 to decimal",
			args: args{"-0.12345678"},
			want: Decimal{Mantissa: -12345678, Exponent: -8},
		},
		{
			name: "0.1234567 to decimal",
			args: args{"0.1234567"},
			want: Decimal{Mantissa: 1234567, Exponent: -7},
		},
		{
			name: "-0.1234567 to decimal",
			args: args{"-0.1234567"},
			want: Decimal{Mantissa: -1234567, Exponent: -7},
		},
		{
			name: "0.123456 to decimal",
			args: args{"0.123456"},
			want: Decimal{Mantissa: 123456, Exponent: -6},
		},
		{
			name: "-0.123456 to decimal",
			args: args{"-0.123456"},
			want: Decimal{Mantissa: -123456, Exponent: -6},
		},
		{
			name: "0.12345 to decimal",
			args: args{"0.12345"},
			want: Decimal{Mantissa: 12345, Exponent: -5},
		},
		{
			name: "-0.12345 to decimal",
			args: args{"-0.12345"},
			want: Decimal{Mantissa: -12345, Exponent: -5},
		},
		{
			name: "0.1234 to decimal",
			args: args{"0.1234"},
			want: Decimal{Mantissa: 1234, Exponent: -4},
		},
		{
			name: "-0.1234 to decimal",
			args: args{"-0.1234"},
			want: Decimal{Mantissa: -1234, Exponent: -4},
		},
		{
			name: "0.123 to decimal",
			args: args{"0.123"},
			want: Decimal{Mantissa: 123, Exponent: -3},
		},
		{
			name: "-0.123 to decimal",
			args: args{"-0.123"},
			want: Decimal{Mantissa: -123, Exponent: -3},
		},
		{
			name: "0.12 to decimal",
			args: args{"0.12"},
			want: Decimal{Mantissa: 12, Exponent: -2},
		},
		{
			name: "-0.12 to decimal",
			args: args{"-0.12"},
			want: Decimal{Mantissa: -12, Exponent: -2},
		},
		{
			name: "0.1 to decimal",
			args: args{"0.1"},
			want: Decimal{Mantissa: 1, Exponent: -1},
		},
		{
			name: "-0.1 to decimal",
			args: args{"-0.1"},
			want: Decimal{Mantissa: -1, Exponent: -1},
		},
		{
			name: "1 to decimal",
			args: args{"1"},
			want: Decimal{Mantissa: 1, Exponent: 0},
		},
		{
			name: "-1 to decimal",
			args: args{"-1"},
			want: Decimal{Mantissa: -1, Exponent: 0},
		},
		{
			name: "1.00000000 to decimal",
			args: args{"1.00000000"},
			want: Decimal{Mantissa: 100000000, Exponent: -8},
		},
		{
			name: "-1.00000000 to decimal",
			args: args{"-1.00000000"},
			want: Decimal{Mantissa: -100000000, Exponent: -8},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseDecimal(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDecimal9(t *testing.T) {
	type args struct {
		d *Decimal
	}
	tests := []struct {
		name    string
		args    args
		want    Decimal9
		wantErr bool
	}{
		{
			name:    "0 Exponent",
			args:    args{&Decimal{Mantissa: 1, Exponent: 0}},
			want:    Decimal9{1_000_000_000, -9},
			wantErr: false,
		},
		{
			name:    "-1 Exponent",
			args:    args{&Decimal{Mantissa: 1, Exponent: -1}},
			want:    Decimal9{100_000_000, -9},
			wantErr: false,
		},
		{
			name:    "-2 Exponent",
			args:    args{&Decimal{Mantissa: 12, Exponent: -2}},
			want:    Decimal9{120_000_000, -9},
			wantErr: false,
		},
		{
			name:    "-3 Exponent",
			args:    args{&Decimal{Mantissa: 123, Exponent: -3}},
			want:    Decimal9{123_000_000, -9},
			wantErr: false,
		},
		{
			name:    "-4 Exponent",
			args:    args{&Decimal{Mantissa: 1234, Exponent: -4}},
			want:    Decimal9{123_400_000, -9},
			wantErr: false,
		},
		{
			name:    "-5 Exponent",
			args:    args{&Decimal{Mantissa: 12345, Exponent: -5}},
			want:    Decimal9{123_450_000, -9},
			wantErr: false,
		},
		{
			name:    "-6 Exponent",
			args:    args{&Decimal{Mantissa: 123456, Exponent: -6}},
			want:    Decimal9{123_456_000, -9},
			wantErr: false,
		},
		{
			name:    "-7 Exponent",
			args:    args{&Decimal{Mantissa: 1234567, Exponent: -7}},
			want:    Decimal9{123_456_700, -9},
			wantErr: false,
		},
		{
			name:    "-8 Exponent",
			args:    args{&Decimal{Mantissa: 12345678, Exponent: -8}},
			want:    Decimal9{123_456_780, -9},
			wantErr: false,
		},
		{
			name:    "-9 Exponent",
			args:    args{&Decimal{Mantissa: 123456789, Exponent: -9}},
			want:    Decimal9{123_456_789, -9},
			wantErr: false,
		},
		{
			name:    "-10 Exponent",
			args:    args{&Decimal{Mantissa: 123456789, Exponent: -10}},
			want:    Decimal9{math.MinInt64, -9},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDecimal9(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDecimal9() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDecimal9() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDecimal8(t *testing.T) {
	type args struct {
		d *Decimal
	}
	tests := []struct {
		name    string
		args    args
		want    Decimal8
		wantErr bool
	}{
		{
			name:    "0 Exponent",
			args:    args{&Decimal{Mantissa: 1, Exponent: 0}},
			want:    Decimal8{100_000_000, -8},
			wantErr: false,
		},
		{
			name:    "-1 Exponent",
			args:    args{&Decimal{Mantissa: 1, Exponent: -1}},
			want:    Decimal8{10_000_000, -8},
			wantErr: false,
		},
		{
			name:    "-2 Exponent",
			args:    args{&Decimal{Mantissa: 12, Exponent: -2}},
			want:    Decimal8{12_000_000, -8},
			wantErr: false,
		},
		{
			name:    "-3 Exponent",
			args:    args{&Decimal{Mantissa: 123, Exponent: -3}},
			want:    Decimal8{12_300_000, -8},
			wantErr: false,
		},
		{
			name:    "-4 Exponent",
			args:    args{&Decimal{Mantissa: 1234, Exponent: -4}},
			want:    Decimal8{12_340_000, -8},
			wantErr: false,
		},
		{
			name:    "-5 Exponent",
			args:    args{&Decimal{Mantissa: 12345, Exponent: -5}},
			want:    Decimal8{12_345_000, -8},
			wantErr: false,
		},
		{
			name:    "-6 Exponent",
			args:    args{&Decimal{Mantissa: 123456, Exponent: -6}},
			want:    Decimal8{12_345_600, -8},
			wantErr: false,
		},
		{
			name:    "-7 Exponent",
			args:    args{&Decimal{Mantissa: 1234567, Exponent: -7}},
			want:    Decimal8{12_345_670, -8},
			wantErr: false,
		},
		{
			name:    "-8 Exponent",
			args:    args{&Decimal{Mantissa: 12345678, Exponent: -8}},
			want:    Decimal8{12_345_678, -8},
			wantErr: false,
		},
		{
			name:    "-9 Exponent",
			args:    args{&Decimal{Mantissa: 123456789, Exponent: -9}},
			want:    Decimal8{math.MinInt64, -8},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDecimal8(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDecimal8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDecimal8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimal_String(t *testing.T) {
	type fields struct {
		Mantissa int64
		Exponent int8
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name:   "Positive number",
			fields: fields(ParseDecimal("12345.6789")),
			want:   "12345.6789",
		},
		{
			name:   "Negative number",
			fields: fields(ParseDecimal("-12345.6789")),
			want:   "-12345.6789",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Decimal{
				Mantissa: tt.fields.Mantissa,
				Exponent: tt.fields.Exponent,
			}
			if got := d.String(); got != tt.want {
				t.Errorf("Decimal.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
