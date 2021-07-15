package utils

import "testing"

func TestShortSourceLoc(t *testing.T) {
	type args struct {
		skip int
	}
	type testCase struct {
		args         args
		wantFileName string
		wantLine     int
		wantFuncName string
		wantOk       bool
	}

	t.Run("when skip is 0 then return utils loc", func(t *testing.T) {
		tc := testCase{
			args: args{
				skip: 0,
			},
			wantOk:       true,
			wantFileName: "sourceutils.go",
			wantLine:     12,
			wantFuncName: "utils.SourceLoc",
		}
		gotFileName, gotLine, gotFuncName, gotOk := ShortSourceLoc(tc.args.skip)
		t.Log(gotFileName, gotLine, gotFuncName, gotOk)
		if gotFileName != tc.wantFileName {
			t.Errorf("ShortSourceLoc() gotFileName = %v, want %v", gotLine, tc.wantLine)
		}

		if gotLine != tc.wantLine {
			t.Errorf("ShortSourceLoc() gotLine = %v, want %v", gotLine, tc.wantLine)
		}

		if gotFuncName != tc.wantFuncName {
			t.Errorf("ShortSourceLoc() gotFuncName = %v, want %v", gotFuncName, tc.wantFuncName)
		}
		if gotOk != tc.wantOk {
			t.Errorf("ShortSourceLoc() gotOk = %v, want %v", gotOk, tc.wantOk)
		}
	})

	t.Run("when skip is 1 then return utils loc", func(t *testing.T) {
		tc := testCase{
			args: args{
				skip: 1,
			},
			wantOk:       true,
			wantFileName: "sourceutils.go",
			wantLine:     22,
			wantFuncName: "utils.ShortSourceLoc",
		}
		gotFileName, gotLine, gotFuncName, gotOk := ShortSourceLoc(tc.args.skip)
		t.Log(gotFileName, gotLine, gotFuncName, gotOk)
		if gotFileName != tc.wantFileName {
			t.Errorf("ShortSourceLoc() gotFileName = %v, want %v", gotLine, tc.wantLine)
		}

		if gotLine != tc.wantLine {
			t.Errorf("ShortSourceLoc() gotLine = %v, want %v", gotLine, tc.wantLine)
		}
		if gotFuncName != tc.wantFuncName {
			t.Errorf("ShortSourceLoc() gotFuncName = %v, want %v", gotFuncName, tc.wantFuncName)
		}
		if gotOk != tc.wantOk {
			t.Errorf("ShortSourceLoc() gotOk = %v, want %v", gotOk, tc.wantOk)
		}
	})

	t.Run("when skip is 2 then return test loc", func(t *testing.T) {
		tc := testCase{
			args: args{
				skip: 2,
			},
			wantOk:       true,
			wantFileName: "sourceutils_test.go",
			wantLine:     82,
			wantFuncName: "utils.TestShortSourceLoc.func3",
		}
		gotFileName, gotLine, gotFuncName, gotOk := ShortSourceLoc(tc.args.skip)
		t.Log(gotFileName, gotLine, gotFuncName, gotOk)
		if gotFileName != tc.wantFileName {
			t.Errorf("ShortSourceLoc() gotFileName = %v, want %v", gotLine, tc.wantLine)
		}

		if gotLine != tc.wantLine {
			t.Errorf("ShortSourceLoc() gotLine = %v, want %v", gotLine, tc.wantLine)
		}
		if gotFuncName != tc.wantFuncName {
			t.Errorf("ShortSourceLoc() gotFuncName = %v, want %v", gotFuncName, tc.wantFuncName)
		}
		if gotOk != tc.wantOk {
			t.Errorf("ShortSourceLoc() gotOk = %v, want %v", gotOk, tc.wantOk)
		}
	})

	t.Run("when skip is 3 then return test loc", func(t *testing.T) {
		tc := testCase{
			args: args{
				skip: 2,
			},
			wantOk:       true,
			wantFileName: "sourceutils_test.go",
			wantLine:     109,
			wantFuncName: "utils.TestShortSourceLoc.func4",
		}
		gotFileName, gotLine, gotFuncName, gotOk := ShortSourceLoc(tc.args.skip)
		t.Log(gotFileName, gotLine, gotFuncName, gotOk)
		if gotFileName != tc.wantFileName {
			t.Errorf("ShortSourceLoc() gotFileName = %v, want %v", gotLine, tc.wantLine)
		}

		if gotLine != tc.wantLine {
			t.Errorf("ShortSourceLoc() gotLine = %v, want %v", gotLine, tc.wantLine)
		}
		if gotFuncName != tc.wantFuncName {
			t.Errorf("ShortSourceLoc() gotFuncName = %v, want %v", gotFuncName, tc.wantFuncName)
		}
		if gotOk != tc.wantOk {
			t.Errorf("ShortSourceLoc() gotOk = %v, want %v", gotOk, tc.wantOk)
		}
	})
}
