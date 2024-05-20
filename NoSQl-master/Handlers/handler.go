package Handlers

import (
	"MongoDb/internal/data"
	"MongoDb/pkg/emailVerification"
	"MongoDb/pkg/logging"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"time"
)

type loginInput struct {
	email    string
	password string
}

func Register(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	if r.Method == http.MethodPost {

		name := r.FormValue("name")
		surname := r.FormValue("surname")
		dobString := r.FormValue("dob") // Retrieve dob as string
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm-password")

		if name == "" || surname == "" || dobString == "" || email == "" || password == "" || confirmPassword == "" {
			http.Error(w, "Not all fields are filled!", http.StatusBadRequest)
			return
		}

		dob, err := time.Parse("2006-01-02", dobString) // Parse dob string to time.Time
		if err != nil {
			http.Error(w, "Invalid date of birth format", http.StatusBadRequest)
			return
		}

		verifiedEmail, err := emailVerification.IsVerifiedEmail(email)
		logger.Infof("verified email: %v", verifiedEmail)
		if err != nil || !verifiedEmail {
			http.Error(w, "Invalid email", http.StatusBadRequest)
			return
		}

		/*var count int64
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		count, err = data.Collection.CountDocuments(ctx, bson.M{"emailVerification": emailVerification})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if count > 0 {
			http.Error(w, "Email already in use", http.StatusBadRequest)
			return
		}*/ //oldMethod
		var recordUser data.User
		recordUser, err = data.GetUser(email)
		if err != nil && err.Error() != "mongo: no documents in result" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if recordUser.ID != primitive.NilObjectID {
			http.Error(w, "Email already in use", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := emailVerification.GenerateToken()
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		recordUser = data.User{
			ID: primitive.NewObjectID(),
			UserInfo: data.UserInfo{
				Name:         name,
				Surname:      surname,
				Dob:          dob,
				Email:        email,
				PasswordHash: hashedPassword,
			},
			VerificationToken: token,
			Verified:          false,
		}

		err = data.CreateUser(recordUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		subject := "Verify your email address" //change domain in body!! !! ! !! !
		body := fmt.Sprintf("Please click the following link to verify your email address: http://localhost:8080/verify?token=%s", token)
		err = emailVerification.SendEmail(email, subject, body)
		if err != nil {
			http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
			return
		}

		logger.Infof("USER WAS CREATED: %s", recordUser)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/registrationForm.html"))
	tmpl.Execute(w, nil)
}

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Invalid verification link", http.StatusBadRequest)
		return
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)

	usersCollection := client.Database("test").Collection("users")
	filter := bson.M{"verification_token": token}
	update := bson.M{"$set": bson.M{"verified": true, "verification_token": ""}}

	result, err := usersCollection.UpdateOne(ctx, filter, update) //TODO make ModifyUser func instead
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		http.Error(w, "Invalid verification token", http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "Email verified successfully")
}

func RecoverPassword(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	token := r.URL.Query().Get("token")
	recovery := r.URL.Query().Get("recovery")
	email := r.FormValue("email")

	if token == "" && recovery == "" {
		tmpl := template.Must(template.ParseFiles("html/passwordRecovery.html"))
		err := tmpl.Execute(w, map[string]interface{}{
			"Token": token,
		})
		if err != nil {
			http.Error(w, "Response Writer Error!", http.StatusInternalServerError)
			return
		}
		return
	}

	if recovery == "linkSent" {
		token, err := emailVerification.GenerateToken()
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		var user data.User
		user, err = data.GetUser(email)
		if user.ID == primitive.NilObjectID {
			logger.Infof("No user with email: %v", email)
			return
		}

		if user.Verified != true {
			logger.Infof("Unverified user %s", email)
			return
		}

		filter := bson.M{"_id": user.ID}
		update := bson.M{"$set": bson.M{"verification_token": token}}

		_, err = data.UpdateUser(filter, update)
		if err != nil {
			logger.Errorf("Error updating user: %v", err)
		}

		subject := "Password recovery"
		body := fmt.Sprintf("Use the following link to create a new password for your account: http://localhost:8080/recoverPassword?recovery=newPassword&token=%s", token)
		err = emailVerification.SendEmail(email, subject, body)
		if err != nil {
			http.Error(w, "Could not send an email. Try later.", http.StatusInternalServerError)
			return
		}
		action := "/shop"
		message := "If your email is correct, we will send a recovery link to it."
		_ = showMessage(action, message, w)
		return
	} else if recovery == "newPassword" {
		tmpl := template.Must(template.ParseFiles("html/passwordRecovery.html"))
		err := tmpl.Execute(w, map[string]interface{}{
			"Token": token,
		})
		if err != nil {
			http.Error(w, "Response Writer Error!", http.StatusInternalServerError)
			return
		}
		return
	}

	newPassword := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if newPassword != confirmPassword {
		action := fmt.Sprintf("/recoverPassword?recovery=newPassword&token=%s", token)
		message := "Password mismatch. Please try again."
		err := showMessage(action, message, w)
		if err != nil {
			http.Error(w, message, http.StatusInternalServerError)
		}
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filter := bson.M{"verification_token": token}
	update := bson.M{"$set": bson.M{"user_info.password": hashedPassword, "verification_token": ""}}

	result, err := data.UpdateUser(filter, update)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		action := "/shop"
		message := "Your password recovery token is not valid anymore."
		err = showMessage(action, message, w)
		return
	}

	action := "/shop"
	message := "Password recovery finished successfully!"
	err = showMessage(action, message, w)
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	if r.Method == http.MethodPost {

		loginData := loginInput{
			email:    r.FormValue("email"),
			password: r.FormValue("password"),
		}

		if loginData.email == "" || loginData.password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		var result data.User

		result, err := data.GetUser(loginData.email)
		if err != nil {
			http.Error(w, "Invalid email", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword(result.UserInfo.PasswordHash, []byte(loginData.password))
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session",
			Value:   result.UserInfo.Email,
			Expires: time.Now().Add(24 * time.Hour),
		})

		http.Redirect(w, r, "/shop", http.StatusSeeOther)

		err = data.SetUser(result)
		if err != nil {
			http.Error(w, "Empty user struct!", http.StatusNotFound)
			return
		}
		logger.Infof("%s LOGGED IN", result.UserInfo.Email)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/loginForm.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Response Writer Error!", http.StatusInternalServerError)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	logger.Infof("%s LOGGED OUT", data.ShowUser())
	data.ClearUser()
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/home.html"))
	tmpl.Execute(w, data.ShowUser())
}

func EditUserInfoForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/editUserInfo.html"))
	tmpl.Execute(w, data.ShowUser())
}

func EditUserInfo(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	err := data.Init("test", "users")
	if err != nil {
		http.Redirect(w, r, "/shop", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		dobString := r.FormValue("dob") // Retrieve dob as string

		if name == "" || surname == "" || dobString == "" {
			http.Error(w, "Not all fields are filled!", http.StatusBadRequest)
			return
		}

		dob, err := time.Parse("2006-01-02", dobString) // Parse dob string to time.Time
		if err != nil {
			http.Error(w, "Invalid date of birth format", http.StatusBadRequest)
			return
		}

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("editUserInfoButton")[10:34])
		filter := bson.M{"_id": ObjID}

		update := bson.M{"$set": bson.M{
			"name":    name,
			"surname": surname,
			"dob":     dob,
		}}

		_, err = data.Collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			logger.Infof("A bulk write error occurred: %v", err)
			return
		} else {
			logger.Infof("User with ID: %s was UPDATED!", ObjID)
		}

		changedUser, _ := data.GetUser(data.ShowUser().UserInfo.Email)
		_ = data.SetUser(changedUser)

		http.Redirect(w, r, "/showUserProfile", http.StatusSeeOther)
	}

}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLogger()

		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if data.ShowUser().UserInfo.Email != cookie.Value {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			http.SetCookie(w, &http.Cookie{
				Name:    "session",
				Value:   "",
				Expires: time.Unix(0, 0),
			})
			logger.Infof("Cookie expired or not got another value.")
			return
		}

		next.ServeHTTP(w, r)
	}
}

func showMessage(action string, message string, w http.ResponseWriter) error {
	tmpl := template.Must(template.ParseFiles("html/message.html"))
	err := tmpl.Execute(w, map[string]interface{}{
		"Action":  action,
		"Message": message,
	})
	if err != nil {
		http.Error(w, "Response Writer Error!", http.StatusInternalServerError)
		return err
	}
	return nil
}

func ShowProfile(w http.ResponseWriter, r *http.Request) {
	//logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/userProfile.html"))
	tmpl.Execute(w, data.ShowUser())
}
