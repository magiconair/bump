package version

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
)

type Version struct {
	Prefix string
	Major  int
	Minor  int
	Patch  int
	Extra  string
}

type Versions []*Version

func (v Versions) Len() int      { return len(v) }
func (v Versions) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v Versions) Less(i, j int) bool {
	a, b := v[i].Ints(), v[j].Ints()

	return a[0] < b[0] ||
		a[0] == b[0] && a[1] < b[1] ||
		a[0] == b[0] && a[1] == b[1] && a[2] < b[2] ||
		a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && v[i].Extra < v[j].Extra
}

var re = regexp.MustCompile(`(v)?([0-9]+)\.([0-9]+)(\.([0-9]+))?(-(.*))?`)

func Read(r io.Reader) (Versions, error) {
	var vv Versions
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		v, err := Parse(sc.Text())
		if v == nil || err != nil {
			continue
		}
		vv = append(vv, v)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	sort.Sort(vv)
	return vv, nil
}

func ParseAll(s []string) (Versions, error) {
	var vv Versions
	for _, x := range s {
		v, err := Parse(x)
		if v == nil || err != nil {
			continue
		}
		vv = append(vv, v)
	}
	sort.Sort(vv)
	return vv, nil
}

func Parse(s string) (*Version, error) {
	var err error
	atoi := func(s string) int {
		if s == "" || err != nil {
			return 0
		}
		var n int
		n, err = strconv.Atoi(s)
		if err != nil {
			return 0
		}
		return n
	}

	m := re.FindStringSubmatch(s)
	if m == nil {
		return nil, nil
	}

	prefix := m[1]
	major := atoi(m[2])
	minor := atoi(m[3])
	patch := atoi(m[5])
	extra := m[7]
	if err != nil {
		return nil, err
	}

	return &Version{
		Prefix: prefix,
		Major:  major,
		Minor:  minor,
		Patch:  patch,
		Extra:  extra,
	}, nil
}

func (v *Version) Ints() []int {
	return []int{v.Major, v.Minor, v.Patch}
}

func (v *Version) Bump() *Version {
	return v.BumpPatch()
}

func (v *Version) BumpMajor() *Version {
	return &Version{
		Prefix: v.Prefix,
		Major:  v.Major + 1,
	}
}

func (v *Version) BumpMinor() *Version {
	return &Version{
		Prefix: v.Prefix,
		Major:  v.Major,
		Minor:  v.Minor + 1,
	}
}

func (v *Version) BumpPatch() *Version {
	return &Version{
		Prefix: v.Prefix,
		Major:  v.Major,
		Minor:  v.Minor,
		Patch:  v.Patch + 1,
	}
}

func (v Version) String() string {
	s := fmt.Sprintf("%s%d.%d", v.Prefix, v.Major, v.Minor)
	if v.Patch > 0 {
		s += fmt.Sprintf(".%d", v.Patch)
	}
	if v.Extra != "" {
		s += "-" + v.Extra
	}
	return s
}
