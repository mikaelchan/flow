package shared

type MediaType uint8

const (
	Movie MediaType = iota
	TvShow
	Anime
	Music
	Adult
	Book
)

func (mt MediaType) IsValid() bool {
	return mt >= Movie && mt <= Book
}

func (mt MediaType) String() string {
	return [...]string{"movie", "tvshow", "anime", "music", "adult", "book"}[mt]
}

type QualityPreference map[string]string
