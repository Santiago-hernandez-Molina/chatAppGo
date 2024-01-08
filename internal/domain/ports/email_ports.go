package ports

import "github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"

type EmailSender interface {
    SendRegisterConfirm(user *models.User, code int) error
}
