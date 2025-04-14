package http

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	quoteapi "github.com/jamoowen/quoteapi/internal"
	"github.com/jamoowen/quoteapi/internal/auth"
	"github.com/jamoowen/quoteapi/internal/email"
	"github.com/jamoowen/quoteapi/internal/problems"
	"github.com/jamoowen/quoteapi/internal/utils"
)

type Handler struct {
	quoteService quoteapi.QuoteService
	authService  auth.AuthService
	logger       *log.Logger
	mailer       email.MailService
}
type ApiKeyResponse struct {
	APIKEY string
}

func (h *Handler) randomQuoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	randomQuote, err := h.quoteService.GetRandomQuote()
	if err != nil {
		h.handleHttpError(w, problems.NewHTTPError(http.StatusInternalServerError, "", err))
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
		h.handleHttpError(w, problems.NewHTTPError(http.StatusBadRequest, "author name must be provided", nil))
		return
	}
	randomQuote, err := h.quoteService.GetQuotesForAuthor(author)
	if err != nil {
		h.handleHttpError(w, problems.NewHTTPError(http.StatusInternalServerError, "", err))
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
			h.handleHttpError(w, problems.NewHTTPError(http.StatusBadRequest, "Request too large", nil))
			return
		}
		h.handleHttpError(w, problems.NewHTTPError(http.StatusBadRequest, "Could not parse request body", nil))
		return
	}
	var newQuote quoteapi.Quote
	err = json.Unmarshal(body, &newQuote)
	if err != nil {
		h.handleHttpError(w, problems.NewHTTPError(http.StatusBadRequest, "Malformed JSON. Expected author & message object", nil))
		return
	}
	// Validate required fields
	if newQuote.Author == "" || newQuote.Message == "" {
		h.handleHttpError(w, problems.NewHTTPError(http.StatusBadRequest, "Author and Message required", nil))
		return
	}
	// Validate message length (100 words)
	words := strings.Fields(newQuote.Message)
	if len(words) > 100 {
		h.handleHttpError(w, problems.NewHTTPError(http.StatusBadRequest, "Message cannot exceed 100 words", nil))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleAuthenticate(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	fmt.Printf("Content type: %v\n", contentType)
	switch contentType {
	case "application/json":
		// to do -> handle authetication via backend requests
	default:
		h.handleFormSubmission(w, r)
	}
}

func (h *Handler) handleFormSubmission(w http.ResponseWriter, r *http.Request) {
	type InjectableData struct {
		Email    string
		Response string
		Error    string
	}

	dataToInjectIntoHtml := InjectableData{}
	t, err := template.ParseFiles("internal/http/templates/authenticate.html")
	if err != nil {
		h.handleHttpError(w, problems.NewHTTPError(http.StatusInternalServerError, "Failed to load html", err))
		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	submittedPin := r.FormValue("pin")

	// no submission, just display the form
	if submittedPin == "" && email == "" {
		t.Execute(w, nil)
		return

	} else if submittedPin == "" {
		httpErr := h.handleOtpRequest(email)
		if httpErr != nil {
			dataToInjectIntoHtml.Error = httpErr.Message
			t.Execute(w, dataToInjectIntoHtml)
			return
		}
		dataToInjectIntoHtml.Email = email
		dataToInjectIntoHtml.Response = "OTP sent! Check your email and enter it here"
		t.Execute(w, dataToInjectIntoHtml)
		return

	} else {
		apiKey, httpErr := h.handleOtpSubmission(email, submittedPin, r.Context())
		if httpErr != nil {
			dataToInjectIntoHtml.Error = httpErr.Message
			t.Execute(w, dataToInjectIntoHtml)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ApiKeyResponse{apiKey})
	}
}

func (h *Handler) handleOtpRequest(email string) *problems.HTTPError {
	if utils.LooksLikeEmail(email) == false {
		return problems.NewHTTPError(http.StatusBadRequest, "Invalid email", nil)
	}
	fmt.Println("gernerating otp")

	pin, err := h.authService.GenerateOtp(email)
	if err != nil {
		return problems.NewHTTPError(http.StatusInternalServerError, "ERROR generating otp", err)
	}
	err = h.mailer.Send(email, "Quote API OTP", fmt.Sprintf("OTP: %v", pin))
	if err != nil {
		return problems.NewHTTPError(http.StatusInternalServerError, "ERROR sending otp email", err)
	}
	return nil
}

func (h *Handler) handleOtpSubmission(email, pin string, ctx context.Context) (string, *problems.HTTPError) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if utils.LooksLikeEmail(email) == false {
		return "", problems.NewHTTPError(http.StatusBadRequest, "Invalid email", nil)
	}
	status, err := h.authService.AuthenticateOtp(email, pin)
	if err != nil {
		return "", problems.NewHTTPError(http.StatusInternalServerError, "", err)
	} else if status == auth.OTPInvalid {
		return "", problems.NewHTTPError(http.StatusBadRequest, "Invalid OTP!", err)
	} else if status == auth.OTPExpired {
		return "", problems.NewHTTPError(http.StatusBadRequest, "Expired OTP!", err)
	} else if status == auth.OTPUserNotFound {
		return "", problems.NewHTTPError(http.StatusBadRequest, "Invalid email!", err)
	}
	newApiKey, err := h.authService.CreateNewApiKeyForUser(email, ctx)
	if err != nil {
		return "", problems.NewHTTPError(http.StatusInternalServerError, "", err)
	}
	return newApiKey, nil
}

func (h *Handler) handleHttpError(w http.ResponseWriter, err *problems.HTTPError) {
	switch err.Code {
	case http.StatusInternalServerError:
		msg := "Internal Server Error"
		if err.Message != "" {
			msg = err.Message
		}
		h.logger.Print(err.Unwrap().Error())
		w.WriteHeader(err.Code)
		http.Error(w, msg, http.StatusInternalServerError)
	default:
		w.WriteHeader(err.Code)
		http.Error(w, err.Message, http.StatusInternalServerError)
	}
}
