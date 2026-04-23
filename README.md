# 📧 Therapie-Termin Massmailer - Web-Edition

Dieser intelligente Assistent hilft dir dabei, zeitnah einen Therapieplatz zu finden. Er automatisiert den Erstkontakt zu Therapeuten und nutzt Künstliche Intelligenz (Google Gemini), um eingehende Emails automatisch zu verstehen und für dich zu sortieren.

---

## ⚡ 1. Direkter Download (Fertige Programme)

Hier findest du die startfertigen Programme. Klicke auf den Link für dein System, um zum Ordner zu gelangen und die Datei herunterzuladen:

*   🍎 **[Mac (neue M1/M2/M3 Chips)](https://github.com/Munnotubbel/therapie-massmailer/raw/main/bin/massmailer_mac_apple_silicon)**
*   🍎 **[Mac (alte Intel Chips)](https://github.com/Munnotubbel/therapie-massmailer/raw/main/bin/massmailer_mac_intel)**
*   🪟 **[Windows PC](https://github.com/Munnotubbel/therapie-massmailer/raw/main/bin/massmailer_windows.exe)**
*   🐧 **[Linux](https://github.com/Munnotubbel/therapie-massmailer/raw/main/bin/massmailer)**

---

## 🚀 2. Installation & Start (Schritt für Schritt)

Das Programm muss nicht installiert werden. Es läuft direkt aus dem Ordner heraus.

### Mac (Apple Silicon & Intel)
1.  **Ordner öffnen:** Gehe in den Ordner `scripts/`.
2.  **Sicherheits-Start:** Mache einen **Rechtsklick** auf die Datei `Start_Massmailer.command` und wähle **"Öffnen"**.
    *   *Wichtig:* Klicke im erscheinenden Fenster auf **"Öffnen"**. Dieser Schritt ist nur beim ersten Mal nötig, damit Apple weiß, dass du dem Programm vertraust.
3.  **Browser:** Es öffnet sich automatisch dein Internet-Browser (z.B. Safari oder Chrome) mit der Benutzeroberfläche.

### Windows
1.  **Starten:** Gehe in den Ordner `bin/` und mache einen Doppelklick auf `massmailer_windows.exe`.
2.  **Konsole:** Es öffnet sich kurz ein schwarzes Fenster im Hintergrund – das ist normal, lass es einfach offen.
3.  **Browser:** Dein Standard-Browser öffnet sich automatisch mit der App.

### Linux
1.  **Starten:** Öffne ein Terminal im Projektordner.
2.  **Befehl:** Führe die Datei mit `./bin/massmailer` aus.
3.  **Berechtigungen:** Falls die Datei nicht startet, gib ihr einmalig Rechte mit `chmod +x bin/massmailer`.
4.  **Browser:** Die App öffnet sich in deinem Browser.

---

## 🛠️ 3. Die Einrichtung (Wizard)

Beim ersten Start hilft dir ein Assistent bei der Einrichtung. Du benötigst zwei Dinge:

### A. Gmail App-Passwort (glasklar erklärt)
Du kannst nicht dein normales Passwort nutzen. Google benötigt ein spezielles 16-stelliges Passwort:
1.  Gehe in dein [Google-Konto](https://myaccount.google.com/).
2.  Klicke links auf **"Sicherheit"**.
3.  Stelle sicher, dass die **Bestätigung in zwei Schritten** aktiviert ist.
4.  Suche oben in der Leiste nach **"App-Passwörter"**.
5.  Erstelle ein Passwort für "Therapie-Tool".
6.  Kopiere den gelben Code (16 Buchstaben) und füge ihn in der App ein.

### B. KI-Schlüssel (Gemini API Key)
Damit die App die Emails für dich lesen kann:
1.  Besuche [Google AI Studio](https://aistudio.google.com/).
2.  Melde dich mit deinem Gmail-Konto an.
3.  Klicke links auf **"Get API key"** und dann auf den blauen Button **"Create API key"**.
4.  Kopiere den langen Code (beginnt mit `AIza...`) in die App.

---

## 📊 4. Bedienung des Dashboards

*   **+ Neu:** Füge einen neuen Therapeuten hinzu.
*   **✏️ Bearbeiten:** Korrigiere Vertipper bei Namen oder Titeln (die Email bleibt fest).
*   **🔄 Refresh:** Nutze diesen Button, wenn du die `contacts.csv` Datei händisch geändert hast.
*   **💾 Download:** Speichere deine gesamte Liste als CSV-Datei zur Sicherung.
*   **▶ Versand starten:** Schickt die Emails an alle Kontakte raus, die noch nichts erhalten haben.

### Was bedeuten die Farben?
*   ⚪ **Grau:** Noch nichts gesendet.
*   🔵 **Blau:** Email wurde erfolgreich verschickt.
*   🔴 **Rot:** Der Therapeut hat abgesagt (KI hat die Email verstanden).
*   🟢 **Grün:** Zusage oder Warteliste erhalten!
    *   *Highlight:* Diese Emails werden in deinem Gmail-Konto automatisch mit dem Label **"Diagnose möglich"** markiert, damit du sie sofort findest.

---

## 🛑 5. App beenden

*   **Automatik:** Schließe einfach den Browser-Tab. Das Programm merkt das und beendet sich nach ca. 20 Sekunden von selbst.
*   **Sofort (⏻):** Klicke oben rechts auf das rote Power-Symbol, um alles sofort zu stoppen.

---

## 📁 6. Projektstruktur
*   `bin/`: Ausführbare Programme.
*   `cmd/`: Quellcode (Go).
*   `ui/`: Benutzeroberfläche (HTML/JS).
*   `therapie_massmailer_daten/`: **Hier liegen deine Daten.** Dieser Ordner enthält deine Einstellungen (`.env`), deine Kontakte (`contacts.csv`) und deine Nachricht (`message.txt`). So bleibt dein Hauptordner immer sauber.
*   `scripts/`: Start-Hilfen für Nutzer.
*   `samples/`: Beispiel-Dateien zur Referenz.

---

## 🏗️ 7. Für Entwickler (Build from Source)

Falls du den Code selbst kompilieren möchtest:

1.  **Go installieren:** Lade Go von [go.dev](https://go.dev/dl/) herunter.
2.  **Projekt laden:** Öffne ein Terminal im Hauptordner.
3.  **Abhängigkeiten:** `go mod tidy`
4.  **Bauen:**
    *   Lokal: `go build -o bin/massmailer ./cmd/massmailer`
    *   Windows: `GOOS=windows GOARCH=amd64 go build -o bin/massmailer_windows.exe ./cmd/massmailer`
    *   Mac M1/M2/M3: `GOOS=darwin GOARCH=arm64 go build -o bin/massmailer_mac_apple_silicon ./cmd/massmailer`
    *   Mac Intel: `GOOS=darwin GOARCH=amd64 go build -o bin/massmailer_mac_intel ./cmd/massmailer`
