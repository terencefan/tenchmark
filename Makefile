image:
	docker build -t stdrickforce/tenchmark .
	docker push stdrickforce/tenchmark

linux:
	docker run -v `pwd`:/tmp stdrickforce/tenchmark cp /usr/bin/tenchmark /tmp/tenchmark
