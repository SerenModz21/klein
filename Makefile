MAIN 	:= main.go

.PHONY: tidy 

default: debug

start: $(MAIN)
	@go build $(MAIN)
	@./main

debug: $(MAIN)
	@systemctl status mongodb.service | grep Active | cut -d: -f2 | xargs -I {} echo "mongo :: {}"
	@systemctl status redis.service | grep Active | cut -d: -f2 | xargs -I {} echo "redis :: {}"
	@echo
	
	@go run $(MAIN)

tidy:
	@go mod tidy
	@rm -rf ./main