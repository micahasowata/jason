# test runs all tests, prepare coverage data and displays it in a web page
test:
	go test . -coverprofile=${name}.out && go tool cover -html=${name}.out 