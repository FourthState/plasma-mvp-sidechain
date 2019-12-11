package plasma

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInputValidation(t *testing.T) {
	type validationCase struct {
		reason string
		Input
	}

	pos, _ := FromPositionString("(1.0.0.0)")
	depositPos, _ := FromPositionString("(0.0.0.1)")
	nilPos, _ := FromPositionString("(0.0.0.0)")

	sampleSig := [65]byte{}
	sampleSig[0] = byte(1)

	invalidInputs := []validationCase{
		validationCase{
			reason: "nil input with a signature",
			Input:  NewInput(nilPos, sampleSig, nil),
		},
		validationCase{
			reason: "nil input with a confirm signature",
			Input:  NewInput(nilPos, [65]byte{}, [][65]byte{sampleSig}),
		},
		validationCase{
			reason: "input with no signature",
			Input:  NewInput(pos, [65]byte{}, [][65]byte{sampleSig}),
		},
		validationCase{
			reason: "transaction input with no confirm signature",
			Input:  NewInput(pos, sampleSig, nil),
		},
		validationCase{
			reason: "deposit input with a no signature",
			Input:  NewInput(depositPos, [65]byte{}, nil),
		},
		validationCase{
			reason: "deposit input with a confirmSignature signature",
			Input:  NewInput(depositPos, sampleSig, [][65]byte{sampleSig}),
		},
	}

	for _, input := range invalidInputs {
		err := input.ValidateBasic()
		require.Error(t, err, input.reason)
	}

	pos, _ = FromPositionString("(1.1.0.0)")
	input := NewInput(pos, sampleSig, [][65]byte{sampleSig})
	err := input.ValidateBasic()
	require.NoError(t, err)

	pos, _ = FromPositionString("(0.0.0.0)")
	input = NewInput(pos, [65]byte{}, nil)
	err = input.ValidateBasic()
	require.NoError(t, err)
}
