package model

import "time"

type ListNotesInput struct {
	Term  string    `json:"term"`
	Tag   string    `json:"tag"`
	Sort  string    `json:"sort"`
	Order OrderType `json:"order"`
	Limit int       `json:"limit"`
}

type AddNoteInput struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type PatchNoteInput struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type Note struct {
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Tags         []string  `json:"tags"`
	LastModified time.Time `json:"lastModified"`
}

type PartialNote struct {
	Title        string    `json:"title"`
	Tags         []string  `json:"tags"`
	LastModified time.Time `json:"lastModified"`
}
