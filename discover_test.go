package mystrom

import (
	"net"
	"reflect"
	"testing"
)

func Test_defaultString(t *testing.T) {
	type args struct {
		s   string
		def string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple", args{"defined", "default"}, "defined"},
		{"default", args{"", "default"}, "default"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultString(tt.args.s, tt.args.def); got != tt.want {
				t.Errorf("defaultString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDevicePayload(t *testing.T) {
	type args struct {
		buf  []byte
		addr net.Addr
	}
	tests := []struct {
		name    string
		args    args
		wantD   Device
		wantErr bool
	}{
		{
			name: "valid broadcast packet",
			args: args{buf: []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, byte(DeviceTypeBulb), 0x01}, addr: nil},
			wantD: Device{
				MAC:        net.HardwareAddr{0x01, 0x23, 0x45, 0x67, 0x89, 0xab},
				Address:    nil,
				Type:       DeviceTypeBulb,
				MeshChild:  true,
				Registered: false,
				Cloud:      false,
			},
			wantErr: false,
		},
		{
			name:    "invalid packet length",
			args:    args{buf: []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0x01}, addr: nil},
			wantD:   Device{},
			wantErr: true,
		},
		{
			name: "valid broadcast packet with all flags set",
			args: args{buf: []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, byte(DeviceTypeGateway), 0x07}, addr: nil},
			wantD: Device{
				MAC:        net.HardwareAddr{0x01, 0x23, 0x45, 0x67, 0x89, 0xab},
				Address:    nil,
				Type:       DeviceTypeGateway,
				MeshChild:  true,
				Registered: true,
				Cloud:      true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotD, err := parseDevicePayload(tt.args.buf, tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDevicePayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("parseDevicePayload() = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}
