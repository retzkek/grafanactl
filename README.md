# Usage

```
SYNOPSIS
    grafanactl is a backup/restore utility for Grafana dashboards.

USAGE
    grafanactl [OPTIONS] COMMAND [COMMAND OPTIONS]

OPTIONS

	-headers=[]
		Comma-separated list of extra headers to pass, e.g. "X-User:foo,X-Grafana-Org-Id:1" (or set GRAFANA_HEADERS)
    -key=[]
        Grafana API key (or set GRAFANA_API_KEY)
    -path=[.]
        path to local dashboard repository (or set GRAFANA_PATH)
    -url=[http://play.grafana.org]
        Grafana base URL (or set GRAFANA_URL)
    -v=[false]
        turn on verbose output

COMMANDS

    get [OPTIONS] [DASHBOARD...]
        Retrieve dashboards and save to file.

    help [COMMAND]
        Print command usage and options.

    list [OPTIONS]
        List dashboards.

    push [OPTIONS] [DASHBOARD...]
        Read dashboards from file and push to Grafana.
```
