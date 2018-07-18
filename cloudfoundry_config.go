package main

import (
	"os"
	"log"

	"github.com/jhunt/vcaptive"
)

func GetAMQPUriToUse() string {
	result := "amqp://localhost:5672/"

	vcapServices, ok := os.LookupEnv("VCAP_SERVICES")
	if ok {
		services, err := vcaptive.ParseServices(vcapServices)
		if err != nil {
			log.Fatalf("VCAP_SERVICES: %s", err)
		}

		instance, found := services.Tagged("amqp", "rabbitmq")
		if !found {
			log.Fatalf("VCAP_SERVICES: No service tagged 'amqp' or 'rabbitmq' bound! %s", err)
		}

		result, found = instance.GetString("uri")
		if !found {
			log.Fatalf("VCAP_SERVICES: No credential 'uri' configured! %s", err)
		}
	}

	return result
}
