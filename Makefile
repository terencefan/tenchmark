alpine:
	docker build -t stdrickforce/tenchmark -f .docker/alpine.df .
	docker tag stdrickforce/tenchmark stdrickforce/tenchmark:alpine
	docker push stdrickforce/tenchmark
	docker push stdrickforce/tenchmark:alpine

centos: # ubuntu
	docker build -t stdrickforce/tenchmark:centos -f .docker/centos.df .
	docker push stdrickforce/tenchmark:centos

image: alpine centos

bin:
	docker run -v `pwd`:/tmp stdrickforce/tenchmark:centos cp /usr/bin/tenchmark /tmp/tenchmark

bin-alpine:
	docker run -v `pwd`:/tmp stdrickforce/tenchmark:alpine cp /usr/bin/tenchmark /tmp/tenchmark
