package bootstrap

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/buildkite/agent/v3/bootstrap/shell"
	"github.com/buildkite/bintest/v3"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGittableURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		url, wantParsed, wantHost string
	}{
		{
			url:        "/home/vagrant/repo",
			wantParsed: "file:///home/vagrant/repo",
			wantHost:   "",
		},
		{
			url:        "file:///C:/Users/vagrant/repo",
			wantParsed: "file:///C:/Users/vagrant/repo",
			wantHost:   "",
		},
		{
			url:        "git@github.com:buildkite/agent.git",
			wantParsed: "ssh://git@github.com/buildkite/agent.git",
			wantHost:   "github.com",
		},
		{
			url:        "git@github.com-alias1:buildkite/agent.git",
			wantParsed: "ssh://git@github.com-alias1/buildkite/agent.git",
			wantHost:   "github.com-alias1",
		},
		{
			url:        "ssh://git@scm.xxx:7999/yyy/zzz.git",
			wantParsed: "ssh://git@scm.xxx:7999/yyy/zzz.git",
			wantHost:   "scm.xxx:7999",
		},
		{
			url:        "ssh://root@git.host.de:4019/var/cache/git/project.git",
			wantParsed: "ssh://root@git.host.de:4019/var/cache/git/project.git",
			wantHost:   "git.host.de:4019",
		},
	}

	for _, test := range tests {
		u, err := parseGittableURL(test.url)
		if err != nil {
			t.Errorf("parseGittableURL(%q) error = %v", test.url, err)
			continue
		}
		if got, want := u.String(), test.wantParsed; got != want {
			t.Errorf("parseGittableURL(%q) u.String() = %q, want %q", test.url, got, want)
		}
		if got, want := u.Host, test.wantHost; got != want {
			t.Errorf("parseGittableURL(%q) u.Host = %q, want %q", test.url, got, want)
		}
	}
}

func TestResolvingGitHostAliasesWithFlagSupport(t *testing.T) {
	t.Parallel()

	sh := shell.NewTestShell(t)

	ssh, err := bintest.NewMock("ssh")
	if err != nil {
		t.Fatalf("bintest.NewMock(ssh) error = %v", err)
	}
	defer ssh.CheckAndClose(t)

	sh.Env.Set("PATH", filepath.Dir(ssh.Path))

	ssh.
		Expect("-G", "github.com-alias1").
		AndWriteToStdout(`user buildkite
hostname github.com
port 22
addkeystoagent false
addressfamily any
batchmode no
canonicalizefallbacklocal yes
canonicalizehostname false
challengeresponseauthentication yes
checkhostip yes
compression no
controlmaster false
enablesshkeysign no
clearallforwardings no
exitonforwardfailure no
fingerprinthash SHA256
forwardagent no
forwardx11 no
forwardx11trusted no
gatewayports no
gssapiauthentication no
gssapidelegatecredentials no
hashknownhosts no
hostbasedauthentication no
identitiesonly no
kbdinteractiveauthentication yes
nohostauthenticationforlocalhost no
passwordauthentication yes
permitlocalcommand no
proxyusefdpass no
pubkeyauthentication yes
requesttty auto
streamlocalbindunlink no
stricthostkeychecking ask
tcpkeepalive yes
tunnel false
verifyhostkeydns false
visualhostkey no
updatehostkeys false
canonicalizemaxdots 1
connectionattempts 1
forwardx11timeout 1200
numberofpasswordprompts 3
serveralivecountmax 3
serveraliveinterval 0
ciphers chacha20-poly1305@openssh.com,aes128-ctr,aes192-ctr,aes256-ctr,aes128-gcm@openssh.com,aes256-gcm@openssh.com
hostkeyalgorithms ecdsa-sha2-nistp256-cert-v01@openssh.com,ecdsa-sha2-nistp384-cert-v01@openssh.com,ecdsa-sha2-nistp521-cert-v01@openssh.com,ssh-ed25519-cert-v01@openssh.com,rsa-sha2-512-cert-v01@openssh.com,rsa-sha2-256-cert-v01@openssh.com,ssh-rsa-cert-v01@openssh.com,ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521,ssh-ed25519,rsa-sha2-512,rsa-sha2-256,ssh-rsa
hostbasedkeytypes ecdsa-sha2-nistp256-cert-v01@openssh.com,ecdsa-sha2-nistp384-cert-v01@openssh.com,ecdsa-sha2-nistp521-cert-v01@openssh.com,ssh-ed25519-cert-v01@openssh.com,rsa-sha2-512-cert-v01@openssh.com,rsa-sha2-256-cert-v01@openssh.com,ssh-rsa-cert-v01@openssh.com,ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521,ssh-ed25519,rsa-sha2-512,rsa-sha2-256,ssh-rsa
kexalgorithms curve25519-sha256,curve25519-sha256@libssh.org,ecdh-sha2-nistp256,ecdh-sha2-nistp384,ecdh-sha2-nistp521,diffie-hellman-group-exchange-sha256,diffie-hellman-group16-sha512,diffie-hellman-group18-sha512,diffie-hellman-group14-sha256,diffie-hellman-group14-sha1
casignaturealgorithms ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521,ssh-ed25519,rsa-sha2-512,rsa-sha2-256,ssh-rsa
loglevel INFO
macs umac-64-etm@openssh.com,umac-128-etm@openssh.com,hmac-sha2-256-etm@openssh.com,hmac-sha2-512-etm@openssh.com,hmac-sha1-etm@openssh.com,umac-64@openssh.com,umac-128@openssh.com,hmac-sha2-256,hmac-sha2-512,hmac-sha1
pubkeyacceptedkeytypes ecdsa-sha2-nistp256-cert-v01@openssh.com,ecdsa-sha2-nistp384-cert-v01@openssh.com,ecdsa-sha2-nistp521-cert-v01@openssh.com,ssh-ed25519-cert-v01@openssh.com,rsa-sha2-512-cert-v01@openssh.com,rsa-sha2-256-cert-v01@openssh.com,ssh-rsa-cert-v01@openssh.com,ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521,ssh-ed25519,rsa-sha2-512,rsa-sha2-256,ssh-rsa
xauthlocation xauth
identityfile ~/.ssh/github_rsa
canonicaldomains
globalknownhostsfile /etc/ssh/ssh_known_hosts /etc/ssh/ssh_known_hosts2
userknownhostsfile ~/.ssh/known_hosts ~/.ssh/known_hosts2
sendenv LANG
sendenv LC_*
connecttimeout none
tunneldevice any:any
controlpersist no
escapechar ~
ipqos af21 cs1
rekeylimit 0 0
streamlocalbindmask 0177
syslogfacility USER`).
		AndExitWith(0)

	assert.Equal(t, "github.com", resolveGitHost(context.Background(), sh, "github.com-alias1"))

	ssh.
		Expect("-G", "blargh-no-alias.com").
		AndWriteToStdout(`user buildkite
hostname blargh-no-alias.com
port 22
addkeystoagent false
addressfamily any
batchmode no
canonicalizefallbacklocal yes
canonicalizehostname false
challengeresponseauthentication yes
checkhostip yes
compression no
controlmaster false
enablesshkeysign no
clearallforwardings no
exitonforwardfailure no
fingerprinthash SHA256
forwardagent no
forwardx11 no
forwardx11trusted no
gatewayports no
gssapiauthentication no
gssapidelegatecredentials no
hashknownhosts no
hostbasedauthentication no
identitiesonly no
kbdinteractiveauthentication yes
nohostauthenticationforlocalhost no
passwordauthentication yes
permitlocalcommand no
proxyusefdpass no
pubkeyauthentication yes
requesttty auto
streamlocalbindunlink no
stricthostkeychecking ask
tcpkeepalive yes
tunnel false
verifyhostkeydns false
visualhostkey no
updatehostkeys false
canonicalizemaxdots 1
connectionattempts 1
forwardx11timeout 1200
numberofpasswordprompts 3
serveralivecountmax 3
serveraliveinterval 0
ciphers chacha20-poly1305@openssh.com,aes128-ctr,aes192-ctr,aes256-ctr,aes128-gcm@openssh.com,aes256-gcm@openssh.com
hostkeyalgorithms ecdsa-sha2-nistp256-cert-v01@openssh.com,ecdsa-sha2-nistp384-cert-v01@openssh.com,ecdsa-sha2-nistp521-cert-v01@openssh.com,ssh-ed25519-cert-v01@openssh.com,rsa-sha2-512-cert-v01@openssh.com,rsa-sha2-256-cert-v01@openssh.com,ssh-rsa-cert-v01@openssh.com,ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521,ssh-ed25519,rsa-sha2-512,rsa-sha2-256,ssh-rsa
hostbasedkeytypes ecdsa-sha2-nistp256-cert-v01@openssh.com,ecdsa-sha2-nistp384-cert-v01@openssh.com,ecdsa-sha2-nistp521-cert-v01@openssh.com,ssh-ed25519-cert-v01@openssh.com,rsa-sha2-512-cert-v01@openssh.com,rsa-sha2-256-cert-v01@openssh.com,ssh-rsa-cert-v01@openssh.com,ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521,ssh-ed25519,rsa-sha2-512,rsa-sha2-256,ssh-rsa
kexalgorithms curve25519-sha256,curve25519-sha256@libssh.org,ecdh-sha2-nistp256,ecdh-sha2-nistp384,ecdh-sha2-nistp521,diffie-hellman-group-exchange-sha256,diffie-hellman-group16-sha512,diffie-hellman-group18-sha512,diffie-hellman-group14-sha256,diffie-hellman-group14-sha1
casignaturealgorithms ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521,ssh-ed25519,rsa-sha2-512,rsa-sha2-256,ssh-rsa
loglevel INFO
macs umac-64-etm@openssh.com,umac-128-etm@openssh.com,hmac-sha2-256-etm@openssh.com,hmac-sha2-512-etm@openssh.com,hmac-sha1-etm@openssh.com,umac-64@openssh.com,umac-128@openssh.com,hmac-sha2-256,hmac-sha2-512,hmac-sha1
pubkeyacceptedkeytypes ecdsa-sha2-nistp256-cert-v01@openssh.com,ecdsa-sha2-nistp384-cert-v01@openssh.com,ecdsa-sha2-nistp521-cert-v01@openssh.com,ssh-ed25519-cert-v01@openssh.com,rsa-sha2-512-cert-v01@openssh.com,rsa-sha2-256-cert-v01@openssh.com,ssh-rsa-cert-v01@openssh.com,ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521,ssh-ed25519,rsa-sha2-512,rsa-sha2-256,ssh-rsa
xauthlocation xauth
identityfile ~/.ssh/github_rsa
canonicaldomains
globalknownhostsfile /etc/ssh/ssh_known_hosts /etc/ssh/ssh_known_hosts2
userknownhostsfile ~/.ssh/known_hosts ~/.ssh/known_hosts2
sendenv LANG
sendenv LC_*
connecttimeout none
tunneldevice any:any
controlpersist no
escapechar ~
ipqos af21 cs1
rekeylimit 0 0
streamlocalbindmask 0177
syslogfacility USER`).
		AndExitWith(0)

	assert.Equal(t, "blargh-no-alias.com", resolveGitHost(context.Background(), sh, "blargh-no-alias.com"))

	ssh.
		Expect("-G", "cool-alias").
		AndWriteToStdout(`user cool-admin
hostname rad-git-host.com
port 443
addkeystoagent false
addressfamily any
batchmode no
canonicalizefallbacklocal yes
canonicalizehostname false
challengeresponseauthentication yes
checkhostip yes
compression no
controlmaster false
enablesshkeysign no
clearallforwardings no
exitonforwardfailure no
fingerprinthash SHA256
forwardagent no
forwardx11 no
forwardx11trusted no
gatewayports no
gssapiauthentication no
gssapidelegatecredentials no
hashknownhosts no
hostbasedauthentication no
identitiesonly no
kbdinteractiveauthentication yes
nohostauthenticationforlocalhost no
passwordauthentication yes
permitlocalcommand no
proxyusefdpass no
pubkeyauthentication yes
requesttty auto
streamlocalbindunlink no
stricthostkeychecking ask
tcpkeepalive yes
tunnel false
verifyhostkeydns false
visualhostkey no
updatehostkeys false
canonicalizemaxdots 1
connectionattempts 1
forwardx11timeout 1200
numberofpasswordprompts 3
serveralivecountmax 3
serveraliveinterval 0
ciphers chacha20-poly1305@openssh.com,aes128-ctr,aes192-ctr,aes256-ctr,aes128-gcm@openssh.com,aes256-gcm@openssh.com
hostkeyalgorithms ecdsa-sha2-nistp256-cert-v01@openssh.com,ecdsa-sha2-nistp384-cert-v01@openssh.com,ecdsa-sha2-nistp521-cert-v01@openssh.com,ssh-ed25519-cert-v01@openssh.com,rsa-sha2-512-cert-v01@openssh.com,rsa-sha2-256-cert-v01@openssh.com,ssh-rsa-cert-v01@openssh.com,ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521,ssh-ed25519,rsa-sha2-512,rsa-sha2-256,ssh-rsa
hostbasedkeytypes ecdsa-sha2-nistp256-cert-v01@openssh.com,ecdsa-sha2-nistp384-cert-v01@openssh.com,ecdsa-sha2-nistp521-cert-v01@openssh.com,ssh-ed25519-cert-v01@openssh.com,rsa-sha2-512-cert-v01@openssh.com,rsa-sha2-256-cert-v01@openssh.com,ssh-rsa-cert-v01@openssh.com,ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521,ssh-ed25519,rsa-sha2-512,rsa-sha2-256,ssh-rsa
kexalgorithms curve25519-sha256,curve25519-sha256@libssh.org,ecdh-sha2-nistp256,ecdh-sha2-nistp384,ecdh-sha2-nistp521,diffie-hellman-group-exchange-sha256,diffie-hellman-group16-sha512,diffie-hellman-group18-sha512,diffie-hellman-group14-sha256,diffie-hellman-group14-sha1
casignaturealgorithms ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521,ssh-ed25519,rsa-sha2-512,rsa-sha2-256,ssh-rsa
loglevel INFO
macs umac-64-etm@openssh.com,umac-128-etm@openssh.com,hmac-sha2-256-etm@openssh.com,hmac-sha2-512-etm@openssh.com,hmac-sha1-etm@openssh.com,umac-64@openssh.com,umac-128@openssh.com,hmac-sha2-256,hmac-sha2-512,hmac-sha1
pubkeyacceptedkeytypes ecdsa-sha2-nistp256-cert-v01@openssh.com,ecdsa-sha2-nistp384-cert-v01@openssh.com,ecdsa-sha2-nistp521-cert-v01@openssh.com,ssh-ed25519-cert-v01@openssh.com,rsa-sha2-512-cert-v01@openssh.com,rsa-sha2-256-cert-v01@openssh.com,ssh-rsa-cert-v01@openssh.com,ecdsa-sha2-nistp256,ecdsa-sha2-nistp384,ecdsa-sha2-nistp521,ssh-ed25519,rsa-sha2-512,rsa-sha2-256,ssh-rsa
xauthlocation xauth
identityfile ~/.ssh/github_rsa
canonicaldomains
globalknownhostsfile /etc/ssh/ssh_known_hosts /etc/ssh/ssh_known_hosts2
userknownhostsfile ~/.ssh/known_hosts ~/.ssh/known_hosts2
sendenv LANG
sendenv LC_*
connecttimeout none
tunneldevice any:any
controlpersist no
escapechar ~
ipqos af21 cs1
rekeylimit 0 0
streamlocalbindmask 0177
syslogfacility USER`).
		AndExitWith(0)

	assert.Equal(t, "rad-git-host.com:443", resolveGitHost(context.Background(), sh, "cool-alias"))
}

func TestResolvingGitHostAliasesWithoutFlagSupport(t *testing.T) {
	t.Parallel()

	sh := shell.NewTestShell(t)

	ssh, err := bintest.NewMock("ssh")
	if err != nil {
		t.Fatalf("bintest.NewMock(ssh) error = %v", err)
	}
	defer ssh.CheckAndClose(t)

	sh.Env.Set("PATH", filepath.Dir(ssh.Path))

	ssh.
		Expect("-G", "github.com-alias1").
		AndWriteToStderr(`unknown option -- G
usage: ssh [-1246AaCfgKkMNnqsTtVvXxYy] [-b bind_address] [-c cipher_spec]
           [-D [bind_address:]port] [-E log_file] [-e escape_char]
           [-F configfile] [-I pkcs11] [-i identity_file]
           [-L [bind_address:]port:host:hostport] [-l login_name] [-m mac_spec]
           [-O ctl_cmd] [-o option] [-p port]
           [-Q cipher | cipher-auth | mac | kex | key]
           [-R [bind_address:]port:host:hostport] [-S ctl_path] [-W host:port]
           [-w local_tun[:remote_tun]] [user@]hostname [command]`).
		AndExitWith(255)

	assert.Equal(t, "github.com", resolveGitHost(context.Background(), sh, "github.com-alias1"))

	ssh.
		Expect("-G", "blargh-no-alias.com").
		AndWriteToStderr(`unknown option -- G
usage: ssh [-1246AaCfgKkMNnqsTtVvXxYy] [-b bind_address] [-c cipher_spec]
           [-D [bind_address:]port] [-E log_file] [-e escape_char]
           [-F configfile] [-I pkcs11] [-i identity_file]
           [-L [bind_address:]port:host:hostport] [-l login_name] [-m mac_spec]
           [-O ctl_cmd] [-o option] [-p port]
           [-Q cipher | cipher-auth | mac | kex | key]
           [-R [bind_address:]port:host:hostport] [-S ctl_path] [-W host:port]
           [-w local_tun[:remote_tun]] [user@]hostname [command]`).
		AndExitWith(255)

	assert.Equal(t, "blargh-no-alias.com", resolveGitHost(context.Background(), sh, "blargh-no-alias.com"))
}

func TestGitCheckRefFormat(t *testing.T) {
	for ref, want := range map[string]bool{
		"hello":          true,
		"hello-world":    true,
		"hello/world":    true,
		"--option":       false,
		" leadingspace":  false,
		"has space":      false,
		"has~tilde":      false,
		"has^caret":      false,
		"has:colon":      false,
		"has\007control": false,
		"has\177del":     false,
		"endswithdot.":   false,
		"two..dots":      false,
		"@":              false,
		"back\\slash":    false,
	} {
		if got := gitCheckRefFormat(ref); got != want {
			t.Errorf("gitCheckRefFormat(%q) = %t, want %t", ref, got, want)
		}
	}
}

func TestGitCheckoutValidatesRef(t *testing.T) {
	sh := new(mockShellRunner)
	defer sh.Check(t)
	err := gitCheckout(context.Background(), &shell.Shell{}, "", "--nope")
	assert.EqualError(t, err, `"--nope" is not a valid git ref format`)
}

func TestGitCheckout(t *testing.T) {
	sh := new(mockShellRunner).Expect("git", "checkout", "-f", "-q", "main")
	defer sh.Check(t)
	err := gitCheckout(context.Background(), sh, "-f -q", "main")
	require.NoError(t, err)
}

func TestGitCheckoutSketchyArgs(t *testing.T) {
	sh := new(mockShellRunner)
	defer sh.Check(t)
	err := gitCheckout(context.Background(), sh, "-f -q", "  --hello")
	assert.EqualError(t, err, `"  --hello" is not a valid git ref format`)
}

func TestGitClone(t *testing.T) {
	sh := new(mockShellRunner).Expect("git", "clone", "-v", "--references", "url", "--", "repo", "dir")
	defer sh.Check(t)
	err := gitClone(context.Background(), sh, "-v --references url", "repo", "dir")
	require.NoError(t, err)
}

func TestGitClean(t *testing.T) {
	sh := new(mockShellRunner).Expect("git", "clean", "--foo", "--bar")
	defer sh.Check(t)
	err := gitClean(context.Background(), sh, "--foo --bar")
	require.NoError(t, err)
}

func TestGitCleanSubmodules(t *testing.T) {
	sh := new(mockShellRunner).Expect("git", "submodule", "foreach", "--recursive", "git clean --foo --bar")
	defer sh.Check(t)
	err := gitCleanSubmodules(context.Background(), sh, "--foo --bar")
	require.NoError(t, err)
}

func TestGitFetch(t *testing.T) {
	sh := new(mockShellRunner).Expect("git", "fetch", "--foo", "--bar", "--", "repo", "ref1", "ref2")
	defer sh.Check(t)
	err := gitFetch(context.Background(), sh, "--foo --bar", "repo", "ref1", "ref2")
	require.NoError(t, err)
}

// mockShellRunner implements shellRunner for testing expected calls.
type mockShellRunner struct {
	got, want [][]string
}

func (r *mockShellRunner) Expect(cmd string, args ...string) *mockShellRunner {
	r.want = append(r.want, append([]string{cmd}, args...))
	return r
}

func (r *mockShellRunner) Run(_ context.Context, cmd string, args ...string) error {
	r.got = append(r.got, append([]string{cmd}, args...))
	return nil
}

func (r *mockShellRunner) Check(t *testing.T) {
	if diff := cmp.Diff(r.got, r.want); diff != "" {
		t.Errorf("mockShellRunner diff (-got +want):\n%s", diff)
	}
}
