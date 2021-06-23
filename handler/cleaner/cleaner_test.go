package cleaner

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestTimePattern(t *testing.T) {
	t.Run("when time pattern then return true", func(t *testing.T) {
		if m := defaultTimePattern.MatchString(defaultTimeFormat); !m {
			t.Errorf("exptect: true,  got: %v", m)
		}
	})
	t.Run("when time pattern with baseName then return true", func(t *testing.T) {
		if m := defaultTimePattern.MatchString("info.2021-06-22.223730.log.zip"); !m {
			t.Errorf("exptect: true,  got: %v", m)
		}
	})

	t.Run("when match then extract", func(t *testing.T) {
		match := defaultTimePattern.FindStringSubmatch("info.2021-06-22.223730.log.zip")
		if len(match) == 0 || match[1] != "2021-06-22.223730" {
			t.Errorf("exptect true,  got: false")
		}
	})

	t.Run("when not match then do nothing", func(t *testing.T) {
		match := defaultTimePattern.FindStringSubmatch("info.2021-096-22.223730.log.zip")
		if len(match) != 0 {
			t.Errorf("exptect not match,  got: match")
		}
	})
}

func TestFileCleaner_expiredFiles(t *testing.T) {
	now := time.Now()

	type fields struct {
		backupDir      string
		backupBaseName string
		backupExt      string
		backupCount    int
		backupDuration time.Duration
		checkInterval  time.Duration
		mutex          *sync.Mutex
	}
	type args struct {
		fileList []FileInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []FileInfo
	}{
		{
			name: "when files empty then empty result",
			fields: fields{
				backupDuration: 1 * time.Second,
			},
			args: args{
				fileList: []FileInfo{},
			},
			want: []FileInfo{},
		},
		{
			name: "when files nil then empty result",
			fields: fields{
				backupDuration: 1 * time.Second,
			},
			args: args{
				fileList: []FileInfo{},
			},
			want: []FileInfo{},
		},
		{
			name: "when files not empty then get result",
			fields: fields{
				backupDuration: 5 * time.Second,
			},
			args: args{
				fileList: []FileInfo{
					{BackupTime: now},
					{BackupTime: now.Add(-10 * time.Second)},
					{BackupTime: now.Add(-11 * time.Second)},
					{BackupTime: now.Add(-2 * time.Second)},
					{BackupTime: now.Add(2 * time.Second)},
				},
			},
			want: []FileInfo{
				{BackupTime: now.Add(-10 * time.Second)},
				{BackupTime: now.Add(-11 * time.Second)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &LogCleaner{
				backupDir:      tt.fields.backupDir,
				backupBaseName: tt.fields.backupBaseName,
				backupExt:      tt.fields.backupExt,
				backupCount:    tt.fields.backupCount,
				backupDuration: tt.fields.backupDuration,
				checkInterval:  tt.fields.checkInterval,
				mutex:          tt.fields.mutex,
			}
			if got := fc.expiredFiles(tt.args.fileList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expiredFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileCleaner_gtBackupLimit(t *testing.T) {
	now := time.Now()
	type fields struct {
		backupDir      string
		backupBaseName string
		backupExt      string
		backupCount    int
		backupDuration time.Duration
		checkInterval  time.Duration
		mutex          *sync.Mutex
	}
	type args struct {
		fileList []FileInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []FileInfo
	}{
		{
			name: "when files greater than return oldest files",
			fields: fields{
				backupCount: 3,
			},
			args: args{
				fileList: []FileInfo{
					{BackupTime: now},
					{BackupTime: now.Add(-10 * time.Second)},
					{BackupTime: now.Add(-11 * time.Second)},
					{BackupTime: now.Add(-2 * time.Second)},
					{BackupTime: now.Add(2 * time.Second)},
				},
			},
			want: []FileInfo{
				{BackupTime: now.Add(-11 * time.Second)},
				{BackupTime: now.Add(-10 * time.Second)},
			},
		},
		{
			name: "when files not greater than limit return empty",
			fields: fields{
				backupCount: 5,
			},
			args: args{
				fileList: []FileInfo{
					{BackupTime: now},
					{BackupTime: now.Add(-10 * time.Second)},
					{BackupTime: now.Add(-11 * time.Second)},
					{BackupTime: now.Add(-2 * time.Second)},
					{BackupTime: now.Add(2 * time.Second)},
				},
			},
			want: []FileInfo{},
		},
		{
			name: "when files is nil return empty",
			fields: fields{
				backupCount: 5,
			},
			args: args{
				fileList: nil,
			},
			want: []FileInfo{},
		},
		{
			name: "when files is empty return empty",
			fields: fields{
				backupCount: 5,
			},
			args: args{
				fileList: nil,
			},
			want: []FileInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &LogCleaner{
				backupDir:      tt.fields.backupDir,
				backupBaseName: tt.fields.backupBaseName,
				backupExt:      tt.fields.backupExt,
				backupCount:    tt.fields.backupCount,
				backupDuration: tt.fields.backupDuration,
				checkInterval:  tt.fields.checkInterval,
				mutex:          tt.fields.mutex,
			}
			if got := fc.gtBackupLimit(tt.args.fileList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("gtBackupLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileCleaner_deduplicateByFilename(t *testing.T) {
	type fields struct {
		backupDir      string
		backupBaseName string
		backupExt      string
		backupCount    int
		backupDuration time.Duration
		checkInterval  time.Duration
		mutex          *sync.Mutex
	}
	type args struct {
		fullList []FileInfo
		sublist  []FileInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []FileInfo
	}{
		{
			name: "when contains duplicated then remove",
			fields: fields{
				backupCount: 3,
			},
			args: args{
				fullList: []FileInfo{
					{Filename: "1"},
					{Filename: "2"},
					{Filename: "3"},
					{Filename: "4"},
					{Filename: "5"},
				},
				sublist: []FileInfo{
					{Filename: "1"},
					{Filename: "2"},
					{Filename: "3"},
					{Filename: "4"},
				},
			},
			want: []FileInfo{
				{Filename: "5"},
			},
		},
		{
			name: "when full is empty then return empty",
			fields: fields{
				backupCount: 3,
			},
			args: args{
				fullList: []FileInfo{},
				sublist: []FileInfo{
					{Filename: "5"},
				},
			},
			want: []FileInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &LogCleaner{
				backupDir:      tt.fields.backupDir,
				backupBaseName: tt.fields.backupBaseName,
				backupExt:      tt.fields.backupExt,
				backupCount:    tt.fields.backupCount,
				backupDuration: tt.fields.backupDuration,
				checkInterval:  tt.fields.checkInterval,
				mutex:          tt.fields.mutex,
			}
			if got := fc.deduplicateByFilename(tt.args.fullList, tt.args.sublist); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deduplicateByFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseName(t *testing.T) {
	t.Run("when filename then get basename", func(t *testing.T) {
		filename := "info.log"
		ext := ".log"
		baseName := filename[:len(filename)-len(ext)]
		t.Log(baseName)
		t.Log(0 * time.Second)
		if baseName != "info" {
			t.Errorf("want; %s, got: %s", "info", baseName)
		}
	})
}
