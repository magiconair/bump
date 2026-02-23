// Copyright 2026 The Bump Authors. All rights reserved. See LICENSE.

package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

func IsEmptyRepository() (bool, error) {
	out, err := exec.Command("git", "count-objects", "-v").Output()
	if err != nil {
		return false, err
	}
	return bytes.Contains(out, []byte("count: 0")), nil
}

func Tags(service string) ([]Version, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("git", "tag")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	if stderr.Len() > 0 {
		return nil, fmt.Errorf("%s", stderr.String())
	}

	vv, err := Read(&stdout, service)
	if err != nil {
		return nil, err
	}
	return vv, nil
}

func Tag(v Version) error {
	if err := signedTag(v); err != nil {
		return annotatedTag(v)
	}
	return nil
}

func annotatedTag(v Version) error {
	return exec.Command("git", "tag", "-a", v.String(), "-m", v.String()).Run()
}

func signedTag(v Version) error {
	return exec.Command("git", "tag", "-s", v.String(), "-m", v.String()).Run()
}
