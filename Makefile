k6: go.mod go.sum *.go server/*.go logger/*.go metrics/*.go
	xk6 build --with github.com/temporalio/xk6-temporalite=.
