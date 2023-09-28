package pkg

type ComboQueue struct {
	combos chan *Combo
	size   int
}

func NewComboQueue(combos []*Combo) *ComboQueue {
	q := &ComboQueue{
		combos: make(chan *Combo, len(combos)),
		size:   len(combos),
	}
	for _, combo := range combos {
		if !combo.IsValid() {
			continue
		}

		q.Enqueue(combo)
	}
	return q
}

func (q *ComboQueue) Enqueue(c *Combo) {
	q.combos <- c
}

func (q *ComboQueue) Dequeue() (*Combo, bool) {
	combo, ok := <-q.combos
	return combo, ok
}

func (q *ComboQueue) Close() {
	close(q.combos)
}
