package templateutil

import (
	"bytes"
	htemplate "html/template"
	"io"
	"io/ioutil"
	ttemplate "text/template"
)

func ParseTextTemplateFileToWriter(filename string, data interface{}, w io.Writer) error {
	b, err := ParseTextTemplateFileToReader(filename, data)
	if err != nil {
		return err
	}
	io.Copy(w, b)
	return nil
}

func ParseTextTemplateFileToString(filename string, data interface{}) (string, error) {
	b, err := ParseTextTemplateFileToReader(filename, data)
	if err != nil {
		return "", err
	}
	s, err := ioutil.ReadAll(b)
	return string(s), err
}

func ParseTextTemplateFileToReader(filename string, data interface{}) (io.Reader, error) {
	b := bytes.Buffer{}

	t, err := ttemplate.ParseFiles(filename)
	if err != nil {
		return &b, err
	}

	t.Execute(&b, data)
	return &b, nil
}

func ParseHtmlTemplateFileToString(filename string, data interface{}) (string, error) {
	b, err := ParseHtmlTemplateFileToReader(filename, data)
	if err != nil {
		return "", err
	}
	s, err := ioutil.ReadAll(b)
	return string(s), err
}

func ParseHtmlTemplateFileToReader(filename string, data interface{}) (io.Reader, error) {
	b := bytes.Buffer{}

	t, err := htemplate.ParseFiles(filename)
	if err != nil {
		return &b, err
	}

	t.Execute(&b, data)
	return &b, nil
}
