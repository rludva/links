package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type link struct {
	keyword   string
	popis     string
	url       string
	kategorie string
}

var links = []link{
	{"g", "Google search engine", "https://www.google.com", "general"},
	{"f", "Facebook", "https://www.facebook.com/", "general"},
	{"y", "YouTube", "https://www.youtube.com/", "general"},
	{"gh", "GitHub", "https://github.com", "development"},
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

	<table border cellpadding="5" cellspacing="0" style="width:80%; border-collapse: collapse;">
	  <tr>
	    <th style="width: 20%;">keyword</th>
	    <th style="width: 50%;">link</th>
			<th style="width: 25%;">kategorie</th>
	  </tr>
	`

	for _, l := range links {
		keyword := strings.ReplaceAll(l.keyword, "%", "%%")
		url := strings.ReplaceAll(l.url, "%", "%%")
		popis := strings.ReplaceAll(l.popis, "%", "%%")
		kategorie := strings.ReplaceAll(l.kategorie, "%", "%%")
		content += fmt.Sprintf("<tr>  <td>%s</td> <td><a href='%s'>%s</a><br><div style=\"color: gray; font-size:8px;\">%s</div></td>  <td>%s</td> </tr>\n", keyword, url, popis, url, kategorie)
	}

	content += `
	</table>
	`

	return content
}

func process_cmd_add() string {
	// Generatr web form to add a new link..
	content := `
	<h2>Add a new link</h2>
	<form action="/add" method="post">
	<label>Keyword:</label>
	<input type="text" name="keyword" required>
	<br>
	<label>Popis:</label>
	<input type="text" name="popis" required>
	<br>
	<label>URL:</label>
	<input type="text" name="url" required>
	<br>
	<label>Kategorie:</label>
	<input type="text" name="kategorie" required>
	<br>
	<input type="submit" value="Submit">
	</form>
	`

	return content
}

func generatePageTop() string {
	content := `
	<h1>Personal homepage</h1>
		
	<hr>
	<form name="links" action="/" method="post">
		<label for="keyword">Odkaz:</label>
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

	//
	if keyword == "add" {
		content = process_cmd_add()
	}

	if content == "" {
		content = process_link(keyword, w, r)
	}

	return content
}

func generatePageBottom() string {
	return `
		<hr>
	<!-- <button class="import-button">Import</button> -->
	<button id="addButton" class="import-button">Add</button>

<script>
  document.getElementById("addButton").addEventListener("click", function () {
    document.getElementById("id_keyword").value = "add";
    document.getElementById("links").submit();
  });
</script>	
	`
}

func handler_csv(w http.ResponseWriter, r *http.Request) {
	// Create a CSV writer that writes to stdout
	writer := csv.NewWriter(os.Stdout)
	writer.Comma = ';'
	defer writer.Flush()

	// Write the header
	err := writer.Write([]string{"keyword", "popis", "url", "kategorie"})
	if err != nil {
		return
	}

	// Write each link as a row in the CSV
	for _, l := range links {
		err = writer.Write([]string{l.keyword, l.popis, l.url, l.kategorie})
		if err != nil {
			return
		}
	}

	return
}

func import_csv(datafile string) {

	// If the datafile is empty, use the default value..
	if datafile == "" {
		datafile = "links.csv"
	}

	fmt.Println("Datafile:", datafile)

	// Open the CSV file
	file, err := os.Open(datafile)
	if err != nil {
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)
	reader.Comma = ';'

	// Read all the records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("ReadAll Error %v", err)
		return
	}
	fmt.Println("ReadAll().. done")

	// Create a new slice to store the links
	links = make([]link, 0, len(records))

	// Loop over the records
	for _, record := range records {
		// Create a new link
		link := link{
			keyword:   record[0],
			popis:     record[1],
			url:       record[2],
			kategorie: record[3],
		}

		// Append the link to the slice
		links = append(links, link)
		fmt.Println("links appended")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	pageTop := generatePageTop()
	pageContent := generatePageContent(w, r)
	pageBottom := generatePageBottom()

	content := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Personal homepage</title>
		<style type="text/css">
		body {font-family: Tahoma, Arial, sans-serif;}
		table {
			border-collapse: collapse;
			font-size: 10px;
		}
		td, th {border: 1px solid #999; padding: 8px;}
		tr:hover {background-color: #F5FFFA; }
		table tr th {background-color: #ffffd7;}
		
		.import-button {
      display: inline-block;
      padding: 10px 20px;
      font-size: 16px;
      color: #fff;
      background-color: #007BFF;
      border: none;
      border-radius: 5px;
      cursor: pointer;
      text-decoration: none;
    }

    .import-button:hover {
      background-color: #0056b3;
    }

		</style>
	</head>
	<body>
	{{.PageTop}}

	{{.PageContent}}

	{{.PageBottom}}

  <script>
      // Detection of pageshow..
      window.addEventListener('pageshow', (event) => {
          const inputField = document.getElementById('id_keyword');
          if (inputField) {
              inputField.focus(); 
          }
      });
  </script>

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
func parseFlags() (string, int, string, string) {
	datafileFlag := flag.String("datafile", "links.csv", "Path to the CSV file")
	portFlag := flag.Int("port", 0, "Port number for this server")
	pFlag := flag.Int("p", 0, "Short port number argument")
	certFlag := flag.String("cert", "", "Path to certificate")
	keyFlag := flag.String("key", "", "Path to certificate key")

	// Parse the flags..
	flag.Parse()

	datafile := manageDatafileFlag(datafileFlag)
	port := managePortFlags(portFlag, pFlag)
	certificate, key := manageCertificateFlags(certFlag, keyFlag)

	return datafile, port, certificate, key
}

func manageDatafileFlag(datafileFlag *string) string {
	// Read the certificate and key files from the environment variables..
	envDatafile := os.Getenv("DATAFILE")

	datafile := ""
	if envDatafile != "" {
		datafile = envDatafile
	}
	if *datafileFlag != "" {
		datafile = *datafileFlag
	}

	return datafile
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
			fmt.Println("Chyba při konverzi:", err)
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

	datafile, port, certificate, key := parseFlags()

	// Read data from the CSV file..
	import_csv(datafile)

	// Build the socket address..
	addr := fmt.Sprintf(":%v", port)
	fmt.Printf("Server operates on portu %d\n", port)

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
