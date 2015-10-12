prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-geojson; then rm -rf src/github.com/whosonfirst/go-whosonfirst-geojson; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-geojson
	cp -r whosonfirst src/github.com/whosonfirst/go-whosonfirst-geojson/whosonfirst

deps:   self
	go get -u "github.com/jeffail/gabs"
	go get -u "github.com/dhconnelly/rtreego"

fmt:
	go fmt bin/dump.go
	go fmt whosonfirst/geojson.go

dump:	self
	go build -o bin/dump bin/dump.go
