package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/obalunenko/getenv"
	"golang.org/x/oauth2"
)

type AppConfig struct {
	Database      *AppConfigDatabase
	ListenPort    int
	LogLevel      string
	Cookies       *AppConfigCookies
	Mode          string
	TLSCert       string
	TLSKey        string
	OAuth         *oauth2.Config // VATSIM Connect
	OAuthUserInfo string
}

type AppConfigDatabase struct {
	Hostname string
	Port     int
	Username string
	Password string
	Database string
	Driver   string
}

type AppConfigCookies struct {
	Secret   string
	Name     string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
}

type AppConfigOAuth struct {
	BaseURL    string
	Client     *AppOAuthClient
	Endpoint   *AppConfigOAuthEndpoint
	MyCallback string
	Scopes     string
}

type AppConfigOAuthEndpoint struct {
	Authorize string
	Token     string
	UserInfo  string
}

type AppOAuthClient struct {
	ID     string
	Secret string
}

var Cfg *AppConfig

func LoadConfig() (*AppConfig, error) {
	// Check if DOTENV is set
	envfile := ".env"
	if os.Getenv("DOTENV") != "" {
		envfile = os.Getenv("DOTENV")
	}

	// If envfile exists, load it
	if _, err := os.Stat(envfile); err == nil {
		err := godotenv.Load(envfile)
		if err != nil {
			return nil, err
		}
	}

	Cfg = &AppConfig{
		Database: &AppConfigDatabase{
			Hostname: getenv.EnvOrDefault("DB_HOSTNAME", "localhost"),
			Port:     getenv.EnvOrDefault("DB_PORT", 3306),
			Username: getenv.EnvOrDefault("DB_USERNAME", "root"),
			Password: getenv.EnvOrDefault("DB_PASSWORD", "root"),
			Database: getenv.EnvOrDefault("DB_DATABASE", "sso"),
			Driver:   getenv.EnvOrDefault("DB_DRIVER", "mysql"),
		},
		ListenPort: getenv.EnvOrDefault("LISTEN_PORT", 3000),
		LogLevel:   getenv.EnvOrDefault("LOG_LEVEL", "info"),
		Mode:       getenv.EnvOrDefault("MODE", "plain"),
		TLSCert:    readFile(getenv.EnvOrDefault("TLS_CERT", "")),
		TLSKey:     readFile(getenv.EnvOrDefault("TLS_KEY", "")),
		Cookies: &AppConfigCookies{
			Secret:   getenv.EnvOrDefault("COOKIE_SECRET", "secret"),
			Name:     getenv.EnvOrDefault("COOKIE_NAME", "training"),
			MaxAge:   getenv.EnvOrDefault("COOKIE_MAX_AGE", 86400),
			Path:     getenv.EnvOrDefault("COOKIE_PATH", "/"),
			Domain:   getenv.EnvOrDefault("COOKIE_DOMAIN", "training.zanartcc.org"),
			Secure:   getenv.EnvOrDefault("COOKIE_SECURE", true),
			HttpOnly: getenv.EnvOrDefault("COOKIE_HTTP_ONLY", true),
		},
		OAuth: &oauth2.Config{
			ClientID:     getenv.EnvOrDefault("OAUTH_CLIENT_ID", ""),
			ClientSecret: getenv.EnvOrDefault("OAUTH_CLIENT_SECRET", ""),
			RedirectURL:  getenv.EnvOrDefault("OAUTH_CALLBACK", ""),
			Scopes: strings.Split(
				getenv.EnvOrDefault(
					"OAUTH_SCOPES",
					"full_name vatsim_details email",
				),
				" ",
			),
			Endpoint: oauth2.Endpoint{
				AuthURL: fmt.Sprintf("%s%s",
					getenv.EnvOrDefault("OAUTH_BASE_URL", "https://auth.vatsim.net"),
					getenv.EnvOrDefault("OAUTH_ENDPOINT_AUTHORIZE", "/oauth/authorize"),
				),
				TokenURL: fmt.Sprintf("%s%s",
					getenv.EnvOrDefault("OAUTH_BASE_URL", "https://auth.vatsim.net"),
					getenv.EnvOrDefault("OAUTH_ENDPOINT_TOKEN", "/oauth/token"),
				),
			},
		},
		OAuthUserInfo: fmt.Sprintf("%s%s",
			getenv.EnvOrDefault("OAUTH_BASE_URL", "https://auth.vatsim.net"),
			getenv.EnvOrDefault("OAUTH_ENDPOINT_USER_INFO", "/api/user"),
		),
	}

	return Cfg, nil
}

func readFile(file string) string {
	if _, err := os.Stat(file); err == nil {
		// Open and read file
		f, err := os.Open(file)
		if err != nil {
			return ""
		}
		defer f.Close()
		body, _ := io.ReadAll(f)
		return string(body)
	}

	return ""
}
