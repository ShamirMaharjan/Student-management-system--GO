package body

type Student struct {
	Id    int
	Name  string `validate:"required"`
	Age   int    `validate:"required"`
	Email string `validate:"required,email"`
}
