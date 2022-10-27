package generate

import "strings"

func withTpName(tp string) string {
	ret := strings.Split(tp, ".")
	return ret[len(ret)-1]
}
