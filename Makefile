NAME := imgeventfilter
VERSION := v0.1

hookurl := "http://127.0.0.1:8080/v1/image/webhook"

all: run

clean:
	rm -rf ./$(NAME)

image: clean
	docker run --rm -v $(shell pwd):/go/src/$(NAME)  -w /go/src/$(NAME)  golang:1.8-alpine go build -o $(NAME)
	docker build -t registry.paas/library/$(NAME):$(VERSION) .

run: image
	docker run -dit --restart always --name $(NAME) -e hookurl=$(hookurl)  -v /etc/localtime:/etc/localtime --net host registry.paas/library/$(NAME):$(VERSION)

package: 
	docker save -o $(NAME)_$(VERSION).tar registry.paas/library/$(NAME):$(VERSION)
kill:
	docker rm -vf $(NAME)
