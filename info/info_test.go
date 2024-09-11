package info_test

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tkrop/go-config/info"
	"github.com/tkrop/go-testing/test"
)

const (
	// revisionHead contains an arbitrary head revision.
	revisionHead = "1b66f320c950b25fa63b81fd4e660c5d1f9d758e"
	// buildPath contains an arbitrary command build path.
	buildPath = "github.com/tkrop/go-config"
	// setupPath contains an arbitrary command setup path.
	setupPath = "github.com/tkrop/go-config/info"
)

type InfoParams struct {
	info       *info.Info
	build      *debug.BuildInfo
	expectInfo *info.Info
}

var testInfoParams = map[string]InfoParams{
	"nil build info": {
		info:       info.New("", "", "", "", "", ""),
		expectInfo: info.New("", "", "", "", "", ""),
	},
	"no build info": {
		info:       info.New("", "", "", "", "", ""),
		build:      &debug.BuildInfo{},
		expectInfo: info.New("", "", "", "", "", ""),
	},

	// Setup build info path.
	"build info setup path": {
		info: info.New(setupPath, "", "", "", "", ""),
		build: &debug.BuildInfo{
			Main: debug.Module{Path: buildPath},
		},
		expectInfo: info.New(setupPath, "", "", "", "", ""),
	},
	"build info build path": {
		info: info.New("", "", "", "", "", ""),
		build: &debug.BuildInfo{
			Main: debug.Module{Path: buildPath},
		},
		expectInfo: info.New(buildPath, "", "", "", "", ""),
	},

	// Setup build info version.
	"build info setup version": {
		info: info.New("", "v2.3.4", "beta.1", "", "", ""),
		build: &debug.BuildInfo{
			Main: debug.Module{Version: "v1.2.3-alpha.1"},
		},
		expectInfo: info.New("", "v1.2.3-alpha.1", "alpha.1", "", "", ""),
	},
	"build info build version": {
		info: info.New("", "", "", "", "", ""),
		build: &debug.BuildInfo{
			Main: debug.Module{Version: "v1.2.3-alpha.1"},
		},
		expectInfo: info.New("", "v1.2.3-alpha.1", "alpha.1", "", "", ""),
	},

	// Setup build info settings.
	"build info revision": {
		info: info.New("", "", "", "", "", ""),
		build: &debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "beta.2"},
			},
		},
		expectInfo: info.New("", "", "beta.2", "", "", ""),
	},
	"build info time": {
		info: info.New("", "", "", "", "", ""),
		build: &debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{Key: "vcs.time", Value: "2023-12-10T18:30:00Z"},
			},
		},
		expectInfo: info.New("", "", "", "", "2023-12-10T18:30:00Z", ""),
	},
	"build info time revision": {
		info: info.New("", "", "", "", "", ""),
		build: &debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "beta.2"},
				{Key: "vcs.time", Value: "2023-12-10T18:30:00Z"},
			},
		},
		expectInfo: info.New("", "", "beta.2", "",
			"2023-12-10T18:30:00Z", ""),
	},
	"build info time hash": {
		info: info.New("", "", "", "", "", ""),
		build: &debug.BuildInfo{
			Settings: []debug.BuildSetting{{
				Key:   "vcs.revision",
				Value: revisionHead,
			}, {Key: "vcs.time", Value: "2023-12-10T18:30:00Z"}},
		},
		expectInfo: info.New("", "", revisionHead, "",
			"2023-12-10T18:30:00Z", ""),
	},
	"build info time hash setup": {
		info: info.New("", "v1.2.3", "", "", "", ""),
		build: &debug.BuildInfo{
			Settings: []debug.BuildSetting{{
				Key:   "vcs.revision",
				Value: revisionHead,
			}, {Key: "vcs.time", Value: "2023-12-10T18:30:00Z"}},
		},
		expectInfo: info.New("", "v1.2.3", revisionHead, "",
			"2023-12-10T18:30:00Z", ""),
	},
	"build info modified": {
		info: info.New("", "", "", "", "", ""),
		build: &debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{Key: "vcs.modified", Value: "true"},
			},
		},
		expectInfo: info.New("", "", "", "", "", "true"),
	},
	"build info unmodified": {
		info: info.New("", "", "", "", "", ""),
		build: &debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{Key: "vcs.modified", Value: "false"},
			},
		},
		expectInfo: info.New("", "", "", "", "", "false"),
	},
}

func TestUseDebug(t *testing.T) {
	test.Map(t, testInfoParams).
		Run(func(t test.Test, param InfoParams) {
			// When
			info := param.info.UseDebug(param.build, true).
				AdjustVersion()

			// Then
			assert.Equal(t, param.expectInfo, info)
			assert.Equal(t, param.expectInfo.String(), info.String())
		})
}

func TestDefault(t *testing.T) {
	// Given
	defaultInfo := info.New("", "", "", "", "", "false")

	// When
	info.SetDefault(defaultInfo)

	// Then
	assert.Equal(t, defaultInfo, info.GetDefault())
}
