package packets

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewPacket(t *testing.T) {

	type args struct {
		header uint8
	}
	tests := []struct {
		name    string
		args    args
		want    *Packet
		wantErr bool
	}{
		{name: "Reserved Test",
			args: args{
				header: 0,
			},
			want: &Packet{
				Type:  Reserved,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Connection Request",
			args: args{
				header: 16,
			},
			want: &Packet{
				Type:  CONNECT,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Connection Acknowledgement",
			args: args{
				header: 32,
			},
			want: &Packet{
				Type:  CONNACK,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Publish Message",
			args: args{
				header: 48,
			},
			want: &Packet{
				Type:  PUBLISH,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Publish Acknowledgement",
			args: args{
				header: 64,
			},
			want: &Packet{
				Type:  PUBACK,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Publish Recieved",
			args: args{
				header: 80,
			},
			want: &Packet{
				Type:  PUBREC,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Publish Release",
			args: args{
				header: 96,
			},
			want: &Packet{
				Type:  PUBREL,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Publish Complete",
			args: args{
				header: 112,
			},
			want: &Packet{
				Type:  PUBCOMP,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Subscribe Request",
			args: args{
				header: 128,
			},
			want: &Packet{
				Type:  SUBSCRIBE,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Subscribe Acknowlegement",
			args: args{
				header: 144,
			},
			want: &Packet{
				Type:  SUBACK,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Unsubscribe Request",
			args: args{
				header: 160,
			},
			want: &Packet{
				Type:  UNSUBSCRIBE,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Unsubscribe Acknowledgement",
			args: args{
				header: 176,
			},
			want: &Packet{
				Type:  UNSUBACK,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Ping Request",
			args: args{
				header: 192,
			},
			want: &Packet{
				Type:  PINGREQ,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Ping Repsonse",
			args: args{
				header: 208,
			},
			want: &Packet{
				Type:  PINGRESP,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Disconnect Notification",
			args: args{
				header: 224,
			},
			want: &Packet{
				Type:  DISCONNECT,
				Flags: FixedHeaderFlags{},
			}},
		{name: "Authentication Exchange",
			args: args{
				header: 240,
			},
			want: &Packet{
				Type:  AUTH,
				Flags: FixedHeaderFlags{},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMQTTPacket([]byte{tt.args.header, 0})
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Type != tt.want.Type {
				t.Errorf("Type = %v, want %v", got.Type, tt.want.Type)
			}
		})
	}
}

func TestDecodeByte(t *testing.T) {
	type args struct {
		b *bytes.Buffer
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
				b: bytes.NewBuffer([]byte{1}),
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{buff: tt.args.b}
			got, err := p.DecodeByte()
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeByte() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeFourByteInt(t *testing.T) {
	type args struct {
		b *bytes.Buffer
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
				b: bytes.NewBuffer([]byte{0, 0, 0, 16}),
			},
			want:    16,
			want1:   4,
			wantErr: false,
		},
		{
			name: "Decode a four byte integer",
			args: args{
				b: bytes.NewBuffer([]byte{0, 0, 0, 32}),
			},
			want:    32,
			want1:   4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{buff: tt.args.b}
			got := p.DecodeFourByteInt()
			if got != tt.want {
				t.Errorf("DecodeFourByteInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeTwoByteInt(t *testing.T) {
	type args struct {
		b *bytes.Buffer
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
				b: bytes.NewBuffer([]byte{0, 16}),
			},
			want:    16,
			want1:   2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{buff: tt.args.b}
			got := p.DecodeTwoByteInt()
			if got != tt.want {
				t.Errorf("DecodeTwoByteInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeBinaryData(t *testing.T) {
	type args struct {
		b *bytes.Buffer
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
				b: bytes.NewBuffer([]byte{0, 4, 2, 3, 4, 5}),
			},
			want:    []byte{2, 3, 4, 5},
			want1:   6,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{buff: tt.args.b}
			got := p.DecodeBinaryData()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeBinaryData() got = %v, want %v", got, tt.want)
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
		b *bytes.Buffer
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
				b: bytes.NewBuffer(testByteArray),
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
			p := &Packet{buff: tt.args.b}
			got := p.DecodeStringPair()
			if got.name != tt.want.name && got.value != tt.want.value {
				t.Errorf("DecodeStringPair() got = %v, want %v", got, tt.want.name)
			}
		})
	}
}

func TestDecodeString(t *testing.T) {
	testByteArray := []byte{0, 16}
	testByteArray = append(testByteArray, []byte("Here is a string")...)
	type args struct {
		b *bytes.Buffer
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
				b: bytes.NewBuffer(testByteArray),
			},
			want:    "Here is a string",
			want1:   18,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Packet{buff: tt.args.b}
			got := p.DecodeString()

			if got != tt.want {
				t.Errorf("DecodeString() got = %v, want %v", got, tt.want)
			}

		})
	}
}
