FROM grafana/grafana:latest

COPY ./grafana/dashboards/default.yml /etc/grafana/provisioning/dashboards/
COPY ./grafana/datasources/default.yml  /etc/grafana/provisioning/datasources/
COPY ./grafana/dashboards/*.json /etc/grafana/provisioning/dashboards/
