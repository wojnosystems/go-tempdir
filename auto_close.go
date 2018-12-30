/*
Copyright 2018 Chris Wojno.
Attribution 4.0 International (CC BY 4.0)
All rights reserved.
You do not have to comply with the license for elements of the material in the public domain or where your use is
permitted by an applicable exception or limitation.
No warranties are given. The license may not give you all of the permissions necessary for your intended use.
For example, other rights such as publicity, privacy, or moral rights may limit how you use the material.
See LICENSE file for the full license
*/

/*
tempdir

Creates a temporary directory using ioutil.TempDir, but also provides a callback (Close() and MustClose()) to clean itself up.
I found this pattern showing up repeatedly in multiple projects and got tired of writing the code and test for this.

@author Chris Wojno
@copyright
*/
package tempdir

import (
	"io/ioutil"
	"os"
)

// AutoClosing tempdir provides the mechanism for cleaning up after creating the temporary directory
//
// @implements Closer with Close()error
//
// Don't create this yourself, you should call TempDir to create the directory and set the value in this struct.
// This type is exported to allow you to pass it around your program in case you need to clean things up later
type AutoClosing struct {
	path string
}

// tempDirFunction allows the create TempDir function to be swapped out for testing
var tempDirFunction = ioutil.TempDir
// removeAllFunction allows the RemoveAll function to be swapped out for testing
var removeAllFunction = os.RemoveAll

// Path the path to the temporary directory as created by calling TempDir
func (a AutoClosing) Path() string {
	return a.path
}

// Close cleans up the temporary directory and all files contained therein
// @return error from RemoveAll or nil if path is empty
func (a *AutoClosing)Close() error {
	if len(a.path) != 0 {
		err := removeAllFunction(a.path)
		a.path = ""
		return err
	}
	return nil
}

// MustClose calls Close and panics if an error is encountered
func (a *AutoClosing)MustClose() {
	err := a.Close()
	if err != nil {
		panic(err)
	}
}

// TempDir creates a new temporary directory in the directory dir
// with a name beginning with prefix and returns the path of the
// new directory. If dir is the empty string, TempDir uses the
// default directory for temporary files (see os.TempDir).
// Multiple programs calling TempDir simultaneously
// will not choose the same directory.
//
// TempDir calls ioutil.TempDir, but saves the returned path to the AutoClosing struct
// To remove the directory, just call "Close()" or "MustClose()" on the struct to delete the directory and all files within.
func TempDir( dir, prefix string ) (ac AutoClosing, err error) {
	ac.path, err = tempDirFunction(dir, prefix)
	if err != nil {
		return ac, err
	}
	return ac, err
}