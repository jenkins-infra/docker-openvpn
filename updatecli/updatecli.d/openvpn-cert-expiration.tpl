{{ range $key, $val := .certificates }}
---
# yamllint disable rule:line-length
name: "Check VPN certificate expiration for {{ $val.username }}"

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
    name: "Extract expiration date from {{ $val.username }}'s certificate"
    kind: shell
    spec:
      command: bash ./updatecli/scripts/cert-expiry-extract.sh {{ $val.cert_file }}
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
  markCertExpired:
    name: "Mark {{ $val.username }}'s certificate as expiring"
    kind: file
    spec:
      file: {{ $val.cert_file }}
      content: "EXPIRED - Certificate expiring soon, please renew"
    scmid: default

actions:
  default:
    kind: github/pullrequest
    scmid: default
    spec:
      draft: true
      title: "[DO NOT MERGE] VPN Certificate Expiring Soon: {{ $val.username }}"
      description: |
        {{ $val.username }} your VPN certificate will expire on **{{ source "certExpiryDate" }}**.
        
        ## Action Required
        
        Your certificate expires in less than 30 days. Please renew it as soon as possible to maintain VPN access.
        
        ### Renewal Instructions
        
        [TODO: Add link to certificate renewal documentation]
        
        ---
        
        **Note**: This is an automated notification PR and will not be merged. Please close this PR after acknowledging.
      labels:
        - vpn
        - certificate-expiration
        - action-required
{{ end }}
