package password

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptor_EncodeString(t *testing.T) {
	tests := []struct {
		key        string
		s          string
		wantString string
	}{
		{
			key:        "",
			s:          "abc",
			wantString: "fd7adb152c05ef80dccf50a1fa4c05d5a3ec6da95575fc312ae7c5d091836351",
		},
		{
			key:        "abcdefghijklmnopqrstvuwxhz",
			s:          "abc",
			wantString: "d54891ec2b68e799a6f9c813b29217183fbb5cf088b80b3f9dd5edce7198dadd",
		},
		{
			key:        "abcdefghijklmnopqrstvuwxhz",
			s:          "abcdefghijklmnopqrstvuwxhz",
			wantString: "d294cc19f820e5047abee86855d8d19f235109a9ef6f8d26523645d5ac447b6c",
		},
		{
			key:        "abcdefghijklmnopqrstvuwxhzabcdefghijklmnopqrstvuwxhzabcdefghijklmnopqrstvuwxhz",
			s:          "abcdefghijklmnopqrstvuwxhz",
			wantString: "98b8b843481781fb8ef1a90d99e3227cbcef6c471db665ef9a6177728f572707",
		},
		{
			key:        "abcdefghijklmnopqrstvuwxhzabcdefghijklmnopqrstvuwxhzabcdefghijklmnopqrstvuwxhz",
			s:          "abcdefghijklmnopqrstvuwxhzabcdefghijklmnopqrstvuwxhzabcdefghijklmnopqrstvuwxhz",
			wantString: "33600f5e901eeef18528645750282adb7a1df1b3fe57a1e0f967a5a67e64f1de",
		},
	}
	for _, tt := range tests {
		t.Run(tt.key+"_"+tt.s, func(t *testing.T) {
			p := NewEncryptorByString(tt.key)
			res, err := p.EncodeString(tt.s)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.wantString, string(res))
			}
		})
	}
}
