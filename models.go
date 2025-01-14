package main

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

const (
	NIP_47_INFO_EVENT_KIND            = 13194
	NIP_47_REQUEST_KIND               = 23194
	NIP_47_RESPONSE_KIND              = 23195
	NIP_47_PAY_INVOICE_METHOD         = "pay_invoice"
	NIP_47_GET_BALANCE_METHOD         = "get_balance"
	NIP_47_GET_INFO_METHOD            = "get_info"
	NIP_47_MAKE_INVOICE_METHOD        = "make_invoice"
	NIP_47_LOOKUP_INVOICE_METHOD      = "lookup_invoice"
	NIP_47_LIST_TRANSACTIONS_METHOD   = "list_transactions"
	NIP_47_PAY_KEYSEND_METHOD         = "pay_keysend"
	NIP_47_ERROR_INTERNAL             = "INTERNAL"
	NIP_47_ERROR_NOT_IMPLEMENTED      = "NOT_IMPLEMENTED"
	NIP_47_ERROR_QUOTA_EXCEEDED       = "QUOTA_EXCEEDED"
	NIP_47_ERROR_INSUFFICIENT_BALANCE = "INSUFFICIENT_BALANCE"
	NIP_47_ERROR_UNAUTHORIZED         = "UNAUTHORIZED"
	NIP_47_ERROR_EXPIRED              = "EXPIRED"
	NIP_47_ERROR_RESTRICTED           = "RESTRICTED"
	NIP_47_OTHER                      = "OTHER"
	NIP_47_CAPABILITIES               = "pay_invoice,pay_keysend,get_balance,get_info,make_invoice,lookup_invoice,list_transactions"
)

const (
	NOSTR_EVENT_STATE_HANDLER_EXECUTED    = "executed"
	NOSTR_EVENT_STATE_HANDLER_ERROR       = "error"
	NOSTR_EVENT_STATE_PUBLISH_CONFIRMED   = "replied"
	NOSTR_EVENT_STATE_PUBLISH_FAILED      = "failed"
	NOSTR_EVENT_STATE_PUBLISH_UNCONFIRMED = "sent"
)

var nip47MethodDescriptions = map[string]string{
	NIP_47_GET_BALANCE_METHOD:       "Read your balance",
	NIP_47_GET_INFO_METHOD:          "Read your node info",
	NIP_47_PAY_INVOICE_METHOD:       "Send payments",
	NIP_47_MAKE_INVOICE_METHOD:      "Create invoices",
	NIP_47_LOOKUP_INVOICE_METHOD:    "Lookup status of invoices",
	NIP_47_LIST_TRANSACTIONS_METHOD: "Read incoming transaction history",
}

var nip47MethodIcons = map[string]string{
	NIP_47_GET_BALANCE_METHOD:       "wallet",
	NIP_47_GET_INFO_METHOD:          "wallet",
	NIP_47_PAY_INVOICE_METHOD:       "lightning",
	NIP_47_MAKE_INVOICE_METHOD:      "invoice",
	NIP_47_LOOKUP_INVOICE_METHOD:    "search",
	NIP_47_LIST_TRANSACTIONS_METHOD: "transactions",
}

// TODO: move to models/Alby
type AlbyMe struct {
	Identifier       string `json:"identifier"`
	NPub             string `json:"nostr_pubkey"`
	LightningAddress string `json:"lightning_address"`
	Email            string `json:"email"`
}

type User struct {
	ID               uint
	AlbyIdentifier   string `validate:"required"`
	AccessToken      string `validate:"required"`
	RefreshToken     string `validate:"required"`
	Email            string
	Expiry           time.Time
	LightningAddress string
	Apps             []App
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type App struct {
	ID          uint
	UserId      uint `validate:"required"`
	User        User
	Name        string `validate:"required"`
	Description string
	NostrPubkey string `validate:"required"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AppPermission struct {
	ID            uint
	AppId         uint `validate:"required"`
	App           App
	RequestMethod string `validate:"required"`
	MaxAmount     int
	BudgetRenewal string
	ExpiresAt     time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type NostrEvent struct {
	ID        uint
	AppId     uint `validate:"required"`
	App       App
	NostrId   string `validate:"required"`
	ReplyId   string
	Content   string
	State     string
	RepliedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Payment struct {
	ID             uint
	AppId          uint `validate:"required"`
	App            App
	NostrEventId   uint `validate:"required"`
	NostrEvent     NostrEvent
	Amount         uint
	PaymentRequest string
	Preimage       *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TODO: move to models/Nip47
type Nip47Transaction struct {
	Type            string      `json:"type"`
	Invoice         string      `json:"invoice"`
	Description     string      `json:"description"`
	DescriptionHash string      `json:"description_hash"`
	Preimage        string      `json:"preimage"`
	PaymentHash     string      `json:"payment_hash"`
	Amount          int64       `json:"amount"`
	FeesPaid        int64       `json:"fees_paid"`
	CreatedAt       int64       `json:"created_at"`
	ExpiresAt       *int64      `json:"expires_at"`
	SettledAt       *int64      `json:"settled_at"`
	Metadata        interface{} `json:"metadata,omitempty"`
}

// TODO: move to models/Alby
type AlbyInvoice struct {
	Amount int64 `json:"amount"`
	// Boostagram AlbyInvoiceBoostagram        `json:"boostagram"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	// CreationDate uint64 `json:"creation_date"`
	Currency string `json:"currency"`
	// custom_records
	DescriptionHash string     `json:"description_hash"`
	ExpiresAt       *time.Time `json:"expires_at"`
	Expiry          uint32     `json:"expiry"`
	// Identifier string
	KeysendMessage string      `json:"keysend_message"`
	Memo           string      `json:"memo"`
	Metadata       interface{} `json:"metadata"`
	PayerName      string      `json:"payer_name"`
	PayerPubkey    string      `json:"payer_pubkey"`
	PaymentHash    string      `json:"payment_hash"`
	PaymentRequest string      `json:"payment_request"`
	Preimage       string      `json:"preimage"`
	// r_hash_str
	Settled   bool       `json:"settled"`
	SettledAt *time.Time `json:"settled_at"`
	State     string     `json:"state"`
	Type      string     `json:"type"`
	// value
}

type PayRequest struct {
	Invoice string `json:"invoice"`
}

// TODO: move to models/Alby
type KeysendRequest struct {
	Amount        int64             `json:"amount"`
	Destination   string            `json:"destination"`
	CustomRecords map[string]string `json:"custom_records,omitempty"`
}

type BalanceResponse struct {
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
	Unit     string `json:"unit"`
}

type PayResponse struct {
	Preimage    string `json:"payment_preimage"`
	PaymentHash string `json:"payment_hash"`
}

type MakeInvoiceRequest struct {
	Amount          int64  `json:"amount"`
	Description     string `json:"description"`
	DescriptionHash string `json:"description_hash"`
}

type MakeInvoiceResponse struct {
	Nip47Transaction
}

type LookupInvoiceResponse struct {
	Nip47Transaction
}

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// TODO: move to models/LNClient
type NodeInfo struct {
	Alias       string
	Color       string
	Pubkey      string
	Network     string
	BlockHeight uint32
	BlockHash   string
}

type Identity struct {
	gorm.Model
	Privkey string
}

// TODO: move to models/Nip47
type Nip47Request struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type Nip47Response struct {
	Error      *Nip47Error `json:"error,omitempty"`
	Result     interface{} `json:"result,omitempty"`
	ResultType string      `json:"result_type"`
}

type Nip47Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type Nip47PayParams struct {
	Invoice string `json:"invoice"`
}
type Nip47PayResponse struct {
	Preimage string `json:"preimage"`
}

type Nip47KeysendParams struct {
	Amount     int64       `json:"amount"`
	Pubkey     string      `json:"pubkey"`
	Preimage   string      `json:"preimage"`
	TLVRecords []TLVRecord `json:"tlv_records"`
}

type TLVRecord struct {
	Type  uint64 `json:"type"`
	Value string `json:"value"`
}

type Nip47BalanceResponse struct {
	Balance       int64  `json:"balance"`
	MaxAmount     int    `json:"max_amount"`
	BudgetRenewal string `json:"budget_renewal"`
}

// TODO: move to models/Nip47
type Nip47GetInfoResponse struct {
	Alias       string   `json:"alias"`
	Color       string   `json:"color"`
	Pubkey      string   `json:"pubkey"`
	Network     string   `json:"network"`
	BlockHeight uint32   `json:"block_height"`
	BlockHash   string   `json:"block_hash"`
	Methods     []string `json:"methods"`
}

type Nip47MakeInvoiceParams struct {
	Amount          int64  `json:"amount"`
	Description     string `json:"description"`
	DescriptionHash string `json:"description_hash"`
	Expiry          int64  `json:"expiry"`
}
type Nip47MakeInvoiceResponse struct {
	Nip47Transaction
}

type Nip47LookupInvoiceParams struct {
	Invoice     string `json:"invoice"`
	PaymentHash string `json:"payment_hash"`
}

type Nip47LookupInvoiceResponse struct {
	Nip47Transaction
}

type Nip47ListTransactionsParams struct {
	From   uint64 `json:"from,omitempty"`
	Until  uint64 `json:"until,omitempty"`
	Limit  uint64 `json:"limit,omitempty"`
	Offset uint64 `json:"offset,omitempty"`
	Unpaid bool   `json:"unpaid,omitempty"`
	Type   string `json:"type,omitempty"`
}

type Nip47ListTransactionsResponse struct {
	Transactions []Nip47Transaction `json:"transactions"`
}
