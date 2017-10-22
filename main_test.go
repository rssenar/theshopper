package main

import (
	"fmt"
	"testing"
)

func TestRemSep(t *testing.T) {
	cases := []struct{ input, expected string }{
		{`-,  Ric-.*(),har(),d`, `Richard`},
		{`"""test"""`, "test"},
		{` g a b i e `, "gabie"},
	}

	for _, c := range cases {
		str := remSep(c.input)
		if str != c.expected {
			t.Log("error should be "+c.expected+" but got", str)
			t.Fail()
		}
	}
}

func TestNewRecord(t *testing.T) {
	cases := []struct {
		rec            []string
		mrc            int
		BundlePerRoute int
		LastBundle     int
	}{
		{[]string{"Corona", "CA", "92882", "C001", "625", "50"}, 0, 13, 25},
		{[]string{"Corona", "CA", "92882", "C001", "650", "50"}, 0, 13, 50},
		{[]string{"Corona", "CA", "92882", "C001", "650", "50"}, 10, 14, 10},
		{[]string{"Corona", "CA", "92882", "C001", "650", "50"}, -10, 13, 40},
		{[]string{"", "CA", "92882", "", "650", "50"}, 10, 14, 10},
	}

	for _, c := range cases {
		rec, err := newRecord(c.rec, c.mrc)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if rec.BundlePerRoute != c.BundlePerRoute {
			t.Log("BundlePerRoute error: should be " + fmt.Sprint(c.BundlePerRoute) + " but got " + fmt.Sprint(rec.BundlePerRoute))
			t.Fail()
		}
		if rec.LastBundle != c.LastBundle {
			t.Log("LastBundle error: should be " + fmt.Sprint(c.LastBundle) + " but got " + fmt.Sprint(rec.LastBundle))
			t.Fail()
		}
	}
}
