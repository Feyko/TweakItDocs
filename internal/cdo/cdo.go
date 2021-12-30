package cdo

import "strings"

const CDOPrefix = "Default__"

func ToClassName(s string) string {
	return strings.TrimPrefix(s, CDOPrefix)
}

func IsCDOOfClass(s, class string) bool {
	return s == CDOPrefix+class
}

func Is(s string) bool {
	return strings.HasPrefix(s, CDOPrefix)
}
