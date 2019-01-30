build:
	docker build -t zeeraw/greeter .

push:
	docker push zeeraw/greeter:latest

run:
	docker run --rm -p 50051:50051 zeeraw/greeter:latest server 0.0.0.0:50051
