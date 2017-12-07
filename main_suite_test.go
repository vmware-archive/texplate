package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTexplateCli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TexplateCli Suite")
}
