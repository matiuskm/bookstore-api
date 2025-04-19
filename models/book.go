package models

type Book struct {
	ID	 		uint 		`json:"id" gorm:"primaryKey"`
	Title		string 		`json:"title"`
	Author		string 		`json:"author"`
	Price		int 		`json:"price"`
	CategoryID	uint 		`json:"categoryId"`
	Category	Category 	`json:"category" gorm:"foreignKey:CategoryID"`
}

type BookPatch struct {
	Title		*string 	`json:"title"`
	Author		*string 	`json:"author"`
	Price		*int 		`json:"price"`
	CategoryID	*uint 		`json:"categoryId"`
}