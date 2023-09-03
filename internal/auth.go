package internal

func auth(userKey string, data *Data) (uint64, error) {
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
