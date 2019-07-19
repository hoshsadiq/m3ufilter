package main

import (
	"fmt"
	"sort"
)

var groupOrder = map[string]int{
	"Group1": 1,
	"Group2": 2,
	"Group3": 3,
	"Group4": 4,
}

type Stream struct {
	Group string
}

type Streams []*Stream

func (s Streams) Len() int {
	return len(s)
}

func (s Streams) Less(i, j int) bool {
	return groupOrder[s[i].Group] < groupOrder[s[j].Group]
}

func (s Streams) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func main() {
	streams := Streams{
		&Stream{"Group2"},
		&Stream{"Group4"},
		&Stream{"Group4"},
		&Stream{"Group4"},
		&Stream{"Group4"},
		&Stream{"Group4"},
		&Stream{"Group1"},
		&Stream{"Group3"},
	}

	sort.Sort(streams)

	fmt.Printf("%v\n", streams)
}
