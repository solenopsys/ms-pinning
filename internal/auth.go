package internal

import "errors"

func auth(userKey string, data *Data) (uint64, error) {
	if userKey == "" {
		return 0, errors.New("empty user key")
	}
	var userId uint64
	user, err := data.GetUserById(userKey)
	if err != nil {
		panic(err)
	}
	if user == nil {
		return data.AddUser(userKey)
	} else {
		userId = user.id
	}
	return userId, nil
}
