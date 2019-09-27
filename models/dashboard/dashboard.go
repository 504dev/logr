package dashboard

var dashboards = Dashboards{
	{
		Id:         1,
		Name:       "vbm",
		PublicKey:  "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1DSXdEUVlKS29aSWh2Y05BUUVCQlFBREVRQXdEZ0lITHNKeEFGSHVzd0lEQVFBQgotLS0tLUVORCBSU0EgUFVCTElDIEtFWS0tLS0tCg==",
		PrivateKey: "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNRGdDQVFBQ0J5N0NjUUJSN3JNQ0F3RUFBUUlIS1g1S0ZFWFFBUUlFQnBsVXN3SUVCeFgrQVFJRUJsa3l3d0lFCkFrVUNBUUlFQU96d2xBPT0KLS0tLS1FTkQgUlNBIFBSSVZBVEUgS0VZLS0tLS0K"},
}

func GetById(id int) *Dashboard {
	for _, item := range dashboards {
		if item.Id == id {
			return item
		}
	}
	return nil
}

func GetByPk(pubk string) *Dashboard {
	// TODO find by public key
	return dashboards[0]
}
