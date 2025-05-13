package ide

import (
	"context"
	"os"

	"dev.khulnasoft.com/e2e/framework"
	"github.com/onsi/ginkgo/v2"
)

var _ = DevSpaceDescribe("devspace ide test suite", func() {
	ginkgo.Context("testing ides", ginkgo.Label("ide"), ginkgo.Ordered, func() {
		var initialDir string

		ginkgo.BeforeEach(func() {
			var err error
			initialDir, err = os.Getwd()
			framework.ExpectNoError(err)
		})

		ginkgo.It("start ides", func() {
			ctx := context.Background()

			f := framework.NewDefaultFramework(initialDir + "/bin")
			tempDir, err := framework.CopyToTempDir("tests/ide/testdata")
			framework.ExpectNoError(err)
			ginkgo.DeferCleanup(framework.CleanupTempDir, initialDir, tempDir)

			_ = f.DevSpaceProviderDelete(ctx, "docker")
			err = f.DevSpaceProviderAdd(ctx, "docker")
			framework.ExpectNoError(err)
			err = f.DevSpaceProviderUse(context.Background(), "docker")
			framework.ExpectNoError(err)

			ginkgo.DeferCleanup(f.DevSpaceWorkspaceDelete, context.Background(), tempDir)

			err = f.DevSpaceUpWithIDE(ctx, tempDir, "--open-ide=false", "--ide=vscode")
			framework.ExpectNoError(err)

			err = f.DevSpaceUpWithIDE(ctx, tempDir, "--open-ide=false", "--ide=openvscode")
			framework.ExpectNoError(err)

			err = f.DevSpaceUpWithIDE(ctx, tempDir, "--open-ide=false", "--ide=jupyternotebook")
			framework.ExpectNoError(err)

			err = f.DevSpaceUpWithIDE(ctx, tempDir, "--open-ide=false", "--ide=fleet")
			framework.ExpectNoError(err)

			// check if ssh works
			err = f.DevSpaceSSHEchoTestString(ctx, tempDir)
			framework.ExpectNoError(err)

			// TODO: test jetbrains ides
		})
	})
})
