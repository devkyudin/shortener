package config

import (
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetAddress_Set(t *testing.T) {
	type fields struct {
		Protocol string
		Host     string
		Port     string
	}
	type args struct {
		rawAddr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := &NetAddress{
				Protocol: tt.fields.Protocol,
				Host:     tt.fields.Host,
				Port:     tt.fields.Port,
			}
			if err := addr.Set(tt.args.rawAddr); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetAddress_String(t *testing.T) {
	type fields struct {
		Protocol string
		Host     string
		Port     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr := &NetAddress{
				Protocol: tt.fields.Protocol,
				Host:     tt.fields.Host,
				Port:     tt.fields.Port,
			}
			if got := addr.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setConfigByFlags(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setConfigByFlags()
		})
	}
}
