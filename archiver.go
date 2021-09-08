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
func (m *MarketDataArchiver) getFilePath() string {
	return fmt.Sprintf("%v/%v_%02d.dat", m.getDir(), m.fnamePrefix, m.currentHour)
}

// Gets the path
func (m *MarketDataArchiver) getDir() string {
	return fmt.Sprintf("%v%02d%02d", m.currentDate.year, m.currentDate.month, m.currentDate.day)
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

	// Ensure directory exists before creating file
	err := os.MkdirAll(m.getDir(), os.ModePerm)
	if err != nil {
		return err
	}

	// Create and open the file for writing
	file, err := os.Create(m.getFilePath())
	if err != nil {
		return err
	}
	m.file = file

	// Buffered writer to wrap the file (buffer size of 4 MB)
	m.writer = bufio.NewWriterSize(m.file, 1024*1024*4)

	return nil
}

// Archives the message
func (m *MarketDataArchiver) ArchiveMessage(timestamp time.Time, msg []byte, sep byte) (int, error) {
	if m.checkShouldRotateFile(timestamp) {
		if err := m.rotateFile(); err != nil {
			return 0, err
		}
	}
	nn, err := m.writer.Write(msg)
	if err != nil {
		return nn, err
	}

	err = m.writer.WriteByte(sep)
	if err != nil {
		return nn, err
	}
	return nn + 1, nil
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
