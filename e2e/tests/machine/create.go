package machine

import (
	"context"
	"os"

	"github.com/google/uuid"
	"dev.khulnasoft.com/e2e/framework"
	"github.com/onsi/ginkgo/v2"
)

var _ = DevSpaceDescribe("devspace machine create", func() {
	ginkgo.Context("testing machine", ginkgo.Label("machine"), ginkgo.Ordered, func() {
		ctx := context.Background()
		initialDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		ginkgo.It("should add simple machine and then delete it", func() {
			tempDir, err := framework.CopyToTempDir("tests/machine/testdata")
			framework.ExpectNoError(err)
			ginkgo.DeferCleanup(framework.CleanupTempDir, initialDir, tempDir)

			f := framework.NewDefaultFramework(initialDir + "/bin")

			// Ensure that mock-provider is deleted
			_ = f.DevSpaceProviderDelete(ctx, "mock-provider")

			ginkgo.By("Add mock provider")
			err = f.DevSpaceProviderAdd(ctx, tempDir+"/mock-provider.yaml")
			framework.ExpectNoError(err)

			ginkgo.By("Use mock provier")
			err = f.DevSpaceProviderUse(context.Background(), "mock-provider")
			framework.ExpectNoError(err)

			machineUUID, _ := uuid.NewRandom()
			machineName := machineUUID.String()

			ginkgo.By("Create test machine with mock provider")
			err = f.DevSpaceMachineCreate([]string{machineName})
			framework.ExpectNoError(err)

			ginkgo.By("Remove test machine")
			err = f.DevSpaceMachineDelete([]string{machineName})
			framework.ExpectNoError(err)
		})
	})
})
