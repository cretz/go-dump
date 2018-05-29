package dump_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cretz/go-dump/dump"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	_, myFile, _, _ := runtime.Caller(0)
	testPkgDir := filepath.Join(filepath.Dir(myFile), "testpkg")
	pkgs, err := dump.LoadDir(testPkgDir)
	require.NoError(t, err)
	require.NotEmpty(t, pkgs)
	// TODO: tests
	// // Dump the JSON
	// jsonable, err := dump.PackagesToJSONMap(pkgs)
	// require.NoError(t, err)
	// byts, err := json.MarshalIndent(jsonable, "", "  ")
	// require.NoError(t, err)
	// println(string(byts))
}
