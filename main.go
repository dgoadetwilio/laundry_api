package main

import (
	"context"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go"
	"fmt"
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"time"
)

type user struct {
	Uid string `json:"uid"`
}

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func getRequestUid(w http.ResponseWriter, r *http.Request) string {

	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return ""
	}
	var u user
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&u)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return ""
	}
	errorResponse(w, "Success", http.StatusOK)

	return u.Uid

}

func pickup(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	fmt.Println("server: pickup handler started")
	defer fmt.Println("server:  pickup handler ended")

	uid := getRequestUid(w, req)
	fmt.Printf("uid passed: %s\n", uid)

	select {
	case <-time.After(3 * time.Second):

		fbConfigEnvVarName := "FIREBASE_CONFIG"
		fbConfigString, ok := os.LookupEnv(fbConfigEnvVarName)
		if !ok {
			fmt.Printf("%s not set\n", fbConfigEnvVarName)
		} else {
			fmt.Printf("%s=%s\n", fbConfigEnvVarName, fbConfigString)
		}

		opt := option.WithCredentialsJSON([]byte(fbConfigString))
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("app: %v\n", app)

		client, err := app.Firestore(ctx)
		if err != nil {
			log.Fatalln(err)
		}

		dsnap, err := client.Collection("users").Doc(uid).Get(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		m := dsnap.Data()
		firstName := m["first"]
		toNumber := m["phone"]

		fmt.Printf("Document data: %#v\n", m)
		defer client.Close()

		twilioClient := twilio.NewRestClient()
		params := &openapi.CreateMessageParams{}
		params.SetTo(toNumber.(string))
		params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
		text := fmt.Sprintf("Hey, %s, a runner is on the way to pickup your dirty laundry!", firstName)
		params.SetBody(text)
		_, err = twilioClient.ApiV2010.CreateMessage(params)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("SMS sent successfully!")
		}

	case <-ctx.Done():
		err := ctx.Err()
		fmt.Println("server:", err)
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
	}

}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/pickup", pickup)

	http.ListenAndServe(":8090", nil)
}