# Fancy Paths

(DRAFT)

Fancy Paths is a url path parser written for Go web apps.

The path simply retreives any part of a http.Request url path based on predefined key values. Keys are defined by a sample url pattern.

The example pattern below defines the keys "category" and "slug" used to retrieve the corresponding values in the http.Request url path:

```
// import paths "github.com/vinm0/fancypaths"
.
.
.
// GET: "https://example.com/blog/poetry/first-post"
// var req *http.Request

pattern := "/*/category/slug"
p := paths.New(req, pattern)

p.Var("category") == "poetry" // true
p.Var("slug") == "first-post" // true
p.Var("*") == ""              // true, astrisks are ignored
```

## Contents

1. [Installation](#installation)
1. [Usage](#usage)
    1. [Initialize](#1-initialize)
    1. [Path Values](#2-path-values)
    1. [Query Values](#3-query-values)
    1. [Fragment Value](#4-fragment-value)
1. [Why not github.com/gorilla/mux?](#why-not-mux)

## Installation

```
go get github.com/vinm0/fancypaths
```

```
import "github.com/vinm0/fancypaths"
```

If you're not feeling very fancy, but you'd still want the paths, I guess you can import an alias (ugh).

```
import paths "github.com/vinm0/fancypaths"
```

## Usage

### 1. Initialize

```
// import paths "github.com/vinm0/fancypaths"
.
.
.
// GET: "https://example.com/blog/poetry/first-post"
// var req *http.Request

pattern := "/*/category/slug"
p := paths.New(req, pattern)
```

### 2. Path Values

```
// GET: "https://example.com/blog/poetry/first-post/"
// var req *http.Request

pattern := "/*/category/slug/id"
p := paths.New(req, pattern)

p.Var("category") == "poetry" // true
p.Var("slug") == "first-post" // true
p.Var("id") == ""             // true, not provided in http.Request
p.Var("*") == ""              // true, astrisks are ignored
```

### 3. Query Values

```
// "https://example.com/edit/article/?item=42"

p := paths.New(req, "/")

p.Query("item") == "42" // true
```

### 4. Fragment Value

```
// "https://example.com/blog/article/#42"

p := New(req, "/")

p.Frag == "42"  // true
```

## Why Not Mux?

Gorilla Mux is great when you want a specific handler to catch a specific pattern.

For instance, consider the following path pattern.

```
gMux.HandleFunc("/edit/{category}/{id:[0-9]+}", specificHandler)
```

Mux will only pass requests to the handler that match the pattern (such as, `"/edit/article/42"`). But it will not match slightly different patterns (such as, `"/edit/article/some-slug"` or `"/edit/profile"`). In which case, you might need more highly specific routes.

Perhaps you don't want to list every possible route pattern. Perhaps, you would prefer one handler for the route `"/edit/"` to make decisions based on the rest of the requested path.

### Solution:

Define the variables within the handler instead of the router.

Keep your path variables while watching for a more general pattern.

With Fancy Paths, could match the pattern `"/edit/"` and let the handler decide what to do with the rest of the request path.
