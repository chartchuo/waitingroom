build: 
	go build
	docker build -t chartchuo/mockserver .

clean:
	go clean
	kubectl delete deployment mockserver
	kubectl delete service mockserver

push:
	docker push chartchuo/mockserver

deploy:
	kubectl apply -f deployment.yml
	kubectl expose deployment mockserver --port=8888 --type=LoadBalancer