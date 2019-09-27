package user

var users = Users{
	{5584596, "DmitryKozhurkin"},
	{3855715, "kumirska"},
}

func GetById(id int) *User {
	for _, item := range users {
		if item.GithubId == id {
			return item
		}
	}
	return nil
}
