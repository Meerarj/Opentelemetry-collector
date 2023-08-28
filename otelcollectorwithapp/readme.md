Created sample golang app using gin and gorm 

Created a collector using config otel.yaml

---- Receivers are used to get data into the collector

---- Here the receivers are http and grpc

---- Exporters are used to export data to an observability backend like SigNoz.

Run the collector using ---> otelcol --config otel.yaml

Start the sample app ---> go run main.go

Trigger from the browser ---> http://localhost:8092/books

Check the app output and collector output

Check the signoz ui 

In signoz ui ,able to see the sample app
