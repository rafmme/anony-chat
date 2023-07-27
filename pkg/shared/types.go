package shared

type Env struct {
	Port       string `barfenv:"key=PORT;required=true"`
	DbHost     string `barfenv:"key=DB_HOST;required=true"`
	DbPort     string `barfenv:"key=DB_PORT;required=true"`
	DbName     string `barfenv:"key=DB_NAME;required=true"`
	DbUser     string `barfenv:"key=DB_USER;required=true"`
	DbPassword string `barfenv:"key=DB_PASSWORD;required=true"`
	SecretKey  string `barfenv:"key=SECRET_KEY;required=true"`
}

type UserSignupData struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type ErrorResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Errors     interface{} `json:"errors"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type AuthData struct {
}
