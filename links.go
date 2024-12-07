package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

type link struct {
	keyword string
	url     string
	popis   string
}

var links = []link{
	{"g", "https://www.google.com", "Google search engine"},
	{"gh", "https://github.com", "GitHub"},
}

func process_link(keyword string, w http.ResponseWriter, r *http.Request) string {
	for _, l := range links {
		if l.keyword == keyword {
			http.Redirect(w, r, l.url, http.StatusFound)
			return ""
		}
	}

	message := fmt.Sprintf("Key word <b>%s</b> was not found..", keyword)
	return message
}

func proces_cmd_links() string {
	content := `
	<h2>Links</h2>

	<table border cellpadding="5" cellspacing="0">
	  <tr>
	    <th>keyword</th>
	    <th>url</th>
	    <th>popis</th>
	  </tr>
	`

	for _, l := range links {
		content += fmt.Sprintf("<tr>  <td>%s</td> <td><a href='%s'>%s</a></td>  <td>%s</td>   </tr>\n", l.keyword, l.url, l.popis, l.url)
	}

	content += `
	</table>
	`

	return content
}

func generatePageTop() string {
	content := `
	<h1>Links</h1>
		
	<hr>
	<form action="/" method="post">
		<label for="keyword">Link:</label>
		 <input type="text" size="50" id="id_keyword" name="keyword" required autofocus>
		 <input type="submit" value="Submit">
	</form>
	<hr>
	`
	return content
}

func generatePageContent(w http.ResponseWriter, r *http.Request) string {

	keyword := r.FormValue("keyword")
	content := ""

	//
	if keyword == "links" || keyword == "" {
		content = proces_cmd_links()
	}

	if content == "" {
		content = process_link(keyword, w, r)
	}

	return content
}

func generatePageBottom() string {
	return `
		<hr>
		<p style="fond-size:40%;">personal homepage</p>	
	`
}

func handler(w http.ResponseWriter, r *http.Request) {

	pageTop := generatePageTop()
	pageContent := generatePageContent(w, r)
	pageBottom := generatePageBottom()

	content := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Links</title>
	</head>
	<body>
	{{.PageTop}}

	{{.PageContent}}

	{{.PageBottom}}
	</body>
	</html>
	`

	// Vytvoření nové šablony a její naparsování
	tmpl, err := template.New("homepage").Parse(content)
	if err != nil {
		panic(err)
	}

	// Data pro vložení do šablony
	data := struct {
		PageTop     string
		PageContent string
		PageBottom  string
	}{
		PageTop:     pageTop,
		PageContent: pageContent,
		PageBottom:  pageBottom,
	}

	// Vytvoření bufferu pro výsledný string
	var output bytes.Buffer

	// Vykreslení šablony s daty
	err = tmpl.Execute(&output, data)
	if err != nil {
		panic(err)
	}

	//
	fmt.Fprintf(w, output.String())
}

// Read flag values and return the port number and certificate with key..
func parseFlags() (int, string, string) {
	portFlag := flag.Int("port", 0, "Port number for this server")
	pFlag := flag.Int("p", 0, "Short port number argument")
	certFlag := flag.String("cert", "", "Path to certificate")
	keyFlag := flag.String("key", "", "Path to certificate key")

	// Parse the flags..
	flag.Parse()

	port := managePortFlags(portFlag, pFlag)
	certificate, key := manageCertificateFlags(certFlag, keyFlag)

	return port, certificate, key
}

func managePortFlags(portFlag *int, pFlag *int) int {
	// Default port number..
	defaultPort := 8080

	// Read the port number from the environment variable..
	envPortStr := os.Getenv("HTTP_PORT")

	// Convert the port number to integer..
	port := defaultPort
	if envPortStr != "" {
		p, err := strconv.Atoi(envPortStr)
		if err != nil {
			fmt.Println("Conversation error for HTTP_PORT env:", err)
		}

		port = p
	}

	// Prioritize the flags and port number..
	if *portFlag != 0 {
		port = *portFlag
	}
	if *pFlag != 0 {
		port = *pFlag
	}

	return port
}

func manageCertificateFlags(certFlag *string, keyFlag *string) (string, string) {
	// Read the certificate and key files from the environment variables..
	envCert := os.Getenv("HTTPS_CERTIFICATE")
	envKey := os.Getenv("HTTPS_CERTIFICATE_KEY")

	certificate := ""
	key := ""
	if envCert != "" && envKey != "" {
		certificate = envCert
		key = envKey
	}
	if *certFlag != "" && *keyFlag != "" {
		certificate = *certFlag
		key = *keyFlag
	}

	return certificate, key
}

func main() {

	port, certificate, key := parseFlags()

	// Build the socket address..
	addr := fmt.Sprintf(":%v", port)
	fmt.Printf("Server operates on port %d\n", port)

	// Register the handler and start the server..
	http.HandleFunc("/", handler)

	// Start the server with or without TLS..
	if certificate == "" || key == "" {
		log.Fatal(http.ListenAndServe(addr, nil))
		return
	}

	// Start the server with TLS..
	log.Fatal(http.ListenAndServeTLS(addr, certificate, key, nil))
}
