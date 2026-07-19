package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RunFridge/wedding-card/internal/session"
)

func TestAdminAuthNoToken(t *testing.T) {
	session.Global = session.NewStore(24 * time.Hour)

	handler := AdminAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestAdminAuthValidCookie(t *testing.T) {
	session.Global = session.NewStore(24 * time.Hour)
	token, err := session.Global.Create()
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	handler := AdminAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/test", nil)
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: token})
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestAdminAuthValidBearer(t *testing.T) {
	session.Global = session.NewStore(24 * time.Hour)
	token, err := session.Global.Create()
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	handler := AdminAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestAdminAuthInvalidToken(t *testing.T) {
	session.Global = session.NewStore(24 * time.Hour)

	handler := AdminAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestAdminAuthRevokedToken(t *testing.T) {
	session.Global = session.NewStore(24 * time.Hour)
	token, err := session.Global.Create()
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}
	session.Global.Revoke(token)

	handler := AdminAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestAdminAuthCookiePrecedence(t *testing.T) {
	session.Global = session.NewStore(24 * time.Hour)
	validToken, err := session.Global.Create()
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	handler := AdminAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Valid cookie + invalid bearer: cookie should take precedence
	req := httptest.NewRequest(http.MethodGet, "/api/admin/test", nil)
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: validToken})
	req.Header.Set("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d (cookie should take precedence)", rr.Code, http.StatusOK)
	}
}

func TestAdminAuthEmptyCookieFallsToBearer(t *testing.T) {
	session.Global = session.NewStore(24 * time.Hour)
	token, err := session.Global.Create()
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	handler := AdminAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Empty cookie should fall through to bearer
	req := httptest.NewRequest(http.MethodGet, "/api/admin/test", nil)
	req.AddCookie(&http.Cookie{Name: "admin_token", Value: ""})
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}
