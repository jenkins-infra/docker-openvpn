{{ range $network, $network_setup := .networks }}
  {{ $subnets := $network_setup.routes }}
  {{ $servers := $network_setup.servers }}
  {{ $multiValuedServers := $network_setup.multivalued_servers }}
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
    {{ if $server_data.report_url }}
      {{ if $server_data.resolve_dns }}
  {{ $server }}-dns:
    transformers:
      - replacers:
          - from: "https://"
            to: ""
          - from: "http://"
            to: ""
          - from: "/"
            to: ""
      {{ else }}
  {{ $server }}-cidr:
    transformers:
      - addsuffix: '/32'
      {{ end }}
    kind: json
    spec:
      file: {{ $server_data.report_url }}
      key: {{ $server_data.report_query }}
    {{ end }}
    {{ if $server_data.resolve_dns }}
  {{ $server }}-cidr:
    name: Get DNS values for {{ $server }}
    kind: shell
      {{ if $server_data.report_url }}
    dependson:
      - {{ $server }}-dns
      {{ end }}
    spec:
      command: >
        dig +short {{ empty $server_data.report_url | ternary $server (print `{{ source "` $server `-dns" }}`) }} | sort -u | grep -v '\.$' | tr '\n' ' ' | sed 's# #/32 #g' | xargs
    {{ end }}
  {{ end }}
targets:
  {{ range $subnet := $subnets }}
  config-{{ $subnet }}:
    name: Update {{ $subnet }} route in the YAML configuration of our OpenVPN CLI
    kind: yaml
    sourceid: {{ $subnet }}-cidr
    scmid: default
    spec:
      file: config.yaml
      key: $.networks.{{ $network }}.routes.{{ $subnet }}
  {{ end }}

  {{ range $server, $server_data := $servers }}
  config-{{ $server }}:
    name: Update {{ $server }} route in the YAML configuration of our OpenVPN CLI
    kind: yaml
    sourceid: {{ $server }}-cidr
    scmid: default
    spec:
      file: config.yaml
      key: $.networks.{{ $network }}.routes.'{{ $server }}'
  {{ end }}

  update-ccd:
    name: Update Client Configuration files
    kind: shell
    disablesourceinput: true
    scmid: default
    # This target should only be executed if one of its dependencies is changed (e.g. only regenerate CCDs if config.yaml changed)
    dependsonchange: true
    # Note: conditional execution only depends on "targets" with a logical OR (e.g. any target changed triggers execution)
    dependson:
    {{ range $server, $server_data := $servers }}
        - "target#config-{{ $server }}:or"
    {{ end }}
    {{ range $subnet := $subnets }}
        - "target#config-{{ $subnet }}:or"
    {{ end }}
    spec:
      command: >
        {
          set -x
          # We need the 'easyvpn' command built (symlinked from the repository root)
          # Requires 'go' installed
          cd ./utils/easyvpn
          go build
          cd -
          ./utils/easyvpn/easyvpn
        } || exit 2

        if [ "${DRY_RUN}" = "true" ];
        then
          echo "DRY_RUN: should run commands"
          exit 1
        else
          # Regenerate client config
          ./utils/easyvpn/easyvpn --commit=false --push=false clientconfig --all
          git status
          git diff --exit-code
        fi
      changedif:
        kind: 'exitcode'
        spec:
          warning: 1
          success: 0
          failure: 2

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
