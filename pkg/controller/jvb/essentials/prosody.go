package essentials

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// prosody deployments

func DeploymentForProsody(namespace string) *appsv1.Deployment {
	replicas := int32(0)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prosody",
			Namespace: namespace,
			Labels: map[string]string{
				"name": "prosody",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": "prosody",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": "prosody",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "prosody-service-account",
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "prosody",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Image: "jitsi/prosody:4101",
							Name:  "prosody",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/config",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "c2s",
									ContainerPort: 5222,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "bosh",
									ContainerPort: 5280,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "component",
									ContainerPort: 5347,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									ConfigMapRef: &corev1.ConfigMapEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "jitsi",
										},
									},
								},
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "jitsi",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return dep
}

func PvcForProsody(namespace string) *corev1.PersistentVolumeClaim {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prosody",
			Namespace: namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: resource.MustParse("100Mi"),
				},
			},
		},
	}
	return pvc
}

func ServiceForProsody(namespace string) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prosody",
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "c2s",
					Port:     5222,
					Protocol: corev1.ProtocolTCP,
					TargetPort: intstr.IntOrString{
						IntVal: 5222,
					},
				},
				{
					Name:     "bosh",
					Port:     5280,
					Protocol: corev1.ProtocolTCP,
					TargetPort: intstr.IntOrString{
						IntVal: 5280,
					},
				},
				{
					Name:     "component",
					Port:     5347,
					Protocol: corev1.ProtocolTCP,
					TargetPort: intstr.IntOrString{
						IntVal: 5347,
					},
				},
			},
			Type: corev1.ServiceType("LoadBalancer"),
			Selector: map[string]string{
				"name": "prosody",
			},
		},
	}
	return service
}

func ServiceAccountForProsody(namespace string) *corev1.ServiceAccount {
	service := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prosody-service-account",
			Namespace: namespace,
		},
	}
	return service
}
