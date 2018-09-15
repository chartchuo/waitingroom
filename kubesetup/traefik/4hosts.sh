#/bin/bash
echo "$(minikube ip) traefik-ui.minikube" | sudo tee -a /etc/hosts