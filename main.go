package main

import (
	"encoding/json"
	"fmt"
	"os"

	functionio "github.com/crossplane/crossplane-functions-go"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CompositeSpec struct {
	NamespaceNames []string `json:"namespace-names"`
}

func main() {
	// Parse the input JSON from Crossplane
	var input functionio.FunctionIO
	err := json.NewDecoder(os.Stdin).Decode(&input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding input: %v\n", err)
		os.Exit(1)
	}

	// Extract the spec from the composite resource
	var compositeSpec CompositeSpec
	if err := json.Unmarshal(input.Resource.Spec.Raw, &compositeSpec); err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding spec: %v\n", err)
		os.Exit(1)
	}

	// Create Namespace resources based on namespace-names array
	var resources []functionio.OutputResource
	for _, name := range compositeSpec.NamespaceNames {
		namespace := &v1.Namespace{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Namespace",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"example": "true",
				},
			},
		}

		// Add the Namespace resource to the output list
		resources = append(resources, functionio.OutputResource{
			Resource: namespace,
		})
	}

	// Write the generated resources back to Crossplane
	output := functionio.FunctionIO{
		Resources: resources,
	}

	// Output the generated resources to stdout
	if err := json.NewEncoder(os.Stdout).Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding output: %v\n", err)
		os.Exit(1)
	}
}
