PROJECT_ID := "fifth-cab-359408"
REGISTRY := "eu.gcr.io"
SERVICE := "notification-manager"
NAMESPACE:= "notification-manager"

dependencies:
	go get -u github.com/Traders-Connect/utils
	go get -u github.com/Traders-Connect/utils/grpc

	go mod vendor
	go mod tidy

test:
	go test -v ./...

run:
	go run cmd/$(SERVICE)/main.go

build:
	go build -o bin/$(SERVICE) cmd/$(SERVICE)/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(SERVICE) cmd/$(SERVICE)/main.go

build-image: build-linux
	docker build . -t traders-connect/$(SERVICE) -t $(REGISTRY)/$(PROJECT_ID)/$(SERVICE)

publish-image: build-image
	docker push $(REGISTRY)/$(PROJECT_ID)/$(SERVICE)

build-dev: build-linux
	docker build . -t traders-connect/$(SERVICE) -t $(REGISTRY)/$(PROJECT_ID)/$(SERVICE):dev

publish-dev: build-dev
	docker push $(REGISTRY)/$(PROJECT_ID)/$(SERVICE):dev

deploy:
	HTTPS_PROXY=localhost:8888 kubectl get namespace $(NAMESPACE) > /dev/null || kubectl create namespace $(NAMESPACE)
	HTTPS_PROXY=localhost:8888 helm upgrade --install -n $(NAMESPACE) -f deployment/$(SERVICE)/values.yaml $(SERVICE) ./deployment/$(SERVICE)
