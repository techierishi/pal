package histm

import (
	"slices"

	"github.com/techierishi/pal/logr"
)

// Histlist is a deduplicated list of cmdLines
type Histlist struct {

	// list of commands lines (deduplicated)
	List []string
	// lookup: cmdLine -> last index
	LastIndex map[string]int
}

// New Histlist
func New() Histlist {
	return Histlist{
		LastIndex: make(map[string]int),
	}
}

// Copy Histlist
func Copy(hl Histlist) Histlist {
	newHl := New()
	// copy list
	newHl.List = make([]string, len(hl.List))
	copy(newHl.List, hl.List)
	// copy map
	for k, v := range hl.LastIndex {
		newHl.LastIndex[k] = v
	}
	return newHl
}

// AddCmdLine to the histlist
func (h *Histlist) AddCmdLine(cmdLine string) {
	logger := logr.GetLogInstance()

	// lenBefore := len(h.List)
	// lookup
	idx, found := h.LastIndex[cmdLine]
	if found {
		// remove duplicate
		if cmdLine != h.List[idx] {
			logger.Error().
				Str("indexKeyCmdLine", cmdLine).
				Str("actualCmdLine", h.List[idx]).
				Msg("Index key is different than actual cmd line in the list")
		}
		h.List = append(h.List[:idx], h.List[idx+1:]...)
		// idx++
		for idx < len(h.List) {
			cmdLn := h.List[idx]
			h.LastIndex[cmdLn]--
			if idx != h.LastIndex[cmdLn] {
				logger.Error().
					Int("actualPosition", idx).
					Int("indexedPosition", h.LastIndex[cmdLn]).
					Msg("Index position is different than actual position of the cmd line")

			}
			idx++
		}
	}
	// update last index
	h.LastIndex[cmdLine] = len(h.List)
	// append new cmdline
	h.List = append(h.List, cmdLine)
}

// AddHistlist contents of another histlist to this histlist
func (h *Histlist) AddHistlist(h2 Histlist) {
	for _, cmdLine := range h2.List {
		h.AddCmdLine(cmdLine)
	}
}

// AddHistlist contents of another histlist to this histlist
func (h *Histlist) Reverse() {
	slices.Reverse[[]string](h.List)
}
