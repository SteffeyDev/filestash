package plg_authenticate_proxy

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	. "github.com/mickael-kerjean/filestash/server/common"
)

func init() {
	Hooks.Register.AuthenticationMiddleware("proxy", Proxy{})
}

type Proxy struct{}

var userLookup map[string]string = make(map[string]string)

func (this Proxy) Setup() Form {
	return Form{
		Elmnts: []FormElement{
			{
				Name:  "type",
				Type:  "hidden",
				Value: "Proxy",
			},
			{
				Name:        "header",
				Type:        "text",
				Placeholder: `X-Remote-User`,
				Default:     "",
				Description: "The name of the header that will contain the username",
			},
			{
				Name:        "password",
				Type:        "password",
				Placeholder: "",
				Default:     "",
				Description: `If the header is set to any value, this password will be used as well (for use with backends that accept a fixed password, such as the admin password)

This plugin exposes {{ .user }}, which will be the value from the specified header, and optionally {{ .password }}`,
			},
		},
	}
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b) // Read random bytes from crypto/rand
	if err != nil {
		return nil, err
	}
	return b, nil
}

func generateRandomString(length int) (string, error) {
	randomBytes, err := generateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil // URL-safe encoding
}

func (this Proxy) EntryPoint(idpParams map[string]string, req *http.Request, res http.ResponseWriter) error {
	username := req.Header.Get(idpParams["header"])

	if username == "" {
		res.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	key, err := generateRandomString(32)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	userLookup[key] = username

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(Page(`
      <form name="form" action="/api/session/auth/" method="post" class="component_middleware">
		<input type="hidden" name="key" value="` + key + `" />
      </form>
	  <script>
	  	document.form.submit();
	  </script>
	  `)))
	return nil
}

func (this Proxy) Callback(formData map[string]string, idpParams map[string]string, res http.ResponseWriter) (map[string]string, error) {
	user, ok := userLookup[formData["key"]]
	if !ok {
		return nil, NewError("", 400)
	}

	delete(userLookup, formData["key"])
	return map[string]string{
		"user":     user,
		"password": idpParams["password"],
	}, nil
}
