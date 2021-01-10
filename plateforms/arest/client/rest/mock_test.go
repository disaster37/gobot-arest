package restClient

import (
	"time"

	"github.com/jarcoal/httpmock"
)

func MockRestClient() *Client {
	client := NewClient("http://localhost", 1*time.Second, true)
	httpmock.ActivateNonDefault(client.Client().GetClient())

	return client
}
