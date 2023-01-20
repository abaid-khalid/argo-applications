package test

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	v1 "k8s.io/api/apps/v1"
	"testing"
)

func TestKafkaConnect(t *testing.T) {

	options := &helm.Options{
		KubectlOptions: k8s.NewKubectlOptions("", "", "confluent"),
	}

	templates := helm.RenderTemplate(t, options, "../", "kafka-connect", []string{"templates/deployment.yaml"})
	var deployment v1.Deployment
	helm.UnmarshalK8SYaml(t, templates, &deployment)

	noOfContainers := len(deployment.Spec.Template.Spec.Containers)
	if noOfContainers != 1 {
		t.Errorf("Expected 1 container, got: %d", noOfContainers)
	}
	checkScrapeAnnotation(t, deployment.Spec.Template.Annotations, "5556")
}

func checkScrapeAnnotation(t *testing.T, annotations map[string]string, expectedScrapePort string) {
	scrapePort, ok := annotations["prometheus.io/port"]
	if !ok {
		t.Errorf("please specify prometheus.io/port annotation")
		return
	}
	if scrapePort != expectedScrapePort {
		t.Errorf("expected scrape port %s, got %s", expectedScrapePort, scrapePort)
	}
}
