package chaperon_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/consul-release/src/confab/chaperon"
	"github.com/cloudfoundry-incubator/consul-release/src/confab/config"
	"github.com/cloudfoundry-incubator/consul-release/src/confab/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		client         chaperon.Client
		timeout        *fakes.Timeout
		controller     *fakes.Controller
		keyringRemover *fakes.KeyringRemover
		configWriter   *fakes.ConfigWriter
		cfg            config.Config
	)

	BeforeEach(func() {
		controller = &fakes.Controller{}
		keyringRemover = &fakes.KeyringRemover{}
		configWriter = &fakes.ConfigWriter{}
		timeout = &fakes.Timeout{}

		cfg = config.Config{
			Node: config.ConfigNode{
				Name: "some-name",
			},
		}

		client = chaperon.NewClient(controller, keyringRemover, configWriter)
	})

	It("writes the consul configuration file", func() {
		err := client.Start(cfg, timeout)
		Expect(err).NotTo(HaveOccurred())
		Expect(configWriter.WriteCall.Receives.Config).To(Equal(cfg))
	})

	It("writes the service definitions", func() {
		err := client.Start(cfg, timeout)
		Expect(err).NotTo(HaveOccurred())
		Expect(controller.WriteServiceDefinitionsCall.CallCount).To(Equal(1))
	})

	It("removes the keyring file", func() {
		err := client.Start(cfg, timeout)
		Expect(err).NotTo(HaveOccurred())
		Expect(keyringRemover.ExecuteCall.CallCount).To(Equal(1))
	})

	It("boots the agent process", func() {
		err := client.Start(cfg, timeout)
		Expect(err).NotTo(HaveOccurred())
		Expect(controller.BootAgentCall.CallCount).To(Equal(1))
		Expect(controller.BootAgentCall.Receives.Timeout).To(Equal(timeout))
	})

	It("configures the client", func() {
		err := client.Start(cfg, timeout)
		Expect(err).NotTo(HaveOccurred())
		Expect(controller.ConfigureClientCall.CallCount).To(Equal(1))
	})

	Context("failure cases", func() {
		Context("when writing the consul config file fails", func() {
			It("returns an error", func() {
				configWriter.WriteCall.Returns.Error = errors.New("failed to write config")

				err := client.Start(cfg, timeout)
				Expect(err).To(MatchError(errors.New("failed to write config")))
			})
		})

		Context("when writing the service definitions fails", func() {
			It("returns an error", func() {
				controller.WriteServiceDefinitionsCall.Returns.Error = errors.New("failed to write service definitions")

				err := client.Start(cfg, timeout)
				Expect(err).To(MatchError(errors.New("failed to write service definitions")))
			})
		})

		Context("when removing the keyring fails", func() {
			It("returns an error", func() {
				keyringRemover.ExecuteCall.Returns.Error = errors.New("failed to remove keyring")

				err := client.Start(cfg, timeout)
				Expect(err).To(MatchError(errors.New("failed to remove keyring")))
			})
		})

		Context("when booting the agent fails", func() {
			It("returns an error", func() {
				controller.BootAgentCall.Returns.Error = errors.New("failed to boot agent")

				err := client.Start(cfg, timeout)
				Expect(err).To(MatchError(errors.New("failed to boot agent")))
			})
		})

		Context("when configuring the client fails", func() {
			It("returns an error", func() {
				controller.ConfigureClientCall.Returns.Error = errors.New("failed to configure client")

				err := client.Start(cfg, timeout)
				Expect(err).To(MatchError(errors.New("failed to configure client")))
			})
		})
	})

	Describe("Stop", func() {
		It("calls stop agent", func() {
			client.Stop()
			Expect(controller.StopAgentCall.CallCount).To(Equal(1))
		})
	})
})
