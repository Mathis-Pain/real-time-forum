package users

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var tmpl = template.Must(template.ParseFiles("./frontend/index.html"))

func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Afficher la page HTML uniquement pour GET
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Error loading page", http.StatusInternalServerError)
		}
		return
	}
	// Ouvrir la base de données
	db, err := sql.Open("sqlite3", "vault/real_time_forum_database.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database connection error"))
		return
	}
	defer db.Close()

	// Vérifier la connexion à la base de données
	if err = db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database connection failed"))
		return
	}

	// Récupérer et nettoyer les données du formulaire
	nickname := strings.TrimSpace(r.FormValue("Nickname"))
	age := strings.TrimSpace(r.FormValue("Age"))
	registerMail := strings.TrimSpace(r.FormValue("Email"))
	gender := r.FormValue("Gender")
	lastName := strings.TrimSpace(r.FormValue("LastName"))
	firstName := strings.TrimSpace(r.FormValue("FirstName"))
	password := r.FormValue("Password")
	confirmPassword := r.FormValue("Confirm-Password")

	// Validation: Champs vides
	if nickname == "" || age == "" || registerMail == "" ||
		lastName == "" || firstName == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("All fields are required"))
		return
	}

	// Validation: Genre par défaut
	if gender == "Default" || gender == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please select a valid gender"))
		return
	}

	// Validation: Mots de passe correspondent
	if password != confirmPassword {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Passwords do not match"))
		return
	}

	// Validation: Longueur du mot de passe
	if len(password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Password must be at least 8 characters"))
		return
	}

	// Vérifier si l'utilisateur existe déjà
	available, err := IsAvailableUser(db, nickname, registerMail)
	if err != nil {
		fmt.Printf("Error checking availability: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error checking user availability"))
		return
	}

	if !available {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email or username already in use"))
		return
	}

	// Hacher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error processing password"))
		return
	}

	// Insérer l'utilisateur dans la base de données
	query := `INSERT INTO users(UserName, Age, Email, Gender, LastName, FirstName, Password) 
          VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = db.Exec(query, nickname, age, registerMail, gender, lastName, firstName, string(hashedPassword))
	if err != nil {
		fmt.Printf("Error inserting user: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating account"))
		return
	}

	// Succès
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Account created successfully!"))
}
