package essentials

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeploymentForJicofo(namespace string) *appsv1.Deployment {
	replicas := int32(0)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jicofo",
			Namespace: namespace,
			Labels: map[string]string{
				"name": "jicofo",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": "jicofo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": "jicofo",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "jicofo-service-account",
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Image: "jitsi/jicofo:4101",
							Name:  "jicofo",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/config",
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
							},
							Env: []corev1.EnvVar{
								{
									Name: "JICOFO_COMPONENT_SECRET",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "jitsi",
											},
											Key: "JICOFO_COMPONENT_SECRET",
										},
									},
								},
								{
									Name: "JICOFO_AUTH_USER",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "jitsi",
											},
											Key: "JICOFO_AUTH_USER",
										},
									},
								},
								{
									Name: "JICOFO_AUTH_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "jitsi",
											},
											Key: "JICOFO_AUTH_PASSWORD",
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

func ServiceAccountForJicofo(namespace string) *corev1.ServiceAccount {
	service := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jicofo-service-account",
			Namespace: namespace,
		},
	}
	return service
}
