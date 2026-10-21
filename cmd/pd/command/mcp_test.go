package command

import (
	"bytes"
	"testing"
)

func TestDecodeMCPData(t *testing.T) {
	for _, test := range []struct {
		name string
		data string
		want []byte
	}{
		{name: "base64", data: "aGVsbG8=", want: []byte("hello")},
		{name: "data URI", data: "data:text/plain;base64,aGVsbG8=", want: []byte("hello")},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := decodeMCPData(test.data)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, test.want) {
				t.Fatalf("decodeMCPData() = %q, want %q", got, test.want)
			}
		})
	}
}
