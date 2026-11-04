package requests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"depot/app/models"
	"depot/test/factories"
)

// contract is the API's OpenAPI document, loaded once for every test.
var contract routers.Router

var _ = BeforeEach(func() {
	if contract != nil {
		return
	}
	doc, err := openapi3.NewLoader().LoadFromFile("../../docs/openapi.yaml")
	Expect(err).NotTo(HaveOccurred())
	Expect(doc.Validate(context.Background())).To(Succeed())
	contract, err = gorillamux.NewRouter(doc)
	Expect(err).NotTo(HaveOccurred())
})

// conforms sends req through the app and fails the test unless the response
// is one the contract describes for that route.
func conforms(req *http.Request) *httptest.ResponseRecorder {
	GinkgoHelper()
	ctx := req.Context()
	route, params, err := contract.FindRoute(req)
	Expect(err).NotTo(HaveOccurred(), "%s %s is not in the contract", req.Method, req.URL.Path)

	in := &openapi3filter.RequestValidationInput{Request: req, PathParams: params, Route: route}

	res := httptest.NewRecorder()
	app.Server.Router.ServeHTTP(res, req)

	err = openapi3filter.ValidateResponse(ctx, &openapi3filter.ResponseValidationInput{
		RequestValidationInput: in,
		Status:                 res.Code,
		Header:                 res.Header(),
		Body:                   io.NopCloser(bytes.NewReader(res.Body.Bytes())),
	})
	Expect(err).NotTo(HaveOccurred(), "the response broke the contract:\n%s", res.Body.String())
	return res
}

// signedIn adds a bearer token for subject, minted with the app's own secret.
func signedIn(req *http.Request, subject string) *http.Request {
	token, err := app.Server.Tokens().Mint(subject, time.Hour)
	Expect(err).NotTo(HaveOccurred())
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

// request builds a request with a JSON body.
func request(method, path, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

var _ = Describe("The contract", func() {
	It("holds for the catalog", func() {
		mug := factories.Product(GinkgoT(), app.Store)

		Expect(conforms(httptest.NewRequest(http.MethodGet, "/products", nil)).Code).To(Equal(http.StatusOK))
		Expect(conforms(httptest.NewRequest(http.MethodGet, "/products?per_page=1000", nil)).Code).To(Equal(http.StatusBadRequest))
		Expect(conforms(httptest.NewRequest(http.MethodGet, fmt.Sprintf("/products/%d", mug.ID), nil)).Code).To(Equal(http.StatusOK))
		Expect(conforms(httptest.NewRequest(http.MethodGet, "/products/999999999", nil)).Code).To(Equal(http.StatusNotFound))
	})

	It("holds for the categories", func() {
		stickers := factories.Category(GinkgoT(), app.Store)
		mug := factories.Product(GinkgoT(), app.Store)
		Expect(app.Store.SetProductCategories(GinkgoT().Context(), fmt.Sprint(mug.ID), []int64{stickers.ID})).To(Succeed())

		Expect(conforms(httptest.NewRequest(http.MethodGet, "/categories", nil)).Code).To(Equal(http.StatusOK))
		Expect(conforms(httptest.NewRequest(http.MethodGet, fmt.Sprintf("/products/%d", mug.ID), nil)).Code).To(Equal(http.StatusOK))
		Expect(conforms(httptest.NewRequest(http.MethodGet, fmt.Sprintf("/products?category_id=%d", stickers.ID), nil)).Code).To(Equal(http.StatusOK))
	})

	It("holds for a guest's cart", func() {
		mug := factories.Product(GinkgoT(), app.Store, func(p *models.Product) { p.Stock = 5 })
		body := strings.NewReader(fmt.Sprintf(`{"product_id": %d, "quantity": 2}`, mug.ID))
		add := signedIn(httptest.NewRequest(http.MethodPost, "/cart_items", body), "guest:1")
		add.Header.Set("Content-Type", "application/json")

		Expect(conforms(add).Code).To(Equal(http.StatusOK))
		Expect(conforms(signedIn(httptest.NewRequest(http.MethodGet, "/cart", nil), "guest:1")).Code).To(Equal(http.StatusOK))
		Expect(conforms(httptest.NewRequest(http.MethodGet, "/cart", nil)).Code).To(Equal(http.StatusUnauthorized))
	})

	It("holds for guest sessions, sign up and sign in", func() {
		Expect(conforms(httptest.NewRequest(http.MethodPost, "/guest_sessions", nil)).Code).To(Equal(http.StatusCreated))

		Expect(conforms(request(http.MethodPost, "/users", `{"email": "ada@example.com", "password": "correct-horse"}`)).Code).To(Equal(http.StatusCreated))
		Expect(conforms(request(http.MethodPost, "/users", `{"email":`)).Code).To(Equal(http.StatusBadRequest))
		Expect(conforms(request(http.MethodPost, "/users", `{"email": "not an email", "password": "x"}`)).Code).To(Equal(http.StatusUnprocessableEntity))

		Expect(conforms(request(http.MethodPost, "/sessions", `{"email": "ada@example.com", "password": "correct-horse"}`)).Code).To(Equal(http.StatusOK))
		Expect(conforms(request(http.MethodPost, "/sessions", `{"email":`)).Code).To(Equal(http.StatusBadRequest))
		Expect(conforms(request(http.MethodPost, "/sessions", `{"email": "ada@example.com", "password": "wrong"}`)).Code).To(Equal(http.StatusUnauthorized))
	})

	It("holds for the cart's items", func() {
		mug := factories.Product(GinkgoT(), app.Store, func(p *models.Product) { p.Stock = 5 })
		as := func(req *http.Request) *http.Request { return signedIn(req, "guest:2") }

		res := conforms(as(request(http.MethodPost, "/cart_items", fmt.Sprintf(`{"product_id": %d, "quantity": 1}`, mug.ID))))
		Expect(res.Code).To(Equal(http.StatusOK))
		Expect(conforms(request(http.MethodPost, "/cart_items", `{}`)).Code).To(Equal(http.StatusUnauthorized))
		Expect(conforms(as(request(http.MethodPost, "/cart_items", `{"product_id":`))).Code).To(Equal(http.StatusBadRequest))
		Expect(conforms(as(request(http.MethodPost, "/cart_items", `{"product_id": 999999999, "quantity": 1}`))).Code).To(Equal(http.StatusNotFound))
		Expect(conforms(as(request(http.MethodPost, "/cart_items", fmt.Sprintf(`{"product_id": %d, "quantity": 99}`, mug.ID)))).Code).To(Equal(http.StatusUnprocessableEntity))

		cart, err := app.Store.FindOrCreateCartByToken(GinkgoT().Context(), "guest:2")
		Expect(err).NotTo(HaveOccurred())
		lines, err := app.Store.CartContents(GinkgoT().Context(), cart.ID)
		Expect(err).NotTo(HaveOccurred())
		item := fmt.Sprintf("/cart_items/%d", lines[0].ItemID)

		Expect(conforms(as(request(http.MethodPatch, item, `{"quantity": 3}`))).Code).To(Equal(http.StatusOK))
		Expect(conforms(request(http.MethodPatch, item, `{"quantity": 3}`)).Code).To(Equal(http.StatusUnauthorized))
		Expect(conforms(as(request(http.MethodPatch, item, `{"quantity":`))).Code).To(Equal(http.StatusBadRequest))
		Expect(conforms(as(request(http.MethodPatch, "/cart_items/999999999", `{"quantity": 1}`))).Code).To(Equal(http.StatusNotFound))
		Expect(conforms(as(request(http.MethodPatch, item, `{"quantity": 99}`))).Code).To(Equal(http.StatusUnprocessableEntity))

		Expect(conforms(request(http.MethodDelete, item, "")).Code).To(Equal(http.StatusUnauthorized))
		Expect(conforms(as(request(http.MethodDelete, "/cart_items/999999999", ""))).Code).To(Equal(http.StatusNotFound))
		Expect(conforms(as(request(http.MethodDelete, item, ""))).Code).To(Equal(http.StatusOK))
	})

	It("holds for orders", func() {
		mug := factories.Product(GinkgoT(), app.Store, func(p *models.Product) { p.Stock = 5 })
		user := factories.User(GinkgoT(), app.Store)
		subject := fmt.Sprintf("user:%d", user.ID)
		cart, err := app.Store.FindOrCreateCartByToken(GinkgoT().Context(), subject)
		Expect(err).NotTo(HaveOccurred())

		Expect(conforms(signedIn(request(http.MethodPost, "/orders", `{}`), subject)).Code).To(Equal(http.StatusUnprocessableEntity)) // empty cart
		Expect(app.Store.AddCartItem(GinkgoT().Context(), cart.ID, mug.ID, 2)).To(Succeed())
		Expect(conforms(request(http.MethodPost, "/orders", `{}`)).Code).To(Equal(http.StatusUnauthorized))
		Expect(conforms(signedIn(request(http.MethodPost, "/orders", `{}`), "guest:3")).Code).To(Equal(http.StatusForbidden))
		Expect(conforms(signedIn(request(http.MethodPost, "/orders", `{"coupon_code":`), subject)).Code).To(Equal(http.StatusBadRequest))
		Expect(conforms(signedIn(request(http.MethodPost, "/orders", `{"coupon_code": "NOPE"}`), subject)).Code).To(Equal(http.StatusUnprocessableEntity))

		res := conforms(signedIn(request(http.MethodPost, "/orders", `{}`), subject))
		Expect(res.Code).To(Equal(http.StatusCreated))
		var placed struct {
			ID int64 `json:"id"`
		}
		Expect(json.Unmarshal(res.Body.Bytes(), &placed)).To(Succeed())

		path := fmt.Sprintf("/orders/%d", placed.ID)
		Expect(conforms(signedIn(httptest.NewRequest(http.MethodGet, path, nil), subject)).Code).To(Equal(http.StatusOK))
		Expect(conforms(httptest.NewRequest(http.MethodGet, path, nil)).Code).To(Equal(http.StatusUnauthorized))
		Expect(conforms(signedIn(httptest.NewRequest(http.MethodGet, "/orders/999999999", nil), subject)).Code).To(Equal(http.StatusNotFound))
	})
})

var _ = Describe("The contract's routes", func() {
	It("are exactly the public routes the app serves", func() {
		doc, err := openapi3.NewLoader().LoadFromFile("../../docs/openapi.yaml")
		Expect(err).NotTo(HaveOccurred())

		documented := []string{}
		for path, item := range doc.Paths.Map() {
			for method := range item.Operations() {
				documented = append(documented, method+" "+path)
			}
		}

		served := []string{}
		for _, r := range app.Server.Router.Routes() {
			if strings.HasPrefix(r.Path, "/admin") || r.Path == "/healthz" || r.Path == "/status" {
				continue // staff and operations: not the public API
			}
			served = append(served, r.Method+" "+openAPIPath(r.Path))
		}

		Expect(served).To(ConsistOf(documented))
	})
})

// openAPIPath writes a Gin path the way OpenAPI does: /products/:id is
// /products/{id}.
func openAPIPath(ginPath string) string {
	parts := strings.Split(ginPath, "/")
	for i, p := range parts {
		if name, ok := strings.CutPrefix(p, ":"); ok {
			parts[i] = "{" + name + "}"
		}
	}
	return strings.Join(parts, "/")
}
