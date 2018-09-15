build: 
	go build
	docker build -t localhost:5000/proxy .

clean:
	go clean
	kubectl delete deployment proxy
	kubectl delete service proxy
	kubectl delete ingress proxy

push:
	docker push localhost:5000/proxy

deploy:
	kubectl apply -f deployment.yml
	kubectl apply -f service.yml 
	kubectl apply -f ingress.yml