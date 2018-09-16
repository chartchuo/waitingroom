service=controller
build: 
	go build
	docker build -t localhost:5000/$(service) .

clean:
	go clean
	kubectl delete deployment $(service)
	kubectl delete service $(service)
	kubectl delete ingress $(service)

cleanhpa:
	kubectl delete hpa $(service)

cleanall: clean cleanhpa

push:
	docker push localhost:5000/$(service)

deploy:
	kubectl apply -f deployment.yml
	kubectl apply -f service.yml 
	kubectl apply -f ingress.yml

autoscale:
	kubectl autoscale deployment $(service) --min=1 --max=10 --cpu-percent=50

a: clean build push deploy