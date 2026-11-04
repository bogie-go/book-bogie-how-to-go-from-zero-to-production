package testapp

import "time"

// As signs a request test in as that token subject: whatever the app's
// sessions mint, usually a user's id as a string. The token is minted with
// the app's own secret, so the real RequireUser accepts it.
//
//	app := testapp.New(GinkgoT())
//	app.Get("/posts").As(strconv.FormatInt(user.ID, 10)).Do().Status(http.StatusOK)
//
// Leave As off to test the 401.
func (r *Request) As(subject string) *Request {
	signed, err := r.app.Server.Tokens().Mint(subject, time.Hour)
	if err != nil {
		r.app.t.Fatal(err)
	}
	r.authorization = "Bearer " + signed
	return r
}
