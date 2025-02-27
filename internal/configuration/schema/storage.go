package schema

import (
	"crypto/tls"
	"net/url"
	"time"
)

// Storage represents the configuration of the storage backend.
type Storage struct {
	Local      *StorageLocal      `koanf:"local" json:"local" jsonschema:"title=Local" jsonschema_description:"The Local SQLite3 Storage configuration settings"`
	MySQL      *StorageMySQL      `koanf:"mysql" json:"mysql" jsonschema:"title=MySQL" jsonschema_description:"The MySQL/MariaDB Storage configuration settings"`
	PostgreSQL *StoragePostgreSQL `koanf:"postgres" json:"postgres" jsonschema:"title=PostgreSQL" jsonschema_description:"The PostgreSQL Storage configuration settings"`

	EncryptionKey string `koanf:"encryption_key" json:"encryption_key" jsonschema:"title=Encryption Key" jsonschema_description:"The Storage Encryption Key used to secure security sensitive values in the storage engine"`
}

// StorageLocal represents the configuration when using local storage.
type StorageLocal struct {
	Path string `koanf:"path" json:"path" jsonschema:"title=Path" jsonschema_description:"The Path for the SQLite3 database file"`
}

// StorageSQL represents the configuration of the SQL database.
type StorageSQL struct {
	Address  *AddressTCP   `koanf:"address" json:"address" jsonschema:"title=Address" jsonschema_description:"The address of the database"`
	Database string        `koanf:"database" json:"database" jsonschema:"title=Database" jsonschema_description:"The database name to use upon a successful connection"`
	Username string        `koanf:"username" json:"username" jsonschema:"title=Username" jsonschema_description:"The username to use to authenticate"`
	Password string        `koanf:"password" json:"password" jsonschema:"title=Password" jsonschema_description:"The password to use to authenticate"`
	Timeout  time.Duration `koanf:"timeout" json:"timeout" jsonschema:"default=5 seconds,title=Timeout" jsonschema_description:"The timeout for the database connection"`

	// Deprecated: use address instead.
	Host string `koanf:"host" json:"host" jsonschema:"deprecated"`

	// Deprecated: use address instead.
	Port int `koanf:"port" json:"port" jsonschema:"deprecated"`
}

// StorageMySQL represents the configuration of a MySQL database.
type StorageMySQL struct {
	StorageSQL `koanf:",squash"`

	TLS *TLS `koanf:"tls" json:"tls"`
}

// StoragePostgreSQL represents the configuration of a PostgreSQL database.
type StoragePostgreSQL struct {
	StorageSQL `koanf:",squash"`
	Schema     string `koanf:"schema" json:"schema" jsonschema:"default=public"`

	TLS *TLS `koanf:"tls" json:"tls"`

	// Deprecated: Use the TLS configuration instead.
	SSL *StoragePostgreSQLSSL `koanf:"ssl" json:"ssl" jsonschema:"deprecated"`
}

// StoragePostgreSQLSSL represents the SSL configuration of a PostgreSQL database.
type StoragePostgreSQLSSL struct {
	Mode            string `koanf:"mode" json:"mode" jsonschema:"deprecated"`
	RootCertificate string `koanf:"root_certificate" json:"root_certificate" jsonschema:"deprecated"`
	Certificate     string `koanf:"certificate" json:"certificate" jsonschema:"deprecated"`
	Key             string `koanf:"key" json:"key"`
}

// DefaultSQLStorageConfiguration represents the default SQL configuration.
var DefaultSQLStorageConfiguration = StorageSQL{
	Timeout: 5 * time.Second,
}

// DefaultMySQLStorageConfiguration represents the default MySQL configuration.
var DefaultMySQLStorageConfiguration = StorageMySQL{
	TLS: &TLS{
		MinimumVersion: TLSVersion{tls.VersionTLS12},
	},
}

// DefaultPostgreSQLStorageConfiguration represents the default PostgreSQL configuration.
var DefaultPostgreSQLStorageConfiguration = StoragePostgreSQL{
	StorageSQL: StorageSQL{
		Address: &AddressTCP{Address{true, false, -1, 5432, &url.URL{Scheme: AddressSchemeTCP, Host: "localhost:5432"}}},
	},
	Schema: "public",
	TLS: &TLS{
		MinimumVersion: TLSVersion{tls.VersionTLS12},
	},
	SSL: &StoragePostgreSQLSSL{
		Mode: "disable",
	},
}
