package service

// StorageProvider единый интерфейс хранилища для упрощения его инциализации.
type StorageProvider interface {
	PutStorageProvider
	GetStorageProvider
	DeleteStorageProvider
	PingStorageProvider
	SupportStorageProvider
}
