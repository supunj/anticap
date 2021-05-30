package unit

import (
	"bytes"
	"crypto/sha256"
	"encoding"
	"fmt"
	"log"
	"testing"

	crypt_util "github.com/supunj/anticap/internal/util/crypt"
)

func TestHash(t *testing.T) {
	const (
		input1 = "The tunneling gopher digs downwards, "
		input2 = "unaware of what he will find."
	)

	first := sha256.New()
	first.Write([]byte(input1))

	marshaler, ok := first.(encoding.BinaryMarshaler)
	if !ok {
		log.Fatal("first does not implement encoding.BinaryMarshaler")
	}
	state, err := marshaler.MarshalBinary()
	if err != nil {
		log.Fatal("unable to marshal hash:", err)
	}

	fmt.Printf("%x\n", first.Sum(nil))

	second := sha256.New()

	unmarshaler, ok := second.(encoding.BinaryUnmarshaler)
	if !ok {
		log.Fatal("second does not implement encoding.BinaryUnmarshaler")
	}
	if err := unmarshaler.UnmarshalBinary(state); err != nil {
		log.Fatal("unable to unmarshal hash:", err)
	}

	first.Write([]byte(input2))
	second.Write([]byte(input2))

	fmt.Printf("%x\n", first.Sum(nil))
	fmt.Println(bytes.Equal(first.Sum(nil), second.Sum(nil)))
}

func TestGetHsh(t *testing.T) {
	hash := crypt_util.GetHash("supun")
	fmt.Println(hash)
}

func TestGetKeys(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		want1   string
		wantErr bool
	}{{
		"t1",
		"",
		"",
		false,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := crypt_util.GetKeys()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetKeys() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetKeys() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
