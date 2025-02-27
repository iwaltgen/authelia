package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/authelia/jsonschema"
	"github.com/spf13/cobra"

	"github.com/authelia/authelia/v4/internal/authentication"
	"github.com/authelia/authelia/v4/internal/configuration/schema"
	"github.com/authelia/authelia/v4/internal/model"
)

func newDocsJSONSchemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "json-schema",
		Short: "Generate docs JSON schema",
		RunE:  rootSubCommandsRunE,

		DisableAutoGenTag: true,
	}

	cmd.AddCommand(newDocsJSONSchemaConfigurationCmd(), newDocsJSONSchemaUserDatabaseCmd(), newDocsJSONSchemaExportsCmd())

	return cmd
}

func newDocsJSONSchemaExportsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exports",
		Short: "Generate docs JSON schema for the various exports",
		RunE:  rootSubCommandsRunE,

		DisableAutoGenTag: true,
	}

	cmd.AddCommand(newDocsJSONSchemaExportsTOTPCmd(), newDocsJSONSchemaExportsWebAuthnCmd(), newDocsJSONSchemaExportsIdentifiersCmd())

	return cmd
}

func newDocsJSONSchemaExportsTOTPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "totp",
		Short: "Generate docs JSON schema for the TOTP exports",
		RunE:  docsJSONSchemaExportsTOTPRunE,

		DisableAutoGenTag: true,
	}

	return cmd
}

func newDocsJSONSchemaExportsWebAuthnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webauthn",
		Short: "Generate docs JSON schema for the WebAuthn exports",
		RunE:  docsJSONSchemaExportsWebAuthnRunE,

		DisableAutoGenTag: true,
	}

	return cmd
}

func newDocsJSONSchemaExportsIdentifiersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identifiers",
		Short: "Generate docs JSON schema for the identifiers exports",
		RunE:  docsJSONSchemaExportsIdentifiersRunE,

		DisableAutoGenTag: true,
	}

	return cmd
}

func newDocsJSONSchemaConfigurationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configuration",
		Short: "Generate docs JSON schema for the configuration",
		RunE:  docsJSONSchemaConfigurationRunE,

		DisableAutoGenTag: true,
	}

	return cmd
}

func newDocsJSONSchemaUserDatabaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-database",
		Short: "Generate docs JSON schema for the user database",
		RunE:  docsJSONSchemaUserDatabaseRunE,

		DisableAutoGenTag: true,
	}

	return cmd
}

func docsJSONSchemaExportsTOTPRunE(cmd *cobra.Command, args []string) (err error) {
	var version *model.SemanticVersion

	if version, err = readVersion(cmd); err != nil {
		return err
	}

	var (
		dir, file, schemaDir string
	)

	if schemaDir, err = getPFlagPath(cmd.Flags(), cmdFlagRoot, cmdFlagDirSchema); err != nil {
		return err
	}

	if dir, file, err = getJSONSchemaOutputPath(cmd, cmdFlagDocsStaticJSONSchemaExportsTOTP); err != nil {
		return err
	}

	return docsJSONSchemaGenerateRunE(cmd, args, version, schemaDir, &model.TOTPConfigurationDataExport{}, dir, file, nil)
}

func docsJSONSchemaExportsWebAuthnRunE(cmd *cobra.Command, args []string) (err error) {
	var version *model.SemanticVersion

	if version, err = readVersion(cmd); err != nil {
		return err
	}

	var (
		dir, file, schemaDir string
	)

	if schemaDir, err = getPFlagPath(cmd.Flags(), cmdFlagRoot, cmdFlagDirSchema); err != nil {
		return err
	}

	if dir, file, err = getJSONSchemaOutputPath(cmd, cmdFlagDocsStaticJSONSchemaExportsWebAuthn); err != nil {
		return err
	}

	return docsJSONSchemaGenerateRunE(cmd, args, version, schemaDir, &model.WebAuthnDeviceDataExport{}, dir, file, nil)
}

func docsJSONSchemaExportsIdentifiersRunE(cmd *cobra.Command, args []string) (err error) {
	var version *model.SemanticVersion

	if version, err = readVersion(cmd); err != nil {
		return err
	}

	var (
		dir, file, schemaDir string
	)

	if schemaDir, err = getPFlagPath(cmd.Flags(), cmdFlagRoot, cmdFlagDirSchema); err != nil {
		return err
	}

	if dir, file, err = getJSONSchemaOutputPath(cmd, cmdFlagDocsStaticJSONSchemaExportsIdentifiers); err != nil {
		return err
	}

	return docsJSONSchemaGenerateRunE(cmd, args, version, schemaDir, &model.UserOpaqueIdentifiersExport{}, dir, file, nil)
}

func docsJSONSchemaConfigurationRunE(cmd *cobra.Command, args []string) (err error) {
	var version *model.SemanticVersion

	if version, err = readVersion(cmd); err != nil {
		return err
	}

	var (
		dir, file, schemaDir string
	)

	if schemaDir, err = getPFlagPath(cmd.Flags(), cmdFlagRoot, cmdFlagDirSchema); err != nil {
		return err
	}

	if dir, file, err = getJSONSchemaOutputPath(cmd, cmdFlagDocsStaticJSONSchemaConfiguration); err != nil {
		return err
	}

	return docsJSONSchemaGenerateRunE(cmd, args, version, schemaDir, &schema.Configuration{}, dir, file, jsonschemaKoanfMapper)
}

func docsJSONSchemaUserDatabaseRunE(cmd *cobra.Command, args []string) (err error) {
	var version *model.SemanticVersion

	if version, err = readVersion(cmd); err != nil {
		return err
	}

	var (
		dir, file, schemaDir string
	)

	if schemaDir, err = getPFlagPath(cmd.Flags(), cmdFlagRoot, cmdFlagDirAuthentication); err != nil {
		return err
	}

	if dir, file, err = getJSONSchemaOutputPath(cmd, cmdFlagDocsStaticJSONSchemaUserDatabase); err != nil {
		return err
	}

	return docsJSONSchemaGenerateRunE(cmd, args, version, schemaDir, &authentication.FileUserDatabase{}, dir, file, jsonschemaKoanfMapper)
}

func docsJSONSchemaGenerateRunE(cmd *cobra.Command, _ []string, version *model.SemanticVersion, schemaDir string, v any, dir, file string, mapper func(reflect.Type) *jsonschema.Schema) (err error) {
	r := &jsonschema.Reflector{
		RequiredFromJSONSchemaTags: true,
		Mapper:                     mapper,
	}

	if runtime.GOOS == windows {
		mapComments := map[string]string{}

		if err = jsonschema.ExtractGoComments(goModuleBase, schemaDir, mapComments); err != nil {
			return err
		}

		if r.CommentMap == nil {
			r.CommentMap = map[string]string{}
		}

		for key, comment := range mapComments {
			r.CommentMap[strings.ReplaceAll(key, `\`, `/`)] = comment
		}
	} else {
		if err = r.AddGoComments(goModuleBase, schemaDir); err != nil {
			return err
		}
	}

	var (
		latest, next bool
	)

	latest, _ = cmd.Flags().GetBool(cmdFlagLatest)
	next, _ = cmd.Flags().GetBool(cmdFlagNext)

	var schemaVersion string

	schemaVersion = fmt.Sprintf("v%d.%d", version.Major, version.Minor)
	if next {
		schemaVersion = fmt.Sprintf("v%d.%d", version.Major, version.Minor+1)
	}

	schema := r.Reflect(v)

	schema.ID = jsonschema.ID(fmt.Sprintf(model.FormatJSONSchemaIdentifier, schemaVersion, file))

	if err = writeJSONSchema(schema, dir, schemaVersion, file); err != nil {
		return err
	}

	if latest {
		if err = writeJSONSchema(schema, dir, "latest", file); err != nil {
			return err
		}
	}

	return nil
}

func writeJSONSchema(schema *jsonschema.Schema, dir, version, file string) (err error) {
	var (
		data []byte
		f    *os.File
	)

	if data, err = json.MarshalIndent(schema, "", "  "); err != nil {
		return err
	}

	if _, err = os.Stat(filepath.Join(dir, version, "json-schema")); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Join(dir, version, "json-schema"), 0755); err != nil {
			return err
		}
	}

	if f, err = os.Create(filepath.Join(dir, version, "json-schema", file+".json")); err != nil {
		return err
	}

	if _, err = f.Write(data); err != nil {
		return err
	}

	return f.Close()
}

func getJSONSchemaOutputPath(cmd *cobra.Command, flag string) (dir, file string, err error) {
	if dir, err = getPFlagPath(cmd.Flags(), cmdFlagRoot, cmdFlagDocs, cmdFlagDocsStatic, cmdFlagDocsStaticJSONSchemas); err != nil {
		return "", "", err
	}

	if file, err = cmd.Flags().GetString(flag); err != nil {
		return "", "", err
	}

	return dir, file, nil
}

func jsonschemaKoanfMapper(t reflect.Type) *jsonschema.Schema {
	switch t.String() {
	case "regexp.Regexp", "*regexp.Regexp":
		return &jsonschema.Schema{
			Type:   jsonschema.TypeString,
			Format: jsonschema.FormatStringRegex,
		}
	case "time.Duration", "*time.Duration":
		return &jsonschema.Schema{
			OneOf: []*jsonschema.Schema{
				{
					Type:    jsonschema.TypeString,
					Pattern: `^\d+\s*(y|M|w|d|h|m|s|ms|((year|month|week|day|hour|minute|second|millisecond)s?))(\s*\d+\s*(y|M|w|d|h|m|s|ms|((year|month|week|day|hour|minute|second|millisecond)s?)))*$`,
				},
				{
					Type:        jsonschema.TypeInteger,
					Description: "The duration in seconds",
				},
			},
		}
	case "schema.CryptographicKey":
		return &jsonschema.Schema{
			Type: jsonschema.TypeString,
		}
	case "schema.CryptographicPrivateKey":
		return &jsonschema.Schema{
			Type:    jsonschema.TypeString,
			Pattern: `^-{5}(BEGIN ((RSA|EC) )?PRIVATE KEY-{5}\n([a-zA-Z0-9/+]{1,64}\n)+([a-zA-Z0-9/+]{1,64}[=]{0,2})\n-{5}END ((RSA|EC) )?PRIVATE KEY-{5}\n?)+$`,
		}
	case "rsa.PrivateKey", "*rsa.PrivateKey", "ecdsa.PrivateKey", "*.ecdsa.PrivateKey":
		return &jsonschema.Schema{
			Type: jsonschema.TypeString,
		}
	case "mail.Address", "*mail.Address":
		return &jsonschema.Schema{
			Type:   jsonschema.TypeString,
			Format: jsonschema.FormatStringEmail,
		}
	case "schema.CSPTemplate":
		return &jsonschema.Schema{
			Type:    jsonschema.TypeString,
			Default: buildCSP(codeCSPProductionDefaultSrc, codeCSPValuesCommon, codeCSPValuesProduction),
		}
	}

	return nil
}
