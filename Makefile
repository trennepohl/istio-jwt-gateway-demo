.DEFAULT_GOAL := help
.PHONY: help istio create-cluster build apply-manifests

export KUBECONFIG=kubeconfig

## This is just a function that lists all targets in this Makefile and pretty prints the help comments
help:
		@IFS=$$'\n' ; \
    	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##/:/'`); \
    	printf "%-30s %s\n" "target" "help" ; \
    	printf "%-30s %s\n" "------" "----" ; \
    	for help_line in $${help_lines[@]}; do \
    		IFS=$$':' ; \
    		help_split=($$help_line) ; \
    		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
    		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
    		printf '\033[36m'; \
    		printf "%-30s %s" $$help_command ; \
    		printf '\033[0m'; \
    		printf "%s\n" $$help_info; \
    	done


install: create-cluster certs build istio apply-manifests ## Setup everything

istio: ## Install istio charts
	@istioctl install -y

create-cluster: ## Create a local kubernetes cluster
	@kind create cluster --config kind-config.yaml
	@kind export kubeconfig --name istio-oauth-poc > kubeconfig

build: ## Build and deploy demo services
	@docker build -t authorization-service -f authorization/Dockerfile authorization/
	@kind load docker-image authorization-service --name istio-oauth-poc
	@docker build -t samplesvc -f samplesvc/Dockerfile samplesvc/
	@kind load docker-image samplesvc --name istio-oauth-poc

apply-manifests: ## Install all components of this demo
	@for m in manifests; do kubectl -f $$m apply; done
	@echo "Waiting for Authorization svc to be ready"
	@kubectl -n default wait --for=condition=ready pod -l app=authorization
	@kubectl -n istio-system rollout restart deploy/istiod

proxy: ## Port-forward the istio ingress
	@kubectl port-forward svc/istio-ingressgateway 8080:http2 -n istio-system

destroy: ## Destroy the cluster
	@kind delete cluster --name istio-oauth-poc

certs:
	@openssl genrsa -out authorization/certs/private-key.pem 2048
	@openssl rsa -in authorization/certs/private-key.pem -pubout -out authorization/certs/public-key.pem