package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/DexScen/WebBook/backend/auth/internal/domain"
	e "github.com/DexScen/WebBook/backend/auth/internal/errors"
	"github.com/gorilla/mux"
)

type Users interface {
	Register(ctx context.Context, user *domain.User) error
	LogIn(ctx context.Context, login, password string) (string, error)
	GetByLogin(ctx context.Context, login string) (*domain.User, error) // Новый метод
}

type Handler struct {
	usersService Users
}

func NewUsers(users Users) *Handler {
	return &Handler{
		usersService: users,
	}
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (h *Handler) OptionsHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(loggingMiddleware)

	links := r.PathPrefix("/users").Subrouter()
	{
		links.HandleFunc("/login", h.LogIn).Methods(http.MethodPost)
		links.HandleFunc("/register", h.Register).Methods(http.MethodPost)
		links.HandleFunc("/check-login", h.CheckLogin).Methods(http.MethodGet) // Новый маршрут
		links.HandleFunc("", h.OptionsHandler).Methods(http.MethodOptions)
	}
	return r
}

func (h *Handler) LogIn(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	var info domain.LoginInfo
	var roleInfo domain.RoleInfo

	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Login error:", err)
		return
	}

	role, err := h.usersService.LogIn(context.TODO(), info.Login, info.Password)
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			role = "unauthorized by user"
			log.Println("Login error2:", err)
		} else if errors.Is(err, e.ErrWrongPassword) {
			role = "unauthorized by password"
			log.Println("Login error3:", err)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Login error1:", err, e.ErrUserNotFound)
			return
		}
	}
	roleInfo.Role = role
	if jsonResp, err := json.Marshal(roleInfo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Login error:", err)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	var user domain.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Register error:", err)
		return
	}

	if err := h.usersService.Register(context.TODO(), &user); err != nil {
		if errors.Is(err, e.ErrUserExists) { // user existed
			if jsonResp, err := json.Marshal("user exists"); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println("Register error:", err)
				return
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonResp)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Register error:", err)
			return
		}
	} else { // user registered success
		if jsonResp, err := json.Marshal("success"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Register error:", err)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResp)
		}
	}
}

// CheckLogin проверяет, доступен ли указанный логин для регистрации
func (h *Handler) CheckLogin(w http.ResponseWriter, r *http.Request) {
    setHeaders(w)
    
    // Получаем логин из query параметров
    login := r.URL.Query().Get("login")
    if login == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "available": false,
            "error": "login parameter is required",
        })
        return
    }

    // Проверяем существование пользователя
    _, err := h.usersService.GetByLogin(r.Context(), login)
    if err != nil {
        if errors.Is(err, e.ErrUserNotFound) {
            // Пользователь не найден - логин доступен
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(map[string]interface{}{
                "available": true,
            })
            return
        }
        // Другая ошибка
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "available": false,
            "error": "internal server error",
        })
        return
    }

    // Пользователь найден - логин занят
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "available": false,
    })
}
