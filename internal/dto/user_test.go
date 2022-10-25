package dto

import (
	"reflect"
	"testing"

	"github.com/Kin-dza-dzaa/userApi/internal/models"
)

func TestUserSignInDto_IntoUser(t *testing.T) {
	tests := []struct {
		dto     UserSignInDto
		want    *models.User
		wantErr bool
	}{
		{
			dto: UserSignInDto{
				Email: "testemail@gmail.com",
				Password: "123456789",
			},
			want: &models.User{
				Email: "testemail@gmail.com",
				Password: "123456789",
			},
			wantErr: false,
		},
		{
			dto: UserSignInDto{
				Email: "bad",
				Password: "bad",
			},
			want: nil,
			wantErr: true,
		},
		{
			dto: UserSignInDto{
				Email: "testemail@gmail.com",
				Password: "bad",
			},
			want: nil,
			wantErr: true,
		},
		{
			dto: UserSignInDto{
			},
			want: nil,
			wantErr: true,
		},
		{
			dto: UserSignInDto{
				Email: "",
				Password: "123123123213",
			},
			want: nil,
			wantErr: true,
		},
		{
			dto: UserSignInDto{
				Email: "testemail@gmail.com",
				Password: "",
			},
			want: nil,
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run("SignInDto", func(t *testing.T) {
			if tc.wantErr {
				user, err := tc.dto.IntoUser()
				if err == nil || user != nil {
					t.FailNow()
				}
				
			} else {
				user, err := tc.dto.IntoUser()
				if err != nil || !reflect.DeepEqual(tc.want, user) {
					t.FailNow()
				}
			}
		})
	}
}

func TestUserSignUpDto_IntoUser(t *testing.T) {
	tests := []struct {
		name string
		dto     UserSignUpDto
		want    *models.User
		wantErr bool
	}{
		{
			name: "good_user",
			dto: UserSignUpDto{
				Email: "testemail@gmail.com",
				UserName: "TestUser",
				Password: "123456789",
			},
			want: &models.User{
				Email: "testemail@gmail.com",
				UserName: "TestUser",
				Password: "123456789",
			},
			wantErr: false,
		},
		{
			name: "bad_email_and_password",
			dto: UserSignUpDto{
				Email: "bad",
				Password: "bad",
			},
			want: nil,
			wantErr: true,
		},
		{
			name: "empty_user_name_and_password",
			dto: UserSignUpDto{
				Email: "testemail@gmail.com",
				UserName: "",
				Password: "",
			},
			want: nil,
			wantErr: true,
		},
		{
			name: "short_user_name",
			dto: UserSignUpDto{
				Email: "testemail@gmail.com",
				UserName: "short",
				Password: "123123123123",
			},
			want: nil,
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantErr {
				user, err := tc.dto.IntoUser()
				if err == nil || user != nil {
					t.FailNow()
				}
				
			} else {
				user, err := tc.dto.IntoUser()
				if err != nil || !reflect.DeepEqual(tc.want, user) {
					t.FailNow()
				}
			}
		})
	}
}

