package uaa_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/pivotalservices/gtils/uaa"
)

var _ = Describe("GetToken()", func() {
	Context("with grant type password", func() {
		var token string
		var err error
		var controlToken = "12345937635"
		var server *ghttp.Server

		BeforeEach(func() {
			expectedFormParams := url.Values{
				"grant_type":    {"password"},
				"response_type": {"token"},
				"username":      {"fakeuser"},
				"password":      {"fakepass"},
				"client_id":     {"opsman"},
				"client_secret": {""},
			}
			server = NewTestServer(ghttp.NewTLSServer(), controlToken, expectedFormParams)
			token, err = uaa.GetToken(server.URL()+"/uaa", "fakeuser", "fakepass", "opsman", "", "password")
		})

		AfterEach(func() {
			server.Close()
		})

		It("Then it should return a valid token", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(token).ShouldNot(BeEmpty())
			Ω(token).Should(Equal(controlToken))
		})
	})

	Context("with grant type client_credentials", func() {
		var token string
		var err error
		var controlToken = "12345937635"
		var server *ghttp.Server

		BeforeEach(func() {
			expectedFormParams := url.Values{
				"grant_type":    {"client_credentials"},
				"response_type": {"token"},
				"username":      {""},
				"password":      {""},
				"client_id":     {"client456"},
				"client_secret": {"clientsecret"},
			}
			server = NewTestServer(ghttp.NewTLSServer(), controlToken, expectedFormParams)
			token, err = uaa.GetToken(server.URL()+"/uaa", "", "", "client456", "clientsecret", "client_credentials")
		})

		AfterEach(func() {
			server.Close()
		})

		It("Then it should return a valid token", func() {
			Ω(err).ShouldNot(HaveOccurred())
			Ω(token).ShouldNot(BeEmpty())
			Ω(token).Should(Equal(controlToken))
		})
	})

	Context("when call to uaa target yields error", func() {
		var token string
		var err error
		var server *ghttp.Server

		BeforeEach(func() {
			server = NewErrorTestServer(ghttp.NewTLSServer())
			token, err = uaa.GetToken(server.URL()+"/uaa", "fakeuser", "fakepass", "opsman", "", "password")
		})

		AfterEach(func() {
			server.Close()
		})

		It("Then it should return the response body as an error message", func() {
			Ω(err).Should(HaveOccurred())
			Ω(token).Should(BeEmpty())
		})
	})
})

func NewErrorTestServer(server *ghttp.Server) *ghttp.Server {
	errTokenHandler := ghttp.RespondWith(http.StatusUnauthorized, "{error:somefailure}")
	server.AppendHandlers(
		errTokenHandler,
	)
	return server
}

func NewTestServer(server *ghttp.Server, token string, expectedParams url.Values) *ghttp.Server {
	tokenJson := getFakeToken("./fixtures/token_response.json", token, "", "")

	successTokenHandler := ghttp.RespondWith(http.StatusOK, tokenJson)
	successFormHandler := ghttp.VerifyForm(expectedParams)

	server.AppendHandlers(
		ghttp.CombineHandlers(
			successFormHandler,
			successTokenHandler,
		),
	)
	return server
}

func getFakeToken(fixturePath, token, refresh, jti string) string {
	b, _ := ioutil.ReadFile(fixturePath)
	return fmt.Sprintf(string(b), token, refresh, jti)
}
