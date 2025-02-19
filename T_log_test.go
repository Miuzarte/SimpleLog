package SimpleLog

import "testing"

func TestLog(t *testing.T) {
	logger := New("Test", true, true)
	for i := range Level(7) {
		logger.Print(i, "Test message ", i)
		logger.Printf(i, "Test message %02d", i)
		logger.Print(i, "Test \n message")
	}
}
