package utils

// utils pour le parsing des erreur HTTP

type HttpErrorsInterface interface {
	Error() string
	GetStatus() int
}

type HttpErrors struct {
	Status  int
	Message string
}

func (h *HttpErrors) Error() string {
	return h.Message
}

func (h *HttpErrors) GetStatus() int {
	return h.Status
}
