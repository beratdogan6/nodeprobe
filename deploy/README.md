# Visualizing nodeprobe with Prometheus + Grafana

This directory is a self-contained Docker stack that scrapes nodeprobe and
draws its metrics. Nothing here is needed to build or run nodeprobe itself.

```
deploy/
├── docker-compose.yml                     Prometheus + Grafana
├── prometheus/prometheus.yml              what to scrape, how often
└── grafana/
    ├── provisioning/datasources/          Prometheus wired up automatically
    ├── provisioning/dashboards/           where Grafana looks for dashboards
    └── dashboards/nodeprobe.json          the dashboard itself
```

## Why nodeprobe stays outside Docker

nodeprobe reads `/proc/loadavg`, `/proc/meminfo` and `statfs("/")`. Inside a
container, `/` is the container's overlay filesystem, so the disk numbers would
describe the image, not the machine. So nodeprobe runs on the host and only
Prometheus and Grafana run in containers. They reach the host through
`host.docker.internal`, which `extra_hosts: host-gateway` maps to the host.

## Run it

1. Start nodeprobe on the host. Use a built binary rather than `go run`:
   the compiler would be using RAM while you measure it.

   ```sh
   go build -o /tmp/nodeprobe ./cmd/nodeprobe && /tmp/nodeprobe
   ```

2. Start the stack:

   ```sh
   cd deploy
   docker compose up -d
   ```

3. Open Grafana at <http://localhost:3000> → Dashboards → **nodeprobe**.
   No login is needed to look at it: anonymous access is on, because this
   stack is meant to listen on localhost only. Grafana pins anonymous users
   to Viewer, so sign in as `admin` / `admin` when you want to edit a panel.
   Prometheus is at <http://localhost:9090>.

Stop with `docker compose down`, or `docker compose down -v` to also throw away
the stored metrics.

## The dashboard

| Panel | Query |
| --- | --- |
| Load average | `nodeprobe_load1` / `load5` / `load15` |
| Memory used | `100 * (1 - nodeprobe_memory_available_bytes / nodeprobe_memory_total_bytes)` |
| Disk used (/) | `100 * nodeprobe_disk_used_bytes / (nodeprobe_disk_used_bytes + nodeprobe_disk_available_bytes)` |
| Memory | used (`total - available`), available, free, total |
| Disk (/) | used, available, free, total |
| Now | `up{job="nodeprobe"}` plus the current load, free RAM and free disk |

Two details worth keeping in mind while reading it:

- **Memory used is `total - available`, not `total - free`.** Free excludes the
  page cache, which Linux hands back the moment a program asks for it, so
  `total - free` makes a healthy machine look nearly full.
- **Disk used % divides by `used + available`, not by total.** ~5% of the
  filesystem is reserved for root, and that is exactly how `df` computes `Use%`.

Edits you make in the Grafana UI (signed in as admin) survive a reload but not
a `docker compose down -v`; to keep one, use *Export → Export as JSON* and write
it back into `grafana/dashboards/nodeprobe.json`.

## If the dashboard is empty

Check <http://localhost:9090/targets>. The `nodeprobe` target should be `UP`.

- **`connection refused`** — nodeprobe is not running, or it is bound to
  something narrower than `:9200`.
- **`no such host` / timeout** — the containers cannot resolve the host.
  On WSL 2 with Docker Desktop, containers live in a separate VM and
  `host.docker.internal` points at Windows, not at your WSL distro. Either
  turn on mirrored networking (`networkingMode=mirrored` in `%UserProfile%\.wslconfig`,
  then `wsl --shutdown`), or run Docker Engine inside the WSL distro itself,
  or hardcode the distro's IP:

  ```sh
  hostname -I | awk '{print $1}'      # run in WSL, put the result in prometheus.yml
  ```

Once the target is `UP` but a panel still reads "No data", query the metric
directly in the Prometheus UI — that separates a scrape problem from a
dashboard problem.
