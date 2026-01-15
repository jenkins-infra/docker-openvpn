{{- range $username, $_ := .certificates }}
---
# yamllint disable rule:line-length
name: "Check VPN certificate expiration for {{ $username }}"

scms:
  default:
    kind: github
    spec:
      user: "{{ $.github.user }}"
      email: "{{ $.github.email }}"
      owner: "{{ $.github.owner }}"
      repository: "{{ $.github.repository }}"
      token: "{{ requiredEnv $.github.token }}"
      branch: "{{ $.github.branch }}"

sources:
  certExpiryDate:
    name: "Extract expiration date from {{ $username }}'s certificate"
    kind: shell
    spec:
      command: >
        bash ./updatecli/scripts/cert-expiry-extract.sh
        cert/pki/issued/{{ $username }}.crt
      environments:
        - name: PATH

conditions:
  checkIfExpiringSoon:
    name: "Check if certificate expires within 30 days"
    kind: shell
    sourceid: certExpiryDate
    spec:
      command: bash ./updatecli/scripts/cert-expiry-check.sh
      environments:
        - name: PATH

targets:
  markCertExpiring:
    name: "Mark {{ $username }}'s certificate as expiring"
    kind: file
    scmid: default
    spec:
      file: cert/pki/issued/{{ $username }}.crt.expiring
      content: |
        Certificate for {{ $username }} expires on {{ source "certExpiryDate" }}.
        Please renew your VPN certificate as soon as possible.

actions:
  default:
    kind: github/pullrequest
    scmid: default
    spec:
      draft: true
      title: "[DO NOT MERGE] VPN Certificate Expiring Soon: {{ $username }}"
      description: |
        @{{ $username }} your VPN certificate will expire on **{{ source "certExpiryDate" }}**.

        ## Action Required

        Your VPN certificate expires in less than **30 days**.
        Please renew it to avoid losing VPN access.

        ---
        **Note:** This is an automated notification PR.
        It is not meant to be merged and can be closed once acknowledged.
      labels:
        - vpn
        - certificate-expiration
        - action-required
{{- end }}
