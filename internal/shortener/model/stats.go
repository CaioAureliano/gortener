package model

type Stats struct {
	Clicks    int
	Sources   map[string]int
	Devices   map[string]int
	Browsers  map[string]int
	Languages map[string]int
	Systems   map[string]int
}
