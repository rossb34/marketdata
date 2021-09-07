package marketdata

import (
	"reflect"
	"testing"
)

func TestParsePrice(t *testing.T) {
	type args struct {
		px string
	}
	tests := []struct {
		name    string
		args    args
		want    Price
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Parse 46128.37",
			args: args{"46128.37"},
			want: Price{Mantissa: 46128370000000, Exponent: -9},
			wantErr: false,
		},
		{
			name: "Parse 28.01938362",
			args: args{"28.01938362"},
			want: Price{Mantissa: 28019383620, Exponent: -9},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePrice(tt.args.px)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
