service=controller
build: 
	go build
	docker build -t localhost:5000/$(service) .

clean:
	go clean
	kubectl delete deployment $(service)
	kubectl delete service $(service)

push:
	docker push localhost:5000/$(service)

deploy:
	kubectl apply -f kube/deployment.yml
	kubectl apply -f kube/service.yml 


a: clean build push deploy