package transformers

import "bytes"

// Permutations merges the values of parts together to find all the
// possible combinations
//
// Input:
//
//	[][]string {
//		[]string{"A", "B", "C"}
//		[]string{"1", "2"}
//	}
//
// Output:
//
//	[]string {"A1", "A2", "B1", "B2", "C1", "C2"}
//
// Resulting length is length of each part passed multiplied by each
// other. So in the example above its 3 x 2 = 6
//
// Is the basis for generating complete table rows from api data.
func Permutations(parts ...[]string) (ret []string) {
	{
		var n = 1
		for _, ar := range parts {
			n *= len(ar)
		}
		ret = make([]string, 0, n)
	}
	var at = make([]int, len(parts))
	var buf bytes.Buffer
loop:
	for {
		// increment position counters
		for i := len(parts) - 1; i >= 0; i-- {
			if at[i] > 0 && at[i] >= len(parts[i]) {
				if i == 0 || (i == 1 && at[i-1] == len(parts[0])-1) {
					break loop
				}
				at[i] = 0
				at[i-1]++
			}
		}
		// construct permutated string
		buf.Reset()
		for i, ar := range parts {
			var p = at[i]
			if p >= 0 && p < len(ar) {
				buf.WriteString(ar[p])
			}
		}
		ret = append(ret, buf.String())
		at[len(parts)-1]++
	}
	return ret
}
