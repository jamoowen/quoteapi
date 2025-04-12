package http

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	quoteapi "github.com/jamoowen/quoteapi/internal"
	"github.com/jamoowen/quoteapi/internal/auth"
	"github.com/jamoowen/quoteapi/internal/email"
	"github.com/jamoowen/quoteapi/internal/utils"
)

type Handler struct {
	quoteService quoteapi.QuoteService
	authService  auth.AuthService
	logger       *log.Logger
	mailer       email.MailService
}

func (h *Handler) randomQuoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	randomQuote, err := h.quoteService.GetRandomQuote()
	if err != nil {
		internalServerError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(randomQuote)
}

// /quote/:quthor
func (h *Handler) getQuotesForAuthorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// need to get author out of the request
	author := r.URL.Query().Get("name")
	if author == "" {
		badRequestError(w, "author name must be provided")

	}
	randomQuote, err := h.quoteService.GetQuotesForAuthor(author)
	if err != nil {
		internalServerError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	if len(randomQuote) == 0 {
		w.Write([]byte("No quotes found for this homie"))
		return
	}
	json.NewEncoder(w).Encode(randomQuote)
}

// POST /quote/add
func (h *Handler) addQuote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Limit request body size to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		if err == http.ErrBodyReadAfterClose {
			badRequestError(w, "Request body too large")
			return
		}
		badRequestError(w, "Could not parse request body")
		return
	}
	var newQuote quoteapi.Quote
	err = json.Unmarshal(body, &newQuote)
	if err != nil {
		badRequestError(w, "Malformed JSON. Expected author & message object")
		return
	}
	// Validate required fields
	if newQuote.Author == "" || newQuote.Message == "" {
		badRequestError(w, "Author and message are required")
		return
	}
	// Validate message length (100 words)
	words := strings.Fields(newQuote.Message)
	if len(words) > 100 {
		badRequestError(w, "Message cannot exceed 100 words")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) newApiKeyRequestForm(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("form").Parse(apiKeyRequestTempl))
	t.Execute(w, nil)
}

func (h *Handler) handleApiKeyRequestFormSubmission(w http.ResponseWriter, r *http.Request) {
	var email, otp string
	contentType:= r.Header.Get("Conent-type")
	switch contentType{
	case "application/json":
		return 
	case "application/x-www-form-urlencoded":
	r.ParseForm()
	email = r.FormValue("email")
	otp = r.FormValue("otp")
	if utils.LooksLikeEmail(email) == false {
		badRequestError(w, "invalid email!")
		return
	}
	type InjectableData struct {
		Email    string
		Response string
		Error    string
	}
	type ApiKeyResponse struct {
		APIKEY string
	}
	dataToInjectIntoHtml := InjectableData{}
	t := template.Must(template.New("form").Parse(apiKeyRequestTempl))
		if err:= h.handleotp 
	switch otp {
	case "":

		if err == nil {
			dataToInjectIntoHtml.Email = email
			dataToInjectIntoHtml.Response = "OTP sent! Check your email and enter it here"
			t.Execute(w, dataToInjectIntoHtml)
			return
		}

	default:
		if isValid == false {

		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ApiKeyResponse{key})
	}
}

func (h *Handler) handleFormSubmission(w http.ResponseWriter, r *http.Request) {

			}

func (h *Handler) handleOtpRequest(email string) error {
	h.logger.Println("Sending otp email...")
	pin, err := h.authService.GenerateOtp(email)
	if err != nil {
		h.logger.Printf("ERROR generating otp: %v", err.Error())
		return fmt.Errorf("ERROR generating otp")
	}
	err = h.mailer.Send(email, "Quote API OTP", fmt.Sprintf("OTP: %v", pin))
	if err != nil {
		h.logger.Printf("ERROR sending email: %v", err.Error())
		return fmt.Errorf("ERROR sending email")
	}
	return nil
}

func (h *Handler) handleOtpSubmission(email, pin string) error {
	status, err := h.authService.AuthenticateOtp(email, pin)
	if err != nil {
		h.logger.Printf("ERROR authenticating OTP")
		return fmt.Errorf("Error authenticating OTP")
	} else if status == auth.OTPInvalid {
		return fmt.Errorf("Invalid OTP!")
	} else if status == auth.OTPExpired {
		return fmt.Errorf("Expired OTP!")
	} else if status == auth.OTPUserNotFound {
		return fmt.Errorf("No OTP found matching email")
	}
	return nil
}
