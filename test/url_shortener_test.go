package test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"url_shortener/internal/http-server/handlers/save"
	"url_shortener/internal/lib/random"
)

const (
	host = "localhost:8082"
)

func Test_url_shortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	alias := random.NewRandomString(6)
	originalURL := gofakeit.URL()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  u.String(),
		Client:   client,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	// 2. SAVE
	e.POST("/url").
		WithJSON(save.Request{
			URL:   originalURL,
			Alias: alias,
		}).
		WithBasicAuth("user", "pass").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("alias")

	// 3. REDIRECT
	e.GET("/" + alias).
		Expect().
		Status(302).
		Header("Location").IsEqual(originalURL)

	// 4. DELETE
	reqDel := e.DELETE("/url/"+alias).
		WithBasicAuth("user", "pass").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("status")

	reqDel.Value("status").String().IsEqual("OK")
}
