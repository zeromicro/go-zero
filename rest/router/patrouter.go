package router

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/core/search"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

const (
	allowHeader          = "Allow"
	allowMethodSeparator = ", "
)

var (
	// ErrInvalidMethod is an error that indicates not a valid http method.
	ErrInvalidMethod = errors.New("not a valid http method")
	// ErrInvalidPath is an error that indicates path is not start with /.
	ErrInvalidPath = errors.New("path must begin with '/'")
)

type patRouter struct {
	trees          map[string]*search.Tree
	notFound       http.Handler
	notAllowed     http.Handler
	fileSystemTree *search.Tree
}

// NewRouter returns a httpx.Router.
func NewRouter() httpx.Router {
	return &patRouter{
		trees: make(map[string]*search.Tree),
	}
}

func (pr *patRouter) Handle(method, reqPath string, handler http.Handler) error {
	if !validMethod(method) {
		return ErrInvalidMethod
	}

	if len(reqPath) == 0 || reqPath[0] != '/' {
		return ErrInvalidPath
	}

	cleanPath := path.Clean(reqPath)
	tree, ok := pr.trees[method]
	if ok {
		return tree.Add(cleanPath, handler)
	}

	tree = search.NewTree()
	pr.trees[method] = tree
	return tree.Add(cleanPath, handler)
}

func (pr *patRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqPath := path.Clean(r.URL.Path)
	if tree, ok := pr.trees[r.Method]; ok {
		if result, ok := tree.Search(reqPath); ok {
			if len(result.Params) > 0 {
				r = pathvar.WithVars(r, result.Params)
			}
			result.Item.(http.Handler).ServeHTTP(w, r)
			return
		}
	}

	if r.Method == http.MethodGet && pr.fileSystemTree != nil {
		var fileSystemHandle interface{}

		fileSystemSearchPath := reqPath
		lastSlashIndex := strings.LastIndexByte(fileSystemSearchPath, '/')
		for lastSlashIndex >= 0 {
			if result, ok := pr.fileSystemTree.Search(fileSystemSearchPath); ok {
				fileSystemHandle = result.Item
				break
			}
			fileSystemSearchPath = fileSystemSearchPath[:lastSlashIndex]
			lastSlashIndex = strings.LastIndexByte(fileSystemSearchPath, '/')
		}

		if fileSystemHandle == nil && fileSystemSearchPath == "" {
			if result, ok := pr.fileSystemTree.Search("/"); ok {
				fileSystemHandle = result.Item
			}
		}

		if fileSystemHandle != nil {
			fileSystemHandle.(http.Handler).ServeHTTP(w, r)
			return
		}
	}

	allows, ok := pr.methodsAllowed(r.Method, reqPath)
	if !ok {
		pr.handleNotFound(w, r)
		return
	}

	if pr.notAllowed != nil {
		pr.notAllowed.ServeHTTP(w, r)
	} else {
		w.Header().Set(allowHeader, allows)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (pr *patRouter) SetNotFoundHandler(handler http.Handler) {
	pr.notFound = handler
}

func (pr *patRouter) SetNotAllowedHandler(handler http.Handler) {
	pr.notAllowed = handler
}

func (pr *patRouter) SetFileSystemHandlerMap(handlerMap map[string]http.Handler) {
	if len(handlerMap) > 0 {
		tree := search.NewTree()
		for k, v := range handlerMap {
			if len(k) == 0 {
				panic(fmt.Errorf("FileSystemHandlerMap key must not be empty"))
			}
			if k[0] != '/' {
				panic(fmt.Errorf("FileSystemHandlerMap key must begin with '/' for %s", k))
			}

			cleanPath := path.Clean(k)
			if cleanPath != k {
				panic(fmt.Errorf("FileSystemHandlerMap key '%s' should be '%s'", k, cleanPath))
			}

			if err := tree.Add(cleanPath, v); err != nil {
				panic(err)
			}
		}

		pr.fileSystemTree = tree
	}
}

func (pr *patRouter) handleNotFound(w http.ResponseWriter, r *http.Request) {
	if pr.notFound != nil {
		pr.notFound.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (pr *patRouter) methodsAllowed(method, path string) (string, bool) {
	var allows []string

	for treeMethod, tree := range pr.trees {
		if treeMethod == method {
			continue
		}

		_, ok := tree.Search(path)
		if ok {
			allows = append(allows, treeMethod)
		}
	}

	if len(allows) > 0 {
		return strings.Join(allows, allowMethodSeparator), true
	}

	return "", false
}

func validMethod(method string) bool {
	return method == http.MethodDelete || method == http.MethodGet ||
		method == http.MethodHead || method == http.MethodOptions ||
		method == http.MethodPatch || method == http.MethodPost ||
		method == http.MethodPut
}
