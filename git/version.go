package git

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

var re = regexp.MustCompile(`(v)?([0-9]+)\.([0-9]+)(\.([0-9]+))?(-(.*))?`)

func Read(r io.Reader) ([]Version, error) {
	var vv byVersion
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		v, err := Parse(sc.Text())
		if v.IsZero() || err != nil {
			continue
		}
		vv = append(vv, v)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	sort.Sort(byVersion(vv))
	return vv, nil
}

func ParseAll(s []string) ([]Version, error) {
	var vv byVersion
	for _, x := range s {
		v, err := Parse(x)
		if v.IsZero() || err != nil {
			continue
		}
		vv = append(vv, v)
	}
	sort.Sort(byVersion(vv))
	return vv, nil
}

func Parse(s string) (Version, error) {
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
		return Version{}, nil
	}

	prefix := m[1]
	major := atoi(m[2])
	minor := atoi(m[3])
	patch := atoi(m[5])
	extra := m[7]
	if err != nil {
		return Version{}, err
	}

	return Version{
		Prefix: prefix,
		Major:  major,
		Minor:  minor,
		Patch:  patch,
		Extra:  extra,
	}, nil
}

func (v Version) Ints() []int {
	return []int{v.Major, v.Minor, v.Patch}
}

func (v Version) IsZero() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0
}

func (v Version) Bump() Version {
	return v.BumpPatch()
}

func (v Version) BumpMajor() Version {
	return Version{
		Prefix: v.Prefix,
		Major:  v.Major + 1,
	}
}

func (v Version) BumpMinor() Version {
	return Version{
		Prefix: v.Prefix,
		Major:  v.Major,
		Minor:  v.Minor + 1,
	}
}

func (v Version) BumpPatch() Version {
	return Version{
		Prefix: v.Prefix,
		Major:  v.Major,
		Minor:  v.Minor,
		Patch:  v.Patch + 1,
	}
}

func (v Version) String() string {
	s := fmt.Sprintf("%s%d.%d.%d", v.Prefix, v.Major, v.Minor, v.Patch)
	if v.Extra != "" {
		s += "-" + v.Extra
	}
	return s
}

type byVersion []Version

func (v byVersion) Len() int      { return len(v) }
func (v byVersion) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v byVersion) Less(i, j int) bool {
	a, b := v[i].Ints(), v[j].Ints()

	return a[0] < b[0] ||
		a[0] == b[0] && a[1] < b[1] ||
		a[0] == b[0] && a[1] == b[1] && a[2] < b[2] ||
		a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && v[i].Extra < v[j].Extra
}
