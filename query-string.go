package urlbuilder

import (
	"fmt"
	"strings"
)

// NewQueryString returns a new QueryString with the provided query string parts
func NewQueryString(mode QueryMode, qp ...QueryPart) *QueryString {
	qs := QueryString{
		mode: mode,
	}
	for _, p := range qp {
		p(&qs)
	}
	return &qs
}

// Nv appends a name-value parameter to query string
func Nv(name string, value any) QueryPart {
	return func(qs *QueryString) error {
		v := fmt.Sprint(value)
		qs.qrs = append(qs.qrs, query{name: name, value: v})
		return nil
	}
}

// Build constructs the query string. Returns an empty string if an error occurred.
func (qs *QueryString) Build() string {
	if len(qs.qrs) == 0 {
		return ""
	}
	var b strings.Builder
	first := true
	if qs.mode == QModeLast || qs.mode == QModeError {
		qmap := make(map[string]string)
		for _, q := range qs.qrs {
			if _, found := qmap[q.name]; found && qs.mode == QModeError {
				qs.err = fmt.Errorf("duplicate query name found")
				return ""
			}
			qmap[q.name] = q.value
		}
		for k, v := range qmap {
			if !first {
				b.WriteByte('&')
			}
			first = false
			b.WriteString(k)
			b.WriteByte('=')
			b.WriteString(escape(v))
		}
	} else {
		for i, q := range qs.qrs {
			if i > 0 {
				b.WriteByte('&')
			}
			b.WriteString(q.name)
			b.WriteByte('=')
			b.WriteString(escape(q.value))
		}
	}
	return b.String()
}

// String implements fmt.Stringer and returns the built query string.
func (qs *QueryString) String() string {
	return qs.Build()
}
