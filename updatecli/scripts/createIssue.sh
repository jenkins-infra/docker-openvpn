#!/bin/bash
# This script log to github and create an issue if not in dry mode
set -eux -o pipefail

command -v "gh" >/dev/null 2>&1 || { echo "ERROR: gh command not found. Exiting."; exit 1; }

cmd=$(cat <<-EOM
    gh issue create --title "[private.vpn.jenkins.io] $1 VPN CRL expires" \
                    --body "follow https://github.com/jenkins-infra/docker-openvpn?tab=readme-ov-file#howto-renew-certificate-revocation-list \
                            See https://github.com/jenkins-infra/helpdesk/issues/4266 for details." \
                    --label crl \
                    --label updatecli \
                    --label triage \
                    --repo jenkins-infra/helpdesk
EOM
)

if test "${DRY_RUN:=true}" == "false"
then
    export GITHUB_TOKEN="${UPDATECLI_GITHUB_TOKEN}"
    alreadyOpened=$(gh issue list --repo jenkins-infra/helpdesk --state open --search "label:crl label:updatecli" | wc -l)
    if test "$alreadyOpened" -eq 0
    then
        eval "${cmd}"
    else
        echo "issue already opened"
    fi
else
    echo "should create an issue on --repo jenkins-infra/helpdesk"
    echo "with title: [private.vpn.jenkins.io] $1 VPN CRL expires"
    echo "${cmd}"
fi

exit 0
