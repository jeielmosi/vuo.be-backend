package adapters

type ShortenBulkRepositoryPigeonhole struct {
	readCount               int
	writeCount              int
	shortenBulkRepositories []*ShortenBulkRepository
}

func NewShortenBulkRepositoryPigeonhole(shortenBulkRepositories []*ShortenBulkRepository) *ShortenBulkRepositoryPigeonhole {
	//Pigeonhole principle
	pigeonsHoles := len(shortenBulkRepositories)
	pigeons := pigeonsHoles + 1

	readCount := pigeons / 2
	writeCount := pigeons - readCount

	return &ShortenBulkRepositoryPigeonhole{
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
