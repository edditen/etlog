package utils

import "testing"

func TestShortSourceLoc(t *testing.T) {
	type args struct {
		skip int
	}
	type testCase struct {
		args         args
		wantLine     string
		wantFuncName string
		wantOk       bool
	}

	t.Run("when skip is 0 then return utils loc", func(t *testing.T) {
		tc := testCase{
			args: args{
				skip: 0,
			},
			wantOk:       true,
			wantLine:     "sourceutils.go:13",
			wantFuncName: "utils.SourceLoc",
		}
		gotLine, gotFuncName, gotOk := ShortSourceLoc(tc.args.skip)
		t.Log(gotLine, gotFuncName, gotOk)
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
			wantLine:     "sourceutils.go:24",
			wantFuncName: "utils.ShortSourceLoc",
		}
		gotLine, gotFuncName, gotOk := ShortSourceLoc(tc.args.skip)
		t.Log(gotLine, gotFuncName, gotOk)
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
			wantLine:     "sourceutils_test.go:69",
			wantFuncName: "utils.TestShortSourceLoc.func3",
		}
		gotLine, gotFuncName, gotOk := ShortSourceLoc(tc.args.skip)
		t.Log(gotLine, gotFuncName, gotOk)
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
			wantLine:     "sourceutils_test.go:92",
			wantFuncName: "utils.TestShortSourceLoc.func4.1",
		}
		gotLine, gotFuncName, gotOk := func() (string, string, bool) {
			return ShortSourceLoc(tc.args.skip)
		}()
		t.Log(gotLine, gotFuncName, gotOk)
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
