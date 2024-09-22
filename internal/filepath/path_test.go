package filepath_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tkrop/go-config/internal/filepath"
	"github.com/tkrop/go-testing/test"
)

var currentDir, _ = os.Getwd()

type testNormalizeParam struct {
	path       string
	setup      func(test.Test)
	clean      func(test.Test)
	expectPath string
}

var testNormalizeParams = map[string]testNormalizeParam{
	"path empty": {
		path:       "",
		expectPath: currentDir,
	},
	"path dot": {
		path:       "",
		expectPath: currentDir,
	},

	"path absolute": {
		path:       "/tmp",
		expectPath: "/tmp",
	},

	"path relative": {
		path:       "tmp",
		expectPath: currentDir + "/tmp",
	},

	"path expand": {
		path: "${HOME}/tmp",
		setup: func(t test.Test) {
			t.Setenv("HOME", "/home/user")
		},
		expectPath: "/home/user/tmp",
	},

	"path expand error": {
		path: "some/${INVALID}/path",
		setup: func(t test.Test) {
			assert.NoError(t, os.Mkdir("_tmp", 0o755), "mkdir")
			assert.NoError(t, os.Chdir("_tmp"), "chdir")
			assert.NoError(t, os.Remove("../_tmp"), "remove")
		},
		expectPath: "some/path",
		clean: func(t test.Test) {
			assert.NoError(t, os.Chdir(".."), "chdir")
		},
	},
}

func TestNormalize(t *testing.T) {
	test.Map(t, testNormalizeParams).
		RunSeq(func(t test.Test, param testNormalizeParam) {
			// Given
			if param.setup != nil {
				param.setup(t)
			}

			// When
			path := filepath.Normalize(param.path)

			// Then
			assert.Equal(t, param.expectPath, path)
			if param.clean != nil {
				param.clean(t)
			}
		})
}
