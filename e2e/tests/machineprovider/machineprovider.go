package machineprovider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dev.khulnasoft.com/e2e/framework"
	"github.com/onsi/ginkgo/v2"
)

var _ = DevSpaceDescribe("devspace machine provider test suite", func() {
	ginkgo.Context("testing machine providers", ginkgo.Label("machineprovider"), ginkgo.Ordered, func() {
		var initialDir string

		ginkgo.BeforeEach(func() {
			var err error
			initialDir, err = os.Getwd()
			framework.ExpectNoError(err)
		})

		ginkgo.It("test start / stop / status", func(ctx context.Context) {
			f := framework.NewDefaultFramework(initialDir + "/bin")

			// copy test dir
			tempDir, err := framework.CopyToTempDirWithoutChdir(initialDir + "/tests/machineprovider/testdata/machineprovider")
			framework.ExpectNoError(err)
			ginkgo.DeferCleanup(func() {
				_ = os.RemoveAll(tempDir)
			})

			tempDirLocation, err := os.MkdirTemp("", "")
			framework.ExpectNoError(err)
			ginkgo.DeferCleanup(func() {
				_ = os.RemoveAll(tempDirLocation)
			})

			// create docker provider
			err = f.DevSpaceProviderAdd(ctx, filepath.Join(tempDir, "provider.yaml"), "-o", "LOCATION="+tempDirLocation)
			framework.ExpectNoError(err)
			ginkgo.DeferCleanup(func() {
				err = f.DevSpaceProviderDelete(context.Background(), "docker123")
				framework.ExpectNoError(err)
			})
			ginkgo.DeferCleanup(f.DevSpaceWorkspaceDelete, context.Background(), tempDir)

			// wait for devspace workspace to come online (deadline: 30s)
			err = f.DevSpaceUp(ctx, tempDir, "--debug")
			framework.ExpectNoError(err)

			// expect workspace
			workspace, err := f.FindWorkspace(ctx, tempDir)
			framework.ExpectNoError(err)

			// check status
			status, err := f.DevSpaceStatus(ctx, tempDir)
			framework.ExpectNoError(err)
			framework.ExpectEqual(strings.ToUpper(status.State), "RUNNING", "workspace status did not match")

			// stop container
			err = f.DevSpaceStop(ctx, tempDir)
			framework.ExpectNoError(err)

			// check status
			status, err = f.DevSpaceStatus(ctx, tempDir)
			framework.ExpectNoError(err)
			framework.ExpectEqual(strings.ToUpper(status.State), "STOPPED", "workspace status did not match")

			// wait for devspace workspace to come online (deadline: 30s)
			err = f.DevSpaceUp(ctx, tempDir)
			framework.ExpectNoError(err)

			// check if ssh works as it should start the container
			out, err := f.DevSpaceSSH(ctx, tempDir, fmt.Sprintf("cat /workspaces/%s/test.txt", workspace.ID))
			framework.ExpectNoError(err)
			framework.ExpectEqual(out, "Test123", "workspace content does not match")

			// delete workspace
			err = f.DevSpaceWorkspaceDelete(ctx, tempDir)
			framework.ExpectNoError(err)
		}, ginkgo.SpecTimeout(framework.GetTimeout()))

		ginkgo.It("test devspace inactivity timeout", func(ctx context.Context) {
			f := framework.NewDefaultFramework(initialDir + "/bin")

			// copy test dir
			tempDir, err := framework.CopyToTempDirWithoutChdir(initialDir + "/tests/machineprovider/testdata/machineprovider2")
			framework.ExpectNoError(err)
			ginkgo.DeferCleanup(func() {
				err = os.RemoveAll(tempDir)
				framework.ExpectNoError(err)
			})

			tempDirLocation, err := os.MkdirTemp("", "")
			framework.ExpectNoError(err)
			ginkgo.DeferCleanup(func() {
				err = os.RemoveAll(tempDirLocation)
				framework.ExpectNoError(err)
			})

			// create provider
			_ = f.DevSpaceProviderDelete(ctx, "docker123")
			err = f.DevSpaceProviderAdd(ctx, filepath.Join(tempDir, "provider.yaml"))
			framework.ExpectNoError(err)
			ginkgo.DeferCleanup(func() {
				err = f.DevSpaceProviderDelete(context.Background(), "docker123")
				framework.ExpectNoError(err)
			})
			ginkgo.DeferCleanup(f.DevSpaceWorkspaceDelete, context.Background(), tempDir)

			// wait for devspace workspace to come online (deadline: 30s)
			err = f.DevSpaceUp(ctx, tempDir, "--debug", "--daemon-interval=3s")
			framework.ExpectNoError(err)
			ginkgo.DeferCleanup(func() {
				// delete workspace
				err = f.DevSpaceWorkspaceDelete(context.Background(), tempDir)
				framework.ExpectNoError(err)
			})

			// check status
			status, err := f.DevSpaceStatus(ctx, tempDir, "--container-status=false")
			framework.ExpectNoError(err)
			framework.ExpectEqual(strings.ToUpper(status.State), "RUNNING", "workspace status did not match")

			// stop container
			err = f.DevSpaceStop(ctx, tempDir)
			framework.ExpectNoError(err)

			// check status
			status, err = f.DevSpaceStatus(ctx, tempDir, "--container-status=false")
			framework.ExpectNoError(err)
			framework.ExpectEqual(strings.ToUpper(status.State), "STOPPED", "workspace status did not match")

			// wait for devspace workspace to come online (deadline: 30s)
			err = f.DevSpaceUp(ctx, tempDir, "--daemon-interval=3s")
			framework.ExpectNoError(err)

			// check status
			status, err = f.DevSpaceStatus(ctx, tempDir, "--container-status=false")
			framework.ExpectNoError(err)
			framework.ExpectEqual(strings.ToUpper(status.State), "RUNNING", "workspace status did not match")

			// wait until workspace is stopped again
			now := time.Now()
			for {
				status, err := f.DevSpaceStatus(ctx, tempDir, "--container-status=false")
				framework.ExpectNoError(err)
				framework.ExpectEqual(time.Since(now) < time.Minute*2, true, "machine did not shutdown in time")
				if status.State == "Stopped" {
					break
				}

				time.Sleep(time.Second * 2)
			}
		}, ginkgo.SpecTimeout(framework.GetTimeout()*5))
	})
})
