package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	port = ":443"
)

type Namespace struct {
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

func UnmarshallConfigMap(file string) []Namespace {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened namespaces.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var namespaces []Namespace
	json.Unmarshal(byteValue, &namespaces)
	return namespaces

}

func ValidateNamespace(w http.ResponseWriter, r *http.Request) {
	admissionReview := v1beta1.AdmissionReview{}
	namespaceControl := false
	limitControl := true
	if err := json.NewDecoder(r.Body).Decode(&admissionReview); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if admissionReview.Request == nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	raw := admissionReview.Request.Object.Raw
	deployment := v1.Deployment{}
	if err := json.Unmarshal(raw, &deployment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if deployment.Name == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	admissionReview.Response = &v1beta1.AdmissionResponse{
		UID:     admissionReview.Request.UID,
		Allowed: false,
	}
	admissionReview.Response.Result = &metav1.Status{
		Message: "Invalid Namespace",
	}

	for _, namespace := range UnmarshallConfigMap("./config/namespaces.json") {
		if namespace.Metadata.Name == deployment.Namespace {
			namespaceControl = true
			break
		}
	}

	for _, container := range deployment.Spec.Template.Spec.Containers {

		if container.Resources.Limits.Cpu == nil {
			limitControl = false
			break
		}

	}

	if namespaceControl == true && limitControl == true {
		admissionReview.Response.Allowed = true
		admissionReview.Response.Result = &metav1.Status{
			Message: "This namespace is allowed and limits are given to Deployment",
		}
	}
	fmt.Println(deployment)
	fmt.Println(admissionReview.Response.Result.Message)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&admissionReview)
}

func main() {
	http.HandleFunc("/validatens", ValidateNamespace)
	log.Fatal(http.ListenAndServeTLS(port, "./certs/cert.pem", "./certs/key.pem", nil))

}
