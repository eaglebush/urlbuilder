// Package urlbuilder provides a flexible and composable way to build URLs in Go.
// It supports customizable parts such as scheme, host, port, path segments, user credentials,
// query parameters, fragments, and query deduplication strategies.
package urlbuilder

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type (
	// UrlPart represents a function that modifies a UrlBuilder.
	UrlPart func(*UrlBuilder) error

	// UrlPart represents a function that modifies a UrlBuilder.
	QueryPart func(*QueryString) error

	// QueryMode determines how query parameters are handled when duplicates are present.
	QueryMode int

	// QueryString holds the components of a query string collection and provides method to construct it.
	//
	// This is a helper struct to create independent query string values to be used as a value of the Query url part function in UrlBuilder.
	// Query strings are sometimes encrypted and encoded to avoid parameter tampering.
	QueryString struct {
		mode QueryMode
		qrs  []query
		err  error
	}

	query struct {
		name  string
		value string
	}

	// UrlBuilder holds components of a URL and provides methods to construct it.
	UrlBuilder struct {
		path     []string
		scheme   string
		host     string
		user     string
		password string
		id       string
		fragment string
		port     uint
		query    []query
		qmode    QueryMode
		err      error
	}
)

const (
	// QModeLast keeps only the last value of duplicate query parameter names. This is the default mode
	QModeLast QueryMode = iota

	// QModeArray keeps all query parameters, allowing duplicates (e.g., ?x=1&x=2).
	QModeArray

	// QModeError triggers an error if duplicate query parameter names are detected.
	QModeError
)

// New creates a new UrlBuilder with the provided UrlPart modifiers.
func New(part ...UrlPart) *UrlBuilder {
	ub := UrlBuilder{
		query: make([]query, 0, 3), // initializing to a capacity minimizes reallocations
		path:  make([]string, 0, 7),
	}
	for _, p := range part {
		p(&ub)
	}
	return &ub
}

// NewUrl returns a UrlBuilder with just a host.
func NewUrl(host string, part ...UrlPart) *UrlBuilder {
	up := make([]UrlPart, 0, 7)
	up = append(up, Host(host))
	up = append(up, part...)
	return New(up...)
}

// NewUrlWithPath returns a UrlBuilder with just a host and a path.
func NewUrlWithPath(host, path string, part ...UrlPart) *UrlBuilder {
	up := make([]UrlPart, 0, 7)
	up = append(up, Host(host))
	up = append(up, Path(path))
	up = append(up, part...)
	return New(up...)
}

// NewUrlWithID returns a UrlBuilder with a host, path, and ID segment.
func NewUrlWithID(host, path string, id any, part ...UrlPart) *UrlBuilder {
	up := make([]UrlPart, 0, 7)
	up = append(up, Host(host))
	up = append(up, Path(path))
	up = append(up, ID(id))
	up = append(up, part...)
	return New(up...)
}

func (ub *UrlBuilder) getHostParts(host string) {
	var (
		scheme, path, id string
		port             int
	)

	host = strings.ReplaceAll(host, "\"", "/")

	// If the host was supplied with a valid url and it has parts, take its result
	if r, err := url.Parse(host); err == nil {
		if r.Host != "" {
			host = r.Host
			if idx := strings.Index(host, ":"); idx != -1 {
				host = host[:idx] // Modify host
			}
		}
		// If it has scheme, this is not a pure host, so flag false
		if r.Scheme == "http" || r.Scheme == "https" {
			scheme = r.Scheme
		}

		// If it has port other than what is standard, flag false
		port, _ = strconv.Atoi(r.Port())
		if port != 0 {
			if scheme == "http" && port == 80 || scheme == "https" && port == 443 {
				port = 0
			}
		}
		// If it has a path, it is not a pure host, flag false
		if r.Path != "" && r.Host != "" {
			path = r.Path

			// If path is just a /, remove it
			if path == "/" {
				path = ""
			}

			// Remove the last path that has no /. It's considered a key in REST
			if !strings.HasSuffix(path, "/") {
				if idx := strings.LastIndex(path, "/"); idx != -1 {
					id = path[idx+1:]
					path = path[:idx]
				}
			}
		}
	}

	// Additional stripping of port
	if idx := strings.Index(host, ":"); idx != -1 {
		pvhost := host
		host = pvhost[:idx]
		port, _ = strconv.Atoi(pvhost[idx+1:])
	}
	ub.host, _ = strings.CutSuffix(host, "/")
	if port != 0 {
		ub.port = uint(port)
	}
	if scheme != "" {
		ub.scheme = scheme
	}
	if path != "" {
		ub.path = append(ub.path, path)
	}
	ub.id = id
}

// Clone returns a new UrlBuilder copied from an existing one and applies additional UrlParts.
func Clone(ub *UrlBuilder, part ...UrlPart) *UrlBuilder {
	cloneUb := *ub
	for _, p := range part {
		p(&cloneUb)
	}
	return &cloneUb
}

// Sch sets the scheme (e.g., "http", "https") of the URL.
func Sch(sch string) UrlPart {
	return func(ub *UrlBuilder) error {
		if sch == "" {
			return nil
		}
		ub.scheme = sch
		return nil
	}
}

// Host sets the host (domain or IP) of the URL.
func Host(h string) UrlPart {
	return func(ub *UrlBuilder) error {
		if h == "" {
			return nil
		}
		ub.getHostParts(h)
		return nil
	}
}

// Usr sets the username for basic authentication.
func Usr(u string) UrlPart {
	return func(ub *UrlBuilder) error {
		if u == "" {
			return nil
		}
		ub.user = u
		return nil
	}
}

// Pwd sets the password for basic authentication.
func Pwd(p string) UrlPart {
	return func(ub *UrlBuilder) error {
		if p == "" {
			return nil
		}
		ub.password = p
		return nil
	}
}

// UsrPwd sets both the username and password for basic authentication.
func UsrPwd(usr, pwd string) UrlPart {
	return func(ub *UrlBuilder) error {
		if usr == "" || pwd == "" {
			return nil
		}
		ub.user = usr
		ub.password = pwd
		return nil
	}
}

// Path appends a path segment to the URL.
func Path(path string) UrlPart {
	return func(ub *UrlBuilder) error {
		if path == "" {
			return nil
		}
		ub.path = append(ub.path, path)
		return nil
	}
}

// ID appends a single ID segment to the end of the URL path.
func ID(id any) UrlPart {
	return func(ub *UrlBuilder) error {
		if id == "" {
			return nil
		}
		ub.id = fmt.Sprint(id)
		return nil
	}
}

// Port sets the port number of the URL.
func Port(port uint) UrlPart {
	return func(ub *UrlBuilder) error {
		if port == 0 {
			return nil
		}
		ub.port = port
		return nil
	}
}

// Mode sets the query parameter deduplication mode.
func Mode(mode QueryMode) UrlPart {
	return func(ub *UrlBuilder) error {
		ub.qmode = mode
		return nil
	}
}

// Query appends a query parameter to the URL.
func Query(name string, value any) UrlPart {
	return func(ub *UrlBuilder) error {
		if name == "" {
			return nil
		}
		v := fmt.Sprint(value)
		// Check for values that may have the same name and value
		// If the keys and values are the same as the one being added,
		// ignore
		if ub.qmode == QModeArray {
			for _, kv := range ub.query {
				if strings.EqualFold(kv.name, name) && strings.EqualFold(kv.value, v) {
					continue
				}
			}
		}
		ub.query = append(ub.query, query{name: name, value: v})
		return nil
	}
}

// Frag sets the URL fragment (part after #).
func Frag(f string) UrlPart {
	return func(ub *UrlBuilder) error {
		if f == "" {
			return nil
		}
		ub.fragment = f
		return nil
	}
}

// Build constructs the URL as a string. Returns an empty string if an error occurred.
// func (ub *UrlBuilder) Build2() string {
// 	if ub.scheme == "" {
// 		ub.scheme = "https"
// 	}
// 	ub.scheme = strings.ToLower(ub.scheme)

// 	// Host: Change slashes and sanitize
// 	if ub.host == "" {
// 		ub.host = "localhost"
// 	}
// 	ub.host = strings.ReplaceAll(ub.host, "\"", "/")
// 	ub.host, _ = strings.CutSuffix(ub.host, "/")

// 	if ub.port == 0 {
// 		switch ub.scheme {
// 		case "https":
// 			ub.port = 443
// 		case "http":
// 			ub.port = 80
// 		}
// 	}

// 	url := ub.scheme + "://"

// 	if ub.user != "" {
// 		auth := ub.user
// 		if ub.password != "" {
// 			auth += ":" + ub.password
// 		}
// 		url += auth + "@"
// 	}
// 	url += ub.host
// 	if ub.port != 80 {
// 		url += ":" + strconv.Itoa(int(ub.port))
// 	}

// 	if len(ub.path) > 0 {
// 		paths := make([]string, len(ub.path))
// 		copy(paths, ub.path)
// 		for i := range paths {
// 			paths[i] = strings.ReplaceAll(paths[i], "\"", "/")
// 			paths[i], _ = strings.CutPrefix(paths[i], "/")
// 			paths[i], _ = strings.CutSuffix(paths[i], "/")
// 		}
// 		url += "/" + strings.Join(paths, "/")
// 	}

// 	pathAppended := false
// 	if ub.id != "" {
// 		url += "/" + ub.id
// 		pathAppended = true
// 	}

// 	if len(ub.query) > 0 {
// 		if !strings.HasSuffix(url, "/") && ub.id == "" {
// 			url += "/"
// 		}
// 		queryParams := []string{}
// 		if ub.qmode == QModeLast || ub.qmode == QModeError {
// 			qmap := make(map[string]string)
// 			for _, q := range ub.query {
// 				if _, found := qmap[q.name]; found && ub.qmode == QModeError {
// 					ub.err = fmt.Errorf("duplicate query name found")
// 					return ""
// 				}
// 				qmap[q.name] = q.value
// 			}
// 			for k, v := range qmap {
// 				queryParams = append(queryParams, k+"="+escape(v))
// 			}
// 		}
// 		if ub.qmode == QModeArray {
// 			for _, q := range ub.query {
// 				queryParams = append(queryParams, q.name+"="+escape(q.value))
// 			}
// 		}
// 		url += "?" + strings.Join(queryParams, "&")
// 		pathAppended = true
// 	}

// 	if ub.fragment != "" {
// 		url += "#" + ub.fragment
// 		pathAppended = true
// 	}

// 	if !pathAppended {
// 		url += "/"
// 	}

// 	return url
// }

// Build constructs the URL as a string. Returns an empty string if an error occurred.
func (ub *UrlBuilder) Build() string {
	if ub.scheme == "" {
		ub.scheme = "https"
	}
	ub.scheme = strings.ToLower(ub.scheme)

	if ub.host == "" {
		ub.host = "localhost"
	}

	switch ub.scheme {
	case "https":
		if ub.port == 0 {
			ub.port = 443
		}
	case "http":
		if ub.port == 0 {
			ub.port = 80
		}
	}

	var b strings.Builder

	b.WriteString(ub.scheme)
	b.WriteString("://")

	if ub.user != "" {
		b.WriteString(ub.user)
		if ub.password != "" {
			b.WriteByte(':')
			b.WriteString(ub.password)
		}
		b.WriteByte('@')
	}

	b.WriteString(ub.host)
	if !((ub.scheme == "http" && ub.port == 80) || (ub.scheme == "https" && ub.port == 443)) {
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(int(ub.port)))
	}

	if len(ub.path) > 0 {
		for _, segment := range ub.path {
			if segment == "" {
				continue
			}
			b.WriteByte('/')
			segment = strings.ReplaceAll(segment, "\"", "/")
			segment, _ = strings.CutPrefix(segment, "/")
			segment, _ = strings.CutSuffix(segment, "/")
			if segment != "" {
				b.WriteString(segment)
			}
		}
	}

	pathAppended := false
	hasTrailingPart := false
	if ub.id != "" {
		b.WriteByte('/')
		b.WriteString(ub.id)
		hasTrailingPart = true
		pathAppended = true
	}

	if len(ub.query) > 0 {
		if b.Len() > 0 && ub.id == "" && b.String()[b.Len()-1] != '/' {
			hasTrailingPart = true
			b.WriteByte('/')
		}
		b.WriteByte('?')

		first := true
		if ub.qmode == QModeLast || ub.qmode == QModeError {
			qmap := make(map[string]string)
			for _, q := range ub.query {
				if _, found := qmap[q.name]; found && ub.qmode == QModeError {
					ub.err = fmt.Errorf("duplicate query name found")
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
				b.WriteString(url.QueryEscape(v))
			}
		} else {
			for i, q := range ub.query {
				if i > 0 {
					b.WriteByte('&')
				}
				b.WriteString(q.name)
				b.WriteByte('=')
				b.WriteString(url.QueryEscape(q.value))
			}
		}
		pathAppended = true
	}

	if ub.fragment != "" {
		b.WriteByte('#')
		b.WriteString(ub.fragment)
		hasTrailingPart = true
		pathAppended = true
	}

	if !pathAppended && hasTrailingPart {
		b.WriteByte('/')
	}

	return b.String()
}

// Clone creates a new UrlBuilder from the current one and applies optional UrlParts.
func (ub *UrlBuilder) Clone(part ...UrlPart) *UrlBuilder {
	return Clone(ub, part...)
}

// Err returns any error that occurred during building.
func (ub *UrlBuilder) Err() error {
	return ub.err
}

// String implements fmt.Stringer and returns the built URL.
func (ub *UrlBuilder) String() string {
	return ub.Build()
}
