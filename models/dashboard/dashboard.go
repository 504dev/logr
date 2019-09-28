package dashboard

var dashboards = Dashboards{
	{
		Id:         1,
		OwnerId:    1,
		Name:       "sandbox",
		PublicKey:  "MCIwDQYJKoZIhvcNAQEBBQADEQAwDgIHLsJxAFHuswIDAQAB",
		PrivateKey: "MDgCAQACBy7CcQBR7rMCAwEAAQIHKX5KFEXQAQIEBplUswIEBxX+AQIEBlkywwIEAkUCAQIEAOzwlA==",
	},
}

func GetById(id int) *Dashboard {
	for _, item := range dashboards {
		if item.Id == id {
			return item
		}
	}
	return nil
}

func GetByPub(pub string) *Dashboard {
	for _, item := range dashboards {
		if item.PublicKey == pub {
			return item
		}
	}
	return nil
}
