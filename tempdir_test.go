package tempdir

import (
	"errors"
	"os"
	"testing"
)

func TestTempDir_OK(t *testing.T) {
	dir, err := TempDir("","")
	if err != nil {
		t.Error("err should be nil")
	}
	if dir.Path() == "" {
		t.Error("TempDir created empty path")
	}
	if _, err := os.Stat(dir.Path()); os.IsNotExist(err) {
		t.Error("path should exist")
	}

	err = dir.Close()
	if err != nil {
		t.Error("Close should not fail")
	}

	if _, err := os.Stat(dir.Path()); !os.IsNotExist(err) {
		t.Error("path should no longer exist")
	}
}

func TestTempDir_OKCloseEmpty(t *testing.T) {
	dir := AutoClosing{}
	err := dir.Close()
	if err != nil {
		t.Error("Close should not fail")
	}
}

func TestTempDir_OKMustClose(t *testing.T) {
	dir, err := TempDir("","")
	if err != nil {
		t.Error("err should be nil")
	}
	didPanic := false
	var m interface{}
	func() {
		defer func() {
			if m = recover(); m != nil {
				didPanic = true
			}
		}()
		dir.MustClose()
	}()

	if didPanic {
		t.Error("MustClose paniced but should not have")
	}
}

func TestTempDir_ErrTempDir(t *testing.T) {
	old := tempDirFunction
	tempDirFunction = func(dir, prefix string) (name string, err error) {
		return "", errors.New("some error")
	}
	_, err := TempDir("","")
	if err.Error() != "some error" {
		t.Error("err expected, but got", err)
	}
	tempDirFunction = old
}

func TestTempDir_ErrClose(t *testing.T) {
	old := removeAllFunction
	removeAllFunction = func( path string ) (err error) {
		return errors.New("some error")
	}
	dir, err := TempDir("","")
	if err != nil {
		t.Error("err should be nil")
	}
	// Clean up for test
	err = os.RemoveAll(dir.Path())
	if err != nil {
		t.Error("err should be nil")
	}

	err = dir.Close()
	if err.Error() != "some error" {
		t.Error("err expected, but got", err)
	}

	removeAllFunction = old
}

func TestTempDir_ErrMustClose(t *testing.T) {
	old := removeAllFunction
	removeAllFunction = func( path string ) (err error) {
		return errors.New("some error")
	}
	dir, err := TempDir("","")
	if err != nil {
		t.Error("err should be nil")
	}
	// Clean up for test
	err = os.RemoveAll(dir.Path())
	if err != nil {
		t.Error("err should be nil")
	}

	didPanic := false
	var m interface{}
	func() {
		defer func() {
			if m = recover(); m != nil {
				didPanic = true
			}
		}()
		dir.MustClose()
	}()

	if !didPanic {
		t.Error("MustClose should panic")
	}

	removeAllFunction = old
}