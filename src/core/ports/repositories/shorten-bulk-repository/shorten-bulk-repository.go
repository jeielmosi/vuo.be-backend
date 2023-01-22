package ports

type ShortenBulkRepository interface {
	GetByHash(hash string) (*ShortenBulkEntity, error)
	GetOldest(size uint) (map[string]*ShortenBulkEntity, error)
	PostAtHash(hash string) error
	IncrementClicks(hash string) error
}
