package options

import (
	"fmt"
	"strings"
)

const (
	SortIteratorKindList    = "list"
	SortIteratorKindHashmap = "hashmap"
	SortIteratorKindAll     = "all"
)

type SortIterators struct {
	SortList    bool
	SortHashmap bool
}

func NewDefaultSortIterators() *SortIterators {
	return &SortIterators{
		SortList:    false,
		SortHashmap: false,
	}
}

func (s *SortIterators) SetFlagFromKind(k string) error {
	switch k {
	case SortIteratorKindList:
		s.SortList = true
	case SortIteratorKindHashmap:
		s.SortHashmap = true
	case SortIteratorKindAll:
		s.SortList = true
		s.SortHashmap = true
	default:
		return fmt.Errorf(
			"the sort iterators kind '%s' is unsupported, it must be '%s', '%s' or '%s'",
			k,
			SortIteratorKindList,
			SortIteratorKindHashmap,
			SortIteratorKindAll,
		)
	}

	return nil
}

func NewSortIteratorsFromLine(line string) (*SortIterators, error) {
	sortIterators := NewDefaultSortIterators()

	kinds := strings.Split(line, ",")

	for _, kind := range kinds {
		err := sortIterators.SetFlagFromKind(kind)
		if err != nil {
			return nil, err
		}
	}

	return sortIterators, nil
}
