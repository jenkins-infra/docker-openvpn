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
    name: "Discover certificates as comma-separated string"
    kind: shell
    spec:
      command: |
        set -eu -o pipefail
        for cert in cert/pki/issued/*.crt; do
          [ -e "$cert" ] || continue
          basename "$cert" .crt
        done | sort | paste -sd ","

targets:
  updateCertificates:
    kind: yaml
    sourceid: certificatesList
    spec:
      file: updatecli/values.yaml
      key: $.certificates
      value: |
        {{ range $index, $cert := splitList "," (source "certificatesList") }}        - {{ $cert }}
        {{ end }}
    # scmigid: default

actions:
  default:
    kind: github/pullrequest
    scmid: default
    spec:
      title: "chore(updatecli): sync OpenVPN certificate inventory"
      description: |
        Synced certificate inventory from cert/pki/issued/
        
        **Certificates:**
        {{ range $index, $cert := splitList "," (source "certificatesList") }}
        - `{{ $cert }}`
        {{- end }}
      labels:
        - vpn
