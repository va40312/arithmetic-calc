package expression

import "testing"

func TestCalculator(t *testing.T) {
	testCases := []struct {
		name        string
		text        string
		expected    float64
		expectError bool
	}{
		{
			name:        "Самый простой случай",
			text:        "2 + 2",
			expected:    4.0,
			expectError: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual, err := Calculate(test.text)

			if !test.expectError && err != nil {
				t.Errorf("not expected error (%s) in Calculate(%s) = %f, expected %f", err, test.text, actual, test.expected)
			} else if test.expectError && err == nil {
				t.Errorf("expected error, but not get it in Calculate(%s) = %f, expected %f", test.text, actual, test.expected)
			}

			if actual != test.expected {
				t.Errorf("Calculate(%s) = %f, expected %f", test.text, actual, test.expected)
			}
		})
	}
}
