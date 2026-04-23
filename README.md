# 📧 Therapie-Termin Massmailer - Web-Edition

Diese Software ist ein intelligenter Assistent für die Suche nach Therapieplätzen. Sie automatisiert den Erstkontakt zu Therapeuten und wertet Rückmeldungen mithilfe von Künstlicher Intelligenz (Google Gemini) automatisch aus. 

Die grafische Benutzeroberfläche ermöglicht eine einfache Bedienung ohne technische Vorkenntnisse.

---

## 🚀 1. Installation & Start

### Mac (Apple Silicon & Intel)
1.  **Entpacken:** Kopiere alle Dateien in einen eigenen Ordner (z.B. auf den Schreibtisch).
2.  **Starten:** Mache einen **Rechtsklick** auf die Datei `Start_Massmailer.command` und wähle **"Öffnen"**. 
    *   *Hinweis:* Beim ersten Start musst du den Rechtsklick nutzen, um die Sicherheitsabfrage von Apple zu bestätigen.
3.  **Browser:** Dein Standard-Browser öffnet sich automatisch mit dem Dashboard.

### Windows
1.  Führe die Datei `massmailer_windows.exe` per Doppelklick aus.
2.  Das Programm öffnet ein schwarzes Konsolenfenster im Hintergrund und startet dann automatisch deinen Browser.

---

## 🛠️ 2. Die Einrichtung (Wizard)

Beim ersten Start führt dich ein Assistent durch die Konfiguration:

### Schritt 1: Gmail App-Passwort
Du kannst nicht dein normales Gmail-Passwort nutzen. Google benötigt ein spezielles 16-stelliges **App-Passwort**:
1. Gehe in dein [Google-Konto](https://myaccount.google.com/).
2. Suche nach **"App-Passwörter"**.
3. Erstelle ein neues Passwort für "Therapie-Tool".
4. Kopiere den Code (z.B. `abcd efgh ijkl mnop`) in die App.

### Schritt 2: KI-Schlüssel (Gemini)
Damit die App Emails verstehen kann, benötigt sie einen API-Key:
1. Besuche [Google AI Studio](https://aistudio.google.com/).
2. Klicke auf **"Get API key"** -> **"Create API key in new project"**.
3. Kopiere den langen Code (beginnt mit `AIza...`) in die App.

### Schritt 3: Email-Vorlage
Hier gestaltest du deine Nachricht. Die Anrede `{{Briefanrede}}` wird automatisch für jeden Kontakt individualisiert (z.B. "Sehr geehrte Frau Müller").

---

## 📊 3. Bedienung des Dashboards

### Kontakte verwalten
*   **+ Neu:** Füge neue Therapeuten direkt über die Oberfläche hinzu.
*   **✏️ Bearbeiten:** Klicke auf das Stift-Symbol, um Daten zu korrigieren. (Die Email-Adresse bleibt dabei fest, um den Status nicht zu verlieren).
*   **🔄 Refresh:** Lade die Liste neu, falls du die `contacts.csv` manuell bearbeitet hast.
*   **💾 Download:** Lade deine aktuelle Kontaktliste als CSV-Datei herunter.

### Status-Anzeigen
*   ⚪ **Grau:** Bereit zum Versand.
*   🔵 **Blau:** Email wurde erfolgreich verschickt.
*   🔴 **Rot:** Absage (KI hat erkannt, dass derzeit kein Platz frei ist).
*   🟢 **Grün:** Zusage oder Warteliste. 
    *   *Highlight:* Bei einer Zusage wird die Email in deinem Gmail-Konto automatisch mit dem Label **"Diagnose möglich"** markiert!

---

## 🛑 4. App beenden

*   **Automatik:** Sobald du den Browser-Tab schließt, beendet sich das Programm nach ca. 20 Sekunden von selbst.
*   **Sofort-Aus (⏻):** Klicke oben rechts auf das rote Power-Symbol, um den Server sofort zu stoppen.

---

## 🔒 Sicherheit & Datenschutz
*   **Lokal:** Alle Passwörter und Daten bleiben auf deinem Rechner in der Datei `.env` gespeichert.
*   **Anonymisierung:** Nur der Inhalt der Antwort-Emails wird zur Analyse verschlüsselt an Google Gemini gesendet.
*   **Transparenz:** Du behältst die volle Kontrolle über dein Gmail-Konto.
