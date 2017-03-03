ALL= backup
clean:
	go clean -i -r
	rm -rf $(ALL) nohup.out
	  
build:
	go build   -o $(ALL)  -ldflags  "-s -w"  -a
run:
	./$(ALL) 
