package types

import "time"

type File struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	ParentID  *string   `json:"parent_id,omitempty"`
	ChannelID int64     `json:"channel_id"`
	MessageID int       `json:"message_id"`
	Parts     int       `json:"parts"`
	Encrypted bool      `json:"encrypted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type UploadSession struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	FileName  string    `json:"file_name"`
	FileSize  int64     `json:"file_size"`
	ChunkSize int64     `json:"chunk_size"`
	Chunks    int       `json:"chunks"`
	Uploaded  int       `json:"uploaded"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
