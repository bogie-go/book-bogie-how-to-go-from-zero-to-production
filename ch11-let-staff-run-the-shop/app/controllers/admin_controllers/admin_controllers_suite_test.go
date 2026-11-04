package admin_controllers_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
	"depot/test/testapp"
)

// app is the whole application on a clean test database, rebuilt before
// every spec, so each spec starts from nothing.
var app *testapp.App

// staff is the subject of a member of staff, for every request in the
// namespace.
var staff string

var _ = BeforeEach(func() {
	app = testapp.New(GinkgoT())
	user := factories.User(GinkgoT(), app.Store, func(u *models.User) { u.Role = models.RoleStaff })
	staff = fmt.Sprintf("user:%d", user.ID)
})

func TestAdmin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Admin")
}
