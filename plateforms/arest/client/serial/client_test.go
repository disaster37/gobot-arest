package serialClient

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/disaster37/gobot-arest/plateforms/arest/client"
	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type ArestTestSuite struct {
	suite.Suite
	client *Client
	mux    sync.Mutex
}

func (s *ArestTestSuite) SetupSuite() {
	// Init logger
	logrus.SetFormatter(new(prefixed.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)

}

func (s *ArestTestSuite) BeforeSuite() {
}

func (s *ArestTestSuite) AfterSuite() {
	httpmock.DeactivateAndReset()
}

func (s *ArestTestSuite) SetupTest() {
	s.client = MockSerialClient()
}

func TestArestTestSuite(t *testing.T) {
	suite.Run(t, new(ArestTestSuite))
}

func (s *ArestTestSuite) TestConnect() {
	s.mux.Lock()
	defer s.mux.Unlock()
	err := s.client.Connect(context.Background())
	assert.NoError(s.T(), err)

	err = s.client.Disconnect(context.Background())
	assert.NoError(s.T(), err)
}

func (s *ArestTestSuite) TestSetMode() {
	s.mux.Lock()
	defer s.mux.Unlock()
	// Error when not yet connected
	err := s.client.SetPinMode(context.Background(), 0, client.ModeOutput)
	assert.Error(s.T(), err)

	// Normal use case
	s.client.Connect(context.Background())
	err = s.client.SetPinMode(context.Background(), 0, client.ModeOutput)
	assert.NoError(s.T(), err)
}

func (s *ArestTestSuite) TestDigitalWrite() {

	s.mux.Lock()
	defer s.mux.Unlock()
	s.client.Connect(context.Background())

	// Return error if pin is not yet setted
	err := s.client.DigitalWrite(context.Background(), 0, client.LevelHigh)
	assert.Error(s.T(), err)

	// Return error if pin is not output mode
	s.client.SetPinMode(context.Background(), 0, client.ModeInput)
	err = s.client.DigitalWrite(context.Background(), 0, client.LevelHigh)
	assert.Error(s.T(), err)

	// Normal use case
	s.client.SetPinMode(context.Background(), 0, client.ModeOutput)
	err = s.client.DigitalWrite(context.Background(), 0, client.LevelHigh)
	assert.NoError(s.T(), err)
}

func (s *ArestTestSuite) TestDigitalRead() {

	s.client.Connect(context.Background())

	var err error

	// Return error if pin is not yet setted

	_, err = s.client.DigitalRead(context.Background(), 0)
	assert.Error(s.T(), err)

	// Return error if pin is not output mode
	s.client.SetPinMode(context.Background(), 0, client.ModeOutput)
	_, err = s.client.DigitalRead(context.Background(), 0)
	assert.Error(s.T(), err)

	// Normal use case
	s.client.SetPinMode(context.Background(), 0, client.ModeInput)

	fixture := map[string]interface{}{
		"return_value": 1,
	}
	s.client.Client().(*MockSerial).ReadData, err = json.Marshal(fixture)
	if err != nil {
		panic(err)
	}

	level, err := s.client.DigitalRead(context.Background(), 0)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), client.LevelHigh, level)
}

func (s *ArestTestSuite) TestReadValue() {

	s.client.Connect(context.Background())

	var err error
	fixture := map[string]interface{}{
		"isRebooted": true,
	}
	s.client.Client().(*MockSerial).ReadData, err = json.Marshal(fixture)
	if err != nil {
		panic(err)
	}

	value, err := s.client.ReadValue(context.Background(), "isRebooted")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), true, value.(bool))

	// Bad
	value, err = s.client.ReadValue(context.Background(), "bad")
	assert.Error(s.T(), err)
}

func (s *ArestTestSuite) TestReadValues() {

	s.client.Connect(context.Background())

	var err error
	fixture := map[string]interface{}{
		"variables": map[string]interface{}{
			"isRebooted": false,
		},
	}
	s.client.Client().(*MockSerial).ReadData, err = json.Marshal(fixture)
	if err != nil {
		panic(err)
	}

	values, err := s.client.ReadValues(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), false, values["isRebooted"].(bool))

}

func (s *ArestTestSuite) TestCallFunction() {

	s.client.Connect(context.Background())

	var err error
	fixture := map[string]interface{}{
		"return_value": 1,
	}
	s.client.Client().(*MockSerial).ReadData, err = json.Marshal(fixture)
	if err != nil {
		panic(err)
	}

	resp, err := s.client.CallFunction(context.Background(), "acknoledgeRebooted", "test")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, resp)

	// Bad
	s.client.Client().(*MockSerial).ReadData = make([]byte, 0)
	resp, err = s.client.CallFunction(context.Background(), "bad", "test")
	assert.Error(s.T(), err)
}
