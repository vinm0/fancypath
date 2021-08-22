package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type path struct {
	// Variable subpath names parsed from pattern string
	//
	pattern []string

	// Subpath values parsed from http request
	//
	request []string

	// Mapped subpaths using pattern subpaths as key
	// and http request subpaths as value.
	//
	// Values are retrieved by calling the Var(key string) method
	//
	vars map[string]string

	// Query values obtained from http request.
	//
	// Values are retrieved by calling the Query(key string) method.
	//
	query url.Values

	// Frag contains the path fragment value
	//
	// Example: "https://example.com/blog/article/#42"
	//
	//		p := NewPath(r, "")
	//		p.Frag == "42"  // true
	//
	Frag string
}

// NewPath parses the request url against the given pattern.
//
// Pattern: The pattern always begins with the root slash.
// Identify variable subpaths within the pattern by providing
// a term in the subpath. Each subpath may contain only one
// variable and the variable must span the entire subpath.
// To ignore a subpath, substitute a subpath with an astrisk, "*".
//
// Example: "https://example.com/blog/random/first-post/"
//
//		// Pattern always begins from the root slash
//		pattern := "/*/category/slug/id"
//		p := NewPath(r, pattern)
//
//		p.Var("category") == "random" // true
//		p.Var("slug") == "first-post" // true
//		p.Var("id") == "" // true
//		p.Var("*") == "" // true
//
func Newpath(r *http.Request, pattern string) *path {
	// Must have root slash
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}

	p := path{}
	p.parsePattern(pattern)
	p.parseReq(r.URL.Path)
	p.mapVars()

	p.query = r.URL.Query()
	p.Frag = r.URL.Fragment

	return &p
}

// Var returns the value corresponding to the path variable
//
// Example: "https://example.com/edit/article/42"
//
// 		pattern = "/*/category/id"
//		p := NewPath(r, pattern)
//
//		p.Var("category") == "article" // true
//		p.Var("id")	== "42" // true
//		p.Var("*") == "" // true
//
func (p *path) Var(key string) (val string) {
	return p.vars[key]
}

// VarInt returns an int representation of val.
// ok returns true on success. Otherwise, false.
//
// val returns 0 if the corresponding value cannot
// be converted to an int or does not exist.
//
func (p *path) VarInt(key string) (val int, ok bool) {
	val, err := strconv.Atoi(p.vars[key])
	return val, (err == nil)
}

// Query returns the query value for the corresponding key.
//
// Example:
// 		"https://example.com/edit/article/?item=42"
//
//		p.Query("item") == "42" // true
//
func (p *path) Query(key string) string {
	return p.query.Get(key)
}

// parseReq loads the request field with the string slice
// representation of the request url path.
//
func (p *path) parseReq(path string) {
	parse(&p.request, path)
}

// parsePattern loads the pattern field with the string slice
// representation of the given pattern
//
func (p *path) parsePattern(pattern string) {
	parse(&p.pattern, pattern)
}

// mapVars identifies variables in the defined pattern and maps
// the variable to the corresponding value in the http request.
//
func (p *path) mapVars() {
	p.vars = map[string]string{}

	for i, v := range p.pattern {
		if i >= len(p.request) {
			break
		}

		k := strings.TrimSpace(v)
		p.vars[k] = p.request[i]
	}

	delete(p.vars, "*")
}

// parse loads the slice with slash separated substrings.
// parse is a helper function used to load the path slice fields
//
func parse(s *[]string, path string) {
	// Exclude first element, always empty
	*s = strings.Split(path, "/")[1:]
}
