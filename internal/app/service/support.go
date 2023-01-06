package service

// SupportStorageProvider интерфейс обслуживающих методов хранилища.
type SupportStorageProvider interface {
	Load()
	Close()
}

// SupportServices структура конструктора ckjz Support бизнес логики обслуживающих методов хранилища.
type SupportServices struct {
	storage SupportStorageProvider
	base    string
}

// NewSupportService конструктор бизнес  логики.
func NewSupportService(s SupportStorageProvider, base string) *SupportServices {
	return &SupportServices{
		s,
		base,
	}
}

// Close метод закрытия соединения с хранилищем.
func (sr *SupportServices) Close() {
	sr.storage.Close()
}

// Load метод загрузки из файлового хранилица.
func (sr *SupportServices) Load() {
	sr.storage.Load()
}
