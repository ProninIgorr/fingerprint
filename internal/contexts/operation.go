package contexts

import (
	"fmt"
	"strings"
)

type Operation string

func (o Operation) String() string {
	return string(o)
}

type Operations struct {
	Stack []Operation
}

func (ops Operations) Empty() bool {
	return len(ops.Stack) == 0
}

func (ops Operations) Equal(ops2 Operations) bool {
	return ops.String() == ops2.String()
}

func (ops *Operations) Add(op Operation) {
	if len(ops.Stack) > 0 && ops.Stack[len(ops.Stack)-1] == op {
		return
	}
	ops.Stack = append(ops.Stack, op)
}

func (ops Operations) First() Operation {
	if len(ops.Stack) > 0 {
		return ops.Stack[0]
	}
	return ""
}

func (ops Operations) Last() Operation {
	if len(ops.Stack) > 0 {
		return ops.Stack[len(ops.Stack)-1]
	}
	return ""
}

func (ops Operations) Contains(wanted Operation) (result bool) {
	for _, op := range ops.Stack {
		if op == wanted {
			return true
		}
	}
	return
}

func (ops Operations) String() string {
	sb := strings.Builder{}
	for i, op := range ops.Stack {
		if i > 0 {
			sb.WriteString("/")
		}
		sb.WriteString(fmt.Sprintf("%s", op))
	}
	return sb.String()
}
