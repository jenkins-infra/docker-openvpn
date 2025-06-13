{{ range $network, $network_setup := .networks }}
  {{ $subnets := $network_setup.routes }}
  {{ $servers := $network_setup.servers }}
---
name: Update the {{ $network | quote }} network YAML configuration of our OpenVPN CLI

scms:
  default:
    kind: github
    spec:
      user: "{{ $.github.user }}"
      email: "{{ $.github.email }}"
      owner: "{{ $.github.owner }}"
      repository: "{{ $.github.repository }}"
      token: "{{ requiredEnv $.github.token }}"
      username: "{{ $.github.username }}"
      branch: "{{ $.github.branch }}"

sources:
  {{ range $subnet := $subnets }}
  {{ $subnet }}-cidr:
    kind: json
    spec:
      file: https://reports.jenkins.io/jenkins-infra-data-reports/azure-net.json
      key: .vnets.{{ $subnet }}-vnet.[0]
  {{ end }}

  {{ range $server, $server_data := $servers }}
  {{ $server }}-cidr:
    kind: json
    spec:
      file: {{ $server_data.report_url }}
      key: {{ $server_data.report_query }}
    transformers:
      - addsuffix: '/32'
  {{ end }}

targets:
  {{ range $subnet := $subnets }}
  config-{{ $subnet }}:
    name: Update {{ $subnet }} route in the YAML configuration of our OpenVPN CLI
    kind: yaml
    sourceid: {{ $subnet }}-cidr
    spec:
      file: config.yaml
      key: $.networks.{{ $network }}.routes.{{ $subnet }}
  {{ end }}

  {{ range $server, $server_data := $servers }}
  config-{{ $server }}:
    name: Update {{ $server }} route in the YAML configuration of our OpenVPN CLI
    kind: yaml
    sourceid: {{ $server }}-cidr
    spec:
      file: config.yaml
      key: $.networks.{{ $network }}.routes.'{{ $server }}'
  {{ end }}

actions:
  default:
    kind: github/pullrequest
    scmid: default
    title: Update list of IPs restricted to VPN access only
    spec:
      labels:
        - enhancement

...
{{ end }}
