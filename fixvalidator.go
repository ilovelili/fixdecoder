package fixdecoder

import (
	"fmt"
	"strconv"
)

// ValidatorFactory validator factory
type ValidatorFactory struct{}

// NewValidatorFactory validator factory
func NewValidatorFactory() *ValidatorFactory {
	return &ValidatorFactory{}
}

// CreateValidators create validators
func (vf *ValidatorFactory) CreateValidators() []Validator {
	return []Validator{BodyLengthValidator{}, CheckSumValidator{}}
}

// Validator field validator. For example, checksum validation, body length validation
type Validator interface {
	Validate(DecodedFields) bool
}

// BodyLengthValidator BodyLength is the character count starting at tag 35 (included, MsgType) all the way to tag 10 (excluded). SOH delimiters do count in body length (length = 1).
type BodyLengthValidator struct{}

// Validate body length validate
func (v BodyLengthValidator) Validate(dfs DecodedFields) bool {
	// pass if no data
	if len(dfs) == 0 {
		return true
	}

	length := 0
	var bodylengthfield *DecodedField
	for _, line := range dfs {
		if line.FieldID == BODYLENGTH {
			bodylengthfield = line
			// exclude body length
			continue
		}

		// Some fields are not part of the FIX message body, exclude them
		if line.FieldID == BEGINSTRING || line.FieldID == CHECKSUM {
			continue
		}

		length += len(line.Raw())
	}

	bodylengthfieldvalue, _ := strconv.Atoi(bodylengthfield.Value)
	if bodylengthfieldvalue == length {
		bodylengthfield.Classes += " Valid"
		bodylengthfield.DecodedValue = "Valid"
		return true
	}

	bodylengthfield.Classes += " Invalid"
	bodylengthfield.DecodedValue = fmt.Sprintf("Invalid (expected %v)", length)
	return false
}

// CheckSumValidator checksum
type CheckSumValidator struct{}

// Validate checksum validate
// https://www.onixs.biz/fix-dictionary/4.2/app_b.html
func (v CheckSumValidator) Validate(dfs DecodedFields) bool {
	// pass if no data
	if len(dfs) == 0 {
		return true
	}

	sum := 0
	var checksumfield *DecodedField
	for _, line := range dfs {
		if line.FieldID == CHECKSUM {
			checksumfield = line
			// exclude checksum
			continue
		}

		for i := 0; i < len(line.Raw()); i++ {
			sum += int((line.Raw())[i])
		}
	}

	// Modulo 256 + pad up to 3 characters with zero
	modulo := "00" + strconv.Itoa(sum%256)
	checksum := modulo[len(modulo)-3:]

	if checksumfield.Value == checksum {
		checksumfield.Classes += " Valid"
		checksumfield.DecodedValue = "Valid"
		return true
	}

	checksumfield.Classes += " Invalid"
	checksumfield.DecodedValue = fmt.Sprintf("Invalid (expected %v)", checksum)
	return false
}
