build: 
	go build
	docker build -t localhost:5000/mockserver .

clean:
	go clean
	kubectl delete deployment mockserver
	kubectl delete service mockserver
	kubectl delete hpa mockserver

push:
	docker push localhost:5000/mockserver

deploy:
	kubectl apply -f deployment.yml
	kubectl apply -f service.yml
	kubectl apply -f ingress.yml

autoscale:
	kubectl autoscale deployment mockserver --min=1 --max=3 --cpu-percent=50
