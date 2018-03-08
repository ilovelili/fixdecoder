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
	expect := `{"ID":"8","Name":"BeginString","Value":"FIX.4.4"}{"ID":"9","Name":"BodyLength","Value":"74","DecodedValue":"Valid"}{"ID":"35","Name":"MsgType","Value":"2","DecodedValue":"ResendRequest"}{"ID":"49","Name":"SenderCompID","Value":"CNX"}{"ID":"34","Name":"MsgSeqNum","Value":"8263336"}{"ID":"52","Name":"SendingTime","Value":"20180126-07:39:59.683"}{"ID":"56","Name":"TargetCompID","Value":"imdstream"}{"ID":"16","Name":"EndSeqNo","Value":"0"}{"ID":"7","Name":"BeginSeqNo","Value":"12812"}{"ID":"10","Name":"CheckSum","Value":"036","DecodedValue":"Valid"}`

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
	expect := `{"ID":"8","Name":"BeginString","Value":"FIX.4.4"}{"ID":"9","Name":"BodyLength","Value":"74","DecodedValue":"Valid"}{"ID":"35","Name":"MsgType","Value":"2","DecodedValue":"ResendRequest"}{"ID":"49","Name":"SenderCompID","Value":"CNX"}{"ID":"34","Name":"MsgSeqNum","Value":"8263336"}{"ID":"52","Name":"SendingTime","Value":"20180126-07:39:59.683"}{"ID":"56","Name":"TargetCompID","Value":"imdstream"}{"ID":"16","Name":"EndSeqNo","Value":"0"}{"ID":"7","Name":"BeginSeqNo","Value":"12812"}{"ID":"10","Name":"CheckSum","Value":"999","DecodedValue":"Invalid(expected036)"}`

	if actual != expect {
		t.Errorf("expect %s, actual %s", expect, actual)
	}
}

func TestFixDecoder_Invalid_BodyLength(t *testing.T) {
	actual := r.Replace(fd.Decode(invalidfixmessage_bodylength).String())
	expect := `{"ID":"8","Name":"BeginString","Value":"FIX.4.4"}{"ID":"9","Name":"BodyLength","Value":"88","DecodedValue":"Invalid(expected74)"}{"ID":"35","Name":"MsgType","Value":"2","DecodedValue":"ResendRequest"}{"ID":"49","Name":"SenderCompID","Value":"CNX"}{"ID":"34","Name":"MsgSeqNum","Value":"8263336"}{"ID":"52","Name":"SendingTime","Value":"20180126-07:39:59.683"}{"ID":"56","Name":"TargetCompID","Value":"imdstream"}{"ID":"16","Name":"EndSeqNo","Value":"0"}{"ID":"7","Name":"BeginSeqNo","Value":"12812"}{"ID":"10","Name":"CheckSum","Value":"036","DecodedValue":"Invalid(expected041)"}`

	if actual != expect {
		t.Errorf("expect %s, actual %s", expect, actual)
	}
}
