package plasma

import (
	"bytes"
	"fmt"
)

// Input represents the input to a spend
type Input struct {
	Position          `json:"position"`
	Signature         [65]byte   `json:"signature"`
	ConfirmSignatures [][65]byte `json:"confirmSignatures"`
}

func NewInput(position Position, sig [65]byte, confirmSigs [][65]byte) Input {
	// nil != empty slice. avoid deserialization issues by forcing empty slices
	if confirmSigs == nil {
		confirmSigs = [][65]byte{}
	}

	return Input{
		Position:          position,
		Signature:         sig,
		ConfirmSignatures: confirmSigs,
	}
}

/*
	So far in the project, serialization for the Input struct was not needed. This can
	be added in here when needed
*/

// ValidateBasic ensures a nil or valid input
func (i Input) ValidateBasic() error {
	var emptySig [65]byte
	if i.Position.IsNilPosition() {
		if !bytes.Equal(i.Signature[:], emptySig[:]) || len(i.ConfirmSignatures) > 0 {
			return fmt.Errorf("nil input should not specifiy a signature nor confirm signatures")
		}
	} else {
		if err := i.Position.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid position { %s }", err)
		}

		if bytes.Equal(i.Signature[:], emptySig[:]) {
			return fmt.Errorf("cannot provide an empty signature")
		}

		confSigLen := len(i.ConfirmSignatures)
		if i.Position.IsDeposit() || i.TxIndex == 1<<16-1 {
			if confSigLen != 0 {
				return fmt.Errorf("deposit or fee inputs must not include confirm signatures")
			}
		} else {
			if confSigLen != 1 && confSigLen != 2 {
				return fmt.Errorf("transaction inputs must specify 1 or 2 confirm signatures")
			}

			for _, sig := range i.ConfirmSignatures {
				if len(sig) != 65 {
					return fmt.Errorf("confirm signatures must be 65 bytes in length")
				}
			}
		}
	}

	return nil
}

func (i Input) String() string {
	if len(i.ConfirmSignatures) != 0 {
		str := fmt.Sprintf("Position: %s, Signature: 0x%x, Confirm Signatures: 0x%x",
			i.Position, i.Signature, i.ConfirmSignatures[0])
		if len(i.ConfirmSignatures) > 1 {
			str = str + fmt.Sprintf(", 0x%x", i.ConfirmSignatures[1])
		}

		return str
	}

	return fmt.Sprintf("Position: %s, Signature: 0x%x, Confirm Signatures: nil", i.Position, i.Signature)
}
