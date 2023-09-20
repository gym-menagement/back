tag=latest

all: server

server: dummy
	buildtool-model ./ 
	buildtool-router ./ > ./router/router.go
	go build -o bin/gym main.go

fswatch:
	fswatch -0 controllers | xargs -0 -n1 build/notify.sh

run:
	gin --port 9000 -a 9004 --bin bin/gym run main.go

allrun:
	fswatch -0 controllers | xargs -0 -n1 build/notify.sh &
	gin --port 9000 -a 9004 --bin bin/gym run main.go

test: dummy
	go test -v ./...

linux:
	env GOOS=linux GOARCH=amd64 go build -o bin/gym.linux main.go

dockerbuild:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o bin/gym.linux main.go

docker: dockerbuild
	docker build -t kobums/gym:$(tag) .

dockerrun:
	docker run -d --name="gym" -p 9004:9004 kobums/gym

push: docker
	docker push kobums/gym:$(tag)

clean:
	rm -f bin/gym

dummy: