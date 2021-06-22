package utils

import (
	"os"
	"testing"
)

func TestZipCompress(t *testing.T) {
	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "when src exist then zip compress",
			args: args{
				src: "../../example/log.yaml",
				dst: "../../example/log.yaml.zip",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ZipCompress(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("ZipCompress() error = %v, wantErr %v", err, tt.wantErr)
			}
			os.Remove(tt.args.dst)
		})
	}
}
