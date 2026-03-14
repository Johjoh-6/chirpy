package utils

import "testing"

func TestRemoveProfanity(t *testing.T) {
	profanity := &Profanity{
		Word:     []string{"badword", "anotherbadword"},
		Replacer: "****",
	}
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "hello world",
			input:    "hello  world",
			expected: "hello  world",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "Unique badword",
			input:    "badword",
			expected: "****",
		},
		{
			name:     "Capitalize badword",
			input:    "Badword Paris Berlin",
			expected: "**** Paris Berlin",
		},
		{
			name:     "String with 2 badwords",
			input:    "It's a long sentence with badword and anotherbadword",
			expected: "It's a long sentence with **** and ****",
		},
	}

	for _, c := range cases {
		actual := profanity.RemoveProfanity(c.input)
		if actual != c.expected {
			t.Errorf("%s: expected %v, got %v", c.name, c.expected, actual)
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("%s: expected %v, got %v", c.name, c.expected, actual)
				break
			}
		}
	}
}
