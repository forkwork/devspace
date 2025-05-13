package up

import (
	"context"
	"os"
	"os/exec"

	"dev.khulnasoft.com/e2e/framework"
	"github.com/onsi/ginkgo/v2"
)

var _ = DevSpaceDescribe("devspace up test suite", func() {
	ginkgo.Context("testing up command", ginkgo.Label("up-podman"), ginkgo.Ordered, func() {
		var initialDir string

		ginkgo.BeforeEach(func() {
			var err error
			initialDir, err = os.Getwd()
			framework.ExpectNoError(err)
		})

		ginkgo.Context("using docker provider", func() {
			ginkgo.Context("with rootfull podman", ginkgo.Ordered, func() {
				ginkgo.It("should setup rootful podman", func(ctx context.Context) {
					wrapper, err := os.Create(initialDir + "/bin/podman-rootful")
					framework.ExpectNoError(err)

					defer wrapper.Close()

					_, err = wrapper.WriteString(`#!/bin/sh
				sudo podman "$@"
				`)
					framework.ExpectNoError(err)

					err = wrapper.Close()
					framework.ExpectNoError(err)

					cmd := exec.Command("sudo", "chmod", "+x", initialDir+"/bin/podman-rootful")
					err = cmd.Run()
					framework.ExpectNoError(err)

					err = exec.Command(initialDir+"/bin/podman-rootful", "ps").Run()
					framework.ExpectNoError(err)
				}, ginkgo.SpecTimeout(framework.GetTimeout()))

				ginkgo.It("should start a new workspace with existing image", func(ctx context.Context) {
					tempDir, err := framework.CopyToTempDir("tests/up/testdata/docker")
					framework.ExpectNoError(err)
					ginkgo.DeferCleanup(framework.CleanupTempDir, initialDir, tempDir)

					f := framework.NewDefaultFramework(initialDir + "/bin")

					_ = f.DevSpaceProviderDelete(ctx, "docker")
					err = f.DevSpaceProviderAdd(ctx, "docker", "-o", "DOCKER_PATH="+initialDir+"/bin/podman-rootful")
					framework.ExpectNoError(err)

					err = f.DevSpaceProviderUse(ctx, "docker")
					framework.ExpectNoError(err)

					ginkgo.DeferCleanup(f.DevSpaceWorkspaceDelete, context.Background(), tempDir)

					// Wait for devspace workspace to come online (deadline: 30s)
					err = f.DevSpaceUp(ctx, tempDir)
					framework.ExpectNoError(err)
				}, ginkgo.SpecTimeout(framework.GetTimeout()))
			})
			ginkgo.Context("with rootless podman", ginkgo.Ordered, func() {
				ginkgo.It("should start a new workspace with existing image", func(ctx context.Context) {
					tempDir, err := framework.CopyToTempDir("tests/up/testdata/docker")
					framework.ExpectNoError(err)
					ginkgo.DeferCleanup(framework.CleanupTempDir, initialDir, tempDir)

					f := framework.NewDefaultFramework(initialDir + "/bin")

					_ = f.DevSpaceProviderDelete(ctx, "docker")
					err = f.DevSpaceProviderAdd(ctx, "docker", "-o", "DOCKER_PATH=podman")
					framework.ExpectNoError(err)
					err = f.DevSpaceProviderUse(ctx, "docker")
					framework.ExpectNoError(err)

					ginkgo.DeferCleanup(f.DevSpaceWorkspaceDelete, context.Background(), tempDir)

					// Wait for devspace workspace to come online (deadline: 30s)
					err = f.DevSpaceUp(ctx, tempDir)
					framework.ExpectNoError(err)
				}, ginkgo.SpecTimeout(framework.GetTimeout()))
			})
		})
	})
})
