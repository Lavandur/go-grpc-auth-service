package users_pb

import (
	"auth-service/internal/grpc/pb"
	"auth-service/internal/models"
	"time"
)

func ToProto(u *models.User) *User {
	return &User{
		Id:    &pb.ID{Id: u.UserID},
		Login: u.Login,
		Person: &User_Person{
			Firstname: u.Person.Firstname,
			Lastname:  u.Person.Lastname,
			Email:     u.Person.Email,
			Gender:    u.Person.Gender,
		},
	}
}

func (f *UserFilter) ToModel() *models.UserFilter {
	if f == nil {
		return nil
	}

	var filter *models.UserFilter
	if f.UserIDs != nil {
		filter.UserID = &f.UserIDs.Value
	}
	if f.Login != nil {
		filter.Login = &f.Login.Value
	}
	if f.Email != nil {
		filter.Email = &f.Email.Value
	}

	return filter
}

func (u *CreateUserRequest) ToModel() *models.UserInput {
	roleIds := make([]string, 0)
	if u.RoleIDs != nil {
		for _, id := range u.RoleIDs.RoleID {
			roleIds = append(roleIds, id.Id)
		}
	}

	return &models.UserInput{
		Login:     u.Login,
		Password:  u.Password,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Birthdate: u.Birthdate.AsTime(),
		Email:     u.Email,
		Gender:    u.Gender,
		RoleIDs:   roleIds,
	}
}

func (u *UpdateUserRequest) ToModel() (string, *models.UserUpdateInput) {
	var birthdate *time.Time = nil
	if u.User.Birthdate != nil {
		t := u.User.Birthdate.AsTime()
		birthdate = &t
	}

	roleIds := make([]string, 0)
	if u.User.RoleIDs != nil {
		for _, id := range u.User.RoleIDs.GetRoleID() {
			roleIds = append(roleIds, id.GetId())
		}
	}

	return u.Id.GetId(), &models.UserUpdateInput{
		Login:     u.User.Login,
		Password:  u.User.Password,
		Firstname: u.User.Firstname,
		Lastname:  u.User.Lastname,
		Birthdate: birthdate, //TODO
		Email:     u.User.Email,
		Gender:    u.User.Gender,
		RoleIDs:   &roleIds,
	}
}
