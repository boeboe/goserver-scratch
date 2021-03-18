# Makefile

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

CHART_DIR := ./goserver
NAMESPACE := example

install: istioctl  ## Install microservices
	kubectl create namespace ${NAMESPACE} || true
	kubectl label namespace ${NAMESPACE} istio-injection=enabled --overwrite
	helm install service1 ${CHART_DIR} --namespace ${NAMESPACE} --values example/service1.yaml || true
	helm install service2 ${CHART_DIR} --namespace ${NAMESPACE} --values example/service2.yaml || true
	helm install service3 ${CHART_DIR} --namespace ${NAMESPACE} --values example/service3.yaml || true
	helm install service4 ${CHART_DIR} --namespace ${NAMESPACE} --values example/service4.yaml || true
	helm install service5 ${CHART_DIR} --namespace ${NAMESPACE} --values example/service5.yaml || true

upgrade:  ## Upgrade microservices
	helm upgrade service1 ${CHART_DIR} --namespace ${NAMESPACE} --values example/service1.yaml || true
	helm upgrade service2 ${CHART_DIR} --namespace ${NAMESPACE} --values example/service2.yaml || true
	helm upgrade service3 ${CHART_DIR} --namespace ${NAMESPACE} --values example/service3.yaml || true
	helm upgrade service4 ${CHART_DIR} --namespace ${NAMESPACE} --values example/service4.yaml || true
	helm upgrade service5 ${CHART_DIR} --namespace ${NAMESPACE} --values example/service5.yaml || true

uninstall:  ## Uninstall microservices 
	helm uninstall service5 --namespace ${NAMESPACE} || true
	helm uninstall service4 --namespace ${NAMESPACE} || true
	helm uninstall service3 --namespace ${NAMESPACE} || true
	helm uninstall service2 --namespace ${NAMESPACE} || true
	helm uninstall service1 --namespace ${NAMESPACE} || true
	kubectl delete namespace ${NAMESPACE} || true

restart: ## Restart all microservices
	kubectl -n ${NAMESPACE} rollout restart deploy

view: ## View all microservices
	kubectl -n ${NAMESPACE} get services -o wide
	kubectl -n ${NAMESPACE} get pods -o wide

k6_install:
	kubectl  -n ${NAMESPACE} apply -f example/k6.yaml

k6_uninstall:
	kubectl  -n ${NAMESPACE} delete -f example/k6.yaml
