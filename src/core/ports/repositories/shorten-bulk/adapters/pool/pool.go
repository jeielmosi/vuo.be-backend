package adapters

type ShortenBulkRepositoryPool struct {
	readCount               int
	writeCount              int
	shortenBulkRepositories []*ShortenBulkRepository
}

func NewShortenBulkPoolRepository(shortenBulkRepositories []*ShortenBulkRepository) *ShortenBulkPoolRepository {
	//Pigeonhole principle
	pigeons := len(shortenBulkRepositories) + 1

	readCount := pigeons / 2
	writeCount := pigeons - readCount

	return &ShortenBulkRepositoryPool{
		readCount,
		writeCount,
		shortenBulkRepositories,
	}

}

/*
func (SBRP *ShortenBulkRepositoryPool) GetByHash(hash string) (*ShortenBulkDTO, error)
func (SBRP *ShortenBulkRepositoryPool) GetOldest(size uint32) (map[string]*ShortenBulkDTO, error)
func (SBRP *ShortenBulkRepositoryPool) PostAtHash(hash string) error
func (SBRP *ShortenBulkRepositoryPool) IncrementClicks(hash string) error
*/
