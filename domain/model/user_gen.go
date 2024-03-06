// Code generated by gen; DO NOT EDIT.
package model

func (user *User) Clone() User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		Name:      user.Name,
		Email:     user.Email,
		Address:   user.Address,
		City:      user.City,
		State:     user.State,
		Zip:       user.Zip,
		BirthDate: user.BirthDate,
		Latitude:  user.Latitude,
		Longitude: user.Longitude,
		Password:  user.Password,
		Source:    user.Source,
		DeletedAt: user.DeletedAt,
	}
}

func (user *User) Field(property string) (string, any) {
	switch property {
	case "id":
		return "id", &user.ID
	case "createdAt":
		return "createdAt", &user.CreatedAt
	case "name":
		return "name", &user.Name
	case "email":
		return "email", &user.Email
	case "address":
		return "address", &user.Address
	case "city":
		return "city", &user.City
	case "state":
		return "state", &user.State
	case "zip":
		return "zip", &user.Zip
	case "birthDate":
		return "birthDate", &user.BirthDate
	case "latitude":
		return "latitude", &user.Latitude
	case "longitude":
		return "longitude", &user.Longitude
	case "password":
		return "password", &user.Password
	case "source":
		return "source", &user.Source
	case "deletedAt":
		return "deletedAt", &user.DeletedAt
	default:
		return "", nil
	}
}

func (user *User) Identity() (string, any) {
	return "id", &user.ID
}

func (user *User) Fields() (string, []any) {
	return `"id", "createdAt", "name", "email", "address", "city", "state", "zip", "birthDate", "latitude", "longitude", "password", "source", "deletedAt"`, []any{&user.ID, &user.CreatedAt, &user.Name, &user.Email, &user.Address, &user.City, &user.State, &user.Zip, &user.BirthDate, &user.Latitude, &user.Longitude, &user.Password, &user.Source, &user.DeletedAt}
}
