package marketdata

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type UTCDate struct {
	year  int
	month time.Month
	day   int
}

type MarketDataArchiver struct {
	fnamePrefix string
	currentDate UTCDate
	currentHour int
	writer      *bufio.Writer
	file        *os.File
}

func NewMarketDataArchiver(filenamePrefix string) *MarketDataArchiver {
	bfw := MarketDataArchiver{fnamePrefix: filenamePrefix}
	return &bfw
}

// Gets the filename
func (m *MarketDataArchiver) getFilename() string {
	return fmt.Sprintf("%v%02d%02d/%v_%02d.dat", m.currentDate.year, m.currentDate.month, m.currentDate.day, m.fnamePrefix, m.currentHour)
}

// Checks if the file should be rotated
func (m *MarketDataArchiver) checkShouldRotateFile(timestamp time.Time) bool {
	year, month, day := timestamp.Date()
	h := timestamp.Hour()
	if h != m.currentHour || year != m.currentDate.year || month != m.currentDate.month || day != m.currentDate.day {
		m.currentDate = UTCDate{year: year, month: month, day: day}
		m.currentHour = h
		return true
	}
	return false
}

// Rotates the file
func (m *MarketDataArchiver) rotateFile() error {
	// Flush the writer
	if m.writer != nil {
		if err := m.writer.Flush(); err != nil {
			return err
		}
		m.writer = nil
	}

	// Close the file
	if m.file != nil {
		if err := m.file.Close(); err != nil {
			return err
		}
		m.file = nil
	}

	// Create and open the file for writing
	file, err := os.Create(m.getFilename())
	if err != nil {
		return err
	}
	m.file = file

	// Buffered writer to wrap the file (buffer size of 4 MB)
	m.writer = bufio.NewWriterSize(m.file, 1024*1024*4)

	return nil
}

// Archives the message
func (m *MarketDataArchiver) ArchiveMessage(timestamp time.Time, msg []byte) (int, error) {
	if m.checkShouldRotateFile(timestamp) {
		if err := m.rotateFile(); err != nil {
			return 0, err
		}
	}

	return m.writer.Write(msg)
}

// Flushes any buffered data and closes the file
func (m *MarketDataArchiver) Close() error {
	// Flush the buffered bytes
	if m.writer != nil {
		if err := m.writer.Flush(); err != nil {
			return err
		}
		m.writer = nil
	}

	// Close the file
	if m.file != nil {
		if err := m.file.Close(); err != nil {
			return err
		}
		m.file = nil
	}
	return nil
}
