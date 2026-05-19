package httpx

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFiles(t *testing.T) {
	t.Run("multipart with file", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.txt")
		assert.Nil(t, err)
		_, err = part.Write([]byte("test content"))
		assert.Nil(t, err)
		writer.Close()

		r := httptest.NewRequest(http.MethodPost, "/upload", body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		files, err := ParseFiles(r)
		assert.Nil(t, err)
		assert.NotNil(t, files["file"])
		assert.Equal(t, "test.txt", files["file"].Filename)
	})

	t.Run("not multipart", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewReader([]byte("hello")))
		r.Header.Set("Content-Type", "application/json")

		files, err := ParseFiles(r)
		assert.Nil(t, err)
		assert.Empty(t, files)
	})

	t.Run("multipart no files", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		err := writer.WriteField("name", "test")
		assert.Nil(t, err)
		writer.Close()

		r := httptest.NewRequest(http.MethodPost, "/upload", body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		files, err := ParseFiles(r)
		assert.Nil(t, err)
		assert.Empty(t, files)
	})
}

func TestParseMultipleFiles(t *testing.T) {
	t.Run("multiple files same field", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part1, err := writer.CreateFormFile("files", "file1.txt")
		assert.Nil(t, err)
		_, err = part1.Write([]byte("content1"))
		assert.Nil(t, err)

		part2, err := writer.CreateFormFile("files", "file2.txt")
		assert.Nil(t, err)
		_, err = part2.Write([]byte("content2"))
		assert.Nil(t, err)

		writer.Close()

		r := httptest.NewRequest(http.MethodPost, "/upload", body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		files, err := ParseMultipleFiles(r)
		assert.Nil(t, err)
		assert.Len(t, files["files"], 2)
		assert.Equal(t, "file1.txt", files["files"][0].Filename)
		assert.Equal(t, "file2.txt", files["files"][1].Filename)
	})

	t.Run("not multipart", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewReader([]byte("hello")))
		r.Header.Set("Content-Type", "application/json")

		files, err := ParseMultipleFiles(r)
		assert.Nil(t, err)
		assert.Empty(t, files)
	})
}

func TestGetFile(t *testing.T) {
	t.Run("existing file", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.txt")
		assert.Nil(t, err)
		_, err = part.Write([]byte("content"))
		assert.Nil(t, err)
		writer.Close()

		r := httptest.NewRequest(http.MethodPost, "/upload", body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		file, err := GetFile(r, "file")
		assert.Nil(t, err)
		assert.NotNil(t, file)
		assert.Equal(t, "test.txt", file.Filename)
	})

	t.Run("non-existing file", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.txt")
		assert.Nil(t, err)
		_, err = part.Write([]byte("content"))
		assert.Nil(t, err)
		writer.Close()

		r := httptest.NewRequest(http.MethodPost, "/upload", body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		file, err := GetFile(r, "nonexistent")
		assert.Nil(t, err)
		assert.Nil(t, file)
	})
}

func TestGetFiles(t *testing.T) {
	t.Run("existing files", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part1, err := writer.CreateFormFile("files", "a.txt")
		assert.Nil(t, err)
		_, err = part1.Write([]byte("a"))
		assert.Nil(t, err)
		part2, err := writer.CreateFormFile("files", "b.txt")
		assert.Nil(t, err)
		_, err = part2.Write([]byte("b"))
		assert.Nil(t, err)
		writer.Close()

		r := httptest.NewRequest(http.MethodPost, "/upload", body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		files, err := GetFiles(r, "files")
		assert.Nil(t, err)
		assert.Len(t, files, 2)
	})
}

func TestParseWithFiles(t *testing.T) {
	t.Run("struct with File field", func(t *testing.T) {
		var req struct {
			Id   string                `form:"id"`
			File *multipart.FileHeader `form:"file,optional"`
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		err := writer.WriteField("id", "123")
		assert.Nil(t, err)
		part, err := writer.CreateFormFile("file", "test.txt")
		assert.Nil(t, err)
		_, err = part.Write([]byte("content"))
		assert.Nil(t, err)
		writer.Close()

		r := httptest.NewRequest(http.MethodPost, "/upload", body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		err = ParseWithFiles(r, &req)
		assert.Nil(t, err)
		assert.Equal(t, "123", req.Id)
		assert.NotNil(t, req.File)
		assert.Equal(t, "test.txt", req.File.Filename)
	})

	t.Run("struct with []File field", func(t *testing.T) {
		var req struct {
			Id    string                  `form:"id"`
			Files []*multipart.FileHeader `form:"files,optional"`
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		err := writer.WriteField("id", "456")
		assert.Nil(t, err)
		part1, err := writer.CreateFormFile("files", "a.txt")
		assert.Nil(t, err)
		_, err = part1.Write([]byte("a"))
		assert.Nil(t, err)
		part2, err := writer.CreateFormFile("files", "b.txt")
		assert.Nil(t, err)
		_, err = part2.Write([]byte("b"))
		assert.Nil(t, err)
		writer.Close()

		r := httptest.NewRequest(http.MethodPost, "/upload", body)
		r.Header.Set("Content-Type", writer.FormDataContentType())

		err = ParseWithFiles(r, &req)
		assert.Nil(t, err)
		assert.Equal(t, "456", req.Id)
		assert.Len(t, req.Files, 2)
		assert.Equal(t, "a.txt", req.Files[0].Filename)
		assert.Equal(t, "b.txt", req.Files[1].Filename)
	})

	t.Run("not multipart", func(t *testing.T) {
		var req struct {
			Id   string                `form:"id"`
			File *multipart.FileHeader `form:"file,optional"`
		}

		r := httptest.NewRequest(http.MethodGet, "/upload?id=123", nil)

		err := ParseWithFiles(r, &req)
		assert.Nil(t, err)
		assert.Equal(t, "123", req.Id)
		assert.Nil(t, req.File)
	})
}

func TestHasFileFields(t *testing.T) {
	t.Run("struct with File field", func(t *testing.T) {
		var v struct {
			File *multipart.FileHeader `form:"file"`
		}
		assert.True(t, hasFileFields(reflect.TypeOf(v)))
	})

	t.Run("struct with []File field", func(t *testing.T) {
		var v struct {
			Files []*multipart.FileHeader `form:"files"`
		}
		assert.True(t, hasFileFields(reflect.TypeOf(v)))
	})

	t.Run("struct without File field", func(t *testing.T) {
		var v struct {
			Name string `form:"name"`
		}
		assert.False(t, hasFileFields(reflect.TypeOf(v)))
	})
}
