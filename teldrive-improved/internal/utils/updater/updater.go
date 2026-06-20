package updater

import (
	"context"
	"fmt"
	"runtime"

	"github.com/Masterminds/semver/v3"
)

type Updater struct {
	os   string
	arch string
}

func New(os, arch string) *Updater {
	return &Updater{os: os, arch: arch}
}

func (u *Updater) GetLatestVersion(ctx context.Context) (*semver.Version, error) {
	return semver.NewVersion("2.0.0")
}

func (u *Updater) DownloadAndInstall(ctx context.Context, version string) error {
	return fmt.Errorf("not implemented")
}
