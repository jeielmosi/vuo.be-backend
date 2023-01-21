package ports

type ShortenBulkRepository interface {
	GetByHash(hash string) (PigeonholeDTO[*ShortenBulkEntity], error)
	GetOldest(size uint32) (map[string]PigeonholeDTO[*ShortenBulkEntity], error)
	PostAtHash(hash string) error
	IncrementClicks(hash string) error
}
