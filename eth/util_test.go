package eth

import (
	"bytes"
	"math/big"
	"testing"
)

func TestPriorityCalc(t *testing.T) {
	position := [3]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
	expected := uint64(1000000*1 + 10*2 + 3)

	priority := calcPriority(position).Uint64()
	if priority != expected {
		t.Fatalf("Position [1,2,3] yielded priority %d. Expected %d",
			priority, expected)
	}

	position = [3]*big.Int{big.NewInt(0), big.NewInt(2), big.NewInt(3)}
	expected = uint64(10*2 + 3)

	priority = calcPriority(position).Uint64()
	if priority != expected {
		t.Fatalf("Position [0,2,3] yielded priority %d. Expected %d",
			priority, expected)
	}

	position = [3]*big.Int{big.NewInt(1), big.NewInt(0), big.NewInt(3)}
	expected = uint64(1000000*1 + 3)

	priority = calcPriority(position).Uint64()
	if priority != expected {
		t.Fatalf("Position [1,0,3] yielded priority %d. Expected %d",
			priority, expected)
	}

	position = [3]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(0)}
	expected = uint64(1000000*1 + 10*2)

	priority = calcPriority(position).Uint64()
	if priority != expected {
		t.Fatalf("Position [1,2,0] yielded priority %d. Expected %d",
			priority, expected)
	}
}

func TestPrefixKey(t *testing.T) {
	expectedKey := []byte("prefix::key")
	key := prefixKey("prefix", []byte("key"))
	if !bytes.Equal(key, expectedKey) {
		t.Fatalf("Actual: %s, Got: %s", string(expectedKey), string(key))
	}
}
