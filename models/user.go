package models

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Image string `json:"image"`
	Description string `json:"description"`
	GitHubURL string `json:"github_url"`
	LinkedInURL string `json:"linkedin_ulr"`
	Password string `json:"-"`
	PasswordHash string `json:"password;omitempty"`
}