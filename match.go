/*
references:
http://swtch.com/~rsc/regexp/regexp2.html
http://research.swtch.com/sparse
*/
package nm

type Program []Inst

type Inst struct {
	Op      Op
	Predict func(*Node) bool
	A, B    int
}

type Op int

const (
	Predict Op = iota
	Ok
	Jump
	Split
)

func Match(node *Node, program Program) []*Node {
	var result []*Node
	walk(node, []*Node{node}, program, &result)
	return result
}

func walk(node *Node, path []*Node, program Program, result *[]*Node) {
	if program.Match(path) {
		*result = append(*result, node)
	}
	for _, c := range node.Children {
		walk(c, append(path, c), program, result)
	}
}

type _Threads struct {
	sparse, dense []int
	n             int
}

func (t *_Threads) clear() {
	t.n = 0
}

func (t *_Threads) add(i int) {
	if t.sparse[i] < t.n && t.dense[t.sparse[i]] == i {
		return
	}
	t.dense[t.n] = i
	t.sparse[i] = t.n
	t.n++
}

func (p Program) Match(path []*Node) bool {
	activeThreads := &_Threads{make([]int, len(p)), make([]int, len(p)), 0}
	nextThreads := &_Threads{make([]int, len(p)), make([]int, len(p)), 0}
	activeThreads.add(0)
	maxMatched := -1
	for n, node := range path {
		for i := 0; i < activeThreads.n; i++ {
			pc := activeThreads.dense[i]
			inst := p[pc]
			switch inst.Op {
			case Predict:
				if inst.Predict(node) {
					nextThreads.add(pc + 1)
					maxMatched = n
				}
			case Ok:
				return maxMatched == len(path)-1 // match exactly
			case Jump:
				activeThreads.add(inst.A)
			case Split:
				activeThreads.add(inst.A)
				activeThreads.add(inst.B)
			}
		}
		activeThreads, nextThreads = nextThreads, activeThreads
		nextThreads.clear()
	}
	for i := 0; i < activeThreads.n; i++ {
		pc := activeThreads.dense[i]
		inst := p[pc]
		switch inst.Op {
		case Ok:
			return true
		case Jump:
			activeThreads.add(inst.A)
		case Split:
			activeThreads.add(inst.A)
			activeThreads.add(inst.B)
		}
	}
	return false
}
