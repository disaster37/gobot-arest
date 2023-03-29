package serialClient

import (
	"sync"
	"time"

	"go.bug.st/serial"
)

func MockSerialClient() *Client {

	mock := NewMockSerial()
	client := NewClient("/dev/null", &serial.Mode{}, 100*time.Second, true)
	client.SetSerial(mock)

	return client

}

type MockSerialBase struct{}

func (m *MockSerialBase) SetMode(mode *serial.Mode) error                      { return nil }
func (m *MockSerialBase) Read(p []byte) (n int, err error)                     { return 0, nil }
func (m *MockSerialBase) Write(p []byte) (n int, err error)                    { return 0, nil }
func (m *MockSerialBase) ResetInputBuffer() error                              { return nil }
func (m *MockSerialBase) ResetOutputBuffer() error                             { return nil }
func (m *MockSerialBase) SetDTR(dtr bool) error                                { return nil }
func (m *MockSerialBase) SetRTS(rts bool) error                                { return nil }
func (m *MockSerialBase) GetModemStatusBits() (*serial.ModemStatusBits, error) { return nil, nil }
func (m *MockSerialBase) Close() error                                         { return nil }
func (m *MockSerialBase) Break(t time.Duration) error                          { return nil }
func (m *MockSerialBase) SetReadTimeout(t time.Duration) error { return nil }

type MockSerial struct {
	MockSerialBase
	read      func(p []byte) (n int, err error)
	write     func(p []byte) (n int, err error)
	close     func() error
	ReadData  []byte
	readData  []byte
	WriteData []byte
	mtx       sync.Mutex
}

func (m *MockSerial) TestRead(f func(p []byte) (n int, err error)) {

	m.read = f
}

func (m *MockSerial) TestWrite(f func(p []byte) (n int, err error)) {

	m.write = f
}

func (m *MockSerial) TestClose(f func() error) {

	m.close = f
}

func (m *MockSerial) Read(p []byte) (n int, err error) {
	return m.read(p)
}

func (m *MockSerial) Write(p []byte) (n int, err error) {
	return m.write(p)
}

func (m *MockSerial) Close() error {
	return m.close()
}

func (m *MockSerial) InitRead() {
	m.mtx.Lock()
	m.read = func(p []byte) (n int, err error) {

		// Simulate wait data
		m.mtx.Lock()

		if len(m.readData) == 0 {
			return 0, nil
		}

		n = len(m.readData)
		for i, b := range m.readData {
			if i < len(p) {
				p[i] = b
			} else {
				n = i
				break
			}
		}

		m.readData = m.readData[n:]
		m.mtx.Unlock()

		return

	}

	m.write = func(p []byte) (n int, err error) {
		m.readData = m.ReadData

		go func() {
			m.mtx.Unlock()
		}()

		return 0, nil
	}
}

func NewMockSerial() serial.Port {
	m := &MockSerial{
		close: func() error { return nil },
	}

	m.mtx.Lock()
	m.read = func(p []byte) (n int, err error) {

		// Simulate wait data
		m.mtx.Lock()

		if len(m.readData) == 0 {
			return 0, nil
		}

		n = len(m.readData)
		for i, b := range m.readData {
			if i < len(p) {
				p[i] = b
			} else {
				n = i
				break
			}
		}

		m.readData = m.readData[n:]
		m.mtx.Unlock()

		return

	}

	m.write = func(p []byte) (n int, err error) {
		m.readData = m.ReadData

		go func() {
			m.mtx.Unlock()
		}()

		return 0, nil
	}

	return m
}
