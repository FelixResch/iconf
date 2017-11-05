
all: copy_assets build

build:
	go-bindata -debug ./assets/... ./templates/...
	go build

deps:
	got get -u github.com/jteeuwen/go-bindata/...

copy_assets:
	bower install -p
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