docker-build:
	docker build -t phanikmullapudi/iscalecc-prometheus-exporter .

docker-push:
	docker push phanikmullapudi/iscalecc-prometheus-exporter

deploy:
	helm upgrade --install iscalecc-prometheus-exporter ./chart --namespace=default