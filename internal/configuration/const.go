package configuration

import (
	"errors"
	"math"
	"time"
)

// DefaultEnvPrefix is the default environment prefix.
const DefaultEnvPrefix = "AUTHELIA_"

// DefaultEnvDelimiter is the default environment delimiter.
const DefaultEnvDelimiter = "_"

const (
	constSecretSuffix = "_FILE"

	constDelimiter = "."

	constWindows = "windows"
)

var (
	errNoValidator = errors.New("no validator provided")
	errNoSources   = errors.New("no sources provided")

	errDecodeNonPtrMustHaveValue = errors.New("must have a non-empty value")
)

const (
	errFmtSecretAlreadyDefined = "secrets: error loading secret into key '%s': it's already defined in other " +
		"configuration sources"
	errFmtSecretOSError         = "secrets: error loading secret path %s into key '%s': %w"
	errFmtSecretOSPermission    = "secrets: error loading secret path %s into key '%s': file permission error occurred: %w"
	errFmtSecretOSNotExist      = "secrets: error loading secret path %s into key '%s': file does not exist error occurred: %w"
	errFmtGenerateConfiguration = "error occurred generating configuration: %+v"

	errFmtDecodeHookCouldNotParse           = "could not decode '%s' to a %s%s: %w"
	errFmtDecodeHookCouldNotParseBasic      = "could not decode to a %s%s: %w"
	errFmtDecodeHookCouldNotParseEmptyValue = "could not decode an empty value to a %s%s: %w"

	errFmtSpecialRemappedKey = "configuration key '%s' is deprecated in %s and has been replaced by '%s' when combined with the '%s' in the format of '%s': this should be automatically mapped for you but you will need to adjust your configuration to remove this message"
)

const (
	durationMax = time.Duration(math.MaxInt64)
)

// IMPORTANT: There is an uppercase copy of this in github.com/authelia/authelia/internal/templates named
// envSecretSuffixes.
// Make sure you update these at the same time.
var (
	secretSuffix          = []string{"key", "secret", "password", "token", "certificate_chain"}
	secretExclusionPrefix = []string{"identity_providers.oidc.lifespans."}
	secretExclusionExact  = []string{"server.tls.key"}
)
