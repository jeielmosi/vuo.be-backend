package ports

type ShortenBulkRepository interface {
	GetByHash(hash string) (*ShortenBulkDTO, error)
	GetOldest(size uint32) (map[string]*ShortenBulkDTO, error)
	PostAtHash(hash string) error
	IncrementClicks(hash string) error
}
