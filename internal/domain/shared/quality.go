package shared

import (
	"strings"
)

type QualityType string

// TODO: add more quality types
const (
	Resolution QualityType = "resolution" // 1080p, 4k, etc.
	Audio      QualityType = "audio"      // 5.1, 7.1, etc.
	Format     QualityType = "format"     // mkv, mp4, etc.
)

var resolutions = map[string]QualityType{
	"1080p": Resolution,
	"4k":    Resolution,
}

type QualityPreference map[QualityType]string

type QualityMatcher func(qualities []string) QualityPreference

func videoQualityMatcher(qualities []string) QualityPreference {
	qualityPreference := QualityPreference{}
	for _, quality := range qualities {
		if resolution, ok := resolutions[quality]; ok {
			qualityPreference[resolution] = quality
		}
	}
	return qualityPreference
}

func Match(typ MediaType, quality string) QualityPreference {
	qualities := strings.Split(quality, ";")
	switch typ {
	case Movie, TvShow:
		return videoQualityMatcher(qualities)
	default:
		return nil
	}
}
