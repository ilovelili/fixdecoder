package fixdecoder_test

import (
	"strings"
	"testing"

	fixdecoder "github.com/min-invastsec/FixDecoder"
)

const (
	validfixmessage              = "8=FIX.4.49=7435=249=CNX34=826333652=20180126-07:39:59.68356=imdstream16=07=1281210=036"
	invalidfixmessage_badformat  = "whosyourdaddy"
	invalidfixmessage_checksum   = "8=FIX.4.49=7435=249=CNX34=826333652=20180126-07:39:59.68356=imdstream16=07=1281210=999"
	invalidfixmessage_bodylength = "8=FIX.4.49=8835=249=CNX34=826333652=20180126-07:39:59.68356=imdstream16=07=1281210=036"
)

var (
	fd *fixdecoder.FixDecoder = fixdecoder.NewFixDecoder()
	r                         = strings.NewReplacer(" ", "", "\n", "")
)

func TestFixDecoder_Valid(t *testing.T) {
	actual := r.Replace(fd.Decode(validfixmessage).String())
	expect := `{"FieldID":"8","Value":"FIX.4.4","Field":{"Name":"BeginString","Type":"STRING"},"Classes":"required-field,header-field"}{"FieldID":"9","Value":"74","Field":{"Name":"BodyLength","Type":"LENGTH"},"DecodedValue":"Valid","Classes":"required-field,header-fieldValid"}{"FieldID":"35","Value":"2","Field":{"Name":"MsgType","Type":"STRING"},"DecodedValue":"ResendRequest","Classes":"required-field,header-field"}{"FieldID":"49","Value":"CNX","Field":{"Name":"SenderCompID","Type":"STRING"},"Classes":"required-field,header-field"}{"FieldID":"34","Value":"8263336","Field":{"Name":"MsgSeqNum","Type":"SEQNUM"},"Classes":"required-field,header-field"}{"FieldID":"52","Value":"20180126-07:39:59.683","Field":{"Name":"SendingTime","Type":"UTCTIMESTAMP"},"Classes":"required-field,header-field"}{"FieldID":"56","Value":"imdstream","Field":{"Name":"TargetCompID","Type":"STRING"},"Classes":"required-field,header-field"}{"FieldID":"16","Value":"0","Field":{"Name":"EndSeqNo","Type":"SEQNUM"},"Classes":""}{"FieldID":"7","Value":"12812","Field":{"Name":"BeginSeqNo","Type":"SEQNUM"},"Classes":""}{"FieldID":"10","Value":"036","Field":{"Name":"CheckSum","Type":"STRING"},"DecodedValue":"Valid","Classes":"system-fieldValid"}`

	if actual != expect {
		t.Errorf("expect %s, actual %s", expect, actual)
	}
}

func TestFixDecoder_Invalid_BadFormat(t *testing.T) {
	actual := r.Replace(fd.Decode(invalidfixmessage_badformat).String())
	expect := ""
	if actual != expect {
		t.Errorf("expect %s, actual %s", expect, actual)
	}
}

func TestFixDecoder_Invalid_CheckSum(t *testing.T) {
	actual := r.Replace(fd.Decode(invalidfixmessage_checksum).String())
	expect := `{"FieldID":"8","Value":"FIX.4.4","Field":{"Name":"BeginString","Type":"STRING"},"Classes":"required-field,header-field"}{"FieldID":"9","Value":"74","Field":{"Name":"BodyLength","Type":"LENGTH"},"DecodedValue":"Valid","Classes":"required-field,header-fieldValid"}{"FieldID":"35","Value":"2","Field":{"Name":"MsgType","Type":"STRING"},"DecodedValue":"ResendRequest","Classes":"required-field,header-field"}{"FieldID":"49","Value":"CNX","Field":{"Name":"SenderCompID","Type":"STRING"},"Classes":"required-field,header-field"}{"FieldID":"34","Value":"8263336","Field":{"Name":"MsgSeqNum","Type":"SEQNUM"},"Classes":"required-field,header-field"}{"FieldID":"52","Value":"20180126-07:39:59.683","Field":{"Name":"SendingTime","Type":"UTCTIMESTAMP"},"Classes":"required-field,header-field"}{"FieldID":"56","Value":"imdstream","Field":{"Name":"TargetCompID","Type":"STRING"},"Classes":"required-field,header-field"}{"FieldID":"16","Value":"0","Field":{"Name":"EndSeqNo","Type":"SEQNUM"},"Classes":""}{"FieldID":"7","Value":"12812","Field":{"Name":"BeginSeqNo","Type":"SEQNUM"},"Classes":""}{"FieldID":"10","Value":"999","Field":{"Name":"CheckSum","Type":"STRING"},"DecodedValue":"Invalid(expected036)","Classes":"system-fieldInvalid"}`

	if actual != expect {
		t.Errorf("expect %s, actual %s", expect, actual)
	}
}

func TestFixDecoder_Invalid_BodyLength(t *testing.T) {
	actual := r.Replace(fd.Decode(invalidfixmessage_bodylength).String())
	expect := `{"FieldID":"8","Value":"FIX.4.4","Field":{"Name":"BeginString","Type":"STRING"},"Classes":"required-field,header-field"}{"FieldID":"9","Value":"88","Field":{"Name":"BodyLength","Type":"LENGTH"},"DecodedValue":"Invalid(expected74)","Classes":"required-field,header-fieldInvalid"}{"FieldID":"35","Value":"2","Field":{"Name":"MsgType","Type":"STRING"},"DecodedValue":"ResendRequest","Classes":"required-field,header-field"}{"FieldID":"49","Value":"CNX","Field":{"Name":"SenderCompID","Type":"STRING"},"Classes":"required-field,header-field"}{"FieldID":"34","Value":"8263336","Field":{"Name":"MsgSeqNum","Type":"SEQNUM"},"Classes":"required-field,header-field"}{"FieldID":"52","Value":"20180126-07:39:59.683","Field":{"Name":"SendingTime","Type":"UTCTIMESTAMP"},"Classes":"required-field,header-field"}{"FieldID":"56","Value":"imdstream","Field":{"Name":"TargetCompID","Type":"STRING"},"Classes":"required-field,header-field"}{"FieldID":"16","Value":"0","Field":{"Name":"EndSeqNo","Type":"SEQNUM"},"Classes":""}{"FieldID":"7","Value":"12812","Field":{"Name":"BeginSeqNo","Type":"SEQNUM"},"Classes":""}{"FieldID":"10","Value":"036","Field":{"Name":"CheckSum","Type":"STRING"},"DecodedValue":"Invalid(expected041)","Classes":"system-fieldInvalid"}`

	if actual != expect {
		t.Errorf("expect %s, actual %s", expect, actual)
	}
}
