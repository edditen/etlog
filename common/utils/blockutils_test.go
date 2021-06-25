package utils

import "testing"

func TestCalculateBlocks(t *testing.T) {
	type args struct {
		totalSize int
		blockSize int
	}
	tests := []struct {
		name       string
		args       args
		wantBlocks int
	}{
		{
			name: "when total is times of blockSize then get total/blockSize",
			args: args{
				totalSize: 10,
				blockSize: 5,
			},
			wantBlocks: 2,
		},
		{
			name: "when total is not times of blockSize then get total/blockSize+1",
			args: args{
				totalSize: 10,
				blockSize: 4,
			},
			wantBlocks: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotBlocks := CalculateBlocks(tt.args.totalSize, tt.args.blockSize); gotBlocks != tt.wantBlocks {
				t.Errorf("CalculateBlocks() = %v, want %v", gotBlocks, tt.wantBlocks)
			}
		})
	}
}
