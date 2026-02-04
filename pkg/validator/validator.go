package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validate is the package-level validator instance.
var Validate *validator.Validate

// Init initializes the validator. Call once on program startup (or tests).
func Init() {
	Validate = validator.New()
}
