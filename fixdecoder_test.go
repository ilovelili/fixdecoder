package fixdecoder_test

import (
	"strings"
	"testing"

	fixdecoder "github.com/ilovelili/FixDecoder"
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
	expect := `{"FieldID":"8","Value":"FIX.4.4","FieldName":"BeginString","FieldType":"STRING","Classes":""}{"FieldID":"9","Value":"74","FieldName":"BodyLength","FieldType":"LENGTH","DecodedValue":"Valid"}{"FieldID":"35","Value":"2","FieldName":"MsgType","FieldType":"STRING","DecodedValue":"ResendRequest"}{"FieldID":"49","Value":"CNX","FieldName":"SenderCompID","FieldType":"STRING","Classes":""}{"FieldID":"34","Value":"8263336","FieldName":"MsgSeqNum","FieldType":"SEQNUM","Classes":""}{"FieldID":"52","Value":"20180126-07:39:59.683","FieldName":"SendingTime","FieldType":"UTCTIMESTAMP","Classes":""}{"FieldID":"56","Value":"imdstream","FieldName":"TargetCompID","FieldType":"STRING","Classes":""}{"FieldID":"16","Value":"0","FieldName":"EndSeqNo","FieldType":"SEQNUM","Classes":""}{"FieldID":"7","Value":"12812","FieldName":"BeginSeqNo","FieldType":"SEQNUM","Classes":""}{"FieldID":"10","Value":"036","FieldName":"CheckSum","FieldType":"STRING","DecodedValue":"Valid"}`

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
	expect := `{"FieldID":"8","Value":"FIX.4.4","FieldName":"BeginString","FieldType":"STRING","Classes":""}{"FieldID":"9","Value":"74","FieldName":"BodyLength","FieldType":"LENGTH","DecodedValue":"Valid"}{"FieldID":"35","Value":"2","FieldName":"MsgType","FieldType":"STRING","DecodedValue":"ResendRequest"}{"FieldID":"49","Value":"CNX","FieldName":"SenderCompID","FieldType":"STRING","Classes":""}{"FieldID":"34","Value":"8263336","FieldName":"MsgSeqNum","FieldType":"SEQNUM","Classes":""}{"FieldID":"52","Value":"20180126-07:39:59.683","FieldName":"SendingTime","FieldType":"UTCTIMESTAMP","Classes":""}{"FieldID":"56","Value":"imdstream","FieldName":"TargetCompID","FieldType":"STRING","Classes":""}{"FieldID":"16","Value":"0","FieldName":"EndSeqNo","FieldType":"SEQNUM","Classes":""}{"FieldID":"7","Value":"12812","FieldName":"BeginSeqNo","FieldType":"SEQNUM","Classes":""}{"FieldID":"10","Value":"999","FieldName":"CheckSum","FieldType":"STRING","DecodedValue":"Invalid(expected036)"}`

	if actual != expect {
		t.Errorf("expect %s, actual %s", expect, actual)
	}
}

func TestFixDecoder_Invalid_BodyLength(t *testing.T) {
	actual := r.Replace(fd.Decode(invalidfixmessage_bodylength).String())
	expect := `{"FieldID":"8","Value":"FIX.4.4","FieldName":"BeginString","FieldType":"STRING","Classes":""}{"FieldID":"9","Value":"88","FieldName":"BodyLength","FieldType":"LENGTH","DecodedValue":"Invalid(expected74)"}{"FieldID":"35","Value":"2","FieldName":"MsgType","FieldType":"STRING","DecodedValue":"ResendRequest"}{"FieldID":"49","Value":"CNX","FieldName":"SenderCompID","FieldType":"STRING","Classes":""}{"FieldID":"34","Value":"8263336","FieldName":"MsgSeqNum","FieldType":"SEQNUM","Classes":""}{"FieldID":"52","Value":"20180126-07:39:59.683","FieldName":"SendingTime","FieldType":"UTCTIMESTAMP","Classes":""}{"FieldID":"56","Value":"imdstream","FieldName":"TargetCompID","FieldType":"STRING","Classes":""}{"FieldID":"16","Value":"0","FieldName":"EndSeqNo","FieldType":"SEQNUM","Classes":""}{"FieldID":"7","Value":"12812","FieldName":"BeginSeqNo","FieldType":"SEQNUM","Classes":""}{"FieldID":"10","Value":"036","FieldName":"CheckSum","FieldType":"STRING","DecodedValue":"Invalid(expected041)"}`

	if actual != expect {
		t.Errorf("expect %s, actual %s", expect, actual)
	}
}
