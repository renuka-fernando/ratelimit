package main

import (
	"fmt"
	"os"
	"text/template"
)

type API struct {
	Name    string
	Version string
	Path    string
}

var routeTemp = `
                        - name: {{.Name}}
                          match:
                            safe_regex:
                              google_re2: {}
                              regex: "^/{{.Name}}/{{.Version}}{{.Path}}[/]{0,1}"
                            headers:
                              - name: ":method"
                                string_match:
                                  safe_regex:
                                    google_re2: {}
                                    regex: "^GET|POST|OPTIONS$"
                          route:
                            cluster_header: x-cluster-header
                            auto_host_rewrite: true
                            regex_rewrite:
                              pattern:
                                google_re2: {}
                                regex: "^/{{.Name}}/{{.Version}}{{.Path}}[/]{0,1}"
                              substitution: "/{{.Name}}/api/{{.Version}}{{.Path}}"
                            timeout: 60s
                            idle_timeout: 300s
                            upgrade_configs:
                              - upgrade_type: websocket
                                enabled: false
                            cors:
                              allow_methods: "GET, PUT, POST, DELETE, PATCH, OPTIONS"
                              allow_headers: "authorization, Access-Control-Allow-Origin, Content-Type, SOAPAction, apikey, testKey, Internal-Key"
                              allow_credentials: false
                              allow_origin_string_match:
                                - safe_regex:
                                    google_re2: {}
                                    regex: ".*"
                          decorator:
                            operation: "cc-envoy:^/{{.Name}}/{{.Version}}{{.Path}}[/]{0,1}"
                          typed_per_filter_config:
                            envoy.filters.http.lua:
                              "@type": "type.googleapis.com/envoy.extensions.filters.http.lua.v3.LuaPerRoute"
                              disabled: true
                            envoy.filters.http.ext_authz:
                              "@type": "type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthzPerRoute"
                              check_settings:
                                context_extensions:
                                  basePath: "/{{.Name}}/{{.Version}}"
                                  method: "GET POST"
                                  name: "{{.Name}}" # API Name
                                  path: "{{.Name}}/{{.Version}}{{.Path}}"
                                  prodClusterName: "clusterProd_cc-envoy_{{.Name}}{{.Version}}" # Cluster Name
                                  sandClusterName: ""
                                  vHost: "cc-envoy"
                                  version: "{{.Version}}"
                                disable_request_body_buffering: true
                            envoy.filters.http.ratelimit:
                              "@type": "type.googleapis.com/envoy.extensions.filters.http.ratelimit.v3.RateLimitPerRoute"
                              vh_rate_limits: INCLUDE
                        - name: {{.Name}}
                          match:
                            safe_regex:
                              google_re2: {}
                              regex: "^/{{.Name}}/{{.Version}}{{.Path}}/([^/]+)[/]{0,1}"
                            headers:
                              - name: ":method"
                                string_match:
                                  safe_regex:
                                    google_re2: {}
                                    regex: "^GET|PUT|DELETE|OPTIONS$"
                          route:
                            cluster_header: x-cluster-header
                            auto_host_rewrite: true
                            regex_rewrite:
                              pattern:
                                google_re2: {}
                                regex: "^/{{.Name}}/{{.Version}}{{.Path}}/([^/]+)[/]{0,1}"
                              substitution: "/{{.Name}}/{{.Version}}{{.Path}}/\\1"
                            timeout: 60s
                            idle_timeout: 300s
                            upgrade_configs:
                              - upgrade_type: websocket
                                enabled: false
                            cors:
                              allow_methods: "GET, PUT, POST, DELETE, PATCH, OPTIONS"
                              allow_headers: "authorization, Access-Control-Allow-Origin, Content-Type, SOAPAction, apikey, testKey, Internal-Key"
                              allow_credentials: false
                              allow_origin_string_match:
                                - safe_regex:
                                    google_re2: {}
                                    regex: ".*"
                          decorator:
                            operation: "cc-envoy:^/{{.Name}}/{{.Version}}{{.Path}}/([^/]+)[/]{0,1}"
                          typed_per_filter_config:
                            envoy.filters.http.lua:
                              "@type": "type.googleapis.com/envoy.extensions.filters.http.lua.v3.LuaPerRoute"
                              disabled: true
                            envoy.filters.http.ext_authz:
                              "@type": "type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthzPerRoute"
                              check_settings:
                                context_extensions:
                                  basePath: "/{{.Name}}/{{.Version}}"
                                  method: "GET PUT DELETE"
                                  name: "{{.Name}}" # API Name
                                  path: "/{{.Name}}/{{.Version}}{{.Path}}/{Id}"
                                  prodClusterName: "clusterProd_cc-envoy_{{.Name}}{{.Version}}" # Cluster Name
                                  sandClusterName: ""
                                  vHost: "cc-envoy"
                                  version: "{{.Version}}"
                                disable_request_body_buffering: true
                            envoy.filters.http.ratelimit:
                              "@type": "type.googleapis.com/envoy.extensions.filters.http.ratelimit.v3.RateLimitPerRoute"
                              vh_rate_limits: INCLUDE
`

var clusterTemp = `
    - name: clusterProd_cc-envoy_{{.Name}}{{.Version}}
      connect_timeout: 20s
      dns_refresh_rate: 5s
      dns_lookup_family: V4_ONLY
      type: STRICT_DNS
      lb_policy: ROUND_ROBIN
      load_assignment:
        cluster_name: clusterProd_cc-envoy_{{.Name}}{{.Version}}
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      address: 127.0.0.1
                      port_value: 9080
`

func main() {
	fmt.Println("Starting")

	fRoutes, err := os.Create("generated-envoy-config-routes.yaml")
	check(err)

	fClusters, err := os.Create("generated-envoy-config-clusters.yaml")
	check(err)

	tRoutes, err := template.New("Routes").Parse(routeTemp)
	check(err)

	tClusters, err := template.New("Clusters").Parse(clusterTemp)
	check(err)

	for apiNumber := 0; apiNumber < 100; apiNumber++ {
		api := &API{
			Name:    fmt.Sprintf("api%02dfoo", apiNumber),
			Version: "4.5.6",
		}

		for resourceNumber := 0; resourceNumber < 10; resourceNumber++ {
			api.Path = fmt.Sprintf("/resource/abcd%02d", resourceNumber)
			tRoutes.Execute(fRoutes, api)
		}
		tClusters.Execute(fClusters, api)
	}

	err = fRoutes.Close()
	check(err)

	err = fClusters.Close()
	check(err)

	fmt.Println("Done")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
