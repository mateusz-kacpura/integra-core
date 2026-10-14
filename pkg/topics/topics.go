package topics

// Topic represents a canonical Kafka topic name in the Integra ecosystem.
type Topic string

// String returns the string representation of the topic.
func (t Topic) String() string {
	return string(t)
}

const (
	// UserProfilesSync is the topic for synchronization of user profile snapshots across services (Key: user_id).
	UserProfilesSync = "user-profiles-sync"

	// UserRolesSync is the topic for IAM user role assignments synchronization across all services (Key: user_id).
	UserRolesSync = "user-roles-sync"

	// UserEvents is the topic for account registration and lifecycle events (Key: user_id).
	UserEvents = "user-events"

	// UserSkillsSync is the topic for user skills updates consumed by Job Offers and RDC (Key: user_id).
	UserSkillsSync = "user-skills-sync"

	// LicenseEvents is the topic for IAM product license events (Key: org_id).
	LicenseEvents = "license-events"

	// OrgEvents is the topic for organization lifecycle events consumed by all services (Key: org_id).
	OrgEvents = "org-events"

	// OrgMembershipEvents is the topic for organization membership changes (Key: org_id).
	OrgMembershipEvents = "org-membership-events"

	// OrgRolesSync is the topic for organization role assignment synchronization (Key: org_id).
	OrgRolesSync = "org-roles-sync"

	// RDCDomainEvents is the topic for RDC Club domain events (Key: aggregate_id).
	RDCDomainEvents = "rdc-domain-events"

	// RDCTapEvents is the topic for RDC member NFC card tap events (Key: member_id).
	RDCTapEvents = "rdc-tap-events"

	// RDCLoyaltyEvents is the topic for RDC loyalty point and reward events (Key: member_id).
	RDCLoyaltyEvents = "rdc-loyalty-events"

	// RDCAnalyticsEvents is the topic for RDC analytics metrics and tracking events (Key: aggregate_id).
	RDCAnalyticsEvents = "rdc-analytics-events"

	// RDCNotifications is the topic for RDC member notification triggers (Key: user_id).
	RDCNotifications = "rdc-notifications"

	// JobOffersEvents is the topic for job postings lifecycle events (Key: offer_id).
	JobOffersEvents = "job-offers-events"

	// JobApplicationsEvents is the topic for candidate job application events (Key: application_id).
	JobApplicationsEvents = "job-applications-events"

	// PartnerNotifications is the topic for notifications destined for business partners (Key: partner_id).
	PartnerNotifications = "partner-notifications"

	// CourseEvents is the topic for course creation and modification events (Key: course_id).
	CourseEvents = "course-events"

	// CourseCompleted is the topic emitted when a user finishes a course (Key: user_id).
	CourseCompleted = "course-completed"

	// CourseProgressEvents is the topic for student learning progress tracking (Key: user_id).
	CourseProgressEvents = "course-progress-events"

	// ContractsEvents is the topic for contracts and employment document events (Key: contract_id).
	ContractsEvents = "contracts-events"

	// TimesheetEvents is the topic for timesheet submissions and approvals (Key: contract_id).
	TimesheetEvents = "timesheet-events"

	// DisputeEvents is the topic for contract dispute events (Key: dispute_id).
	DisputeEvents = "dispute-events"

	// BillingPaymentSuccess is the topic for successful payment transactions (Key: tx_id).
	BillingPaymentSuccess = "billing-payment-success"

	// SubscriptionEvents is the topic for subscription lifecycle and renewal events (Key: user_id).
	SubscriptionEvents = "subscription-events"

	// WalletEvents is the topic for user and organization wallet balance events (Key: wallet_id).
	WalletEvents = "wallet-events"

	// CRMDealClosed is the topic emitted when a CRM deal is won or lost (Key: deal_id).
	CRMDealClosed = "crm-deal-closed"

	// CRMEvents is the topic for customer relationship management lead events (Key: lead_id).
	CRMEvents = "crm-events"

	// NotificationDelivered is the topic emitted when a notification is successfully delivered (Key: notification_id).
	NotificationDelivered = "notification-delivered"

	// SupportTicketEvents is the topic for customer support ticket lifecycle events (Key: ticket_id).
	SupportTicketEvents = "support-ticket-events"

	// FeatureIdeaEvents is the topic for user product feature suggestions (Key: idea_id).
	FeatureIdeaEvents = "feature-idea-events"

	// ProductEvents is the topic for product and marketplace catalog events (Key: product_id).
	ProductEvents = "product-events"

	// CalendarEvents is the topic for calendar schedules and booking events (Key: event_id).
	CalendarEvents = "calendar-events"

	// KRSEvents is the topic for National Court Register (KRS) updates (Key: nip).
	KRSEvents = "krs-events"

	// LedgerUpdated is the topic for double-entry accounting ledger balance changes (Key: wallet_id).
	LedgerUpdated = "ledger-updated"

	// PayoutEvents is the topic for withdrawal and payout processing events (Key: request_id).
	PayoutEvents = "payout-events"

	// AccountingInvoiceIssued is the topic emitted when an official invoice is generated (Key: invoice_id).
	AccountingInvoiceIssued = "accounting-invoice-issued"

	// AdminEvents is the topic for administration and moderation actions (Key: admin_id).
	AdminEvents = "admin-events"

	// SecurityAuditEvents is the topic for immutable, append-only security audit log events (Key: event_id).
	SecurityAuditEvents = "security-audit-events"
)

// AllTopics returns a slice of all 38 registered topic names.
func AllTopics() []string {
	return []string{
		UserProfilesSync, UserRolesSync, UserEvents, UserSkillsSync, LicenseEvents,
		OrgEvents, OrgMembershipEvents, OrgRolesSync,
		RDCDomainEvents, RDCTapEvents, RDCLoyaltyEvents, RDCAnalyticsEvents, RDCNotifications,
		JobOffersEvents, JobApplicationsEvents, PartnerNotifications,
		CourseEvents, CourseCompleted, CourseProgressEvents,
		ContractsEvents, TimesheetEvents, DisputeEvents,
		BillingPaymentSuccess, SubscriptionEvents, WalletEvents,
		CRMDealClosed, CRMEvents,
		NotificationDelivered,
		SupportTicketEvents, FeatureIdeaEvents,
		ProductEvents,
		CalendarEvents,
		KRSEvents,
		LedgerUpdated, PayoutEvents, AccountingInvoiceIssued,
		AdminEvents,
		SecurityAuditEvents,
	}
}

// TopicForAggregate maps entity/aggregate types to canonical Kafka topics in the ecosystem.
func TopicForAggregate(aggregateType string) string {
	switch aggregateType {
	case "session", "user":
		return UserEvents
	case "role":
		return UserRolesSync
	case "license":
		return LicenseEvents
	default:
		return UserEvents
	}
}
