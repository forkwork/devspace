package provider

import (
	"context"
	"os"
	"time"

	"dev.khulnasoft.com/e2e/framework"
	"github.com/onsi/ginkgo/v2"
)

var _ = DevSpaceDescribe("devspace provider test suite", func() {
	ginkgo.Context("testing non-machine providers", ginkgo.Label("provider"), ginkgo.Ordered, func() {
		ctx := context.Background()
		initialDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		ginkgo.It("should add simple provider and delete it", func() {
			tempDir, err := framework.CopyToTempDir("tests/provider/testdata/simple-k8s-provider")
			framework.ExpectNoError(err)
			ginkgo.DeferCleanup(framework.CleanupTempDir, initialDir, tempDir)

			f := framework.NewDefaultFramework(initialDir + "/bin")

			// Ensure that provider 1 is deleted
			err = f.DevSpaceProviderDelete(ctx, "provider1", "--ignore-not-found")
			framework.ExpectNoError(err)

			// Add provider 1
			err = f.DevSpaceProviderAdd(ctx, tempDir+"/provider1.yaml")
			framework.ExpectNoError(err)

			// Ensure provider 1 exists but not provider X
			err = f.DevSpaceProviderUse(context.Background(), "provider1")
			framework.ExpectNoError(err)
			err = f.DevSpaceProviderUse(context.Background(), "providerX")
			framework.ExpectError(err)

			// Cleanup: delete provider 1
			err = f.DevSpaceProviderDelete(ctx, "provider1")
			framework.ExpectNoError(err)

			// Cleanup: ensure provider 1 is deleted
			err = f.DevSpaceProviderUse(context.Background(), "provider1")
			framework.ExpectError(err)
		})

		ginkgo.It("should add simple provider and update it", func() {
			tempDir, err := framework.CopyToTempDir("tests/provider/testdata/simple-k8s-provider")
			framework.ExpectNoError(err)
			defer framework.CleanupTempDir(initialDir, tempDir)

			f := framework.NewDefaultFramework(initialDir + "/bin")

			// Ensure that provider 2 is deleted
			err = f.DevSpaceProviderDelete(ctx, "provider2", "--ignore-not-found")
			framework.ExpectNoError(err)

			// Add provider 2 and use it
			err = f.DevSpaceProviderAdd(ctx, tempDir+"/provider2.yaml")
			framework.ExpectNoError(err)
			err = f.DevSpaceProviderUse(context.Background(), "provider2")
			framework.ExpectNoError(err)

			// Ensure provider 2 namespace parameter has the default value
			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
			err = f.DevSpaceProviderOptionsCheckNamespaceDescription(ctx, "provider2", "The namespace to use")
			framework.ExpectNoError(err)
			cancel()

			// Update provider 2 (change the namespace description value)
			err = f.DevSpaceProviderUpdate(context.Background(), "provider2", tempDir+"/provider2-update.yaml")
			framework.ExpectNoError(err)

			// Ensure that provider 2 was updated
			ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
			err = f.DevSpaceProviderOptionsCheckNamespaceDescription(ctx, "provider2", "Updated namespace parameter")
			framework.ExpectNoError(err)
			cancel()

			// Cleanup: delete provider 2
			err = f.DevSpaceProviderDelete(context.Background(), "provider2")
			framework.ExpectNoError(err)

			// Cleanup: ensure provider 2 is deleted
			err = f.DevSpaceProviderUse(context.Background(), "provider2")
			framework.ExpectError(err)
		})

		ginkgo.It("should list all providers", func() {
			tempDir, err := framework.CopyToTempDir("tests/provider/testdata/simple-k8s-provider")
			framework.ExpectNoError(err)
			ginkgo.DeferCleanup(framework.CleanupTempDir, initialDir, tempDir)

			f := framework.NewDefaultFramework(initialDir + "/bin")

			// Ensure that provider 1 is deleted
			err = f.DevSpaceProviderDelete(ctx, "provider1", "--ignore-not-found")
			framework.ExpectNoError(err)

			// Add provider 1
			err = f.DevSpaceProviderAdd(ctx, tempDir+"/provider1.yaml")
			framework.ExpectNoError(err)
			// Ensure provider 1 exists
			err = f.DevSpaceProviderUse(context.Background(), "provider1")
			framework.ExpectNoError(err)

			// Add .DS_Store file to tempDir
			err = os.Mkdir(tempDir+"/.DS_Store", 0755)
			framework.ExpectNoError(err)

			// List providers
			err = f.DevSpaceProviderList(context.Background())
			framework.ExpectNoError(err)

			// Cleanup: delete provider 1
			err = f.DevSpaceProviderDelete(ctx, "provider1")
			framework.ExpectNoError(err)

			// Cleanup: ensure provider 1 is deleted
			err = f.DevSpaceProviderUse(context.Background(), "provider1")
			framework.ExpectError(err)
		})

		ginkgo.It("should parse options", func() {
			tempDir, err := framework.CopyToTempDir("tests/provider/testdata/simple-k8s-provider")
			framework.ExpectNoError(err)
			ginkgo.DeferCleanup(framework.CleanupTempDir, initialDir, tempDir)

			f := framework.NewDefaultFramework(initialDir + "/bin")

			// Ensure that provider is deleted
			err = f.DevSpaceProviderDelete(ctx, "provider3", "--ignore-not-found")
			framework.ExpectNoError(err)

			podManifest := `
apiVersion: v1
kind: Pod
metadata:
	name: test
spec:
	containers:
	- name: devspace
`
			// Add provider
			err = f.DevSpaceProviderAdd(ctx, tempDir+"/provider3.yaml", "--option=TEMPLATE="+podManifest)
			framework.ExpectNoError(err)
			// Ensure provider exists
			err = f.DevSpaceProviderUse(context.Background(), "provider3")
			framework.ExpectNoError(err)

			// look for template option
			err = f.DevSpaceProviderFindOption(context.Background(), "provider3", podManifest)
			framework.ExpectNoError(err)

			// Cleanup: delete provider
			err = f.DevSpaceProviderDelete(ctx, "provider3")
			framework.ExpectNoError(err)

			// Cleanup: ensure provider is deleted
			err = f.DevSpaceProviderUse(context.Background(), "provider3")
			framework.ExpectError(err)
		})
	})
})
