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

func FromString(s string) (MediaType, error) {
	for i, v := range [...]string{"movie", "tvshow", "anime", "music", "adult", "book"} {
		if v == s {
			return MediaType(i), nil
		}
	}
	return 0, ErrInvalidMediaType
}
