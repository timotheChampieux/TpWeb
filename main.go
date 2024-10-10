package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"
)

func main() {
	listTemp, tempErr := template.ParseGlob("./templates/*.html")
	if tempErr != nil {
		fmt.Printf("erreur avec le temp : %v", tempErr.Error())
		os.Exit(02)
	}

	// PageConditionSimple struct {
	//	Check bool
	//}

	type Eleve struct {
		Nom    string
		Prenom string
		Age    int
		Sexe   string
	}
	type Classe struct {
		NomClasse string
		Filiere   string
		NbrEt     int
		Etu       []Eleve
	}

	fileServer := http.FileServer(http.Dir("./assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/ynov", func(w http.ResponseWriter, r *http.Request) {
		//dataPage := PageConditionSimple{true}
		dataClasse := Classe{"b1", "info", 34, []Eleve{{"romain", "bourdot", 20, "m"}, {"Timothé", "Champieux", 18, "m"}, {"Tomy", "grospd", 19, "m"}, {"Tom", "Amaru", 18, "f"}}}

		//listTemp.ExecuteTemplate(w, "condition", dataPage)
		listTemp.ExecuteTemplate(w, "Classe", dataClasse)
	})

	http.ListenAndServe("localhost:8080", nil)
}
