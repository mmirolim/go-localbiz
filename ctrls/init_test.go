package ctrls

import "testing"

func TestGetUrl(t *testing.T) {
	data := []struct {
		in  []string // input
		out string   // expected result
	}{
		{[]string{"123", "456", "789"}, "/123/456/789"},
		{[]string{"1 2 3", "45 6", "abc"}, "/1_2_3/45_6/abc"},
	}
	for _, v := range data {
		got := GetUrl(v.in)
		if got != v.out {
			t.Errorf("expected %v got %v", v.out, got)
		}
	}

}
