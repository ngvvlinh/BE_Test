package service

type (
	ControlMessage struct {
		Prio int
		Data []byte
	}

	ControlMessages []ControlMessage
)

func (h ControlMessages) Len() int            { return len(h) }
func (h ControlMessages) Less(i, j int) bool  { return h[i].Prio < h[j].Prio }
func (h ControlMessages) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *ControlMessages) Push(x interface{}) { *h = append(*h, x.(ControlMessage)) }

func (h *ControlMessages) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
