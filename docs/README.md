`grafana.png` is the screenshot used in the main README: the **nodeprobe**
dashboard from `deploy/`, taken at 1280px wide after the stack has been
collecting for a while, so the graphs have some history in them.

To retake it: start nodeprobe, run `docker compose up -d` in `deploy/`, open
<http://localhost:3000>, wait for the time range to fill, and save the
dashboard screenshot here as `grafana.png`.
