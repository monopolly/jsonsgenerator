package main

import (
	"sort"
)

// SortFieldsByGoMemoryPadding sorts fields to *usually* minimize struct padding on 64-bit Go.
// Heuristic: larger alignment first, then larger size, then name.
func Optimize(list []*Field) (sum int) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Go.Bits > list[j].Go.Bits
	})

	return StructSize(list)

	// for _, x := range list {
	// 	sum += x.Go.Bits
	// }
	// return
}

// StructSize calculates total struct size including padding.
func StructSize(fields []*Field) int {
	offset := 0
	maxAlign := 1

	for _, f := range fields {
		if f == nil {
			continue
		}
		if f.Go.Align <= 0 {
			f.Go.Align = 1
		}
		if f.Go.Align > maxAlign {
			maxAlign = f.Go.Align
		}

		// padding before field
		pad := padding(offset, f.Go.Align)
		offset += pad

		// field itself
		offset += f.Go.Bits
	}

	// tail padding: struct size must be multiple of maxAlign
	offset += padding(offset, maxAlign)

	return offset
}

func padding(offset, align int) int {
	if align <= 1 {
		return 0
	}
	rem := offset % align
	if rem == 0 {
		return 0
	}
	return align - rem
}
