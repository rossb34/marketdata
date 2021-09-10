package marketdata

import (
	"bufio"
	"os"
	"testing"
	"time"
)

func TestMarketDataArchiver_getFilename(t *testing.T) {
	type fields struct {
		fnamePrefix string
		archiveDir  string
		currentDate UTCDate
		currentHour int
		writer      *bufio.Writer
		file        *os.File
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "single digit hour",
			fields: fields{fnamePrefix: "foo", archiveDir: "archive", currentDate: UTCDate{year: 2021, month: time.January, day: 1}, currentHour: 1},
			want:   "archive/20210101/foo_01.dat",
		},
		{
			name:   "double digit hour",
			fields: fields{fnamePrefix: "foo", archiveDir: "archive", currentDate: UTCDate{year: 2021, month: time.January, day: 1}, currentHour: 23},
			want:   "archive/20210101/foo_23.dat",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MarketDataArchiver{
				fnamePrefix: tt.fields.fnamePrefix,
				archiveDir:  tt.fields.archiveDir,
				currentDate: tt.fields.currentDate,
				currentHour: tt.fields.currentHour,
				writer:      tt.fields.writer,
				file:        tt.fields.file,
			}
			if got := m.getFilePath(); got != tt.want {
				t.Errorf("MarketDataArchiver.getFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarketDataArchiver_checkShouldRotateFile(t *testing.T) {
	type fields struct {
		fnamePrefix string
		archiveDir  string
		currentDate UTCDate
		currentHour int
		writer      *bufio.Writer
		file        *os.File
	}
	type args struct {
		timestamp time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Should not rotate file",
			fields: fields{fnamePrefix: "foo", archiveDir: "archive", currentDate: UTCDate{year: 2021, month: time.January, day: 1}, currentHour: 1},
			args:   args{time.Date(2021, time.January, 1, 1, 59, 59, 0, time.UTC)},
			want:   false,
		},
		{
			name:   "Should rotate file",
			fields: fields{fnamePrefix: "foo", archiveDir: "archive", currentDate: UTCDate{year: 2021, month: time.January, day: 1}, currentHour: 1},
			args:   args{time.Date(2021, time.January, 1, 2, 0, 0, 0, time.UTC)},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MarketDataArchiver{
				fnamePrefix: tt.fields.fnamePrefix,
				archiveDir:  tt.fields.archiveDir,
				currentDate: tt.fields.currentDate,
				currentHour: tt.fields.currentHour,
				writer:      tt.fields.writer,
				file:        tt.fields.file,
			}
			if got := m.checkShouldRotateFile(tt.args.timestamp); got != tt.want {
				t.Errorf("MarketDataArchiver.checkShouldRotateFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
