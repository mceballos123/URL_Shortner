package api
import (
	"time"
	"net/http"
	"encoding/json"
	"log"
	"github.com/gin-gonic/gin"
)

type CreateUrlRequest struct{
	Alias string `json:"alias"`
	Link string `json:"link"`
	ExpiredAt *time.Time `json:"expired_at"`

}

type CreateUserRequest struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

