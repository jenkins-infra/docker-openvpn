---
name: "Sync OpenVPN certificate inventory from cert/pki/issued"
scms:
  default:
    kind: github
    spec:
      user: "{{ .github.user }}"
      email: "{{ .github.email }}"
      owner: "{{ .github.owner }}"
      repository: "{{ .github.repository }}"
      token: "{{ requiredEnv .github.token }}"
      branch: "{{ .github.branch }}"

sources:
  certificatesList:
    name: "Discover OpenVPN certificate usernames"
    kind: shell
    spec:
      command: |
        set -eu
        for cert in cert/pki/issued/*.crt; do
          [ -e "$cert" ] || continue
          basename "$cert" .crt
        done | sort | paste -sd "," -

targets:
  updateCertificates:
    name: "Update certificate inventory in updatecli/values.yaml"
    kind: yaml
    sourceid: certificatesList
    spec:
      file: updatecli/values.yaml
      key: $.certificates
    scmid: default

actions:
  default:
    kind: github/pullrequest
    scmid: default
    spec:
      draft: true
      title: "chore(updatecli): sync OpenVPN certificate inventory"
      description: |
        This PR synchronizes the OpenVPN certificate inventory from the
        contents of `cert/pki/issued/`.
        
        **Certificates discovered:**
        {{- range $cert := splitList "," (source "certificatesList") }}
        - `{{ $cert }}`
        {{- end }}
      labels:
        - vpn
        - updatecli
