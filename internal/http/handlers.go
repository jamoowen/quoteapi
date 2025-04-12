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
	"github.com/jamoowen/quoteapi/internal/email"
	"github.com/jamoowen/quoteapi/internal/utils"
)

type Handler struct {
	quoteService quoteapi.QuoteService
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

	// now we need to
	// a) append to list (cache)
	// b) write to csv
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) newApiKeyRequestForm(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("form").Parse(apiKeyRequestTempl))
	t.Execute(w, nil)
}

func (h *Handler) handleApiKeyRequestFormSubmission(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	otp := r.FormValue("otp")
	if utils.LooksLikeEmail(email) == false {
		badRequestError(w, "invalid email!")
		return
	}

	switch otp {
	case "":
		//send email here
		type InjectableData struct {
			Email    string
			Response string
			Error    string
		}

		h.logger.Println("Sending email...")
		dataToInjectIntoHtml := InjectableData{}
		err := h.mailer.Send(email, "Quote API OTP", "OTP: 289347y234r")
		if err == nil {
			dataToInjectIntoHtml.Email = email
			dataToInjectIntoHtml.Response = "OTP sent! Check your email and enter it here"
			// need to store in otpcache
		} else {
			h.logger.Printf("ERROR sending email: %v", err.Error())
			dataToInjectIntoHtml.Error = fmt.Sprintf("Unable to send email to %v", email)
		}
		t := template.Must(template.New("form").Parse(apiKeyRequestTempl))
		t.Execute(w, dataToInjectIntoHtml)

	default:
		//verify otp
		// easiest way to do this is create an in mem cache mapping email to otp. if expired, throw away
		// must match otp -> email
		h.logger.Print("Sending email...")
		type ApiKeyResponse struct {
			APIKEY string
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ApiKeyResponse{"20983rhnjfw2iuh"})
	}
}
