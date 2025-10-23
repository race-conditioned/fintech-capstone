package contracts

// TransferResponse is the standard response body for a transfer.
// It is the Result returned from the usecase
type TransferResponse struct {
	TransactionID string `json:"transaction_id"` // uuid
	Status        string `json:"status"`         // success | rejected | rate_limited | duplicate
	Message       string `json:"message,omitempty"`
}
