package httpx

import (
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"

	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/core/validation"
)

// ParseFiles parses files from multipart form and returns a map of field names
// to their first file. For multiple files under the same field name, use ParseMultipleFiles.
// It returns an empty map if no files are present or the form is not multipart.
func ParseFiles(r *http.Request) (map[string]*multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		// If not multipart, return empty map without error
		if err == http.ErrNotMultipart {
			return map[string]*multipart.FileHeader{}, nil
		}
		return nil, err
	}

	files := make(map[string]*multipart.FileHeader)
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return files, nil
	}

	for name, fileHeaders := range r.MultipartForm.File {
		if len(fileHeaders) > 0 {
			files[name] = fileHeaders[0]
		}
	}

	return files, nil
}

// ParseMultipleFiles parses files from multipart form and returns a map of field names
// to slices of file headers. This is useful when multiple files are uploaded under
// the same field name.
// It returns an empty map if no files are present or the form is not multipart.
func ParseMultipleFiles(r *http.Request) (map[string][]*multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		// If not multipart, return empty map without error
		if err == http.ErrNotMultipart {
			return map[string][]*multipart.FileHeader{}, nil
		}
		return nil, err
	}

	files := make(map[string][]*multipart.FileHeader)
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return files, nil
	}

	for name, fileHeaders := range r.MultipartForm.File {
		files[name] = fileHeaders
	}

	return files, nil
}

// GetFile returns the first file for the given form field name.
// Returns nil if the field doesn't exist or has no files.
func GetFile(r *http.Request, name string) (*multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		if err == http.ErrNotMultipart {
			return nil, nil
		}
		return nil, err
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return nil, nil
	}

	fileHeaders := r.MultipartForm.File[name]
	if len(fileHeaders) == 0 {
		return nil, nil
	}

	return fileHeaders[0], nil
}

// GetFiles returns all files for the given form field name.
// Returns nil if the field doesn't exist or has no files.
func GetFiles(r *http.Request, name string) ([]*multipart.FileHeader, error) {
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		if err == http.ErrNotMultipart {
			return nil, nil
		}
		return nil, err
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return nil, nil
	}

	return r.MultipartForm.File[name], nil
}

// ParseWithFiles parses the request including file fields.
// It extends Parse to automatically populate *multipart.FileHeader and
// []*multipart.FileHeader fields from multipart form data.
// The form tag is used to match the field name in the multipart form.
// For structs with file fields, this function handles both regular form
// values and file fields, skipping the standard form unmarshaler for
// file-type fields that it cannot handle.
func ParseWithFiles(r *http.Request, v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return Parse(r, v)
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return Parse(r, v)
	}

	// If the struct has no file fields, fall back to standard Parse
	if !hasFileFields(elem.Type()) {
		return Parse(r, v)
	}

	// Parse multipart form first
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		if err == http.ErrNotMultipart {
			return Parse(r, v)
		}
		return err
	}

	// Parse path and headers
	kind := mapping.Deref(reflect.TypeOf(v)).Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		if err := ParsePath(r, v); err != nil {
			return err
		}

		if err := ParseHeaders(r, v); err != nil {
			return err
		}
	}

	if err := ParseJsonBody(r, v); err != nil {
		return err
	}

	// Parse regular form fields (non-file fields)
	params, err := GetFormValues(r)
	if err != nil {
		return err
	}

	// Get file field names to exclude from form unmarshaler
	fileFieldNames := getFileFieldNames(elem.Type())

	// Remove file field keys from params to prevent type mismatch errors
	for _, name := range fileFieldNames {
		delete(params, name)
	}

	// Create a temporary struct without file fields for form unmarshaler
	// by building a map with only non-file field values
	if len(params) > 0 {
		if err := formUnmarshaler.Unmarshal(params, v); err != nil {
			return err
		}
	}

	// Populate file fields from multipart form
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		typ := elem.Type()
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			fieldValue := elem.Field(i)

			if !fieldValue.CanSet() {
				continue
			}

			tag := field.Tag.Get("form")
			if tag == "" || tag == "-" {
				continue
			}

			tagParts := strings.Split(tag, ",")
			fieldName := tagParts[0]

			fileHeaders := r.MultipartForm.File[fieldName]
			if len(fileHeaders) == 0 {
				continue
			}

			switch fieldValue.Type() {
			case reflect.TypeOf(&multipart.FileHeader{}):
				fieldValue.Set(reflect.ValueOf(fileHeaders[0]))
			case reflect.TypeOf([]*multipart.FileHeader{}):
				fieldValue.Set(reflect.ValueOf(fileHeaders))
			}
		}
	}

	// Run validation
	if valid, ok := v.(validation.Validator); ok {
		return valid.Validate()
	} else if v := getValidator(); v != nil {
		return v.Validate(r, v)
	}

	return nil
}

// SetFileField sets a file field value in a struct using reflection.
// The field must be of type *multipart.FileHeader.
// This is a helper function for custom request parsing.
func SetFileField(v reflect.Value, fieldName string, file *multipart.FileHeader) bool {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return false
	}

	if !field.CanSet() {
		return false
	}

	fieldType := field.Type()
	if fieldType == reflect.TypeOf(&multipart.FileHeader{}) {
		field.Set(reflect.ValueOf(file))
		return true
	}

	return false
}

// SetFilesField sets a multiple files field value in a struct using reflection.
// The field must be of type []*multipart.FileHeader.
// This is a helper function for custom request parsing.
func SetFilesField(v reflect.Value, fieldName string, files []*multipart.FileHeader) bool {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return false
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return false
	}

	if !field.CanSet() {
		return false
	}

	fieldType := field.Type()
	if fieldType == reflect.TypeOf([]*multipart.FileHeader{}) {
		field.Set(reflect.ValueOf(files))
		return true
	}

	return false
}

// hasFileFields checks if a struct type has any file fields
func hasFileFields(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return false
	}

	fileHeaderType := reflect.TypeOf(&multipart.FileHeader{})
	fileHeadersType := reflect.TypeOf([]*multipart.FileHeader{})

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Type == fileHeaderType || field.Type == fileHeadersType {
			return true
		}
	}

	return false
}

// getFileFieldNames returns the form tag names of all file fields in a struct type.
// This is used to exclude file fields from the standard form unmarshaler.
func getFileFieldNames(typ reflect.Type) []string {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil
	}

	fileHeaderType := reflect.TypeOf(&multipart.FileHeader{})
	fileHeadersType := reflect.TypeOf([]*multipart.FileHeader{})

	var names []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Type == fileHeaderType || field.Type == fileHeadersType {
			tag := field.Tag.Get("form")
			if tag != "" && tag != "-" {
				tagParts := strings.Split(tag, ",")
				names = append(names, tagParts[0])
			}
		}
	}

	return names
}

// shouldParseFiles determines if we should attempt to parse files based on the Content-Type header
func shouldParseFiles(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	return strings.Contains(contentType, "multipart/form-data")
}
