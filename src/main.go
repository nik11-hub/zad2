package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	// HEALTHCHECK
	if len(os.Args) > 1 && os.Args[1] == "-health" {

		// Prosty self-check aplikacji - sprawdzenie, czy serwer odpowiada na żądania.
		resp, err := http.Get("http://127.0.0.1:8080/")
		if err != nil || resp.StatusCode >= 400 {
			os.Exit(1)
		}

		os.Exit(0)
	}

	port := "8080"

	log.Printf("Data uruchomienia: %s", time.Now().Format(time.RFC3339))
	log.Println("Autor: Mikita Liaiko")
	log.Printf("Aplikacja nasłuchuje na porcie TCP: %s", port)

	// routowanie endpointow
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/weather", weatherHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	// Statyczny HTML z formularzem wyboru miasta.
	html := `
<!DOCTYPE html>
<html lang="pl">

<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">

	<title>Weather App</title>

	<style>

		body {
			margin: 0;
			padding: 0;
			font-family: Arial, sans-serif;
			background: linear-gradient(135deg, #4facfe, #00f2fe);
			height: 100vh;
			display: flex;
			justify-content: center;
			align-items: center;
		}

		.container {
			background: white;
			padding: 40px;
			border-radius: 20px;
			box-shadow: 0 10px 30px rgba(0,0,0,0.2);
			width: 350px;
			text-align: center;
		}

		h1 {
			margin-bottom: 25px;
			color: #333;
		}

		select {
			width: 100%;
			padding: 12px;
			border-radius: 10px;
			border: 1px solid #ccc;
			font-size: 16px;
			margin-bottom: 20px;
		}

		button {
			width: 100%;
			padding: 12px;
			border: none;
			border-radius: 10px;
			background: #4facfe;
			color: white;
			font-size: 16px;
			cursor: pointer;
			transition: 0.2s;
		}

		button:hover {
			background: #2196f3;
		}

		.footer {
			margin-top: 20px;
			font-size: 12px;
			color: #666;
		}

	</style>
</head>

<body>

	<div class="container">

		<h1>🌤 Weather App</h1>

		<form action="/weather" method="get">

			<select name="city">

				<option value="Lublin">
					Polska — Lublin
				</option>

				<option value="Wroclaw">
					Polska — Wrocław
				</option>

				<option value="Minsk">
					Białoruś — Mińsk
				</option>

				<option value="New York">
					USA — New York
				</option>

			</select>

			<button type="submit">
				Sprawdź pogodę
			</button>

		</form>


	</div>

</body>
</html>
`

	fmt.Fprint(w, html)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {

	// Odczyt miasta z query stringa i przekazanie krotkiego podsumowania pogody.
	city := r.URL.Query().Get("city")

	if city == "" {
		http.Error(w, "Nie wybrano miasta", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	url := fmt.Sprintf("https://wttr.in/%s?format=3", city)

	// Wywolanie zewnetrznej uslugi i wyrenderowanie odpowiedzi jako HTML.
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Błąd pobierania pogody", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	// Usunięcie powielonej nazwy miasta z odpowiedzi API
	weatherData := string(body)
	if parts := strings.SplitN(weatherData, ": ", 2); len(parts) == 2 {
		weatherData = parts[1]
	}

	fmt.Fprintf(w, `
<!DOCTYPE html>
<html lang="pl">

<head>
	<meta charset="UTF-8">

	<title>Pogoda</title>

	<style>

		body {
			margin: 0;
			padding: 0;
			font-family: Arial, sans-serif;
			background: linear-gradient(135deg, #4facfe, #00f2fe);
			height: 100vh;
			display: flex;
			justify-content: center;
			align-items: center;
		}

		.card {
			background: white;
			padding: 40px;
			border-radius: 20px;
			box-shadow: 0 10px 30px rgba(0,0,0,0.2);
			width: 350px;
			text-align: center;
		}

		h1 {
			color: #333;
		}

		p {
			font-size: 32px;
			color: #444;
			margin: 20px 0;
		}

		a {
			display: inline-block;
			padding: 10px 20px;
			background: #4facfe;
			color: white;
			text-decoration: none;
			border-radius: 10px;
		}

		a:hover {
			background: #2196f3;
		}

	</style>
</head>

<body>

	<div class="card">

		<h1>📍 %s</h1>

		<p>%s</p>

		<a href="/">
			Powrót
		</a>

	</div>

</body>
</html>
`, city, weatherData)
}
