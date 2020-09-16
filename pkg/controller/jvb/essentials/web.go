package essentials

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//jitsi web deployments

func DeploymentForWeb(namespace string) *appsv1.Deployment {
	replicas := int32(0)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "web",
			Namespace: namespace,
			Labels: map[string]string{
				"name": "web",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": "web",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": "web",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "web-service-account",
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "nginx-default",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "jitsi-web",
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "default",
											Path: "default",
										},
									},
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Image: "jitsi/web:4101",
							Name:  "web",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/config",
								},
								{
									Name:      "nginx-default",
									MountPath: "/config/nginx/site-confs",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
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
							},
							Env: []corev1.EnvVar{
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
									Name: "JIBRI_XMPP_USER",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "jitsi",
											},
											Key: "JIBRI_XMPP_USER",
										},
									},
								},
								{
									Name: "JIBRI_XMPP_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "jitsi",
											},
											Key: "JIBRI_XMPP_PASSWORD",
										},
									},
								},
								{
									Name: "JIBRI_RECORDER_USER",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "jitsi",
											},
											Key: "JIBRI_RECORDER_USER",
										},
									},
								},
								{
									Name: "JIBRI_RECORDER_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "jitsi",
											},
											Key: "JIBRI_RECORDER_PASSWORD",
										},
									},
								},
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.IntOrString{
											IntVal: 8080,
										},
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.IntOrString{
											IntVal: 8080,
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

func ServiceForWeb(namespace string) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "web",
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.IntOrString{IntVal: 8080},
					Protocol:   "TCP",
				},
			},
			Type: corev1.ServiceType("ClusterIP"),
			Selector: map[string]string{
				"name": "web",
			},
		},
	}
	return service
}

func ServiceAccountForWeb(namespace string) *corev1.ServiceAccount {
	service := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "web-service-account",
			Namespace: namespace,
		},
	}
	return service
}
