package v1alpha1

import (
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BackupSpec defines the desired state of Backup
// +k8s:openapi-gen=true
type BackupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Name of the PostgreSQL CR applied which this backup will work with
	PostgresqlCRName string `json:"postgresqlCRName,omitempty"`

	// Schedule period for the CronJob.
	// Default Value: <0 0 * * *> daily at 00:00
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	Schedule string `json:"schedule,omitempty"`

	// Image:tag used to do the backup.
	// Default Value: <quay.io/integreatly/backup-container:1.0.8>
	// More Info: https://github.com/integr8ly/backup-container-image
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Image:tag"
	Image string `json:"image,omitempty"`

	// Database version. (E.g 9.6).
	// Default Value: <9.6>
	// IMPORTANT: Just the first 2 digits should be used.
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="PostgreSQLversion"
	DatabaseVersion string `json:"databaseVersion,omitempty"`

	// Used to create the directory where the files will be stored
	// Default Value: <postgresql>
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="AWS tag name"
	ProductName string `json:"productName,omitempty"`

	// Name of AWS S3 storage.
	// Default Value: nil
	// Required to create the Secret with the AWS data to allow send the backup files to AWS S3 storage.
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="AWS S3 Bucket name"
	AwsS3BucketName string `json:"awsS3BucketName,omitempty"`

	// Key ID of AWS S3 storage.
	// Default Value: nil
	// Required to create the Secret with the data to allow send the backup files to AWS S3 storage.
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="AWS S3 accessKey/token ID"
	AwsAccessKeyId string `json:"awsAccessKeyId,omitempty"`

	// Secret/Token of AWS S3 storage.
	// Default Value: nil
	// Required to create the Secret with the data to allow send the backup files to AWS S3 storage.
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="AWS S3 accessKey/token"
	AwsSecretAccessKey string `json:"awsSecretAccessKey,omitempty"`

	// Name of the secret with the AWS data credentials pre-existing in the cluster
	// Default Value: nil
	// See here the template: https://github.com/integr8ly/backup-container-image/blob/master/templates/openshift/sample-config/s3-secret.yaml
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="AWS Secret name:"
	AwsSecretName string `json:"awsSecretName,omitempty"`

	// Namespace of the secret with the AWS data credentials pre-existing in the cluster
	// Default Value: nil
	// NOTE: If the namespace be not informed then the operator will try to find it in the same namespace where it is applied
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="AWS Secret namespace:"
	AwsSecretNamespace string `json:"awsSecretNamespace,omitempty"`

	// Name of the secret with the Encrypt data pre-existing in the cluster
	// Default Value: nil
	// See here the template: https://github.com/integr8ly/backup-container-image/blob/master/templates/openshift/sample-config/gpg-secret.yaml
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="EncryptKey Secret name:"
	EncryptKeySecretName string `json:"encryptKeySecretName,omitempty"`

	// Namespace of the secret with the Encrypt data pre-existing in the cluster
	// Default Value: nil
	// NOTE: If the namespace be not informed then the operator will try to find it in the same namespace where it is applied
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="EncryptKey Secret namespace:"
	EncryptKeySecretNamespace string `json:"encryptKeySecretNamespace,omitempty"`

	// GPG public key to create the EncryptionKeySecret with this data
	// Default Value: nil
	// See here how to create this key : https://help.github.com/en/articles/generating-a-new-gpg-key
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Gpg public key:"
	GpgPublicKey string `json:"gpgPublicKey,omitempty"`

	// GPG email to create the EncryptionKeySecret with this data
	// Default Value: nil
	// See here how to create this key : https://help.github.com/en/articles/generating-a-new-gpg-key
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Gpg public email:"
	GpgEmail string `json:"gpgEmail,omitempty"`

	// GPG trust model to create the EncryptionKeySecret with this data. the default value is true when it is empty.
	// Default Value: nil
	// See here how to create this key : https://help.github.com/en/articles/generating-a-new-gpg-key
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Gpg trust model:"
	GpgTrustModel string `json:"gpgTrustModel,omitempty"`
}

// BackupStatus defines the observed state of Backup
// +k8s:openapi-gen=true
type BackupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Will be as "OK when all objects are created successfully
	BackupStatus string `json:"backupStatus"`

	// Name of the CronJob object created and managed by it to schedule the backup job
	CronJobName string `json:"cronJobName"`

	// Name of the secret object created with the database data to allow the backup image connect to the database
	DBSecretName string `json:"dbSecretName"`

	// Data  of the secret object created with the database data to allow the backup image connect to the database
	DBSecretData map[string]string `json:"dbSecretData"`

	// Name  of the secret object with the Aws data to allow send the backup files to the AWS storage
	AWSSecretName string `json:"awsSecretName"`

	// Data  of the secret object with the Aws data to allow send the backup files to the AWS storage
	AWSSecretData map[string]string `json:"awsSecretData"`

	// Namespace  of the secret object with the Aws data to allow send the backup files to the AWS storage
	AwsCredentialsSecretNamespace string `json:"awsCredentialsSecretNamespace"`

	// Name  of the secret object with the Encryption GPG Key
	EncryptKeySecretName string `json:"encryptKeySecretName"`

	// Namespace of the secret object with the Encryption GPG Key
	EncryptKeySecretNamespace string `json:"encryptKeySecretNamespace"`

	// Data of the secret object with the Encryption GPG Key
	EncryptKeySecretData map[string]string `json:"encryptKeySecretData"`

	// Boolean value which has true when it has an EncryptionKey to be used to send the backup files
	HasEncryptKey bool `json:"hasEncryptKey"`

	// Boolean value which has true when the Database Pod was found in order to create the secret with the database data to allow the backup image connect into it.
	IsDatabasePodFound bool `json:"isDatabasePodFound"`

	// Boolean value which has true when the Service Database Pod was found in order to create the secret with the database data to allow the backup image connect into it.
	IsDatabaseServiceFound bool `json:"isDatabaseServiceFound"`

	// Status of the CronJob object
	CronJobStatus v1beta1.CronJobStatus `json:"cronJobStatus"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Backup represents the backup service in the cluster
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +operator-sdk:gen-csv:customresourcedefinitions.displayName="Backup"
// +operator-sdk:gen-csv:customresourcedefinitions.resources="CronJob,v1beta1,\"A Kubernetes CronJob\""
// +operator-sdk:gen-csv:customresourcedefinitions.resources="Secret,v1,\"A Kubernetes Secret\""
// +operator-sdk:gen-csv:customresourcedefinitions.resources="Service,v1,\"A Kubernetes Service\""
// +operator-sdk:gen-csv:customresourcedefinitions.resources="Pod,v1,\"A Kubernetes Pod\""
type Backup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupSpec   `json:"spec,omitempty"`
	Status BackupStatus `json:"status,omitempty"`
}

// BackupList contains a list of Backup
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Backup{}, &BackupList{})
}
