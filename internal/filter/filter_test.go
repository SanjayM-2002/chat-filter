package filter

import "testing"

func TestMask(t *testing.T) {
	f := New()
	for _, w := range []string{"shit", "ass", "fuck", "bastard", "son of a bitch"} {
		f.AddWord(w)
	}

	tests := []struct {
		in, want string
	}{
		{"mondays are shit", "mondays are ****"},
		{"mondays are SHIT!", "mondays are ****!"},
		{"this is sh1t", "this is ****"},
		{"$hit happens", "**** happens"},
		{"shiiiiit man", "******** man"},
		{"s.h.i.t", "*******"},
		{"s h i t", "* * * *"},
		{"f-u-c-k you", "******* you"},
		{"you ass", "you ***"},
		{"I passed the class", "I passed the class"}, // must NOT be masked
		{"assistant manager", "assistant manager"},   // must NOT be masked
		{"this hit hard", "this hit hard"},           // must NOT be masked (t-h-i-s h-i-t)
		{"you son of a bitch", "you *** ** * *****"}, // phrase
		{"fuckfuck", "fuckfuck"},                     // no boundary between -> not caught (limitation)
	}

	for _, tt := range tests {
		if got := f.Mask(tt.in); got != tt.want {
			t.Errorf("Mask(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
