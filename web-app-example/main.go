package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/xxxxxxxxxx/misc/datastore"
	"github.com/xxxxxxxxxx/misc/routes"
	"github.com/xxxxxxxxxx/misc/tmpls"
	"github.com/xxxxxxxxxx/simplegostringtools"
)

const (
	host = "localhost"
	port = 4321
)

type HTMLOutput struct {
	Output []string `json:"output"`
}

type IndexInfo struct {
	HTMLOutput `json:"htmlOutput"`
	Time       time.Time `json"time"`
	Info       string    `json"info"`
	Layout     string    `json"layout"`
}

type CustomersData struct {
	Customers          []*datastore.Customer
	CurrentSearchQuery string
}

// Make files sortable.
type SortableFiles []os.FileInfo

func (s SortableFiles) Len() int {
	return len(s)
}

func (s SortableFiles) Less(i, j int) bool {
	return s[i].Name() < s[j].Name()
}

func (s SortableFiles) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var (
	exitable chan bool
)

var router *mux.Router

func init() {
	os.Setenv("PORT", fmt.Sprintf("%d", port))
	router = mux.NewRouter()
	tmpls.SetRouter(router)
}

func main() {
	// Assets
	//
	// Alternatively, static assets could be served using `FileServer',
	// but this would also serve directory listings for paths that
	// represent directories instead of files:
	//
	// http.Handle("/assets/", http.StripPrefix("/assets/",
	// 	http.FileServer(http.Dir("./assets"))))
	//
	http.HandleFunc(routes.URL("assets"), func(
		w http.ResponseWriter, r *http.Request) {
		fileOrDir, err := os.Open(r.URL.Path[1:])
		defer fileOrDir.Close()
		if err != nil {
			fmt.Printf("File not found: %s", r.URL.Path[1:])
		}
		fileInfo, err := fileOrDir.Stat()
		if err != nil {
			fmt.Printf("File not found: %s", r.URL.Path[1:])
		}
		if fileInfo.IsDir() {
			http.Redirect(w, r, routes.URL("404"), 302)
			return
		}
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	// Regular Routes
	//
	// Instead of using a custom muxer like Gorilla, the standard
	// library way would look like this:
	//
	// http.HandleFunc(routes.URL("home"), handleIndex)
	// http.HandleFunc(routes.URL("redirect"), handleRedirect)
	// http.HandleFunc(routes.URL("404"), handle404)
	// http.HandleFunc(routes.URL("file"), handleFile)
	// http.HandleFunc(routes.URL("exit"), handleExit)
	// http.HandleFunc(routes.URL("customers"), handleCustomers)

	router.StrictSlash(true)

	router.HandleFunc(routes.URL("home"), handle(handleIndex)).Name("home")
	router.HandleFunc(routes.URL("redirect"), handleRedirect)
	router.HandleFunc(routes.URL("404"), handle404)
	router.HandleFunc(routes.URL("file"), handleFile)
	router.HandleFunc(routes.URL("exit"), handleExit)
	router.HandleFunc("/customers", handleCustomersSearch).Methods("POST").Name("customers-search")
	router.HandleFunc(routes.URL("customers"), handle(handleCustomers)).Name("customers")
	router.HandleFunc(routes.URL("customer"), handleCustomer).Name("customer")
	router.HandleFunc("/customer", handleCustomerDelete).Methods("DELETE").Name("customer-delete")
	router.HandleFunc("/customer", handleCustomerDelete).Methods("POST").Queries("delete", "true")
	router.HandleFunc("/customer", handleCustomerPostPartial).Queries("partial", "true").Methods("POST").Name("customer-post-partial")
	router.HandleFunc("/customer", handleCustomerPost).Methods("POST").Name("customer-post")
	router.HandleFunc("/say/{message}/{n:[0-9]+}", handleSay).Name("say")
	router.HandleFunc("/headers", handleHeaders).Name("headers")
	router.HandleFunc("/json", handleJSON).Name("json")
	router.HandleFunc("/edit", handleFileEditDelete).Methods("POST").Queries("delete", "true").Name("edit-file-delete")
	router.HandleFunc("/edit", handleFileEditCreate).Methods("POST").Queries("create", "true").Name("edit-file-create")
	router.HandleFunc("/edit", handleFileEdit).Name("edit-files")
	router.HandleFunc("/edit/{id:.*}", handleFileEditPost).Methods("POST").Name("edit-file-post")
	router.HandleFunc("/edit/{id:.*}", handleFileEdit).Name("edit-file")
	http.Handle("/", router)

	// Open the web browser after a second so the server has some time
	// to start.
	//
	// go func() {
	// 	time.Sleep(100 * time.Millisecond)
	// 	openInBrowser("http://localhost:4455/")
	// }()

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		panic(err)
	}

	// Cleanly close the database connection.
	defer datastore.DB.Close()
}

func runManualTests() {
	fmt.Printf(simplegostringtools.Hline("=", 20))
	fmt.Printf(simplegostringtools.Sep())
	fmt.Println(simplegostringtools.Reverse("!dlroW ,olleH"))
	fmt.Println(simplegostringtools.Frame("This is a string.\n\nOne\nTwo\nThree", "#"))
	fmt.Println(simplegostringtools.Frame("Onetwo\n\nthree", "*"))
	fmt.Println(simplegostringtools.Frame("This is a test.\nHere is another line.\n... yet another one.", "@"))
	fmt.Println(simplegostringtools.Frame(simplegostringtools.Reverse("Hello"), "@"))
	fmt.Println(simplegostringtools.Studly("This is a test string."))
	fmt.Println(simplegostringtools.RandomStrings(10))

	fmt.Printf(simplegostringtools.Sep())

	x := 3
	fmt.Println(x)
	ChangeInt(&x)
	fmt.Println(x)
	ChangeInt(&x)
	fmt.Println(x)
	ChangeInt(&x)
	fmt.Println(x)

	fmt.Printf(simplegostringtools.Sep())

	var y float64 = 3
	fmt.Println("reflect.typeOf(y) = ", reflect.TypeOf(y))
	fmt.Println("reflect.valueOf(y) = ", reflect.ValueOf(y))
}

func ChangeInt(i *int) {
	*i = *i * *i
}

func makeHTTPRequestGET(url string) string {
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func serveFiles() {
	http.Handle(routes.URL("home"), http.FileServer(http.Dir("./")))
	err := http.ListenAndServe("localhost:4455", nil)
	if err != nil {
		panic(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		handleIndexFlashPOST(w, r)
		return
	}
	if r.URL.Path != routes.URL("home") {
		http.NotFound(w, r)
		return
	}
	info := IndexInfo{Time: time.Now().UTC(),
		Info: "This is some info."}

	printToHTML(&info.HTMLOutput, "---------------------------------------------------------------------------")
	printToHTML(&info.HTMLOutput, "Something")
	printToHTML(&info.HTMLOutput, "Something else")
	printToHTML(&info.HTMLOutput, "---------------------------------------------------------------------------")
	printToHTML(&info.HTMLOutput, fmt.Sprintf("Status code 404: %s", http.StatusText(404)))
	printToHTML(&info.HTMLOutput, fmt.Sprintf("Status code 202: %s", http.StatusText(202)))
	printToHTML(&info.HTMLOutput, fmt.Sprintf("Status code 201: %s", http.StatusText(201)))
	printToHTML(&info.HTMLOutput, fmt.Sprintf("Status code 200: %s", http.StatusText(200)))
	printToHTML(&info.HTMLOutput, fmt.Sprintf("Status code 505: %s", http.StatusText(505)))
	printToHTML(&info.HTMLOutput, "---------------------------------------------------------------------------")
	printToHTML(&info.HTMLOutput, r.UserAgent())
	printToHTML(&info.HTMLOutput, "---------------------------------------------------------------------------")
	printToHTML(&info.HTMLOutput, fmt.Sprintf(`Form value ("q"): %v`, r.FormValue("q")))

	tmpls.Render(w, info, "index", "layout")
}

func handle(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request")
		fn(w, r)
	}
	return handler
}

func handleIndexFlashPOST(w http.ResponseWriter, r *http.Request) {
	tmpls.Flash = r.FormValue("flash-message")
	http.Redirect(w, r, routes.URL("home"), 303)
}

func handle404(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Info": "There is no page here. Sorry.",
	}
	tmpls.Render(w, data, "404", "layout")
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, routes.URL("404"), 303)
}

func handleFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./main.go")
}

func handleExit(w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Exit program")
		os.Exit(0)
	}()
	tmpls.Render(w, map[string]string{}, "exit", "layout")
}

func handleCustomers(w http.ResponseWriter, r *http.Request) {
	customers := datastore.QueryCustomers()
	data := CustomersData{customers, ""}
	tmpls.Render(w, data, "customers", "layout")
}

func handleCustomersSearch(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	customers := datastore.QueryCustomersBySearchTerm(query)
	data := CustomersData{customers, query}
	tmpls.Render(w, data, "customers", "layout")
}

func handleCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	customer := datastore.QueryCustomerById(id)
	log.Println(reflect.TypeOf(customer).Elem().Field(0).Tag.Get("name"))
	data := struct {
		Customer *datastore.Customer
	}{Customer: customer}
	tmpls.Render(w, data, "customer", "layout")
}

func handleCustomerPost(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	company := r.FormValue("company")
	address := r.FormValue("address")
	city := r.FormValue("city")
	state := r.FormValue("state")
	country := r.FormValue("country")
	postalCode := r.FormValue("postalCode")
	phone := r.FormValue("phone")
	fax := r.FormValue("fax")
	email := r.FormValue("email")
	err := datastore.UpdateCustomer(id, firstName, lastName, company,
		address, city, state, country, postalCode, phone, fax,
		email)
	if err != nil {
		log.Fatal(err)
		tmpls.Flash = "Could not update user data"
	} else {
		route, err := router.Get("customer").URL("id", id)
		if err != nil {
			log.Fatal(err)
		}
		url := route.Path
		http.Redirect(w, r, url, 303)
	}
}

func handleCustomerPostPartial(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")

	err := datastore.CreateCustomer(firstName, lastName)

	if err != nil {
		log.Fatal(err)
		tmpls.Flash = "Could not create user"
	} else {
		tmpls.Flash = fmt.Sprintf("Successfully created new user: %s %s",
			firstName, lastName)
	}

	route, err := router.Get("customers").URL()
	if err != nil {
		log.Fatal(err)
	}
	url := route.Path
	http.Redirect(w, r, url, 303)
}

func handleCustomerDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	err := datastore.DeleteCustomer(id)
	if err != nil {
		log.Fatal(err)
		tmpls.Flash = "Could not delete user"
	} else {
		tmpls.Flash = "Successfully deleted user"
	}

	route, err := router.Get("customers").URL()
	if err != nil {
		log.Fatal(err)
	}
	url := route.Path
	http.Redirect(w, r, url, 303)
}

func handleHeaders(w http.ResponseWriter, r *http.Request) {
	headers := map[string]string{}
	var s string
	fields := []string{}
	for k, v := range r.Header {
		for _, field := range v {
			fields = append(fields, field)
		}
		s = strings.Join(fields, ", ")
		headers[k] = fmt.Sprintf("%s", s)
		s = ""
		fields = []string{}
	}
	data := struct {
		Headers interface{}
	}{headers}
	tmpls.Render(w, data, "headers", "layout")
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	data := &IndexInfo{
		HTMLOutput: HTMLOutput{[]string{"one", "two", "three"}},
		Time:       time.Now().UTC(),
		Info:       "Here is some info",
		Layout:     "left, center, right",
	}
	j, err := json.Marshal(map[string]interface{}{"data": data})
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(j))
}

func handleFileEdit(w http.ResponseWriter, r *http.Request) {
	// List all files
	files, err := ioutil.ReadDir("./files")
	if err != nil {
		log.Fatal(err)
	}
	// Sort the files (file systems are not guaranteed to return lists
	// of files sorted lexically/by code point.).
	sort.Sort(SortableFiles(files))
	data := struct {
		ID      string
		Files   []interface{}
		Name    string
		Content string
		Lines   int
		Chars   int
	}{ID: string(mux.Vars(r)["id"])}
	var s []byte
	var content string
	for _, f := range files {
		// Generate the MD5 hashes
		file, err := ioutil.ReadFile("./files/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		s = append(file, f.Name()...)
		hash := md5.Sum(s)
		// Set current file name + content.
		if fmt.Sprintf("%x", hash) == mux.Vars(r)["id"] {
			content = fmt.Sprintf("%s", file)
			data.Content = content
			data.Name = f.Name()
			data.Lines = len(strings.Split(content, "\n"))
			data.Chars = len(content)
		}
		if err != nil {
			log.Fatal(err)
		}
		// Gather the file names + MD5 hash sums
		data.Files = append(data.Files, struct {
			Name string
			ID   string
		}{f.Name(), fmt.Sprintf("%x", hash)})
	}
	tmpls.Render(w, data, "layout", "edit-file")
}

func handleFileEditPost(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	content := r.FormValue("content")

	// Find the correnct file based on its hash (file name + file content).
	files, err := ioutil.ReadDir("./files")
	if err != nil {
		log.Fatal(err)
	}
	var file []byte
	for _, f := range files {
		file, err = ioutil.ReadFile("./files/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		s := append(file, f.Name()...)
		hash := md5.Sum(s)
		// Write the file to disk.
		if fmt.Sprintf("%x", hash) == id {
			err := ioutil.WriteFile(fmt.Sprintf("./files/%s", f.Name()), []byte(content), 0644)
			if err != nil {
				log.Fatal(err)
			} else {
				tmpls.Flash = "Successfully saved file."
			}
		}
	}

	url, err := router.Get("edit-files").URL()
	if err != nil {
		log.Fatal(err)
	}
	route := url.Path
	http.Redirect(w, r, route, 303)
}

func handleFileEditCreate(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	// Create the new file.
	err := ioutil.WriteFile("./files/"+name, []byte{}, 0644)
	if err != nil {
		log.Fatal(err)
	} else {
		tmpls.Flash = "Successfully created file."
	}
	hash := fmt.Sprintf("%x", md5.Sum(append([]byte{}, name...)))
	// Redirect to files page.
	url, err := router.Get("edit-file").URL("id", hash)
	if err != nil {
		log.Fatal(err)
	}
	path := url.Path
	http.Redirect(w, r, path, 303)
}

func handleFileEditDelete(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	// Find and delete the file.
	files, err := ioutil.ReadDir("./files")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		file, err := ioutil.ReadFile("./files/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		s := append(file, f.Name()...)
		hash := md5.Sum(s)
		if fmt.Sprintf("%x", hash) == id {
			err := os.Remove("./files/" + f.Name())
			if err != nil {
				log.Fatal(err)
			} else {
				tmpls.Flash = "Successfully deleted file."
			}
		}
	}

	// Redirect to files page.
	url, err := router.Get("edit-files").URL()
	if err != nil {
		log.Fatal(err)
	}
	path := url.Path
	http.Redirect(w, r, path, 303)
}

func handleSay(w http.ResponseWriter, r *http.Request) {
	// route := router.Get("say")

	// url, err := router.Get("say").URL("message", "MESSAGE", "n", "123")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	message := mux.Vars(r)["message"]
	n, err := strconv.Atoi(mux.Vars(r)["n"])
	if err != nil {
		log.Fatal(err)
	}
	s := ""
	for i := 0; i < n; i++ {
		s = fmt.Sprintf("%s%s\n", s, message)
	}

	incremented, err := router.Get("say").URL("message", message, "n",
		fmt.Sprintf("%d", n+1))
	if err != nil {
		log.Fatal(err)
	}
	s = fmt.Sprintf("%s\nIncremented: %s\n", s, incremented)

	fmt.Fprintf(w, "%s", s)
}

// Helpers

func printToHTML(o *HTMLOutput, s string) {
	o.Output = append(o.Output, s)
}

func openInBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
