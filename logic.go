package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiFormat struct {
	IsRelevant                 bool   `json:"is_relevant"`
	MatchedEmail               string `json:"matched_email"`
	AnsweredMe                 bool   `json:"answered_me"`
	HasAppointmentOrWaitingList bool   `json:"has_appointment_or_waitinglist"`
	Dec                        bool   `json:"dec"`
	DecReason                  string `json:"dec_reason"`
}

type Contact struct {
	Bezirk          string `json:"bezirk"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	AnredeTitel     string `json:"anrede_titel"`
	Nachname        string `json:"nachname"`
	Geschlecht      string `json:"geschlecht"`
	Massmailed      bool   `json:"massmailed"`
	AnsweredMe      bool   `json:"answered_me"`
	HasAppointment  bool   `json:"has_appointment"`
	Dec             bool   `json:"dec"`
	DecDate         string `json:"dec_date"`
	DecReason       string `json:"dec_reason"`
	OriginalRecord  []string `json:"-"`
}

func readCSV(path string) ([]Contact, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	if len(records) < 1 {
		return nil, nil, fmt.Errorf("CSV is empty")
	}

	header := records[0]
	colIndices := make(map[string]int)
	for i, col := range header {
		colIndices[strings.ToLower(col)] = i
	}

	var contacts []Contact
	for i := 1; i < len(records); i++ {
		r := records[i]
		for len(r) < len(header) {
			r = append(r, "")
		}

		c := Contact{
			OriginalRecord: r,
		}
		if idx, ok := colIndices["bezirk"]; ok && idx < len(r) { c.Bezirk = r[idx] }
		if idx, ok := colIndices["name"]; ok && idx < len(r) { c.Name = r[idx] }
		if idx, ok := colIndices["email"]; ok && idx < len(r) { c.Email = r[idx] }
		if idx, ok := colIndices["anrede_titel"]; ok && idx < len(r) { c.AnredeTitel = r[idx] }
		if idx, ok := colIndices["nachname"]; ok && idx < len(r) { c.Nachname = r[idx] }
		if idx, ok := colIndices["geschlecht"]; ok && idx < len(r) { c.Geschlecht = r[idx] }
		if idx, ok := colIndices["massmailed"]; ok && idx < len(r) { c.Massmailed = r[idx] == "true" }
		if idx, ok := colIndices["answered_me"]; ok && idx < len(r) { c.AnsweredMe = r[idx] == "x" }
		if idx, ok := colIndices["has_appointment_or_waitinglist"]; ok && idx < len(r) { c.HasAppointment = r[idx] == "x" }
		if idx, ok := colIndices["dec"]; ok && idx < len(r) { c.Dec = r[idx] == "x" }
		if idx, ok := colIndices["dec_date"]; ok && idx < len(r) { c.DecDate = r[idx] }
		if idx, ok := colIndices["dec_reason"]; ok && idx < len(r) { c.DecReason = r[idx] }

		contacts = append(contacts, c)
	}

	return contacts, header, nil
}

func saveContacts(path string, header []string, contacts []Contact) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	
	records := [][]string{header}
	colIndices := make(map[string]int)
	for i, col := range header {
		colIndices[strings.ToLower(col)] = i
	}

	for _, c := range contacts {
		r := make([]string, len(header))
		copy(r, c.OriginalRecord)
		if len(r) < len(header) {
		    r = append(r, make([]string, len(header)-len(r))...)
		}
		
		if idx, ok := colIndices["massmailed"]; ok { 
			if c.Massmailed { r[idx] = "true" } else { r[idx] = "false" }
		}
		if idx, ok := colIndices["answered_me"]; ok {
			if c.AnsweredMe { r[idx] = "x" } else { r[idx] = "" }
		}
		if idx, ok := colIndices["has_appointment_or_waitinglist"]; ok {
			if c.HasAppointment { r[idx] = "x" } else { r[idx] = "" }
		}
		if idx, ok := colIndices["dec"]; ok {
			if c.Dec { r[idx] = "x" } else { r[idx] = "" }
		}
		if idx, ok := colIndices["dec_date"]; ok { r[idx] = c.DecDate }
		if idx, ok := colIndices["dec_reason"]; ok { r[idx] = c.DecReason }
		
		records = append(records, r)
	}
	return writer.WriteAll(records)
}

func sendEmail(auth smtp.Auth, from, to, subject, body, host, port string) error {
	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := host + ":" + port
	return smtp.SendMail(addr, auth, from, []string{to}, msg)
}

func getSalutation(c Contact) string {
	geschlecht := strings.ToLower(c.Geschlecht)
	titel := c.AnredeTitel
	nachname := c.Nachname

	salutation := "Sehr geehrte(r)"
	if geschlecht == "f" || geschlecht == "w" || geschlecht == "frau" {
		salutation = "Sehr geehrte Frau"
	} else if geschlecht == "m" || geschlecht == "h" || geschlecht == "herr" {
		salutation = "Sehr geehrter Herr"
	}

	result := salutation
	if titel != "" {
		result += " " + titel
	}
	if nachname != "" {
		result += " " + nachname
	}
	return result
}

func extractBody(msg *imap.Message) string {
	section := &imap.BodySectionName{Peek: true}
	r := msg.GetBody(section)
	if r == nil {
		r = msg.GetBody(&imap.BodySectionName{})
	}
	if r == nil {
		return ""
	}
	mr, err := mail.CreateReader(r)
	if err != nil {
		return ""
	}
	bodyText := ""
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			break
		}
		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			contentType, _, _ := h.ContentType()
			if strings.HasPrefix(contentType, "text/plain") {
				b, _ := io.ReadAll(p.Body)
				bodyText += string(b)
			}
		}
	}
	return bodyText
}

func checkReplies(email, pass, geminiKey, startDate string, contacts []Contact, header []string, sendProgress func(string)) {
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil { 
		sendProgress(fmt.Sprintf("IMAP Fehler: %v", err))
		return 
	}
	defer c.Logout()

	if err := c.Login(email, pass); err != nil { 
		sendProgress(fmt.Sprintf("Login Fehler: %v", err))
		return 
	}
	mbox, err := c.Select("INBOX", false)
	if err != nil || mbox.Messages == 0 { return }

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.FlaggedFlag}
	
	// Apply START_DATE filter
	if startDate != "" {
		t, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			criteria.Since = t
		}
	}

	uids, err := c.Search(criteria)
	if err != nil || len(uids) == 0 { return }

	ctx := context.Background()
	genaiClient, err := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
	if err != nil { return }
	defer genaiClient.Close()
	model := genaiClient.GenerativeModel("gemini-2.0-flash")

	seqset := new(imap.SeqSet)
	seqset.AddNum(uids...)
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchItem("BODY.PEEK[]"), imap.FetchUid}, messages)
	}()

	var contactsEmails []string
	for _, contact := range contacts {
		if contact.Email != "" { contactsEmails = append(contactsEmails, contact.Email) }
	}

	for msg := range messages {
		bodyText := extractBody(msg)
		
		prompt := fmt.Sprintf(`Analyze this email to determine if it is a response to a therapy appointment inquiry.
Match the sender to one of our target email addresses: %v

Return a JSON object with exactly these fields:
- "is_relevant": boolean (true if it's a response to our inquiry)
- "matched_email": string (the email address from our list that matched)
- "answered_me": boolean (true if they actually wrote back)
- "has_appointment_or_waitinglist": boolean (true if they offer a session, a first meeting, or a spot on a waiting list)
- "dec": boolean (true if they declined/rejected the inquiry)
- "dec_reason": string (briefly explain why they declined in German, e.g., "Keine Kapazitäten", "Nimmt nur Privatpatienten")

Email Subject: %s
Email Body: %s`, contactsEmails, msg.Envelope.Subject, bodyText)

		resp, err := model.GenerateContent(ctx, genai.Text(prompt))
		if err == nil && len(resp.Candidates) > 0 {
			text := ""
			for _, part := range resp.Candidates[0].Content.Parts {
				if t, ok := part.(genai.Text); ok { text += string(t) }
			}
			
			cleanJSON := strings.TrimSpace(text)
			if strings.HasPrefix(cleanJSON, "```json") {
				cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
				cleanJSON = strings.TrimSuffix(cleanJSON, "```")
				cleanJSON = strings.TrimSpace(cleanJSON)
			}

			var result GeminiFormat
			if err := json.Unmarshal([]byte(cleanJSON), &result); err == nil && result.IsRelevant {
				for i := range contacts {
					if strings.EqualFold(contacts[i].Email, result.MatchedEmail) {
						contacts[i].AnsweredMe = true
						if result.HasAppointmentOrWaitingList { 
							contacts[i].HasAppointment = true 
							// LABEL THE EMAIL IN GMAIL
							seq := new(imap.SeqSet)
							seq.AddNum(msg.Uid)
							// X-GM-LABELS is the extension for Gmail labels
							_ = c.UidStore(seq, "+X-GM-LABELS", []interface{}{"Diagnose möglich"}, nil)
							sendProgress(fmt.Sprintf("Zusage von %s erkannt! Label 'Diagnose möglich' gesetzt.", contacts[i].Email))
						}
						if result.Dec {
							contacts[i].Dec = true
							contacts[i].DecReason = result.DecReason
							contacts[i].DecDate = time.Now().Format("02.01.2006")
						}
						saveContacts(csvPath, header, contacts)
						break
					}
				}
			}
		}
		
		// Flag as seen (FlaggedFlag used to avoid re-processing)
		item := imap.FormatFlagsOp(imap.AddFlags, true)
		seq := new(imap.SeqSet)
		seq.AddNum(msg.Uid)
		_ = c.UidStore(seq, item, []interface{}{imap.FlaggedFlag}, nil)
	}
	<-done
}
