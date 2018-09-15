setup:
	kubectl apply -f traefik-rbac.yaml
	kubectl apply -f traefik-ds.yaml
	kubectl apply -f ui.yaml
	echo "$(minikube ip) traefik-ui.minikube" | sudo tee -a /etc/hosts

