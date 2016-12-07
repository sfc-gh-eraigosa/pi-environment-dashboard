all: build run

build:
	docker build \
	--build-arg http_proxy=$(http_proxy) \
	--build-arg https_proxy=$(https_proxy) \
	-t cpumonitor .

run:
	docker run -it --rm --hostname samplehost --name samplehost --net envirodashboard_default \
		-e username=admin -e password=admin -e address=http://influxdb:8086 cpumonitor
