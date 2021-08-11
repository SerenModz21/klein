MAIN 	:= main.go

.PHONY: tidy 

default: debug

start: $(MAIN)
	@go build $(MAIN)
	@./main &
	@cd web && yarn build |> /dev/null
	@./tools/go-live -dir web/dist &
	@echo To stop the process, close the terminal

debug: $(MAIN)
	@systemctl status mongodb.service | grep Active | cut -d: -f2 | xargs -I {} echo "mongo :: {}"
	@systemctl status redis.service | grep Active | cut -d: -f2 | xargs -I {} echo "redis :: {}"
	@echo
	
	@go run $(MAIN) &
	@cd web && yarn serve

tidy:
	@go mod tidy
	@rm -rf ./main ./web/dist