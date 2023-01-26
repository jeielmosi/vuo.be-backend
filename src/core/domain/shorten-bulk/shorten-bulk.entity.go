package entities

type ShortenBulkEntity struct {
	URL    string
	Clicks uint64
}

func NewShortenBulkEntity(
	url string,
	clicks uint64,
) *ShortenBulkEntity {
	return &ShortenBulkEntity{
		URL:    url,
		Clicks: clicks,
	}
}
