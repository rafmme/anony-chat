package shared

import (
	"regexp"
)

func SignUpDataValidator(userSignupData UserSignupData) []map[string]string {
	var errList []map[string]string

	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	uppercasePattern := `[A-Z]`
	numberPattern := `[0-9]`
	specialCharPattern := `[!@#$%^&*()_+{}|:;<>,.?/~[\]\\]`

	validEmail := regexp.MustCompile(emailPattern)
	uppercaseRegex := regexp.MustCompile(uppercasePattern)
	numberRegex := regexp.MustCompile(numberPattern)
	specialCharRegex := regexp.MustCompile(specialCharPattern)

	isEmailValid := validEmail.MatchString(userSignupData.Email)
	hasUppercase := uppercaseRegex.MatchString(userSignupData.Password)
	hasNumber := numberRegex.MatchString(userSignupData.Password)
	hasSpecialChar := specialCharRegex.MatchString(userSignupData.Password)
	isPasswordValid := hasUppercase && hasNumber && hasSpecialChar

	if !isEmailValid {
		errList = append(errList, map[string]string{
			"email": "You entered an invalid email address.",
		})
	}

	if !isPasswordValid {
		errList = append(errList, map[string]string{
			"password": "You entered an invalid password. Password must have at least one uppercase letter, at least one numeric digit, and at least one special character.",
		})
	}

	if len(userSignupData.Password) < 8 || len(userSignupData.Password) > 26 {
		errList = append(errList, map[string]string{
			"password": "Password length must be between 8 and 26 characters.",
		})
	}

	if userSignupData.ConfirmPassword != userSignupData.Password {
		errList = append(errList, map[string]string{
			"confirmPassword": "Password field and Confirm Password field must match.",
		})
	}

	return errList
}
