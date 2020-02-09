package packets

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewFixedHeader(t *testing.T) {

	type args struct {
		header uint8
	}
	tests := []struct {
		name    string
		args    args
		want    *FixedHeader
		wantErr bool
	}{
		{name: "Reserved Test",
			args: args{
				header: 0,
			},
			want: &FixedHeader{
				Type:  Reserved,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Connection Request",
			args: args{
				header: 16,
			},
			want: &FixedHeader{
				Type:  CONNECT,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Connection Acknowledgement",
			args: args{
				header: 32,
			},
			want: &FixedHeader{
				Type:  CONNACK,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Publish Message",
			args: args{
				header: 48,
			},
			want: &FixedHeader{
				Type:  PUBLISH,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Publish Acknowledgement",
			args: args{
				header: 64,
			},
			want: &FixedHeader{
				Type:  PUBACK,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Publish Recieved",
			args: args{
				header: 80,
			},
			want: &FixedHeader{
				Type:  PUBREC,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Publish Release",
			args: args{
				header: 96,
			},
			want: &FixedHeader{
				Type:  PUBREL,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Publish Complete",
			args: args{
				header: 112,
			},
			want: &FixedHeader{
				Type:  PUBCOMP,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Subscribe Request",
			args: args{
				header: 128,
			},
			want: &FixedHeader{
				Type:  SUBSCRIBE,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Subscribe Acknowlegement",
			args: args{
				header: 144,
			},
			want: &FixedHeader{
				Type:  SUBACK,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Unsubscribe Request",
			args: args{
				header: 160,
			},
			want: &FixedHeader{
				Type:  UNSUBSCRIBE,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Unsubscribe Acknowledgement",
			args: args{
				header: 176,
			},
			want: &FixedHeader{
				Type:  UNSUBACK,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Ping Request",
			args: args{
				header: 192,
			},
			want: &FixedHeader{
				Type:  PINGREQ,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Ping Repsonse",
			args: args{
				header: 208,
			},
			want: &FixedHeader{
				Type:  PINGRESP,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Disconnect Notification",
			args: args{
				header: 224,
			},
			want: &FixedHeader{
				Type:  DISCONNECT,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Authentication Exchange",
			args: args{
				header: 240,
			},
			want: &FixedHeader{
				Type:  AUTH,
				Flags: FixedHeaderFlags{},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFixedHeader(tt.args.header, tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFixedHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFixedHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeByte(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    byte
		want1   int
		wantErr bool
	}{
		{
			name: "Decode a byte",
			args: args{
				b: []byte{1},
			},
			want:    1,
			want1:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DecodeByte(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeByte() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DecodeByte() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDecodeFourByteInt(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    uint32
		want1   int
		wantErr bool
	}{
		{
			name: "Decode a four byte integer",
			args: args{
				b: []byte{0, 0, 0, 16},
			},
			want:    16,
			want1:   4,
			wantErr: false,
		},
		{
			name: "Decode a four byte integer",
			args: args{
				b: []byte{0, 0, 0, 32},
			},
			want:    32,
			want1:   4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DecodeFourByteInt(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeFourByteInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeFourByteInt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DecodeFourByteInt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDecodeTwoByteInt(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		want1   int
		wantErr bool
	}{
		{
			name: "Decode a two byte integer",
			args: args{
				b: []byte{0, 16},
			},
			want:    16,
			want1:   2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DecodeTwoByteInt(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeTwoByteInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeTwoByteInt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DecodeTwoByteInt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDecodeBinaryData(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		want1   int
		wantErr bool
	}{
		{
			name: "Decode binary data",
			args: args{
				b: []byte{0, 4, 2, 3, 4, 5},
			},
			want:    []byte{2, 3, 4, 5},
			want1:   6,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DecodeBinaryData(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeBinaryData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeBinaryData() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DecodeBinaryData() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDecodeStringPair(t *testing.T) {
	testByteArray := []byte{0, 16}
	testByteArray = append(testByteArray, []byte("Here is a string")...)
	testByteArray = append(testByteArray, 0, 23)
	testByteArray = append(testByteArray, []byte("Here is another string")...)
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *StringPair
		want1   int
		wantErr bool
	}{
		{
			name: "Decode a two byte integer",
			args: args{
				b: testByteArray,
			},
			want: &StringPair{
				name:  "Here is a string",
				value: "Here is another string",
			},
			want1:   16 + 23 + 4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DecodeStringPair(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeStringPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.name != tt.want.name && got.value != tt.want.value {
				t.Errorf("DecodeStringPair() got = %v, want %v", got, tt.want.name)
			}
			if got1 != tt.want1 {
				t.Errorf("DecodeStringPair() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDecodeString(t *testing.T) {
	testByteArray := []byte{0, 16}
	testByteArray = append(testByteArray, []byte("Here is a string")...)
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   int
		wantErr bool
	}{
		{
			name: "Decode a two byte integer",
			args: args{
				b: testByteArray,
			},
			want:    "Here is a string",
			want1:   18,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DecodeString(tt.args.b)
			fmt.Println(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DecodeString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
