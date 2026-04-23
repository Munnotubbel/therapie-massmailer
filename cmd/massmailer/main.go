package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/browser"
)

//go:embed ui/index.html defaults/*
var staticContent embed.FS

type Config struct {
	Email     string `json:"email"`
	Pass      string `json:"pass"`
	Key       string `json:"key"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	StartDate string `json:"start_date"`
}

var dataDir = "therapie_massmailer_daten"
var csvPath = dataDir + "/contacts.csv"
var envPath = dataDir + "/.env"
var msgPath = dataDir + "/message.txt"
var lastHeartbeat time.Time

func main() {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.Mkdir(dataDir, 0755)
	}

	// Initialize default files if missing
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		content, _ := staticContent.ReadFile("defaults/contacts.csv")
		os.WriteFile(csvPath, content, 0644)
	}
	if _, err := os.Stat(msgPath); os.IsNotExist(err) {
		content, _ := staticContent.ReadFile("defaults/message.txt")
		os.WriteFile(msgPath, content, 0644)
	}

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/api/config", handleConfig)
	http.HandleFunc("/api/contacts", handleContacts)
	http.HandleFunc("/api/contacts/add", handleAddContact)
	http.HandleFunc("/api/contacts/update", handleUpdateContact)
	http.HandleFunc("/api/download", handleDownload)
	http.HandleFunc("/api/quit", handleQuit)
	http.HandleFunc("/api/heartbeat", handleHeartbeat)
	http.HandleFunc("/api/run", handleRun)
	
	_ = godotenv.Load(envPath)
	if os.Getenv("START_DATE") == "" {
		now := time.Now().Format("2006-01-02")
		f, err := os.OpenFile(envPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			fmt.Fprintf(f, "START_DATE=%s\n", now)
			f.Close()
			_ = godotenv.Load(envPath)
		}
	}

	port := "8080"
	url := "http://localhost:" + port
	
	fmt.Printf("Therapie Massmailer läuft auf %s\n", url)
	
	// Shutdown logic if no heartbeat
	lastHeartbeat = time.Now()
	go func() {
		for {
			time.Sleep(5 * time.Second)
			if time.Since(lastHeartbeat) > 20*time.Second {
				fmt.Println("Kein Browser-Fenster mehr offen. Beende App...")
				os.Exit(0)
			}
		}
	}()

	go func() {
		time.Sleep(1 * time.Second)
		browser.OpenURL(url)
	}()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	content, _ := staticContent.ReadFile("ui/index.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_ = godotenv.Load(envPath)
		cfg := Config{
			Email:     os.Getenv("GMAIL_USER"),
			Pass:      os.Getenv("GMAIL_PASS"),
			Key:       os.Getenv("GEMINI_API_KEY"),
			StartDate: os.Getenv("START_DATE"),
		}
		
		if cfg.StartDate == "" {
			cfg.StartDate = time.Now().Format("2006-01-02")
		}
		
		bodyBytes, _ := os.ReadFile(msgPath)
		content := string(bodyBytes)
		lines := strings.SplitN(content, "\n", 2)
		if len(lines) > 0 && strings.HasPrefix(strings.ToLower(lines[0]), "betreff:") {
			cfg.Subject = strings.TrimSpace(strings.TrimPrefix(lines[0][8:], " "))
			if len(lines) > 1 {
				cfg.Body = strings.TrimLeft(lines[1], "\r\n")
			}
		} else {
			cfg.Body = content
		}
		
		json.NewEncoder(w).Encode(cfg)
	} else if r.Method == "POST" {
		var cfg Config
		json.NewDecoder(r.Body).Decode(&cfg)
		
		if cfg.StartDate == "" {
			cfg.StartDate = time.Now().Format("2006-01-02")
		}

		envContent := fmt.Sprintf("GMAIL_USER=%s\nGMAIL_PASS=%s\nGEMINI_API_KEY=%s\nSTART_DATE=%s\n", cfg.Email, cfg.Pass, cfg.Key, cfg.StartDate)
		os.WriteFile(envPath, []byte(envContent), 0644)
		
		fullMsg := fmt.Sprintf("Betreff: %s\n\n%s", cfg.Subject, cfg.Body)
		os.WriteFile(msgPath, []byte(fullMsg), 0644)
		
		w.WriteHeader(http.StatusOK)
	}
}

func handleContacts(w http.ResponseWriter, r *http.Request) {
	contacts, _, err := readCSV(csvPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(contacts)
}

func handleAddContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var c Contact
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	contacts, header, err := readCSV(csvPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Prepare original record to match header length
	c.OriginalRecord = make([]string, len(header))
	contacts = append(contacts, c)
	
	if err := saveContacts(csvPath, header, contacts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func handleUpdateContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var c Contact
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	contacts, header, err := readCSV(csvPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	found := false
	for i := range contacts {
		if strings.EqualFold(contacts[i].Email, c.Email) {
			contacts[i].Bezirk = c.Bezirk
			contacts[i].Name = c.Name
			contacts[i].AnredeTitel = c.AnredeTitel
			contacts[i].Nachname = c.Nachname
			contacts[i].Geschlecht = c.Geschlecht
			found = true
			break
		}
	}
	
	if !found {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}
	
	if err := saveContacts(csvPath, header, contacts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "attachment; filename=contacts.csv")
	w.Header().Set("Content-Type", "text/csv")
	http.ServeFile(w, r, csvPath)
}

func handleQuit(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}

func handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	lastHeartbeat = time.Now()
	w.WriteHeader(http.StatusOK)
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	sendProgress := func(msg string) {
		fmt.Fprintf(w, "data: %s\n\n", `{"type":"progress", "message":"`+msg+`"}`)
		w.(http.Flusher).Flush()
	}

	_ = godotenv.Load(envPath)
	email := os.Getenv("GMAIL_USER")
	pass := os.Getenv("GMAIL_PASS")
	key := os.Getenv("GEMINI_API_KEY")

	contacts, header, _ := readCSV(csvPath)
	
	// Mailing
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	auth := smtp.PlainAuth("", email, pass, smtpHost)
	
	bodyBytes, _ := os.ReadFile(msgPath)
	bodyContent := string(bodyBytes)
	subject := "Anfrage"
	lines := strings.SplitN(bodyContent, "\n", 2)
	if len(lines) > 0 && strings.HasPrefix(strings.ToLower(lines[0]), "betreff:") {
		subject = strings.TrimSpace(strings.TrimPrefix(lines[0][8:], " "))
		bodyContent = strings.TrimLeft(lines[1], "\r\n")
	}

	for i := range contacts {
		if contacts[i].Massmailed || contacts[i].Email == "" {
			continue
		}
		
		sendProgress(fmt.Sprintf("Sende an %s...", contacts[i].Email))
		
		personalizedBody := strings.ReplaceAll(bodyContent, "{{Briefanrede}}", getSalutation(contacts[i]))
		err := sendEmail(auth, email, contacts[i].Email, subject, personalizedBody, smtpHost, smtpPort)
		if err == nil {
			contacts[i].Massmailed = true
			saveContacts(csvPath, header, contacts)
		} else {
			sendProgress(fmt.Sprintf("Fehler bei %s: %v", contacts[i].Email, err))
		}
		time.Sleep(1 * time.Second)
	}

	// Reply check
	sendProgress("Prüfe auf neue Antworten...")
	startDate := os.Getenv("START_DATE")
	checkReplies(email, pass, key, startDate, contacts, header, sendProgress)
	
	fmt.Fprintf(w, "data: %s\n\n", `{"type":"done"}`)
	w.(http.Flusher).Flush()
}
