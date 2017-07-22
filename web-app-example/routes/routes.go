package routes

var Routes map[string]string

func init() {
	Routes = map[string]string{
		"home":      "/",
		"404":       "/404/",
		"assets":    "/assets/",
		"file":      "/file/",
		"redirect":  "/redirect/",
		"exit":      "/exit/",
		"customers": "/customers/",
		"customer":  "/customer/{id:[0-9]+}",
	}
}

func URL(name string) string {
	return Routes[name]
}
