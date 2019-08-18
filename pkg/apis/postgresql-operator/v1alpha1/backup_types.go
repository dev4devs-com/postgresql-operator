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

	// Schedule period for the CronJob  "0 0 * * *" # daily at 00:00.
	Schedule string `json:"schedule,omitempty"`

	// Image:tag used to do the backup.
	// More Info: https://github.com/integr8ly/backup-container-image
	Image string `json:"image,omitempty"`

	// Database version. (E.g 9.6).
	// IMPORTANT: Just the first 2 digits should be used.
	DatabaseVersion string `json:"databaseVersion,omitempty"`

	// Used to create the directory where the files will be stored
	ProductName string `json:"productName,omitempty"`

	// Name of AWS S3 storage.
	// Required to create the Secret with the data to allow send the backup files to AWS S3 storage.
	AwsS3BucketName string `json:"awsS3BucketName,omitempty"`

	// Key ID of AWS S3 storage.
	// Required to create the Secret with the data to allow send the backup files to AWS S3 storage.
	AwsAccessKeyId string `json:"awsAccessKeyId,omitempty"`

	// Secret/Token of AWS S3 storage.
	// Required to create the Secret with the data to allow send the backup files to AWS S3 storage.
	AwsSecretAccessKey string `json:"awsSecretAccessKey,omitempty"`

	// Name of the secret with the AWS data credentials already created in the cluster
	AwsCredentialsSecretName string `json:"awsCredentialsSecretName,omitempty"`

	// Name of the namespace where the scret with the AWS data credentials is in the cluster
	AwsCredentialsSecretNamespace string `json:"awsCredentialsSecretNamespace,omitempty"`

	// Name of the secret with the EncryptKey data already created in the cluster
	EncryptKeySecretName string `json:"encryptKeySecretName,omitempty"`

	// Name of the namespace where the secret with the EncryptKey data is in the cluster
	EncryptKeySecretNamespace string `json:"encryptKeySecretNamespace,omitempty"`

	// GPG public key to create the EncryptionKeySecret with this data
	// See here how to create this key : https://help.github.com/en/articles/generating-a-new-gpg-key
	GpgPublicKey string `json:"gpgPublicKey,omitempty"`

	// GPG email to create the EncryptionKeySecret with this data
	// See here how to create this key : https://help.github.com/en/articles/generating-a-new-gpg-key
	GpgEmail string `json:"gpgEmail,omitempty"`

	// GPG trust model to create the EncryptionKeySecret with this data. the default value is true when it is empty.
	// See here how to create this key : https://help.github.com/en/articles/generating-a-new-gpg-key
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
	EncryptionKeySecretName string `json:"encryptionKeySecretName"`

	// Namespace of the secret object with the Encryption GPG Key
	EncryptionKeySecretNamespace string `json:"encryptionKeySecretNamespace"`

	// Data of the secret object with the Encryption GPG Key
	EncryptionKeySecretData map[string]string `json:"encryptionKeySecretData"`

	// Boolean value which has true when it has an EncryptionKey to be used to send the backup files
	HasEncryptionKey bool `json:"hasEncryptionKey"`

	// Boolean value which has true when the Database Pod was found in order to create the secret with the database data to allow the backup image connect into it.
	DatabasePodFound bool `json:"databasePodFound"`

	// Boolean value which has true when the Service Database Pod was found in order to create the secret with the database data to allow the backup image connect into it.
	DatabaseServiceFound bool `json:"databaseServiceFound"`

	// Status of the CronJob object
	CronJobStatus v1beta1.CronJobStatus `json:"cronJobStatus"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Backup is the Schema for the backups API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Backup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackupSpec   `json:"spec,omitempty"`
	Status BackupStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BackupList contains a list of Backup
type BackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Backup{}, &BackupList{})
}
