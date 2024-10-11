package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"
	"unicode"
)

func main() {
	listTemp, tempErr := template.ParseGlob("./templates/*.html")
	if tempErr != nil {
		fmt.Printf("erreur avec le temp : %v", tempErr.Error())
		os.Exit(02)
	}
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

	http.HandleFunc("/promo", func(w http.ResponseWriter, r *http.Request) {
		dataClasse := Classe{"b1", "info", 34, []Eleve{{"romain", "bourdot", 20, "m"}, {"Timothé", "Champieux", 18, "m"}, {"Tomy", "faible", 16, "m"}, {"Tom", "Amaru", 18, "f"}}}
		listTemp.ExecuteTemplate(w, "Promo", dataClasse)

	})

	type Change struct {
		Count    int
		ViewType string
	}
	var count int
	http.HandleFunc("/change", func(w http.ResponseWriter, r *http.Request) {
		count++
		viewType := "impaire"
		if count%2 == 0 {
			viewType = "paire"
		}
		dataChange := Change{count, viewType}
		listTemp.ExecuteTemplate(w, "change", dataChange)
	})

	type User struct {
		Nom           string
		Prenom        string
		DateNaissance string
		Sexe          string
	}
	http.HandleFunc("/user/form", func(w http.ResponseWriter, r *http.Request) {
		listTemp.ExecuteTemplate(w, "form", nil)
	})
	http.HandleFunc("/user/treatment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/user/form", http.StatusSeeOther)
			return
		}

		nom := r.FormValue("nom")
		prenom := r.FormValue("prenom")
		dateNaissance := r.FormValue("dateNaissance")
		sexe := r.FormValue("sexe")

		validName := len(nom) >= 1 && len(nom) <= 32
		for _, char := range nom {
			if !unicode.IsLetter(char) {
				validName = false
			}
		}
		validPrenom := len(prenom) >= 1 && len(prenom) <= 32
		for _, char := range prenom {
			if !unicode.IsLetter(char) {
				validName = false
			}
		}
		validSexe := sexe == "masculin" || sexe == "féminin" || sexe == "autre"

		if !validName || !validPrenom || !validSexe {
			http.Redirect(w, r, "/user/error", http.StatusSeeOther)
			return
		}

		os.Setenv("NOM", nom)
		os.Setenv("PRENOM", prenom)
		os.Setenv("DATE_NAISSANCE", dateNaissance)
		os.Setenv("SEXE", sexe)

		http.Redirect(w, r, "/user/display", http.StatusSeeOther)
	})

	http.HandleFunc("/user/display", func(w http.ResponseWriter, r *http.Request) {
		nom, _ := os.LookupEnv("NOM")
		prenom, _ := os.LookupEnv("PRENOM")
		dateNaissance, _ := os.LookupEnv("DATE_NAISSANCE")
		sexe, _ := os.LookupEnv("SEXE")

		user := User{Nom: nom, Prenom: prenom, DateNaissance: dateNaissance, Sexe: sexe}

		listTemp.ExecuteTemplate(w, "display", user)
	})

	http.HandleFunc("/user/error", func(w http.ResponseWriter, r *http.Request) {
		listTemp.ExecuteTemplate(w, "error", nil)
	})
	http.ListenAndServe("localhost:8080", nil)
}
