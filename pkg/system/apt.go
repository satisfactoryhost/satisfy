package system

import (
	"errors"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func NewAptPackageClient(distro string, logger *log.Entry) *PackageClient {
	return &PackageClient{
		logger: logger.WithField("package", "apt"),
		distro: distro,
	}
}

func (p *PackageClient) InstallPrereqs() error {
	p.logger.Info("installing pre-requisites")

	if err := p.updatePackages(); err != nil {
		return err
	}
	p.logger.Info("✅ done updating packages")

	if err := p.installPackages("software-properties-common"); err != nil {
		return err
	}
	p.logger.Info("✅ done installing software-properties-common")

	if err := p.addMultiverse(); err != nil {
		return err
	}
	p.logger.Info("✅ done adding multiverse")

	if err := p.add32BitSupport(); err != nil {
		return err
	}
	p.logger.Info("✅ done adding 32-bit support")

	if err := p.updatePackages(); err != nil {
		return err
	}
	p.logger.Info("✅ done updating packages")

	p.logger.Info("✅ done installing pre-requisites")
	return nil
}

func (p *PackageClient) InstallSteam() error {
	p.logger.Info("installing SteamCMD")

	if err := p.installPackages("lib32gcc1"); err != nil {
		return err
	}

	cmdText := "echo \"I AGREE\" | apt-get install -y steamcmd"
	cmd := exec.Command("bash", "-c", cmdText)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	p.logger.Info("creating /usr/games/steamcmd symlink")
	cmd = exec.Command("ln", "-s", "/usr/games/steamcmd", "steamcmd")
	if err := cmd.Run(); err != nil {
		return err
	}

	p.logger.Info("✅ done installing SteamCMD")
	return nil
}

func (p *PackageClient) CheckInstall() error {
	p.logger.Info("checking install, this could take a few minutes")
	cmd := exec.Command("/usr/games/steamcmd", "+info", "+quit")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	p.logger.Info("✅ done checking install")
	return nil
}

func (p *PackageClient) updatePackages() error {
	p.logger.Info("updating packages")
	cmd := exec.Command("apt-get", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *PackageClient) installPackages(packages ...string) error {
	p.logger.Infof("installing packages: %s", packages)
	args := []string{"install", "-y"}
	args = append(args, packages...)
	cmd := exec.Command("apt-get", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *PackageClient) addMultiverse() error {
	var cmd *exec.Cmd
	p.logger.Info("adding multiverse")
	if p.distro == "ubuntu" {
		cmd = exec.Command("add-apt-repository", "multiverse")
	} else if p.distro == "debian" {
		cmd = exec.Command("add-apt-repository", "non-free")
	} else {
		return errors.New("unsupported distro")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *PackageClient) add32BitSupport() error {
	p.logger.Info("adding 32-bit support")
	cmd := exec.Command("dpkg", "--add-architecture", "i386")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
