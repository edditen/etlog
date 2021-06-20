package core

import "testing"

func TestFields_String(t *testing.T) {
	tests := []struct {
		name string
		f    Fields
		want string
	}{
		{
			name: "when nil then return blank",
			f:    nil,
			want: "",
		},
		{
			name: "when empty then return blank",
			f:    map[string]interface{}{},
			want: "",
		},
		{
			name: "when fields then return json",
			f:    map[string]interface{}{"hello": "world"},
			want: "{\"hello\":\"world\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
