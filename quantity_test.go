package marketdata

import (
	"reflect"
	"testing"
)

func TestParseQuantity(t *testing.T) {
	type args struct {
		qty string
	}
	tests := []struct {
		name    string
		args    args
		want    Quantity
		wantErr bool
	}{
		{
			name: "Parse 245812",
			args: args{"245812"},
			want: Quantity{Mantissa: 24581200000000, Exponent: -8},
			wantErr: false,
		},
		{
			name: "Parse 10",
			args: args{"10"},
			want: Quantity{Mantissa: 1000000000, Exponent: -8},
			wantErr: false,
		},
		{
			name: "Parse 0.00000001",
			args: args{"0.00000001"},
			want: Quantity{Mantissa: 1, Exponent: -8},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseQuantity(tt.args.qty)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuantity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQuantity() = %v, want %v", got, tt.want)
			}
		})
	}
}
