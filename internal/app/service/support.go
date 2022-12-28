package service

// интерфейс методов хранилища
type SupportStorageProvider interface {
	Load()
	Close()
}

// структура конструктора бизнес логики
type SupportServices struct {
	storage SupportStorageProvider
	base    string
}

// конструктор бизнес  логики
func NewSupportService(s SupportStorageProvider, base string) *SupportServices {
	return &SupportServices{
		s,
		base,
	}
}

func (sr *SupportServices) Close() {
	sr.storage.Close()
}

func (sr *SupportServices) Load() {
	sr.storage.Load()
}
