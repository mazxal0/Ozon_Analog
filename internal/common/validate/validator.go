package validate

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()

	// Регистрация кастомных валидаторов
	registerCustomValidators()
}

func registerCustomValidators() {
	// Валидатор для пользовательских ролей
	Validate.RegisterValidation("userrole", func(fl validator.FieldLevel) bool {
		role := fl.Field().String()
		validRoles := map[string]bool{
			"student": true,
			"tutor":   true,
			"admin":   true,
		}
		return validRoles[role]
	})

	// Валидатор для урока статусов
	Validate.RegisterValidation("lessonstatus", func(fl validator.FieldLevel) bool {
		status := fl.Field().String()
		validStatuses := map[string]bool{
			"free":     true,
			"booked":   true,
			"finished": true,
			"canceled": true,
		}
		return validStatuses[status]
	})

	// Валидатор для телефона
	Validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		// Международный формат E.164
		matched, _ := regexp.MatchString(`^\+[1-9]\d{1,14}$`, phone)
		return matched
	})

	// Валидатор для пароля
	Validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 {
			return false
		}
		// Должен содержать хотя бы одну цифру, одну букву и один спецсимвол
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
		hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

		return hasNumber && hasLetter && hasSpecial
	})

	Validate.RegisterValidation("starttime", func(fl validator.FieldLevel) bool {
		t, ok := fl.Field().Interface().(time.Time)
		if !ok {
			return false
		}
		return !t.Before(time.Now()) // t >= now
	})

	Validate.RegisterValidation("endtime", func(fl validator.FieldLevel) bool {
		endTime, ok := fl.Field().Interface().(time.Time)
		if !ok {
			return false
		}
		lesson := fl.Top().Interface()
		m, ok := lesson.(map[string]interface{})
		if ok {
			if start, ok := m["StartTime"].(time.Time); ok {
				return endTime.After(start)
			}
		}

		// fallback: cannot validate
		return true
	})
}
