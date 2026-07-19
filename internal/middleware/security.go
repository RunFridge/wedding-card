package middleware

import (
	"net/http"
	"strings"
	"sync"
)

var (
	cspMu       sync.RWMutex
	cspMapExtra string
)

func SetMapCSP(embedProvider string) {
	cspMu.Lock()
	defer cspMu.Unlock()
	cspMapExtra = embedProvider
}

func SecurityHeaders() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			imgSrc := "'self' data: blob:"
			scriptSrc := "'self'"
			frameSrc := ""
			connectSrc := "'self'"

			cspMu.RLock()
			mapEmbed := cspMapExtra
			cspMu.RUnlock()

			switch mapEmbed {
			case "google":
				frameSrc = "https://www.google.com"
			case "kakao":
				scriptSrc += " https://dapi.kakao.com *.daumcdn.net"
				imgSrc += " https://dapi.kakao.com *.daumcdn.net"
				connectSrc += " https://dapi.kakao.com"
			case "naver":
				scriptSrc += " https://oapi.map.naver.com https://openapi.map.naver.com"
				imgSrc += " https://nrbe.pstatic.net https://simg.pstatic.net https://*.pstatic.net"
				connectSrc += " https://oapi.map.naver.com https://openapi.map.naver.com"
			case "tmap":
				scriptSrc += " https://apis.openapi.sk.com https://topopentile1.tmap.co.kr"
				imgSrc += " https://topopentile1.tmap.co.kr https://topopentile2.tmap.co.kr https://topopentile3.tmap.co.kr https://topopentile4.tmap.co.kr"
				connectSrc += " https://apis.openapi.sk.com"
			}

			directives := []string{
				"default-src 'none'",
				"script-src " + scriptSrc,
				"style-src 'self' 'unsafe-inline'",
				"img-src " + imgSrc,
				"font-src 'self'",
				"connect-src " + connectSrc,
				"object-src 'none'",
				"base-uri 'self'",
				"form-action 'self'",
				"frame-ancestors 'none'",
				"manifest-src 'self'",
			}
			if frameSrc != "" {
				directives = append(directives, "frame-src "+frameSrc)
			}
			csp := strings.Join(directives, "; ")

			h := w.Header()
			h.Set("Content-Security-Policy", csp)
			h.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("X-XSS-Protection", "0")
			h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
			h.Set("X-DNS-Prefetch-Control", "off")
			h.Set("X-Permitted-Cross-Domain-Policies", "none")
			h.Set("Cross-Origin-Opener-Policy", "same-origin")
			h.Set("Cross-Origin-Resource-Policy", "same-origin")
			h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
			next.ServeHTTP(w, r)
		})
	}
}
