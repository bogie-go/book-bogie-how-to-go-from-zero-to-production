package order_confirmation_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOrderConfirmation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OrderConfirmation")
}
