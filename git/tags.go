package git

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/magiconair/bump/version"
)

func Tags() (version.Versions, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("git", "tag")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	if stderr.Len() > 0 {
		return nil, fmt.Errorf(stderr.String())
	}

	vv, err := version.Read(&stdout)
	if err != nil {
		return nil, err
	}
	return vv, nil
}

func Tag(v *version.Version) error {
	if err := signedTag(v); err != nil {
		return annotatedTag(v)
	}
	return nil
}

func annotatedTag(v *version.Version) error {
	return exec.Command("git", "tag", "-a", v.String(), "-m", v.String()).Run()
}

func signedTag(v *version.Version) error {
	return exec.Command("git", "tag", "-s", v.String(), "-m", v.String()).Run()
}
