package domain

import "time"

type Report struct {
	ID           string   `json:"id"`
	LayoutID     string   `json:"layout_id"`
	Title        string   `json:"title,omitempty"`
	TemplateKind string   `json:"template_kind,omitempty"`
	Extension    string   `json:"extension,omitempty"`
	Subscribers  []string `json:"subscribers,omitempty"`
	FilesCount   int      `json:"files_count,omitempty"`
}

type Reports []Report

func (rs *Reports) AddItem(r Report) {
	(*rs) = append((*rs), r)
}

type ReportFiles []ReportFile

type ReportFile struct {
	FileID    string    `json:"file_id"`
	Name      string    `json:"name"`
	MimeType  string    `json:"mime_type"`
	Content   []byte    `json:"content,omitempty"`
	Size      int       `json:"size"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	SentAt    time.Time `json:"sent_at,omitempty"`
}
