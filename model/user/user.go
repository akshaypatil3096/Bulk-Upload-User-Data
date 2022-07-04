package user

import (
	"github.com/akshaypatil3096/Bulk-Upload-User-Data/model/payment"
	"github.com/akshaypatil3096/Bulk-Upload-User-Data/model/subscription"
)

type User struct {
	ID                    int64                `json:"id"`
	UID                   string               `json:"uid"`
	Password              string               `json:"password"`
	FirstName             string               `json:"first_name"`
	LastName              string               `json:"last_name"`
	Username              string               `json:"username"`
	Email                 string               `json:"email"`
	Avatar                string               `json:"avatar"`
	Gender                string               `json:"gender"`
	PhoneNumber           string               `json:"phone_number"`
	SocialInsuranceNumber string               `json:"social_insurance_number"`
	DateOfBirth           string               `json:"date_of_birth"`
	Employment            Employment           `json:"employment"`
	Address               Address              `json:"address"`
	CreditCard            payment.CreditCard   `json:"credit_card"`
	Subscription          subscription.Details `json:"subscription"`
}
