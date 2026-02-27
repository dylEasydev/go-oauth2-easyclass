package service

import "errors"

var (
	ErrNotCode = errors.New("mauvais code de verification !")
	ErrDestroy = errors.New("impossible de supprimer l'utilisateur temporaire")
)
