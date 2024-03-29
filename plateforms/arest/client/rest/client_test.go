package restClient

import (
	"context"
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
	s.client = MockRestClient()
	httpmock.Reset()
}

func TestArestTestSuite(t *testing.T) {
	suite.Run(t, new(ArestTestSuite))
}

func (s *ArestTestSuite) TestConnect() {
	fixture := `{"message": "Pin D0 set to output", "id": "002", "name": "TFP", "hardware": "arduino", "connected": true}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeURL := "http://localhost/id"
	httpmock.RegisterResponder("GET", fakeURL, responder)

	err := s.client.Connect(context.Background())
	assert.NoError(s.T(), err)
	assert.True(s.T(), s.client.connected.Load().(bool))
}

func (s *ArestTestSuite) TestDisconnect() {

	err := s.client.Disconnect(context.Background())
	assert.NoError(s.T(), err)
	assert.False(s.T(), s.client.connected.Load().(bool))
}

func (s *ArestTestSuite) TestReconnect() {
	fixture := `{"message": "Pin D0 set to output", "id": "002", "name": "TFP", "hardware": "arduino", "connected": true}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeURL := "http://localhost/id"
	httpmock.RegisterResponder("GET", fakeURL, responder)

	err := s.client.Reconnect(context.Background())
	assert.NoError(s.T(), err)
	assert.True(s.T(), s.client.connected.Load().(bool))
}

func (s *ArestTestSuite) TestSetMode() {

	fixture := `{"message": "Pin D0 set to output", "id": "002", "name": "TFP", "hardware": "arduino", "connected": true}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeURL := "http://localhost/mode/0/o"
	httpmock.RegisterResponder("POST", fakeURL, responder)

	err := s.client.SetPinMode(context.Background(), 0, client.ModeOutput)
	assert.NoError(s.T(), err)

}

func (s *ArestTestSuite) TestDigitalWrite() {

	fixture := `{"message": "Pin D0 set to 1", "id": "002", "name": "TFP", "hardware": "arduino", "connected": true}`
	responder := httpmock.NewStringResponder(200, fixture)
	fakeURL := "http://localhost/digital/0/1"
	httpmock.RegisterResponder("POST", fakeURL, responder)
	httpmock.RegisterResponder("POST", "http://localhost/mode/0/i", responder)
	httpmock.RegisterResponder("POST", "http://localhost/mode/0/o", responder)

	// Return error if pin is not yet setted
	err := s.client.DigitalWrite(context.Background(), 0, client.LevelHigh)
	assert.Error(s.T(), err)

	// Return error if pin is not output mode
	if err := s.client.SetPinMode(context.Background(), 0, client.ModeInput); err != nil {
		s.T().Fatal(err)
	}
	err = s.client.DigitalWrite(context.Background(), 0, client.LevelHigh)
	assert.Error(s.T(), err)

	// Normal use case
	if err := s.client.SetPinMode(context.Background(), 0, client.ModeOutput); err != nil {
		s.T().Fatal(err)
	}
	err = s.client.DigitalWrite(context.Background(), 0, client.LevelHigh)
	assert.NoError(s.T(), err)
}

func (s *ArestTestSuite) TestDigitalRead() {

	//fixture := `{"return_value": 1, "id": "002", "name": "TFP", "hardware": "arduino", "connected": true}`
	fixture := map[string]interface{}{
		"return_value": 1,
	}
	responder := httpmock.NewJsonResponderOrPanic(200, fixture)
	fakeURL := "http://localhost/digital/0"
	httpmock.RegisterResponder("GET", fakeURL, responder)
	httpmock.RegisterResponder("POST", "http://localhost/mode/0/o", responder)
	httpmock.RegisterResponder("POST", "http://localhost/mode/0/i", responder)

	// Return error if pin is not yet setted
	_, err := s.client.DigitalRead(context.Background(), 0)
	assert.Error(s.T(), err)

	// Return error if pin is not output mode
	if err := s.client.SetPinMode(context.Background(), 0, client.ModeOutput); err != nil {
		s.T().Fatal(err)
	}
	_, err = s.client.DigitalRead(context.Background(), 0)
	assert.Error(s.T(), err)

	// Normal use case
	if err := s.client.SetPinMode(context.Background(), 0, client.ModeInput); err != nil {
		s.T().Fatal(err)
	}
	level, err := s.client.DigitalRead(context.Background(), 0)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), client.LevelHigh, level)
}

func (s *ArestTestSuite) TestReadValue() {

	//fixture := `{"isRebooted": true, "id": "002", "name": "TFP", "hardware": "arduino", "connected": true}`
	fixture := map[string]interface{}{
		"isRebooted": true,
	}
	responder := httpmock.NewJsonResponderOrPanic(200, fixture)

	fakeURL := "http://localhost/isRebooted"
	httpmock.RegisterResponder("GET", fakeURL, responder)

	value, err := s.client.ReadValue(context.Background(), "isRebooted")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), true, value.(bool))

	// Bad
	_, err = s.client.ReadValue(context.Background(), "bad")
	assert.Error(s.T(), err)
}

func (s *ArestTestSuite) TestReadValues() {

	//fixture := `{"variables": {"isRebooted": false}, "id": "002", "name": "TFP", "hardware": "arduino", "connected": true}`
	fixture := map[string]interface{}{
		"variables": map[string]interface{}{
			"isRebooted": false,
		},
	}
	responder := httpmock.NewJsonResponderOrPanic(200, fixture)
	fakeURL := "http://localhost/"
	httpmock.RegisterResponder("GET", fakeURL, responder)

	values, err := s.client.ReadValues(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), false, values["isRebooted"].(bool))

}

func (s *ArestTestSuite) TestCallFunction() {

	//fixture := `{"return_value": 1, "id": "002", "name": "TFP", "hardware": "arduino", "connected": true}`
	fixture := map[string]interface{}{
		"return_value": 1,
	}
	responder := httpmock.NewJsonResponderOrPanic(200, fixture)
	fakeURL := "http://localhost/acknoledgeRebooted?params=test"
	httpmock.RegisterResponder("POST", fakeURL, responder)

	resp, err := s.client.CallFunction(context.Background(), "acknoledgeRebooted", "test")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, resp)

	// Bad
	_, err = s.client.CallFunction(context.Background(), "bad", "test")
	assert.Error(s.T(), err)
}
