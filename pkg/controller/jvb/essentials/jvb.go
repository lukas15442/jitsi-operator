package essentials

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func ServiceAccountForJvb(namespace string) *corev1.ServiceAccount {
	service := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: "jvb-service-account",
			Namespace: namespace,
		},
	}
	return service
}

func ServiceForJvb(namespace string) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "jvb",
			Labels: map[string]string{},
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "jvb",
					Port:       10000,
					TargetPort: intstr.IntOrString{IntVal: 10000},
					Protocol:   "UDP",
				},
			},
			Type: corev1.ServiceType("LoadBalancer"),
			Selector: map[string]string{},
		},
	}
	return service
}

func ServiceForJvbHttp(namespace string) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "jvb-http",
			Labels: map[string]string{},
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "jvb-http",
					Port:       80,
					TargetPort: intstr.IntOrString{IntVal: 8080},
					Protocol:   "TCP",
				},
			},
			Type: corev1.ServiceType("ClusterIP"),
			Selector: map[string]string{},
		},
	}
	return service
}

func ServiceForJvbTcp(namespace string) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "jvb-tcp",
			Labels: map[string]string{},
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "jvb-tcp",
					Port:       443,
					TargetPort: intstr.IntOrString{IntVal: 443},
					Protocol:   "TCP",
				},
			},
			Type: corev1.ServiceType("LoadBalancer"),
			Selector: map[string]string{},
		},
	}
	return service
}

func DeploymentForJvb(namespace string) *appsv1.Deployment {
	replicas := int32(0)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "jvb",
			Labels: map[string]string{},
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{},
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "jvb-service-account",
					Containers: []corev1.Container{
						{
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/about/health",
										Port: intstr.IntOrString{
											IntVal: 8080,
										},
										Scheme: corev1.URISchemeHTTP,
									},
								},
								TimeoutSeconds:   1,
								PeriodSeconds:    10,
								SuccessThreshold: 1,
								FailureThreshold: 3,
							},
							Name: "jvb",
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/about/health",
										Port: intstr.IntOrString{
											IntVal: 8080,
										},
										Scheme: corev1.URISchemeHTTP,
									},
								},
								TimeoutSeconds:   1,
								PeriodSeconds:    10,
								SuccessThreshold: 1,
								FailureThreshold: 3,
							},
							Env: []corev1.EnvVar{
								{
									Name: "JVB_AUTH_USER",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "jitsi",
											},
											Key: "JVB_AUTH_USER",
										},
									},
								},
								{
									Name: "JVB_AUTH_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "jitsi",
											},
											Key: "JVB_AUTH_PASSWORD",
										},
									},
								},
								{
									Name:  "DOCKER_HOST_ADDRESS",
									Value: "WILL_BE_REPLACED",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "jvb",
									ContainerPort: 10000,
									Protocol:      "UDP",
								},
								{
									Name:          "jvb-tcp",
									ContainerPort: 443,
									Protocol:      "TCP",
								},
							},
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
							Image: "jitsi/jvb:4101",
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config",
						},
					},
				},
			},
		},
	}
	return dep
}
