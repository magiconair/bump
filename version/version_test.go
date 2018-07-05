package version

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		s string
		v *Version
	}{
		{"1.0", &Version{Major: 1}},
		{"v1.0", &Version{Prefix: "v", Major: 1}},
		{"v0.28", &Version{Prefix: "v", Major: 0, Minor: 28}},
		{"v1.0.1", &Version{Prefix: "v", Major: 1, Patch: 1}},
		{"v1.1.1", &Version{Prefix: "v", Major: 1, Minor: 1, Patch: 1}},
		{"v1.1", &Version{Prefix: "v", Major: 1, Minor: 1}},
		{"v1.1-foo", &Version{Prefix: "v", Major: 1, Minor: 1, Extra: "foo"}},
		{"v1.1.1-foo", &Version{Prefix: "v", Major: 1, Minor: 1, Patch: 1, Extra: "foo"}},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			v, err := Parse(tt.s)
			if err != nil {
				t.Fatalf("got error %s want nil", err)
			}
			if got, want := v, tt.v; *got != *want {
				t.Fatalf("\ngot  %#v\nwant %#v", *got, *want)
			}
			if got, want := v.String(), tt.s; got != want {
				t.Fatalf("got version %s want %s", got, want)
			}
		})
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		in, out []string
	}{
		{[]string{"1.0", "2.0"}, []string{"1.0", "2.0"}},
		{[]string{"2.0", "1.0"}, []string{"1.0", "2.0"}},
		{[]string{"1.2", "1.0"}, []string{"1.0", "1.2"}},
		{[]string{"1.0.1", "1.0"}, []string{"1.0", "1.0.1"}},
		{[]string{"v0.28", "v0.9"}, []string{"v0.9", "v0.28"}},
		{[]string{"v0.24", "v0.22-test"}, []string{"v0.22-test", "v0.24"}},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%#v", tt.in)
		t.Run(name, func(t *testing.T) {
			vv, err := ParseAll(tt.in)
			if err != nil {
				t.Fatalf("got error %s want nil", err)
			}
			var ss []string
			for _, v := range vv {
				ss = append(ss, v.String())
			}
			if got, want := ss, tt.out; !reflect.DeepEqual(got, want) {
				t.Fatalf("\ngot  %#v\nwant %#v", got, want)
			}
		})
	}
}
