# AutoClosing TempDir

Because os.TempDir leaves it to the implementer to clean up the directory,
this was written to avoid writing boiler plate code to do clean ups.

Library to create temporary directories, treat them like objects and call Close using defers to ensure that the directory is cleaned up.

# How to use

```go
// Basic imports
import (
    "github.com/wojnosystems/tempdir"
)

func main() {
	dir1, err := tempdir.TempDir("/tmp/thing", "prefix")
	if err != nil {
		panic(err)
	}
	defer dir.MustClose()
	
	dir2, err := tempdir.TempDir("/tmp/thing", "prefix")
	if err != nil {
		panic(err)
	}
	err = dir.Close()
	if err != nil {
		panic("Unable to delete the files")
	}
}
```

You can call MustClose, which will panic if os.RemoveAll throws an error, or call Close which will return the error if you want it. You can wrap both in defer's to allow them to auto-close when you're done inside the method, or just call Close whenever you want to clean up the directory.