package integration_test

import "testing"

func TestKurzd(t *testing.T) {
	for _, test := range []struct {
		name string
	}{
		{
			name: "first",
		},
		{
			name: "second",
		},
	} {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
