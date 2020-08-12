package cookies

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

//SetCookie sets cookie(session id here) for a given method
func SetCookie(session string, method string) *http.Client {
	jar, _ := cookiejar.New(nil)

	var cookies []*http.Cookie

	firstCookie := &http.Cookie{
		Name:   "session",
		Value:  session,
		Path:   "/",
		Domain: ".irccloud.com",
	}

	cookies = append(cookies, firstCookie)

	// URL for cookies to remember. i.e reply when encounter this URL
	cookieURL, _ := url.Parse("https://www.irccloud.com/chat/" + method)

	jar.SetCookies(cookieURL, cookies)
	client := &http.Client{
		Jar: jar,
	}

	return client

}
