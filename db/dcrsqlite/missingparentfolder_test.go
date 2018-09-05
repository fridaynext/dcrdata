package dcrsqlite

import (
	"path/filepath"
	"testing"

	"github.com/fridaynext/dcrdata/testutil"
)

// TestMissingParentFolder ensures InitDB() is able to create
// a new DB-file parent directory if necessary
// See https://github.com/fridaynext/dcrdata/issues/515
func TestMissingParentFolder(t *testing.T) {
	testutil.BindCurrentTestSetup(t)
	testMissingParentFolder()
}

func testMissingParentFolder() {
	testName := testutil.TestName()
	testutil.ResetTempFolder(&testName)
	// Specify DB file in non-existent path
	target := filepath.Join(testName, "x", "y", "z", testutil.DefaultDBFileName)
	targetDBFile := testutil.FilePathInsideTempDir(target)

	InitTestDB(targetDBFile)
}
