package repository

type KeyOriginalURL struct {
	Key       string `json:"short_url" gorm:"primaryKey;size:32"`
	Original  string `json:"original_url"`
	UserID    string `json:"user_id" gorm:"primaryKey;size:36"`
	IsDeleted bool   `json:"is_deleted" gorm:"is_deleted"`
}
