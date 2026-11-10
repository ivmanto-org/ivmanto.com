package email

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	"ivmanto.com/backend/internal/config"
	"ivmanto.com/backend/internal/ical"
)

// SmtpService is a concrete implementation of the email Service using SMTP.
type SmtpService struct {
	cfg    *config.EmailConfig
	auth   smtp.Auth
	logger *slog.Logger
}

// NewSmtpService creates a new SMTP email service.
// It requires a valid EmailConfig. The SMTP password should be loaded from a secure source.
func NewSmtpService(cfg *config.EmailConfig, logger *slog.Logger) *SmtpService {
	// In a real app, you'd check if cfg.SmtpPass is empty and handle it.
	auth := smtp.PlainAuth("", cfg.SendFrom, cfg.SmtpPass, cfg.SmtpHost)
	return &SmtpService{
		cfg:    cfg,
		auth:   auth,
		logger: logger,
	}
}

// stripPlusAlias removes the +alias part of an email address, which is sometimes
// required for the SMTP RCPT TO command.
// e.g., "user+alias@example.com" becomes "user@example.com".
func stripPlusAlias(email string) string {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return email // Not a standard email format, return as is.
	}
	localPart, domainPart := parts[0], parts[1]

	plusIndex := strings.Index(localPart, "+")
	if plusIndex == -1 {
		return email // No plus alias, return as is.
	}

	return fmt.Sprintf("%s@%s", localPart[:plusIndex], domainPart)
}

// send is a private helper to construct and dispatch an email.
func (s *SmtpService) send(to, cc []string, subject, htmlBody string, attachment *ical.Attachment) error {
	// Build the message body using multipart writer for attachments.
	var bodyBuffer bytes.Buffer
	writer := multipart.NewWriter(&bodyBuffer)

	// HTML part
	htmlHeaders := textproto.MIMEHeader{}
	htmlHeaders.Set("Content-Type", "text/html; charset=utf-8")
	part, err := writer.CreatePart(htmlHeaders)
	if err != nil {
		return fmt.Errorf("failed to create html part: %w", err)
	}
	_, err = part.Write([]byte(htmlBody))
	if err != nil {
		return fmt.Errorf("failed to write html body: %w", err)
	}

	// Attachment part
	if attachment != nil {
		attachment.Headers.Set("Content-Transfer-Encoding", "base64")
		part, err := writer.CreatePart(attachment.Headers)
		if err != nil {
			return fmt.Errorf("failed to create attachment part: %w", err)
		}
		b64Writer := base64.NewEncoder(base64.StdEncoding, part)
		_, _ = b64Writer.Write(attachment.Body)
		b64Writer.Close()
	}

	writer.Close()

	// Build the full message with top-level headers.
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.cfg.SendFromAlias, s.cfg.SendFrom))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))
	if len(cc) > 0 {
		msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(cc, ", ")))
	}
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary()))
	msg.WriteString("\r\n")
	msg.Write(bodyBuffer.Bytes())

	addr := fmt.Sprintf("%s:%s", s.cfg.SmtpHost, s.cfg.SmtpPort)

	allRecipients := append(to, cc...)

	s.logger.Info("Connecting to SMTP server", "address", addr)
	c, err := smtp.Dial(addr)
	if err != nil {
		s.logger.Error("Failed to connect to SMTP server", "error", err)
		return err
	}
	defer c.Close()

	// It's good practice to say hello to the server.
	// The hostname can be anything, but localhost is common for clients.
	if err = c.Hello("localhost"); err != nil {
		s.logger.Error("Failed to send HELO to SMTP server", "error", err)
		return err
	}

	// Use STARTTLS. Gmail requires this on port 587.
	if ok, _ := c.Extension("STARTTLS"); ok {
		s.logger.Info("Server supports STARTTLS. Upgrading connection...")
		config := &tls.Config{ServerName: s.cfg.SmtpHost}
		if err = c.StartTLS(config); err != nil {
			s.logger.Error("Failed to start TLS", "error", err)
			return err
		}
	}

	// Authenticate.
	if s.auth != nil {
		s.logger.Info("Authenticating with SMTP server")
		if err = c.Auth(s.auth); err != nil {
			s.logger.Error("SMTP authentication failed", "error", err)
			return err
		}
	}

	// Set the sender.
	s.logger.Debug("Setting SMTP sender", "from", s.cfg.SendFrom)
	if err = c.Mail(s.cfg.SendFrom); err != nil {
		s.logger.Error("SMTP MAIL command failed", "error", err)
		return err
	}

	// Set the recipients.
	for _, rcpt := range allRecipients {
		// Strip +alias for the RCPT TO command, as some servers require the base address.
		baseRcpt := stripPlusAlias(rcpt)
		s.logger.Debug("Adding SMTP recipient", "recipient", rcpt, "base_recipient", baseRcpt)
		if err = c.Rcpt(baseRcpt); err != nil {
			s.logger.Error("SMTP RCPT command failed", "recipient", rcpt, "error", err)
			return err
		}
	}

	// Get the writer for the data and write the message.
	s.logger.Debug("Sending email body")
	wc, err := c.Data()
	if err != nil {
		s.logger.Error("SMTP DATA command failed", "error", err)
		return err
	}
	_, err = wc.Write(msg.Bytes())
	if err != nil {
		s.logger.Error("Failed to write email body to SMTP connection", "error", err)
		return err
	}
	err = wc.Close()
	if err != nil {
		s.logger.Error("Failed to close SMTP data writer", "error", err)
		return err
	}

	// Quit the session.
	s.logger.Debug("Quitting SMTP session")
	err = c.Quit()
	if err != nil {
		s.logger.Warn("Failed to quit SMTP session cleanly", "error", err)
		return err
	}

	s.logger.Info("Email sent successfully", "recipients", strings.Join(allRecipients, ", "))
	return nil
}

// buildBookingConfirmationHTML renders the user-facing confirmation email
// body. Extracted from SendBookingConfirmation so tests can assert the
// rendered HTML without going through the SMTP transport.
func buildBookingConfirmationHTML(details BookingConfirmationDetails) string {
	var meetLinkHTML string
	if details.MeetLink != "" {
		meetLinkHTML = fmt.Sprintf(`<li><strong>Google Meet Link:</strong> <a href="%s">%s</a></li>`, details.MeetLink, details.MeetLink)
	}

	body := fmt.Sprintf(`
		<p>Hi %s,</p>
		<p>Your 30-minute consultation is confirmed. Here are the details:</p>
		<ul>
		<li><strong>Date:</strong> %s</li>
		<li><strong>Time:</strong> %s - %s (%s)</li>
		%s
		</ul>
		<p>A calendar invitation (.ics file) is attached to this email. Please open it to add the event to your calendar.</p>
		<p>We look forward to speaking with you!</p>
		<p>Thanks,<br>The IVMANTO Team</p>`,
		details.ToName,
		details.StartTime.Format("Monday, January 2, 2006"),
		details.StartTime.Format("3:04 PM"),
		details.EndTime.Format("3:04 PM"),
		details.Timezone,
		meetLinkHTML)

	if details.CancellationURL != "" {
		body += fmt.Sprintf(`<p style="font-size: small; color: #666;">Need to make a change? <a href="%s">Cancel this booking</a>.</p>`, details.CancellationURL)
	}
	return body
}

// SendBookingConfirmation sends a confirmation email to the user.
func (s *SmtpService) SendBookingConfirmation(details BookingConfirmationDetails) error {
	subject := "Your consultation is confirmed!"
	htmlBody := buildBookingConfirmationHTML(details)

	// Generate the .ics file content
	icsContent := ical.Generate(ical.EventDetails{
		UID:         details.IcsUID,
		StartTime:   details.StartTime,
		EndTime:     details.EndTime,
		Summary:     details.IcsSummary,
		Description: details.IcsDescription,
		Location:    details.MeetLink,
		Name:        details.ToName,
		Email:       details.ToEmail,
		Timezone:    details.IcsTimezone,
	})

	// Create the attachment
	attachment := &ical.Attachment{
		Headers: textproto.MIMEHeader{
			"Content-Type":        {`text/calendar; charset="utf-8"; method=REQUEST`},
			"Content-Disposition": {`attachment; filename="invite.ics"`},
		},
		Body: []byte(icsContent),
	}

	return s.send([]string{details.ToEmail}, nil, subject, htmlBody, attachment)
}

// SendBookingNotificationToAdmin sends a notification email to the admin.
func (s *SmtpService) SendBookingNotificationToAdmin(name, clientEmail string, startTime time.Time, notes string) error {
	// Use a '+booking' alias to ensure delivery to the admin's inbox.
	parts := strings.Split(s.cfg.SendFrom, "@")
	var adminEmail string
	if len(parts) == 2 {
		adminEmail = fmt.Sprintf("%s+booking@%s", parts[0], parts[1])
	} else {
		adminEmail = s.cfg.SendFrom // Fallback for non-standard emails
	}

	subject := "New Consultation Booked!"
	body := fmt.Sprintf("New booking with:<br>Name: %s<br>Email: %s<br>Time: %s<br>Notes: %s", name, clientEmail, startTime.Format(time.RFC1123), notes)
	return s.send([]string{adminEmail}, nil, subject, body, nil)
}

// SendContactMessage sends the contact form message to the admin.
func (s *SmtpService) SendContactMessage(msg ContactMessage) error {
	// Trick to help Gmail deliver to inbox: add a suffix to the username.
	// e.g., nikolay.tonev@ivmanto.com -> nikolay.tonev+contact@ivmanto.com
	parts := strings.Split(s.cfg.SendFrom, "@")
	var adminEmail string
	if len(parts) == 2 {
		adminEmail = fmt.Sprintf("%s+contact@%s", parts[0], parts[1])
	} else {
		adminEmail = s.cfg.SendFrom // Fallback for non-standard emails
	}

	subject := fmt.Sprintf("New Contact Message from %s", msg.Name)
	body := fmt.Sprintf("From: %s &lt;%s&gt;<br><br>Message:<br>%s", msg.Name, msg.Email, msg.Message)

	var ccList []string
	if msg.SendCopyToSelf {
		ccList = append(ccList, msg.Email)
	}

	return s.send([]string{adminEmail}, ccList, subject, body, nil)
}

// SendLeadNotification emails an advisor lead (epic #112) to the configured
// lead recipient (LEADS_TO, default SendFrom). The recipient uses the
// +lead username alias for Gmail inbox delivery, mirroring SendContactMessage.
func (s *SmtpService) SendLeadNotification(lead LeadInfo) error {
	toEmail := s.cfg.LeadTo
	if toEmail == "" {
		toEmail = s.cfg.SendFrom
	}
	parts := strings.Split(toEmail, "@")
	var leadRecipient string
	if len(parts) == 2 {
		leadRecipient = fmt.Sprintf("%s+lead@%s", parts[0], parts[1])
	} else {
		leadRecipient = toEmail
	}

	subject := fmt.Sprintf("New advisor lead: %s (%s)", lead.Company, lead.Name)
	return s.send([]string{leadRecipient}, nil, subject, buildLeadNotificationHTML(lead), nil)
}

// buildLeadNotificationHTML renders the lead email body. Consent is already
// checked true by the handler before this is called. Source, package and
// insights opt-in are rendered so Nick knows which surface produced the
// lead and what future consent the visitor gave.
func buildLeadNotificationHTML(lead LeadInfo) string {
	packageLine := ""
	if lead.PackageID != "" {
		packageLine = fmt.Sprintf("<li><strong>Package:</strong> %s</li>\n", htmlEscape(lead.PackageID))
	}
	insights := "no"
	if lead.InsightsOptIn {
		insights = "yes"
	}
	return fmt.Sprintf(`<p><strong>New lead (source: %s) from the ivmanto.com AI service advisor funnel.</strong></p>
<ul>
<li><strong>Name:</strong> %s</li>
<li><strong>Email:</strong> <a href="mailto:%s">%s</a></li>
<li><strong>Company:</strong> %s</li>
%s</ul>
<p><strong>Request / conversation summary:</strong></p>
<p>%s</p>
<p><em>Consent given (Art. 6(1)(a)): yes. Insights opt-in: %s.</em></p>`,
		htmlEscape(lead.Source),
		htmlEscape(lead.Name),
		htmlEscape(lead.Email), htmlEscape(lead.Email),
		htmlEscape(lead.Company),
		packageLine,
		htmlEscape(lead.Message),
		insights)
}

var htmlEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&#34;",
	"'", "&#39;",
)

func htmlEscape(s string) string { return htmlEscaper.Replace(s) }

// SendBookingCancellationToClient sends a cancellation confirmation to the user.
// The visitor's IANA zone and a display label are forwarded from the booking
// handler (read off the calendar event's private properties) so the slot
// time renders in the visitor's local time, not the calendar owner's.
func (s *SmtpService) SendBookingCancellationToClient(toName, toEmail string, startTime time.Time, visitorLoc *time.Location, visitorTZLabel string) error {
	subject := "Your consultation has been cancelled"

	// Localise the rendered time to the visitor's zone. If we have no
	// visitor location (legacy event without a stored TZ, or unknown
	// IANA name), fall back to the start time's own location and the
	// empty label.
	renderedTime := startTime
	if visitorLoc != nil {
		renderedTime = startTime.In(visitorLoc)
	}
	slotLabel := renderedTime.Format("Monday, January 2, 2006 at 3:04 PM")
	if visitorTZLabel != "" {
		slotLabel = slotLabel + " " + visitorTZLabel
	}

	htmlBody := fmt.Sprintf(`
		<p>Hi %s,</p>
		<p>This is a confirmation that your consultation scheduled for <strong>%s</strong> has been successfully cancelled.</p>
		<p>If you wish to book another time, please feel free to visit our <a href="https://ivmanto.com/booking"><strong>booking page</strong></a> again.</p>
		<p>Thanks,<br>The IVMANTO Team</p>`,
		toName,
		slotLabel)

	return s.send([]string{toEmail}, nil, subject, htmlBody, nil)
}

// SendBookingCancellationToAdmin sends a notification to the admin about a client cancellation.
func (s *SmtpService) SendBookingCancellationToAdmin(clientName, clientEmail string, startTime time.Time) error {
	parts := strings.Split(s.cfg.SendFrom, "@")
	var adminEmail string
	if len(parts) == 2 {
		// Use a '+cancellation' alias to help with filtering in the admin's inbox.
		adminEmail = fmt.Sprintf("%s+cancellation@%s", parts[0], parts[1])
	} else {
		adminEmail = s.cfg.SendFrom // Fallback for non-standard emails
	}

	subject := "Consultation Cancelled by Client"
	body := fmt.Sprintf("The consultation with <strong>%s (%s)</strong> for <strong>%s</strong> has been cancelled by the client.", clientName, clientEmail, startTime.Format(time.RFC1123))
	return s.send([]string{adminEmail}, nil, subject, body, nil)
}

// SendGeneratedIdeas sends an email with the list of generated ideas.
func (s *SmtpService) SendGeneratedIdeas(toEmail, topic string, ideasBody string) error {
	subject := fmt.Sprintf("Your generated ideas for \"%s\"", topic)

	htmlBody := fmt.Sprintf(`
		<p>Hi there,</p>
		<p>As requested, here are the blog post ideas we generated for the topic "<strong>%s</strong>":</p>
		%s
		<p>If these ideas spark your interest, imagine what we could achieve with a dedicated consultation. We can help you turn these concepts into a full-fledged data strategy.</p>
		<p>Ready to take the next step? <a href="https://ivmanto.com/booking"><strong>Book a free consultation today!</strong></a></p> 
		<p>Best,<br>The IVMANTO Team</p>`,
		topic,
		ideasBody, // Now directly use the pre-formatted ideasBody
	)

	return s.send([]string{toEmail}, nil, subject, htmlBody, nil)
}
