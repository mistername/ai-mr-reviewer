package gitlab

import "testing"

func TestClientImplementsPort(t *testing.T) {
	var _ interface{} = (*Client)(nil)
}
