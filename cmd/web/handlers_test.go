package main

import (
	"bytes"
	"net/http"
	"net/url"
	"regexp"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestShowSnippet(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func TestCreateSnippet(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)
	ts.LoginForTest(t, csrfToken)

	longTitle := ""
	for i := 0; i < 102; {
		longTitle = longTitle + "1"
		i++
	}

	createdRegex, err := regexp.Compile("/snippet/[0-9]+$")
	if err != nil {
		t.Errorf("regex didn't compile")
	}

	tests := []struct {
		name       string
		title      string
		content    string
		expires    string
		csrfToken  string
		wantCode   int
		wantBody   []byte
		wantHeader bool
	}{
		{"Valid submission (one day expire)", "validTitle", "validContent", "1", csrfToken, http.StatusSeeOther, []byte(""), true},
		{"Valid submission (one week expire)", "validTitle", "validContent", "7", csrfToken, http.StatusSeeOther, []byte(""), true},
		{"Valid submission (one year expire)", "validTitle", "validContent", "365", csrfToken, http.StatusSeeOther, []byte(""), true},
		{"Empty title", "", "validContent", "1", csrfToken, http.StatusOK, []byte("This field cannot be blank"), false},
		{"Long Title", longTitle, "validContent", "1", csrfToken, http.StatusOK, []byte("This field is too long (maximum is 100 characters)"), false},
		{"Empty Content", "validTitle", "", "1", csrfToken, http.StatusOK, []byte("This field cannot be blank"), false},
		{"Empty Expire", "validTitle", "validContent", "", csrfToken, http.StatusOK, []byte("This field cannot be blank"), false},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("title", tt.title)
			form.Add("content", tt.content)
			form.Add("expires", tt.expires)
			form.Add("csrf_token", tt.csrfToken)
			form.Add("contextKeyIsAuthenticated", "true")

			code, header, body := ts.postForm(t, "/snippet/create", form)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}

			if tt.wantHeader {
				if !createdRegex.MatchString(header.Get("location")) {
					t.Errorf("want header %s to contain %q", header.Get("location"), createdRegex.String())
				}
			}
		})
	}
}

func TestSignupUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid submission", "Bob", "bob@example.com", "validPa$$word", csrfToken, http.StatusSeeOther, nil},
		{"Empty name", "", "bob@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty password", "Bob", "bob@example.com", "", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Short password", "Bob", "bob@example.com", "pa$$word", csrfToken, http.StatusOK, []byte("This field is too short (minimum is 10 characters)")},
		{"Duplicate email", "Bob", "dupe@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("Address is already in use")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name         string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid Credentials", "alice@example.com", "my plain text password", csrfToken, http.StatusSeeOther, nil},
		{"Empty Email", "", "validPassword", csrfToken, http.StatusOK, []byte("Email or Password is incorrect")},
		{"Empty Password", "valid@email.com", "", csrfToken, http.StatusOK, []byte("Email or Password is incorrect")},
		{"Empty Fields", "", "", csrfToken, http.StatusOK, []byte("Email or Password is incorrect")},
		{"Wrong Credentials", "wrongEmail", "wrongPassword", csrfToken, http.StatusOK, []byte("Email or Password is incorrect")},
		{"Invalid CSRF Token", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/login", form)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/login")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name      string
		csrfToken string
		login     bool
		wantCode  int
	}{
		{"User logged in", csrfToken, true, http.StatusSeeOther},
		{"User not logged in", csrfToken, false, http.StatusSeeOther},
		{"Invalid CSRF Token", "wrongToken", false, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.login {
				_, _, body := ts.get(t, "/user/login")
				logincsrfToken := extractCSRFToken(t, body)
				ts.LoginForTest(t, logincsrfToken)
			}

			form := url.Values{}
			form.Add("csrf_token", tt.csrfToken)

			code, _, _ := ts.postForm(t, "/user/logout", form)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}
		})
	}
}
