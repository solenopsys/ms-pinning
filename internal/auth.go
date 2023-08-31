package internal

func auth(userKey string, data *Data) (int64, error) {
	var userId int64
	user := data.GetUserById(userKey)
	if user == nil {
		return data.AddUser(userKey)
	} else {
		userId = user.id
	}
	return userId, nil
}
