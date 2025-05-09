// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repo

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Message struct {
	ID        string           `json:"id"`
	ThreadID  string           `json:"thread_id"`
	Content   string           `json:"content"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type Order struct {
	ID        string           `json:"id"`
	Amount    string           `json:"amount"`
	Number    string           `json:"number"`
	Status    string           `json:"status"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

type Thread struct {
	ID        string           `json:"id"`
	Title     string           `json:"title"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}
