package agent

import (
	"regexp"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

// AgentConfiguration is the run-time configuration for an agent that
// has been loaded from the config file and command-line params
type AgentConfiguration struct {
	ConfigPath            string
	BootstrapScript       string
	BuildPath             string
	HooksPath             string
	SocketsPath           string
	GitMirrorsPath        string
	GitMirrorsLockTimeout int
	GitMirrorsSkipUpdate  bool
	PluginsPath           string
	GitCheckoutFlags      string
	GitCloneFlags         string
	GitCloneMirrorFlags   string
	GitCleanFlags         string
	GitFetchFlags         string
	GitSubmodules         bool
	AllowedRepositories   []*regexp.Regexp
	AllowedPlugins        []*regexp.Regexp
	SSHKeyscan            bool
	CommandEval           bool
	PluginsEnabled        bool
	PluginValidation      bool
	LocalHooksEnabled     bool
	StrictSingleHooks     bool
	RunInPty              bool

	SigningJWKSFile  string // Where to find the key to sign pipeline uploads with (passed through to jobs, they might be uploading pipelines)
	SigningJWKSKeyID string // The key ID to sign pipeline uploads with

	VerificationJWKS             jwk.Set // The set of keys to verify jobs with
	VerificationFailureBehaviour string  // What to do if job verification fails (one of `block` or `warn`)

	ANSITimestamps             bool
	TimestampLines             bool
	HealthCheckAddr            string
	DisconnectAfterJob         bool
	DisconnectAfterIdleTimeout int
	CancelGracePeriod          int
	SignalGracePeriod          time.Duration
	EnableJobLogTmpfile        bool
	JobLogPath                 string
	WriteJobLogsToStdout       bool
	LogFormat                  string
	Shell                      string
	Profile                    string
	RedactedVars               []string
	AcquireJob                 string
	TracingBackend             string
	TracingServiceName         string
}
