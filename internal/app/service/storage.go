package services

// Единый интерфейс хранилища для упрощения его инциализации
type StorageProvider interface {
	PutStorageProvider
	GetStorageProvider
	DeleteStorageProvider
	PingStorageProvider
	SupportStorageProvider
} 