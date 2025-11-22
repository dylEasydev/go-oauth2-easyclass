package utils

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/gomail.v2"
)

// SendVerificationCode envoie un code de vérification par Gmail
func SendVerificationCode(dest string, code string, name string) error {
	from := os.Getenv("COMPANING_MAIl")
	appPassword := os.Getenv("PASSWORD_MAIL")
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", dest)
	m.SetHeader("Subject", "Votre code de vérification")

	expiredAt := time.Now().Add(1 * time.Hour)

	plain := fmt.Sprintf("Bonjour,\n\nVotre code de vérification est : %s\n\nCe code expirera :%s .\n", code, expiredAt.Format("16h04"))

	html := fmt.Sprintf(`
	<!doctype html>
	<html lang="fr">
	<head>
	  <meta charset="utf-8">
	  <style>
	    body { font-family: Arial, sans-serif; background:#f9f9f9; padding:20px; }
	    .box { max-width:500px; margin:0 auto; background:white; padding:20px; border-radius:8px; box-shadow:0 2px 8px rgba(0,0,0,0.1);}
	    h1 { color:#333; font-size:20px; }
	    .code { font-size:28px; font-weight:bold; color:#2d89ef; letter-spacing:4px; margin:20px 0; }
	    p { color:#555; font-size:14px; }
	  </style>
	</head>
	<body>
	  <div class="box">
	    <h1>Code de Vérification</h1>
	    <p>Bonjour,%s</p>
	    <p>Voici votre code de vérification :</p>
	    <div class="code">%s</div>
	    <p>Ce code est valable jusqu'à %s.</p>
	    <p style="margin-top:20px; font-size:12px; color:#888;">&copy; 2025 MonApplication</p>
	  </div>
	</body>
	</html>`, name, code, expiredAt.Format("16h04"))

	m.SetBody("text/plain", plain)
	m.AddAlternative("text/html", html)

	d := gomail.NewDialer(smtpHost, smtpPort, from, appPassword)

	return d.DialAndSend(m)
}
