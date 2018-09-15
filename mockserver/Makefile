build: 
	go build
	docker build -t localhost:5000/mockserver .

clean:
	go clean
	kubectl delete deployment mockserver
	kubectl delete service mockserver

push:
	docker push localhost:5000/mockserver

deploy:
	kubectl apply -f deployment.yml
	kubectl apply -f service.yml
	kubectl apply -f ingress.yml
