package utils

import "testing"

func TestLastSubstring(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "when multiple string then return last one",
			args: args{
				s:   "com.github.etlog",
				sep: ".",
			},
			want: "etlog",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LastSubstring(tt.args.s, tt.args.sep); got != tt.want {
				t.Errorf("LastSubstring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSize(t *testing.T) {
	type args struct {
		fileSize string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "when 1k then return 1024",
			args: args{
				fileSize: "1k",
			},
			want:    1024,
			wantErr: false,
		},
		{
			name: "when 1K then return 1024",
			args: args{
				fileSize: "1K",
			},
			want:    1024,
			wantErr: false,
		},
		{
			name: "when 12K then return 1024",
			args: args{
				fileSize: "12K",
			},
			want:    12 * 1024,
			wantErr: false,
		},
		{
			name: "when 1m then return 1024*1024",
			args: args{
				fileSize: "1m",
			},
			want:    1024 * 1024,
			wantErr: false,
		},
		{
			name: "when 1M then return 1024*1024",
			args: args{
				fileSize: "1M",
			},
			want:    1024 * 1024,
			wantErr: false,
		},
		{
			name: "when 100M then return 1024*1024",
			args: args{
				fileSize: "100M",
			},
			want:    100 * 1024 * 1024,
			wantErr: false,
		},
		{
			name: "when 1g then return 1024*1024*1024",
			args: args{
				fileSize: "1g",
			},
			want:    1024 * 1024 * 1024,
			wantErr: false,
		},
		{
			name: "when 1G then return 1024*1024",
			args: args{
				fileSize: "1G",
			},
			want:    1024 * 1024 * 1024,
			wantErr: false,
		},
		{
			name: "when 100G then return 1024*1024",
			args: args{
				fileSize: "100G",
			},
			want:    100 * 1024 * 1024 * 1024,
			wantErr: false,
		},
		{
			name: "when 12 then return error",
			args: args{
				fileSize: "12",
			},
			want:    -1,
			wantErr: true,
		},
		{
			name: "when 1KB then return error",
			args: args{
				fileSize: "1KB",
			},
			want:    -1,
			wantErr: true,
		},
		{
			name: "when 1T then return error",
			args: args{
				fileSize: "1T",
			},
			want:    -1,
			wantErr: true,
		},
		{
			name: "when blank then return error",
			args: args{
				fileSize: "",
			},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSize(tt.args.fileSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseSize() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSeconds(t *testing.T) {
	type args struct {
		interval string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "when 1m then return 60",
			args: args{
				interval: "1m",
			},
			want:    60,
			wantErr: false,
		},
		{
			name: "when 1M then return 60",
			args: args{
				interval: "1M",
			},
			want:    60,
			wantErr: false,
		},
		{
			name: "when 12M then return 12*60",
			args: args{
				interval: "12M",
			},
			want:    12 * 60,
			wantErr: false,
		},

		{
			name: "when 1h then return 60*60",
			args: args{
				interval: "1h",
			},
			want:    60 * 60,
			wantErr: false,
		},
		{
			name: "when 1H then return 60*60",
			args: args{
				interval: "1H",
			},
			want:    60 * 60,
			wantErr: false,
		},
		{
			name: "when 12H then return 12*60*60",
			args: args{
				interval: "12H",
			},
			want:    12 * 60 * 60,
			wantErr: false,
		},

		{
			name: "when 1d then return 24*60*60",
			args: args{
				interval: "1d",
			},
			want:    24 * 60 * 60,
			wantErr: false,
		},
		{
			name: "when 1D then return 24*60*60",
			args: args{
				interval: "1D",
			},
			want:    24 * 60 * 60,
			wantErr: false,
		},
		{
			name: "when 12D then return 12*24*60*60",
			args: args{
				interval: "12D",
			},
			want:    12 * 24 * 60 * 60,
			wantErr: false,
		},
		{
			name: "when blank then return error",
			args: args{
				interval: "",
			},
			want:    -1,
			wantErr: true,
		},
		{
			name: "when 12 then return error",
			args: args{
				interval: "12",
			},
			want:    -1,
			wantErr: true,
		},
		{
			name: "when 1w then return error",
			args: args{
				interval: "1w",
			},
			want:    -1,
			wantErr: true,
		},
		{
			name: "when 1ms then return error",
			args: args{
				interval: "1ms",
			},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSeconds(tt.args.interval)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSeconds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseSeconds() got = %v, want %v", got, tt.want)
			}
		})
	}
}
