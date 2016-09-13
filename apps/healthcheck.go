package apps

import (
	. "github.com/cloudfoundry/cf-acceptance-tests/cats_suite_helpers"
	"github.com/cloudfoundry/cf-acceptance-tests/helpers/app_helpers"
	"github.com/cloudfoundry/cf-acceptance-tests/helpers/assets"
	"github.com/cloudfoundry/cf-acceptance-tests/helpers/skip_messages"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers"
	"github.com/cloudfoundry/cf-acceptance-tests/helpers/random_name"
)

var _ = AppsDescribe("Healthcheck", func() {
	var appName string

	BeforeEach(func() {
		if Config.Backend != "diego" {
			Skip(skip_messages.SkipDiegoMessage)
		}

		appName = random_name.CATSRandomName("APP")
	})

	AfterEach(func() {
		app_helpers.AppReport(appName, DEFAULT_TIMEOUT)
		Eventually(cf.Cf("delete", appName, "-f"), DEFAULT_TIMEOUT).Should(Exit(0))
	})

	Describe("when the healthcheck is set to none", func() {
		It("starts up successfully", func() {
			By("pushing it")
			Eventually(cf.Cf(
				"push", appName,
				"-p", assets.NewAssets().WorkerApp,
				"--no-start",
				"-b", "go_buildpack",
				"-m", DEFAULT_MEMORY_LIMIT,
				"-d", Config.AppsDomain,
				"-i", "1",
				"-u", "none"),
				CF_PUSH_TIMEOUT,
			).Should(Exit(0))

			By("staging and running it")
			app_helpers.SetBackend(appName)
			Eventually(cf.Cf("start", appName), CF_PUSH_TIMEOUT).Should(Exit(0))

			By("verifying it's up")
			Eventually(func() *Session {
				appLogsSession := cf.Cf("logs", "--recent", appName)
				Expect(appLogsSession.Wait(DEFAULT_TIMEOUT)).To(Exit(0))
				return appLogsSession
			}, DEFAULT_TIMEOUT).Should(gbytes.Say("I am working at"))
		})
	})

	Describe("when the healthcheck is set to port", func() {
		It("starts up successfully", func() {
			By("pushing it")
			Eventually(cf.Cf(
				"push", appName,
				"-p", assets.NewAssets().Dora,
				"--no-start",
				"-b", Config.RubyBuildpackName,
				"-m", DEFAULT_MEMORY_LIMIT,
				"-d", Config.AppsDomain,
				"-i", "1",
				"-u", "port"),
				DEFAULT_TIMEOUT,
			).Should(Exit(0))

			By("staging and running it")
			app_helpers.SetBackend(appName)
			Eventually(cf.Cf("start", appName), CF_PUSH_TIMEOUT).Should(Exit(0))

			By("verifying it's up")
			Eventually(helpers.CurlingAppRoot(appName)).Should(ContainSubstring("Hi, I'm Dora!"))
		})
	})
})