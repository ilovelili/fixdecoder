// Package fixdecoder fix message decoder
package fixdecoder

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

const (
	// CHECKSUM Three byte, simple checksum. ALWAYS LAST FIELD IN MESSAGE; i.e. serves, with the trailing <SOH>, as the end-of-message delimiter. Always defined as three characters. (Always unencrypted)
	CHECKSUM = "10"
	// BEGINSTRING Identifies beginning of new message and protocol version. ALWAYS FIRST FIELD IN MESSAGE. (Always unencrypted)
	BEGINSTRING = "8"
	// BODYLENGTH Message length, in bytes, forward to the CheckSum <10> field. ALWAYS SECOND FIELD IN MESSAGE. (Always unencrypted)
	BODYLENGTH = "9"
)

// FieldMetaData meta data of a field
type FieldMetaData struct {
	Name string
	Type string
}

// DecodedField decoded field
type DecodedField struct {
	FieldID      string
	Value        string
	Field        *FieldMetaData
	DecodedValue string
	Classes      string
	Decoded      bool // Whether decoding succeeded or not
}

// DecodedFields alias of DecodedField slice
type DecodedFields []*DecodedField

// Raw parse the raw message
func (df *DecodedField) Raw() string {
	// "\x01" stands for ascii SOH, which is used as a delimiter in FIX protocol
	return df.FieldID + "=" + df.Value + "\x01"
}

// String decode to string
func (dfs DecodedFields) String() string {
	result := make([]string, 0)

	validators := NewValidatorFactory().CreateValidators()
	for _, v := range validators {
		v.Validate(dfs)
	}

	for _, line := range dfs {
		if !line.Decoded {
			return "Decoding failed"
		}

		var output interface{}
		if line.DecodedValue != "" {
			output = struct {
				FieldID      string
				Value        string
				FieldName    string
				FieldType    string
				DecodedValue string
			}{
				FieldID:      line.FieldID,
				Value:        line.Value,
				FieldName:    line.Field.Name,
				FieldType:    line.Field.Type,
				DecodedValue: line.DecodedValue,
			}
		} else {
			output = struct {
				FieldID   string
				Value     string
				FieldName string
				FieldType string
				Classes   string
			}{
				FieldID:   line.FieldID,
				Value:     line.Value,
				FieldName: line.Field.Name,
				FieldType: line.Field.Type,
			}
		}

		binary, _ := json.Marshal(output)
		result = append(result, string(binary))
	}

	return strings.Join(result, "\n")
}

// FixDecoder the main struct
type FixDecoder struct{}

// NewFixDecoder new fix decoder instance
func NewFixDecoder() *FixDecoder {
	return &FixDecoder{}
}

// parseVersionFromBeginString get the FIX protocol version
func (f *FixDecoder) parseVersionFromBeginString(beginStr string) string {
	return beginStr[len("FIX."):]
}

// Decode the main decode function
func (f *FixDecoder) Decode(message string) (decodedfields DecodedFields) {
	regex := regexp.MustCompile("([0-9]+)=([^|;\x01]*)")
	decodedfields = make([]*DecodedField, 0)

	fixVersion := "unknown"
	fields := Fields()
	systemFieldIDs := SystemFieldIDs()

	for i, result := 0, regex.FindAllString(message, -1); i < len(result); i++ {
		// {{fieldId}}={{value}}
		if parsed := regex.FindStringSubmatch(result[i]); len(parsed) == 3 {
			fieldID := parsed[1]
			value := parsed[2]
			field := gjson.Get(fields, fieldID)
			decodedValue := ""

			if found, values := hasField(field, "values"); found {
				decodedValue = gjson.Get(values, value).String()
			}

			if BEGINSTRING == fieldID {
				fixVersion = f.parseVersionFromBeginString(value)
			}

			classes := make([]string, 0)

			if _, contain := contains(systemFieldIDs, fieldID); contain {
				classes = append(classes, "system-field")
			}

			if found, _ := hasField(field, "isRequired"); found {
				classes = append(classes, "required-field")
			}

			if found, _ := hasField(field, "isHeaderField"); found {
				classes = append(classes, "header-field")
			}

			if found, value := hasField(field, "deprecatedSince"); found {
				if value <= fixVersion {
					classes = append(classes, "deprecated-field")
				}
			}

			decodedfields = append(decodedfields, &DecodedField{
				FieldID: fieldID,
				Value:   value,
				Field: &FieldMetaData{
					Name: gjson.Get(field.String(), "name").String(),
					Type: gjson.Get(field.String(), "type").String(),
				},
				Classes:      strings.Join(classes, ","),
				DecodedValue: decodedValue,
				Decoded:      true,
			})
		} else {
			// parsing failed
			decodedfields = append(decodedfields, &DecodedField{
				Decoded: false,
			})
		}
	}

	return DecodedFields(decodedfields)
}

// Array.contains
func contains(source []string, target string) (index int, contains bool) {
	index = -1
	contains = false
	for i, item := range source {
		if item == target {
			index = i
			contains = true
			return
		}
	}

	return
}

// help method to wrap gjson exists
func hasField(g gjson.Result, field string) (found bool, value string) {
	value = ""
	searchField := gjson.Get(g.String(), field)
	found = searchField.Exists()
	if found {
		value = searchField.String()
	}

	return
}
