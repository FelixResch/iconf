
all: copy_assets build

build-debug:
	go-bindata -debug ./assets/... ./templates/...
	go build

build:
	go-bindata ./assets/... ./templates/...
	go build

deps:
	go get -u github.com/jteeuwen/go-bindata/...

copy_assets:
	bower --allow-root install -p
	mkdir -p assets/js
	mkdir -p assets/css
	mkdir -p assets/fonts
	mkdir -p assets/templates
	cp -r bower_components/jquery/dist/jquery.min.* assets/js
	cp -r bower_components/font-awesome/css/*.min.css assets/css
	cp -r bower_components/font-awesome/css/*.css.map assets/css
	cp -r bower_components/font-awesome/fonts/* assets/fonts

run:
	@sudo ./iconf