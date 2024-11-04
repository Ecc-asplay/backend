package db

import "testing"

func CreateAdminUser() CreateAdminUserParams {
	Data := CreateAdminUserParams{
		Email:        "",
		Hashpassword: "",
		StaffName:    "",
		Department:   "",
	}
	return Data
}

func TestCreateAdminUser(t *testing.T) {

}
