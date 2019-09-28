package dashboard

var dashboards = Dashboards{
	{
		Id:         1,
		OwnerId:    1,
		Name:       "sandbox",
		PublicKey:  "MDwwDQYJKoZIhvcNAQEBBQADKwAwKAIhAJ/dD/71Ak7I0e0Q9R9M+yQr45dk3FxbSLT6PQ3vPwLLAgMBAAE=",
		PrivateKey: "MIGpAgEAAiEAn90P/vUCTsjR7RD1H0z7JCvjl2TcXFtItPo9De8/AssCAwEAAQIgO8gzfyiooEXBG2JICFiFYb4dArQbN+TJgkJdqTHthcECEQDFJHDL3Px0ZlPNfgt6bBYJAhEAz5doRIkWyGkNx1dk5lFnMwIQLEeDDkcUbRCOwhkNevHMAQIQGE7IXM4Yptr2TBNFs1pw1QIQYjC4bb5B1l4cGi2zHH8a2Q==",
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
