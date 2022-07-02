package user

type User struct {
	ID                    int64        `json:"id"`
	UID                   string       `json:"uid"`
	Password              string       `json:"password"`
	FirstName             string       `json:"first_name"`
	LastName              string       `json:"last_name"`
	Username              string       `json:"username"`
	Email                 string       `json:"email"`
	Avatar                string       `json:"avatar"`
	Gender                string       `json:"gender"`
	PhoneNumber           string       `json:"phone_number"`
	SocialInsuranceNumber string       `json:"social_insurance_number"`
	DateOfBirth           string       `json:"date_of_birth"`
	Employment            Employment   `json:"employment"`
	Address               Address      `json:"address"`
	CreditCard            CreditCard   `json:"credit_card"`
	Subscription          Subscription `json:"subscription"`
}
