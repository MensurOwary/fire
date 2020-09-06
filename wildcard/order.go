package wildcard

import (
	"github.com/mensurowary/fire/arg"
	log "github.com/sirupsen/logrus"
	r "regexp"
	"strconv"
)

// OrderCounter keeps track of the next value for the order wildcard
type OrderCounter struct {
	numeric int
	text    string
}

// ReplaceIfNeeded replaces the wildcard with the correct next value
func (w *OrderCounter) ReplaceIfNeeded(input string) string {
	if w.contains(input) {
		return w.replace(input)
	}
	return input
}

func (w *OrderCounter) contains(input string) bool {
	res := orderRegex.MatchString(input)
	log.Debug("Ordering match result with %s : %v\n", input, res)
	return res
}

func (w *OrderCounter) replace(input string) string {
	return orderRegex.ReplaceAllString(input, w.advance())
}

func (w *OrderCounter) advance() string {
	if w.text == "" {
		v := w.numeric
		w.numeric = v + 1
		return strconv.Itoa(v)
	}
	v := next(w.text)
	w.text = v
	return v
}

func next(text string) string {
	chars := []rune(text)
	incNext := false
	for i := len(chars) - 1; i >= 0; i-- {
		char := chars[i]
		if char == 'z' {
			chars[i] = 'a'
			incNext = true
		} else {
			chars[i] = rune(int(chars[i]) + 1)
			incNext = false
			break
		}
	}
	if incNext {
		chars = append([]rune{'a'}, chars...)
	}
	return string(chars)
}

// ordering wildcard regex
// #{o:number} => #{o:12} will start counting from 12
// #{o:ascii} => #{o:a} will start counting from a until z...z
var orderRegex = r.MustCompile("#{o(:(.+))?}")

// MakeOrdering initializes the order counter
func MakeOrdering(arg arg.Arguments) OrderCounter {
	subMatches := orderRegex.FindStringSubmatch(arg.To)
	wildCardArg := ""
	if len(subMatches) == 3 {
		wildCardArg = subMatches[2]
	}

	if wildCardArg == "" {
		return OrderCounter{numeric: 1}
	} else {
		if val, err := strconv.Atoi(wildCardArg); err != nil {
			return OrderCounter{text: wildCardArg}
		} else {
			return OrderCounter{numeric: val}
		}
	}
}
