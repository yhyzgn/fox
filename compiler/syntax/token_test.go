// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-08-25 16:40
// version: 1.0.0
// desc   :

package syntax

import "testing"

func TestToken(t *testing.T) {
	tokens := []token{Identifier, If, And, When, AddAssign, DotDotDot}
	wants := []string{"IDENTIFIER", "if", "&", "when", "+=", "..."}

	for i, tk := range tokens {
		if tk.String() != wants[i] {
			t.Errorf("Want: %s, but found: %s", tk.String(), wants[i])
		}
	}
}
