package email

import "time"

// Service defines the interface for sending all types of application emails.
type Service interface {
	SendContactMessage(msg ContactMessage) error
	SendBookingConfirmation(details BookingConfirmationDetails) error
	SendBookingNotificationToAdmin(name, clientEmail string, startTime time.Time, notes string) error
	// SendBookingCancellationToClient renders the cancellation email using the
	// visitor's timezone (visitorLoc + visitorTZLabel) when available, so the
	// slot time is in the visitor's local zone rather than the calendar owner's.
	SendBookingCancellationToClient(toName, toEmail string, startTime time.Time, visitorLoc *time.Location, visitorTZLabel string) error
	SendBookingCancellationToAdmin(clientName, clientEmail string, startTime time.Time) error
	SendGeneratedIdeas(toEmail, topic string, ideasBody string) error
	// SendLeadNotification emails the advisor lead to the configured lead
	// recipient (LEADS_TO). No storage happens in the handler.
	SendLeadNotification(lead LeadInfo) error
}
