package main

import "testing"

func TestCharSetHas(t *testing.T) {
	for ch := rune('a'); ch <= 'z'; ch++ {
		var cs charSet
		if cs.Has(ch) {
			t.Errorf("%c in charSet before it is added", ch)
		}
		cs.Add(ch)
		if !cs.Has(ch) {
			t.Errorf("%c not in charSet after it is added", ch)
		}
	}
}

func TestCharSetAddWouldFill(t *testing.T) {
	tests := []struct {
		ch            rune
		existingChars []rune
		want          bool
	}{
		{},
		{
			ch:            'f',
			existingChars: []rune{'a', 'b', 'c', 'd', 'e'},
		},
		{
			ch:            'c',
			existingChars: []rune{'a', 'b', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'},
			want:          true,
		},
	}
	for i, test := range tests {
		cs := newCharSetHelper(t, test.existingChars...)
		if want, got := test.want, cs.AddWouldFill(test.ch); want != got {
			t.Errorf("test %v: addWouldFill not equal: wanted %v, got %v", i, want, got)
		}
	}
}

func TestCharSetString(t *testing.T) {
	cs := newCharSetHelper(t, 'f', 'y', 'r', 'o', 't')
	if want, got := "[forty]", cs.String(); want != got {
		t.Errorf("wanted %q, got %q", want, got)
	}
}

func TestCharSetBadChars(t *testing.T) {
	badChars := []rune{'?', 'A', 'Z', ' ', '!', '`', '\n', 0, 0x7F, 0xFF}
	for i, ch := range badChars {
		t.Run("bad-add-#"+string(rune('0'+i)), func(t *testing.T) {
			var cs charSet
			if cs.Has(ch) {
				t.Errorf("bad character 0x%x in charSet", ch)
			}
			defer func() {
				r := recover()
				if _, ok := r.(error); r == nil || !ok {
					t.Errorf("expected panic error adding bad character")
				}
			}()
			cs.Add(ch)
		})
	}
}

func newCharSetHelper(t *testing.T, chars ...rune) charSet {
	t.Helper()
	var cs charSet
	for _, ch := range chars {
		cs.Add(ch)
	}
	return cs
}
