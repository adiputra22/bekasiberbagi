package user

type UserFormatter struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Token      string `json:"token"`
	ImageUrl   string `json:"image_url"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{}
	formatter.ID = user.ID
	formatter.Name = user.Name
	formatter.Occupation = user.Occupation
	formatter.Email = user.Email
	formatter.Token = token
	formatter.ImageUrl = user.AvatarFileName

	return formatter
}
