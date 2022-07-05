package model

type Stats struct {
	Clicks    int            `json:"clicks"`
	Sources   map[string]int `json:"sources"`
	Devices   map[string]int `json:"devices"`
	Browsers  map[string]int `json:"browsers"`
	Languages map[string]int `json:"languages"`
	Systems   map[string]int `json:"systems"`
}

func (s *Stats) Initialize() {
	s.Clicks = 0
	s.Sources = make(map[string]int)
	s.Devices = make(map[string]int)
	s.Browsers = make(map[string]int)
	s.Languages = make(map[string]int)
	s.Systems = make(map[string]int)
}

func (s *Stats) IncrementIfExists(click Click) {
	if click.Browser != "" {
		s.Browsers[click.Browser] += 1
	}

	if click.Source != "" {
		s.Sources[click.Source] += 1
	}

	if click.Device != "" {
		s.Devices[click.Device] += 1
	}

	if click.Language != "" {
		s.Languages[click.Language] += 1
	}

	if click.System != "" {
		s.Systems[click.System] += 1
	}
}
