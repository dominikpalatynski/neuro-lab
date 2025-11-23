deploy-mqtt-broker:
	helm install mqtt ./charts/mqtt-broker \
	--namespace mqtt \
	--create-namespace \
	--set auth.enabled=false \
	--set persistence.storageClass=local-path

uninstall-mqtt-broker:
	helm uninstall mqtt

deploy-otel-collector:
	helm install opentelemetry-collector ./charts/opentelemetry -f charts/opentelemetry/values.yaml

uninstall-otel-collector:
	helm uninstall opentelemetry-collector

deploy-prometheus:
	helm install prometheus ./charts/prometheus -f charts/prometheus/values.yaml

upgrade-prometheus:
	helm upgrade prometheus ./charts/prometheus -f charts/prometheus/values.yaml

uninstall-prometheus:
	helm uninstall prometheus

deploy-metallb:
	helm repo add metallb https://metallb.github.io/metallb
	helm repo update
	helm install metallb metallb/metallb --set overrideNamespace=metallb