package Handlers

import (
	"MongoDb/internal/data"
	"MongoDb/pkg/emailVerification"
	"MongoDb/pkg/logging"
	"MongoDb/pkg/session"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		dobString := r.FormValue("dob") //string
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm-password")

		if name == "" || surname == "" || dobString == "" || email == "" || password == "" || confirmPassword == "" {
			_ = showMessage("/shop", "Not all fields are filled!", w)
			return
		}

		dob, err := time.Parse("2006-01-02", dobString) //Parse dob string to time.Time
		if err != nil {
			HandleError(errors.New("invalid date of birth format"), logger, w)
			return
		}

		verifiedEmail, err := emailVerification.IsVerifiedEmail(email)
		logger.Infof("verified email: %v", verifiedEmail)
		if err != nil || !verifiedEmail {
			HandleError(errors.New("invalid email"), logger, w)
			return
		}

		var recordUser data.User
		recordUser, err = data.GetUser(email)
		if err != nil && err.Error() != "mongo: no documents in result" {
			HandleError(errors.New("could not check user email for duplication"), logger, w)
			return
		}

		if recordUser.ID != primitive.NilObjectID {
			HandleError(errors.New("email already in use"), logger, w)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			HandleError(err, logger, w)
			return
		}

		token, err := emailVerification.GenerateToken()
		if err != nil {
			HandleError(err, logger, w)
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
			SessionToken:      "",
			VerificationToken: token,
			Verified:          false,
		}

		err = data.CreateUser(recordUser)
		if err != nil {
			HandleError(err, logger, w)
		}

		subject := "Verify your email address" //change domain in body!! !! ! !! !
		body := fmt.Sprintf("Please click the following link to verify your email address: http://localhost:8080/verify?token=%s", token)
		err = emailVerification.SendEmail(email, subject, body)
		if err != nil {
			HandleError(errors.New("failed to send verification email"), logger, w)
			return
		}

		logger.Infof("USER WAS CREATED: %s", recordUser)
		messageText := "Your account has been created! Please, check your email to verify your account!"
		_ = showMessage("/shop", messageText, w)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/registrationForm.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		_ = showMessage("/shop", err.Error(), w)
		return
	}
}

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	var err error

	token := r.URL.Query().Get("token")
	if token == "" {
		HandleError(errors.New("invalid verification link"), logger, w)
		return
	}

	filter := bson.M{"verification_token": token}
	update := bson.M{"$set": bson.M{"verified": true, "verification_token": ""}}

	result, err := data.UpdateUser(filter, update)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if result.MatchedCount == 0 {
		HandleError(errors.New("invalid verification token"), logger, w)
		return
	}
	_ = showMessage("/shop", "Email verified successfully!", w)
}

func RecoverPassword(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	token := r.URL.Query().Get("token")
	recovery := r.URL.Query().Get("recovery")
	email := r.FormValue("email")
	var err error

	if token == "" && recovery == "" {
		tmpl := template.Must(template.ParseFiles("html/passwordRecovery.html"))
		err = tmpl.Execute(w, map[string]interface{}{
			"Token": token,
		})
		if err != nil {
			http.Error(w, "Response Writer Error!", http.StatusInternalServerError)
			return
		}
		return
	}

	if recovery == "linkSent" {
		token, err = emailVerification.GenerateToken()
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

		var user data.User

		user, err := data.GetUser(loginData.email)
		if err != nil {
			HandleError(errors.New("invalid email"), logger, w)
			return
		}

		err = bcrypt.CompareHashAndPassword(user.UserInfo.PasswordHash, []byte(loginData.password))
		if err != nil {
			HandleError(errors.New("invalid password"), logger, w)
			return
		}

		sessionToken, err := session.GenerateSessionToken()
		if err != nil {
			HandleError(errors.New("internal server error"), logger, w)
			return
		}

		filter := bson.M{"_id": user.ID}
		update := bson.M{"$set": bson.M{"session_token": sessionToken}}
		_, err = data.UpdateUser(filter, update)
		if err != nil {
			HandleError(errors.New("internal server error"), logger, w)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session",
			Value:   sessionToken,
			Expires: time.Now().Add(24 * time.Hour),
		})

		http.Redirect(w, r, "/shop", http.StatusSeeOther)

		logger.Infof("%s LOGGED IN", user.UserInfo.Email)
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
	logger.Infof("%v LOGGED OUT", data.ShowUser(r).UserInfo.Email)

	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Value:   "",
		Expires: time.Unix(0, 0),
	})

	filter := bson.M{"session_token": session.GetSessionTokenFromCookie(r)}
	update := bson.M{"$set": bson.M{"session_token": ""}}
	_, err := data.UpdateUser(filter, update)
	if err != nil {
		HandleError(errors.New("internal server error"), logger, w)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/home.html"))
	err := tmpl.Execute(w, data.ShowUser(r))
	if err != nil {
		HandleError(err, logging.GetLogger(), w)
		return
	}
}

func EditUserInfoForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/editUserInfo.html"))

	user, _ := data.GetUserBySessionToken(session.GetSessionTokenFromCookie(r))

	err := tmpl.Execute(w, user)
	if err != nil {
		HandleError(err, logging.GetLogger(), w)
		return
	}
}

func EditUserInfo(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	if data.UsersCollection == nil {
		HandleError(errors.New("empty collection"), logger, w)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		dobString := r.FormValue("dob") //Retrieve dob as string

		if name == "" || surname == "" || dobString == "" {
			HandleError(errors.New("not all fields are filled"), logger, w)
			return
		}

		dob, err := time.Parse("2006-01-02", dobString) //Convert dob string to time.Time
		if err != nil {
			HandleError(errors.New("invalid date of birth format"), logger, w)
			return
		}

		ObjID, err := primitive.ObjectIDFromHex(r.FormValue("editUserInfoButton")[10:34])
		filter := bson.M{"_id": ObjID}

		update := bson.M{"$set": bson.M{
			"user_info.name":    name,
			"user_info.surname": surname,
			"user_info.dob":     dob,
		}}

		_, err = data.UpdateUser(filter, update)
		if err != nil {
			logger.Errorf("A bulk write error occurred: %v", err)
			if err != nil {
				HandleError(err, logger, w)
				return
			}
			return
		} else {
			logger.Infof("User with ID: %s was UPDATED!", ObjID)
		}

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

		sessionToken := cookie.Value

		var user data.User
		user, err = data.GetUserBySessionToken(sessionToken)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			http.SetCookie(w, &http.Cookie{
				Name:    "session",
				Value:   "",
				Expires: time.Unix(0, 0),
			})
			logger.Errorf("Invalid session token.")
		}

		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
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
	logger := logging.GetLogger()
	tmpl := template.Must(template.ParseFiles("html/userProfile.html"))
	user := data.ShowUser(r)
	ordersList, err := data.GetOrdersByUserID(user.ID, true)
	if err != nil {
		logger.Errorf("Could not Get Orders By User ID: %v", err)
	}
	dataToSend := struct {
		User       data.User
		OrdersList []data.Order
	}{
		User:       user,
		OrdersList: ordersList,
	}
	err = tmpl.Execute(w, dataToSend)
	if err != nil {
		http.Error(w, "Response Writer Error!", http.StatusInternalServerError)
		logger.Errorf("Could not execute template: %v", err)
	}
}

func HandleError(err error, logger *logging.Logger, w http.ResponseWriter) {
	if err != nil {
		logger.Errorf("error: %v", err)
		_ = showMessage("/shop", "Error occurred! "+err.Error(), w)
		return
	}
	return
}
