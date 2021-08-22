package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type path struct {
	pattern []string
	request []string
	vars    map[string]string
	query   url.Values
	frag    string
}

// NewPath parses the request url against the given pattern.
// Returns a pointer to the parsed path object.
//
// Pattern: The pattern always begins with the root slash.
// Identify variable subpaths within the pattern by
// enclosing them in curly braces, "{varName}". All other
// subpaths are ignored. Each subpath may contain only one
// variable and the variable must span the entire subpath.
//
// Example: "https://example.com/blog/random/first-post/"
//
//		// Pattern always begins from the root slash
//		pattern := "/anything/{category}/{slug}/{id}"
//		p := NewPath(r, pattern)
//
//		p.Var("category") ==> "random"
//		p.var("slug") ==> "first-post"
//		p.var("id") ==> ""
func Newpath(r *http.Request, pattern string) *path {
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}

	p := path{}
	p.parsePattern(pattern)
	p.parseReq(r.URL.Path)

	p.vars = map[string]string{}
	p.mapVars()

	p.query = r.URL.Query()
	p.frag = r.URL.Fragment

	return &p
}

// Var returns the value corresponding to the path variable
//
// Example:
//		"https://example.com/edit/article/42"
// 		Pattern: "/anything/{category}/{id}"
//
//		p.Var("category") ==> "article"
//		p.Var("id")	==> "42"
func (p *path) Var(key string) (val string) {
	return p.vars[key]
}

// VarInt returns an int representation of val.
// ok returns true on success. Otherwise, false.
//
// val returns 0 if the corresponding value cannot
// be converted to an int or does not exist.
func (p *path) VarInt(key string) (val int, ok bool) {
	val, err := strconv.Atoi(p.vars[key])
	return val, (err == nil)
}

// Query returns the query value for the corresponding key.
//
// Example:
// 		"https://example.com/edit/article/?item=42"
//
//		p.Query("item")	==> "42"
func (p *path) Query(key string) string {
	return p.query.Get(key)
}

// Frag returns the path fragment value
//
// Example:
//		"https://example.com/blog/article/#42"
//
//		p.Frag() ==> "42"
func (p *path) Frag() string {
	return p.frag
}

// parseReq loads the request field with the string slice
// representation of the request url path.
func (p *path) parseReq(path string) {
	parse(&p.request, path)
}

// parsePattern loads the pattern field with the string slice
// representation of the given pattern
func (p *path) parsePattern(pattern string) {
	parse(&p.pattern, pattern)
}

// mapVars identifies variables in the defined pattern and maps
// the variable to the corresponding value in the http request.
func (p *path) mapVars() {
	for i, v := range p.pattern {
		if i >= len(p.request) {
			break
		}

		if isVar(v, true) { // v == "{var}"
			k := strings.TrimSpace(v[1 : len(v)-1])
			p.vars[k] = p.request[i]
		}

	}
}

// parse loads the slice with slash separated substrings.
// parse is a helper function used to load the path slice fields
func parse(s *[]string, path string) {
	// Exclude first element, always empty
	*s = strings.Split(path, "/")[1:]
}

// isVar returns true if val is a string formatted as a variable
// (i.e., "{var}"). Otherwise, returns false.
//
// Set strict to true if val must be validated as a variable.
// Set strict to false if val only needs to be recognized as a variable.
func isVar(val string, strict bool) bool {
	if !strict {
		return strings.HasPrefix(val, "{")
	}

	return strings.HasPrefix(val, "{") &&
		strings.HasSuffix(val, "}") &&
		strings.Index(val[1:], "{") == -1
}
