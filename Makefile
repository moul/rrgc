GOPKG ?=	moul.io/rrgc
DOCKER_IMAGE ?=	moul/rrgc
GOBINS ?=	.
NPM_PACKAGES ?=	.

include rules.mk

generate: install _gen-logs-dir
	GO111MODULE=off go get github.com/campoy/embedmd
	mkdir -p .tmp
	echo 'foo@bar:~$$ rrgc -h' > .tmp/usage.txt
	(rrgc -h || true) 2>> .tmp/usage.txt
	echo 'foo@bar:~$$ ls logs' >> .tmp/usage.txt
	ls logs >> .tmp/usage.txt
	echo 'foo@bar:~$$ rrgc 24h,5 1h,5 -- ./logs/*.log | xargs rm -v' >> .tmp/usage.txt
	rrgc 24h,5 1h,5 -- "./logs/*.log" | xargs rm -v >> .tmp/usage.txt
	echo 'foo@bar:~$$ rrgc 24h,5 1h,5 -- ./logs/*.log' >> .tmp/usage.txt
	rrgc 24h,5 1h,5 -- "./logs/*.log" >> .tmp/usage.txt
	echo 'foo@bar:~$$ ls logs' >> .tmp/usage.txt
	ls logs >> .tmp/usage.txt
	go doc -all ./rrgc > .tmp/godoc.txt
	embedmd -w README.md
	rm -rf .tmp
.PHONY: generate

_gen-logs-dir:
	rm -rf logs
	mkdir -p logs
	touch -t 8001031305 logs/A.log
	touch -t 8001031335 logs/B.log
	touch -t 8001031404 logs/C.log
	touch -t 8001031406 logs/D.log
	touch -t 8001041706 logs/E.log
	touch -t 8001041606 logs/F.log
	touch -t 8001041506 logs/G.log
	touch -t 8001041436 logs/H.log
	touch -t 8001041406 logs/I.log
	touch -t 8001051406 logs/J.log
	touch -t 8001061406 logs/K.log
	touch -t 8001131406 logs/L.log
	touch -t 8001161406 logs/M.log
	touch -t 8001231406 logs/N.log
	touch -t 8001261406 logs/O.log
	touch -t 8002011406 logs/P.log
	touch -t 8003011406 logs/Q.log
	touch -t 8004011406 logs/R.log
