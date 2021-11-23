package system

import (
	"os"
	"os/exec"

	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

var (
	packageServices = map[string]func(distro string, logger *log.Entry) *PackageClient{
		"debian": NewAptPackageClient,
		"ubuntu": NewAptPackageClient,
	}
)

type PackageClient struct {
	logger *log.Entry
	distro string
}

type PackageService interface {
	InstallPrereqs() error
	InstallSteam() error
	CheckInstall() error
}

type Client struct {
	ID             string
	PackageService PackageService
}

func NewClient(logger *log.Entry) *Client {
	distroInfo, err := readOSRelease("/etc/os-release")
	if err != nil {
		logger.Fatal(err)
	}

	packageClientFunc, exists := packageServices[distroInfo["ID"]]
	if !exists {
		logger.Fatal("unsupported OS")
	}
	return &Client{
		ID:             distroInfo["ID"],
		PackageService: packageClientFunc(distroInfo["ID"], logger),
	}
}

func readOSRelease(configfile string) (map[string]string, error) {
	cfg, err := ini.Load(configfile)
	if err != nil {
		return nil, err
	}

	ConfigParams := make(map[string]string)
	ConfigParams["ID"] = cfg.Section("").Key("ID").String()

	return ConfigParams, nil
}

func (c *Client) CreateUser() error {
	cmd := exec.Command("adduser", "--system", "--group", "satisfactory")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Client) CreateService() error {
	cmd := exec.Command("mkdir", "-p", "/home/satisfactory/satisfactory")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("chown", "satisfactory:satisfactory", "/home/satisfactory/satisfactory")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmdText := `cat > /etc/systemd/system/satisfactory.service << EOF
[Unit]
Description=Satisfactory dedicated server
Wants=network-online.target
After=syslog.target network.target nss-lookup.target network-online.target

[Service]
Environment="LD_LIBRARY_PATH=./linux64"
ExecStartPre=/usr/games/steamcmd +login anonymous +force_install_dir "/home/satisfactory/satisfactory" +app_update 1690800 validate +quit
ExecStart=/home/satisfactory/satisfactory/FactoryServer.sh
User=satisfactory
Group=satisfactory
StandardOutput=journal
Restart=on-failure
KillSignal=SIGINT
WorkingDirectory=/home/satisfactory/satisfactory

[Install]
WantedBy=multi-user.target
EOF
`
	cmd = exec.Command("bash", "-c", cmdText)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
