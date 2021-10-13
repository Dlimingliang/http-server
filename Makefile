project=httpserver
version=v1.0

version:
	echo ${version}

build: version
	echo "building httpserver binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64/httpserver .

build-mac: version
	echo "building httpserver binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/amd64/httpserver-mac .

release: build
	echo "building httpserver container"
	docker build -t limingliang/${project}:${version} .

push: release
	echo "pushing limingliang/httpserver"
	docker push limingliang/${project}:${version}



