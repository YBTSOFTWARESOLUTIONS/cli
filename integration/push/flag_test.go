package push

import (
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/cli/integration/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("flag combinations", func() {
	var (
		appName string
	)

	BeforeEach(func() {
		appName = helpers.NewAppName()
	})

	Context("when the --docker-username is provided without an image", func() {
		It("errors with usage", func() {
			session := helpers.CF(PushCommandName, "--docker-username", "some-docker-username", appName)
			Eventually(session.Err).Should(Say("Incorrect Usage: '--docker-image, -o' and '--docker-username' must be used together."))
			Eventually(session).Should(Say("NAME:"))
			Eventually(session).Should(Say("Push a new app or sync changes to an existing app"))
			Eventually(session).Should(Exit(1))
		})
	})

	Context("when the -p and -o flags are used together", func() {
		var path string

		BeforeEach(func() {
			tempFile, err := ioutil.TempFile("", "integration-push-path")
			Expect(err).ToNot(HaveOccurred())
			path = tempFile.Name()
			Expect(tempFile.Close())
		})

		AfterEach(func() {
			err := os.Remove(path)
			Expect(err).ToNot(HaveOccurred())
		})

		It("tells the user that they cannot be used together, displays usage and fails", func() {
			session := helpers.CF(PushCommandName, appName, "-o", PublicDockerImage, "-p", path)

			Eventually(session.Err).Should(Say("Incorrect Usage: The following arguments cannot be used together: --docker-image, -o, -p"))
			Eventually(session).Should(Say("FAILED"))
			Eventually(session).Should(Say("USAGE:"))

			Eventually(session).Should(Exit(1))
		})
	})

	Context("when the -b and -o flags are used together", func() {
		It("tells the user that they cannot be used together, displays usage and fails", func() {
			session := helpers.CF(PushCommandName, appName, "-o", PublicDockerImage, "-b", "some-buildpack")

			Consistently(session.Out).ShouldNot(Say("Creating app"))
			Eventually(session.Err).Should(Say("Incorrect Usage: The following arguments cannot be used together: -b, --docker-image, -o"))
			Eventually(session).Should(Say("FAILED"))
			Eventually(session).Should(Say("USAGE:"))

			Eventually(session).Should(Exit(1))
		})
	})

	Context("when the --no-hostname and --no-route flags are used together", func() {
		It("tells the user that they cannot be used together, displays usage and fails", func() {
			session := helpers.CF(PushCommandName, appName, "--no-hostname", "--no-route")

			Eventually(session.Err).Should(Say("Incorrect Usage: The following arguments cannot be used together: --no-hostname, --no-route"))
			Eventually(session).Should(Say("FAILED"))
			Eventually(session).Should(Say("USAGE:"))

			Eventually(session).Should(Exit(1))
		})
	})
})
