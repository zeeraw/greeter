build:
	docker build -t zeeraw/greeter .

push:
	docker push zeeraw/greeter:latest

run:
	docker run --rm -p 50051:50051 -p 5117:5117 zeeraw/greeter:latest server 0.0.0.0:50051
