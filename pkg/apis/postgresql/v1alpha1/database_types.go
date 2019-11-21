package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DatabaseSpec defines the desired state of Database
// +k8s:openapi-gen=true
type DatabaseSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Value for the Database Environment Variable (spec.databaseNameKeyEnvVar).
	// Default value: example
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Database Name"
	DatabaseName string `json:"databaseName,omitempty"`

	// Value for the Database Environment Variable (spec.databasePasswordKeyEnvVar).
	// Default value: postgres
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Database Password"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:password"
	DatabasePassword string `json:"databasePassword,omitempty"`

	// Value for the Database Environment Variable (spec.databaseUserKeyEnvVar).
	// Default value: postgres
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Database User"
	DatabaseUser string `json:"databaseUser,omitempty"`

	// Key Value for the Database Environment Variable in order to inform the database mame
	// Note that each database version/image can expected a different value for it.
	// Default value: nil
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="EnvVar Key (Database Name)"
	DatabaseNameKeyEnvVar string `json:"databaseNameKeyEnvVar,omitempty"`

	// Key Value for the Database Environment Variable in order to inform the database password
	// Note that each database version/image can expected a different value for it.
	// Default value: nil
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="EnvVar Key (Database Password)"
	DatabasePasswordKeyEnvVar string `json:"databasePasswordKeyEnvVar,omitempty"`

	// Key Value for the Database Environment Variable in order to inform the database user
	// Note that each database version/image can expected a different value for it.
	// Default value: nil
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="EnvVar Key (Database User)"
	DatabaseUserKeyEnvVar string `json:"databaseUserKeyEnvVar,omitempty"`

	// Value for the Database Environment Variable in order to define the port which it should use. It will be used in its container as well
	DatabasePort int32 `json:"databasePort,omitempty"`

	// Quantity of instances
	// Default value: 1
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Size int32 `json:"size,omitempty"`

	// Database image:tag
	// Default value: centos/postgresql-96-centos7
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Image:tag"
	Image string `json:"image,omitempty"`

	// Name to create the Database container
	ContainerName string `json:"containerName,omitempty"`

	// Limit of Memory which will be available for the database container
	// Default value: 512Mi
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Database Memory Limit"
	DatabaseMemoryLimit string `json:"databaseMemoryLimit,omitempty"`

	// Limit of Memory Request which will be available for the database container
	// Default value: 128Mi
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Database Memory Request"
	DatabaseMemoryRequest string `json:"databaseMemoryRequest,omitempty"`

	// Limit of Storage Request which will be available for the database container
	// Default value: 1Gi
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Database Storage Request"
	DatabaseStorageRequest string `json:"databaseStorageRequest,omitempty"`

	// CPU resource request which will be available for the database container
	// Default value: 10Mi
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Database CPU"
	DatabaseCpu string `json:"databaseCpu,omitempty"`

	// Limit of CPU request which will be available for the database container
	// Default value: 20Mi
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Database CPU Limit"
	DatabaseCpuLimit string `json:"databaseCpuLimit,omitempty"`

	// Policy definition to pull the Database Image
	// More info: https://kubernetes.io/docs/concepts/containers/images/
	// Default value:
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Container Image Pull Policy"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:imagePullPolicy"
	ContainerImagePullPolicy v1.PullPolicy `json:"containerImagePullPolicy,omitempty"`

	// Name of the ConfigMap where the operator should looking for the EnvVars keys and/or values only
	// Default value: nil
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="ConfigMap name"
	ConfigMapName string `json:"configMapName,omitempty"`

	// Name of the configMap key where the operator should looking for the value for the database name for its env var
	// Default value: nil
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="ConfigMap Database Key"
	ConfigMapDatabaseNameKey string `json:"configMapDatabaseNameKey,omitempty"`

	// Name of the configMap key where the operator should looking for the value for the database user for its env var
	// Default value: nil
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="ConfigMap Password Key"
	ConfigMapDatabasePasswordKey string `json:"configMapDatabasePasswordKey,omitempty"`

	// Name of the configMap key where the operator should looking for the value for the database password for its env var
	// Default value: nil
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="ConfigMap User Key"
	ConfigMapDatabaseUserKey string `json:"configMapDatabaseUserKey,omitempty"`
}

// DatabaseStatus defines the observed state of Database
// +k8s:openapi-gen=true
type DatabaseStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Name of the PersistentVolumeClaim created and managed by it
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="v1.PersistentVolumeClaimStatus"
	PVCStatus v1.PersistentVolumeClaimStatus `json:"pvcStatus"`

	// Status of the Database Deployment created and managed by it
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="appsv1.DeploymentStatus"
	DeploymentStatus appsv1.DeploymentStatus `json:"deploymentStatus"`

	// Status of the Database Service created and managed by it
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="v1.ServiceStatus"
	ServiceStatus v1.ServiceStatus `json:"serviceStatus"`

	// It will be as "OK when all objects are created successfully
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Database Status"
	DatabaseStatus string `json:"databaseStatus"`
}

// Database is the Schema for the the Database Database API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +operator-sdk:gen-csv:customresourcedefinitions.displayName="Database Database"
// +operator-sdk:gen-csv:customresourcedefinitions.resources="Deployment,v1,\"A Kubernetes Deployment\""
// +operator-sdk:gen-csv:customresourcedefinitions.resources="Service,v1,\"A Kubernetes Service\""
// +operator-sdk:gen-csv:customresourcedefinitions.resources="PersistentVolumeClaim,v1,\"A Kubernetes PersistentVolumeClaim\""
type Database struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseSpec   `json:"spec,omitempty"`
	Status DatabaseStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DatabaseList contains a list of Database
type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Database `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Database{}, &DatabaseList{})
}
