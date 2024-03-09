package server

func MustListenPROM() {
	if err := ListenPROM(); err != nil {
		panic(err)
	}
}

func ListenPROM() error {
	//address := config.Get().Bind.Prom
	//if address == "" {
	//	return nil
	//}
	////use separated ServeMux to prevent handling on the global Mux
	//mux := http.NewServeMux()
	//mux.Handle("/metrics", promhttp.Handler())
	//
	//return http.ListenAndServe(address, mux)
	return nil
}
