package star

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// SiteArchive describes entrypoints and metadata for a Meteor application.  It
// is generally located in the `star.json` file of the `.meteor/local/build` of
// a Meteor application directory.
type SiteArchive struct {

	// Format describes the format or schema of the data.  For now, the only valid value is "site-archive-pre1".
	Format string

	// BuiltBy indicates the Meteor tool which built the site archive
	BuildBy string

	// Programs describes the program components of the Meteor application
	Programs []ProgramDescription

	// MeteorRelease indicates the version of Meteor for which this application was designed
	MeteorRelease string

	// NodeVersion indicates the version of Node for which this application was designed
	NodeVersion string

	// NPMVersion indicates the version of NPM for which this application was designed
	NPMVersion string
}

// MeteorVersion parses the MeteorRelease to pull the raw version string
func (star *SiteArchive) MeteorVersion() string {
	if star == nil {
		return ""
	}

	pieces := strings.Split(star.MeteorRelease, "@")
	if len(pieces) != 2 {
		return ""
	}
	return pieces[1]
}

// A ProgramDescription describes a program of a Meteor application
type ProgramDescription struct {

	// Name is the name of the program
	Name string

	// Arch is the architecture of the program
	Arch string

	// Path is the path of the entrypoint, relative to the application bundle root
	Path string
}

// Find attempts to locate and read the SiteArchive (star.json) file of a Meteor application
func Find(baseDir string) (*SiteArchive, error) {
	var star *SiteArchive

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "star.json" {
			star, err = ReadFile(path)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if star == nil {
		return nil, errors.New("not found")
	}
	return star, nil
}

// ReadFile reads a SiteArchive from the given file path
func ReadFile(path string) (*SiteArchive, error) {
	if path == "" {
		return nil, errors.New("not found")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	star := new(SiteArchive)
	err = json.NewDecoder(f).Decode(star)
	return star, err
}
