package validator

import( 
  "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator() *validator.Validate{
	validate = validator.New()
	return validate
}