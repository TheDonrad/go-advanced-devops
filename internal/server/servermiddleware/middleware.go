package servermiddleware

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"goAdvancedTpl/internal/fabric/encryption"
	"goAdvancedTpl/internal/fabric/logs"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

// GzipHandle сжимает ответ в случе наличия заголовка "Accept-Encoding": "gzip"
func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil && !errors.Is(err, io.EOF) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			err = gz.Close()
			if err != nil && err != io.EOF {
				log.Print(err.Error())
			}
		}()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func Decryption(key string) func(next http.Handler) http.Handler {

	decoder := Decoder{
		Key: key,
	}

	return decoder.Handler
}

type Decoder struct {
	Key string
}

func (d *Decoder) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if d.Key == "" {
			next.ServeHTTP(w, r)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logs.Logger().Println(err.Error())
			next.ServeHTTP(w, r)
		}

		if b, err := encryption.Decrypt(d.Key, body); err == nil {
			r.Body = io.NopCloser(bytes.NewReader(b))
		}

		next.ServeHTTP(w, r)

	})
}

func CheckIP(trustedSubnet string) func(next http.Handler) http.Handler {

	ipChecker := ipChecker{
		trustedSubnet: trustedSubnet,
	}

	return ipChecker.Handler
}

type ipChecker struct {
	trustedSubnet string
}

func (i *ipChecker) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if i.trustedSubnet == "" {
			next.ServeHTTP(w, r)
			return
		}

		_, ipNet, err := net.ParseCIDR(i.trustedSubnet)
		if err != nil {
			logs.Logger().Println(err.Error())
			next.ServeHTTP(w, r)
			return
		}

		ip := net.ParseIP(r.Header.Get("X-Real-IP"))
		if ip == nil || ipNet.Contains(ip) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)

	})
}
