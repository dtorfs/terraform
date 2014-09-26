package module

import (
	"bytes"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

// Getter defines the interface that schemes must implement to download
// and update modules.
type Getter interface {
	// Get downloads the given URL into the given directory. This always
	// assumes that we're updating and gets the latest version that it can.
	//
	// The directory may already exist (if we're updating). If it is in a
	// format that isn't understood, an error should be returned. Get shouldn't
	// simply nuke the directory.
	Get(string, *url.URL) error
}

// Getters is the mapping of scheme to the Getter implementation that will
// be used to get a dependency.
var Getters map[string]Getter

// forcedRegexp is the regular expression that finds forced getters. This
// syntax is schema::url, example: git::https://foo.com
var forcedRegexp = regexp.MustCompile(`^([A-Za-z]+)::(.+)$`)

func init() {
	httpGetter := new(HttpGetter)

	Getters = map[string]Getter{
		"file":  new(FileGetter),
		"git":   new(GitGetter),
		"hg":    new(HgGetter),
		"http":  httpGetter,
		"https": httpGetter,
	}
}

// Get downloads the module specified by src into the folder specified by
// dst. If dst already exists, Get will attempt to update it.
//
// src is a URL, whereas dst is always just a file path to a folder. This
// folder doesn't need to exist. It will be created if it doesn't exist.
func Get(dst, src string) error {
	var force string
	force, src = getForcedGetter(src)

	u, err := url.Parse(src)
	if err != nil {
		return err
	}
	if force == "" {
		force = u.Scheme
	}

	g, ok := Getters[force]
	if !ok {
		return fmt.Errorf(
			"module download not supported for scheme '%s'", force)
	}

	err = g.Get(dst, u)
	if err != nil {
		err = fmt.Errorf("error downloading module '%s': %s", src, err)
	}

	return err
}

// getRunCommand is a helper that will run a command and capture the output
// in the case an error happens.
func getRunCommand(cmd *exec.Cmd) error {
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	if err == nil {
		return nil
	}
	if exiterr, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return fmt.Errorf(
				"%s exited with %d: %s",
				cmd.Path,
				status.ExitStatus(),
				buf.String())
		}
	}

	return fmt.Errorf("error running %s: %s", cmd.Path, buf.String())
}

// getDirSubdir takes a source and returns a tuple of the URL without
// the subdir and the URL with the subdir.
func getDirSubdir(src string) (string, string) {
	// Calcaulate an offset to avoid accidentally marking the scheme
	// as the dir.
	var offset int
	if idx := strings.Index(src, "://"); idx > -1 {
		offset = idx + 3
	}

	// First see if we even have an explicit subdir
	idx := strings.Index(src[offset:], "//")
	if idx == -1 {
		return src, ""
	}

	idx += offset
	subdir := src[idx+2:]
	src = src[:idx]

	// Next, check if we have query parameters and push them onto the
	// URL.
	if idx = strings.Index(subdir, "?"); idx > -1 {
		query := subdir[idx:]
		subdir = subdir[:idx]
		src += query
	}

	return src, subdir
}

// getForcedGetter takes a source and returns the tuple of the forced
// getter and the raw URL (without the force syntax).
func getForcedGetter(src string) (string, string) {
	var forced string
	if ms := forcedRegexp.FindStringSubmatch(src); ms != nil {
		forced = ms[1]
		src = ms[2]
	}

	return forced, src
}