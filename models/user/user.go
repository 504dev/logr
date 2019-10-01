package user

var users = Users{
	{1, 5584596, "DmitryKozhurkin"},
	{2, 3855715, "kumirska"},
}

func GetAll() Users {
	return users
}

func GetById(id int) *User {
	for _, item := range users {
		if item.GithubId == id {
			return item
		}
	}
	return nil
}
