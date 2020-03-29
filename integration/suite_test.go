package integration_test

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var imap []uuid.UUID

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	data, err := ioutil.ReadFile("fixture/guid.txt")
	Expect(err).NotTo(HaveOccurred())

	scanner := bufio.NewScanner(bytes.NewBuffer(data))

	for scanner.Scan() {
		id, err := uuid.Parse(scanner.Text())
		Expect(err).NotTo(HaveOccurred())
		imap = append(imap, id)
	}
})
