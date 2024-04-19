package Handlers

import (
	"MongoDb/internal/data"
	"MongoDb/pkg/logging"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	if r.Method == http.MethodPost {

		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" || name == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		var count int64
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		count, err := data.Collection.CountDocuments(ctx, bson.M{"email": email})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if count > 0 {
			http.Error(w, "Email already in use", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
		_, err = data.Collection.InsertOne(ctx, bson.M{"email": email, "name": name, "password": hashedPassword})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Infof("Create a new user %s", email)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/register.html"))
	tmpl.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	if r.Method == http.MethodPost {

		//name := r.FormValue("name")
		email := r.FormValue("email")
		data.Gmail = email
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		var result data.User
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := data.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
		if err != nil {
			http.Error(w, "Invalid email", http.StatusUnauthorized)
			return
		}

		/*err = data.Collection.FindOne(ctx, bson.M{"name": name}).Decode(&result)
		if err != nil {
			http.Error(w, "Invalid name", http.StatusUnauthorized)
			return
		}*/

		err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session",
			Value:   email,
			Expires: time.Now().Add(24 * time.Hour),
		})

		http.Redirect(w, r, "/home", http.StatusSeeOther)
		logger.Infof("%s log in", email)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/loginForm.html"))
	tmpl.Execute(w, nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	logger.Infof("%s do log out", data.Gmail)
	data.Gmail = ""
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/home.html"))
	tmpl.Execute(w, data.GetUser(w))
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		var count int64
		count, err = data.Collection.CountDocuments(ctx, bson.M{"email": cookie.Value})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if count == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}
}
