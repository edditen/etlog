package opt

type LogFunc func(e *LogE)

type LogE struct {
	Time     int64                  `json:"time,omitempty"`
	Level    string                 `json:"level,omitempty"`
	SrcFile  string                 `json:"srcf,omitempty"`
	Line     int                    `json:"line,omitempty"`
	FuncName string                 `json:"func,omitempty"`
	Msg      string                 `json:"msg,omitempty"`
	Marker   string                 `json:"marker,omitempty"`
	Err      error                  `json:"error,omitempty"`
	Fields   map[string]interface{} `json:"fields,omitempty"`
}
