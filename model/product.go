package model

type Product struct {
	ID          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Price       uint32 `db:"price" json:"price"`
	Description string `db:"description" json:"description"`
	Command     string `db:"command" json:"-"` // why client would need command lol
}
