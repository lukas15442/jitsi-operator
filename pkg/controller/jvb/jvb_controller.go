package jvb

import (
	"context"
	"fmt"
	fbiv1alpha1 "jitsi-operator/pkg/apis/fbi/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var JvbName = []string{"app.kubernetes.io/name", "jvb"}
var JvbPreInstanceName = []string{"app.kubernetes.io/instance", "jvb-test"}

var log = logf.Log.WithName("controller_jvb")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new JVB Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileJVB{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("jvb-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource JVB
	err = c.Watch(&source.Kind{Type: &fbiv1alpha1.JVB{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner JVB
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &fbiv1alpha1.JVB{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileJVB implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileJVB{}

// ReconcileJVB reconciles a JVB object
type ReconcileJVB struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a JVB object and makes changes based on the state read
// and what is in the JVB.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileJVB) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling JVB")

	// Fetch the JVB instance
	instance := &fbiv1alpha1.JVB{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	size := int(instance.Spec.Size)
	reqLogger.Info(fmt.Sprint(size))

	deployments := &appsv1.DeploymentList{}
	opts := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace),
		client.MatchingLabels{JvbName[0]: JvbName[1]},
	}
	ctx := context.TODO()
	err = r.client.List(ctx, deployments, opts...)

	if len(deployments.Items) < size {
		r.createDeploymentsAndServices(size-len(deployments.Items), request)
	} else if len(deployments.Items) > size {
		r.deleteDeploymentsAndServices(len(deployments.Items)-size, request)
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileJVB) createDeploymentsAndServices(size int, request reconcile.Request) {
	for i := 0; i < size; i++ {
		random := RandomString(10)
		name := JvbPreInstanceName[1] + "-" + random
		nameHttp := JvbPreInstanceName[1] + "-http-" + random
		nameTcp := JvbPreInstanceName[1] + "-tcp-" + random
		ctx := context.TODO()

		dep := r.deploymentForJvb(request.NamespacedName.Namespace, name)
		_ = r.client.Create(ctx, dep)

		service := r.serviceForJvb(request.NamespacedName.Namespace, name)
		_ = r.client.Create(ctx, service)

		serviceHttp := r.serviceForJvbHttp(request.NamespacedName.Namespace, nameHttp, name)
		_ = r.client.Create(ctx, serviceHttp)

		serviceTcp := r.serviceForJvbTcp(request.NamespacedName.Namespace, nameTcp, name)
		_ = r.client.Create(ctx, serviceTcp)
	}
}

func (r *ReconcileJVB) deleteDeploymentsAndServices(size int, request reconcile.Request) {
	for i := 0; i < size; i++ {
		services := &corev1.ServiceList{}
		deployments := &appsv1.DeploymentList{}
		opts := []client.ListOption{
			client.InNamespace(request.NamespacedName.Namespace),
			client.MatchingLabels{JvbName[0]: JvbName[1]},
		}
		ctx := context.TODO()
		_ = r.client.List(ctx, services, opts...)
		_ = r.client.List(ctx, deployments, opts...)

		dep := deployments.Items[0]
		name := dep.ObjectMeta.Name
		random := strings.Split(name, "-")[len(strings.Split(name, "-"))-1]

		for _, currentService := range services.Items {
			currentRandom := strings.Split(currentService.Name, "-")[len(strings.Split(currentService.Name, "-"))-1]
			if currentRandom == random {
				_ = r.client.Delete(ctx, &currentService, client.GracePeriodSeconds(2))
			}
		}

		_ = r.client.Delete(ctx, &dep, client.GracePeriodSeconds(2))
	}
}

func (r *ReconcileJVB) deploymentForJvb(namespace string, name string) *appsv1.Deployment {
	replicas := int32(1)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				JvbPreInstanceName[0]: name,
				JvbName[0]:            JvbName[1],
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					JvbPreInstanceName[0]: name,
					JvbName[0]:            JvbName[1],
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						JvbPreInstanceName[0]: name,
						JvbName[0]:            JvbName[1],
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "jvb",
					Containers: []corev1.Container{
						{
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: nil,
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/about/health",
										Port: intstr.IntOrString{
											IntVal: 8080,
										},
										Scheme: corev1.URISchemeHTTP,
									},
									TCPSocket: nil,
								},
								TimeoutSeconds:   1,
								PeriodSeconds:    10,
								SuccessThreshold: 1,
								FailureThreshold: 3,
							},
							Name: "jvb",
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									Exec: nil,
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/about/health",
										Port: intstr.IntOrString{
											IntVal: 8080,
										},
										Scheme: corev1.URISchemeHTTP,
									},
									TCPSocket: nil,
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
									Name: "DOCKER_HOST_ADDRESS",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "jitsi",
											},
											Key: "JVB0_PUBLIC_ADDR",
										},
									},
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

func (r *ReconcileJVB) serviceForJvb(namespace string, name string) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				JvbPreInstanceName[0]: name,
				JvbName[0]:            JvbName[1],
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				JvbPreInstanceName[0]: name,
				JvbName[0]:            JvbName[1],
			},
			Type: corev1.ServiceType("LoadBalancer"),
			Ports: []corev1.ServicePort{
				{
					Name:       "jvb",
					Port:       10000,
					TargetPort: intstr.IntOrString{IntVal: 10000},
					Protocol:   "UDP",
				},
			},
		},
	}
	return service
}

func (r *ReconcileJVB) serviceForJvbHttp(namespace string, httpName string, name string) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      httpName,
			Namespace: namespace,
			Labels: map[string]string{
				JvbPreInstanceName[0]: httpName,
				JvbName[0]:            JvbName[1],
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				JvbPreInstanceName[0]: name,
				JvbName[0]:            JvbName[1],
			},
			Type: corev1.ServiceType("ClusterIP"),
			Ports: []corev1.ServicePort{
				{
					Name:       "jvb-http",
					Port:       80,
					TargetPort: intstr.IntOrString{IntVal: 8080},
					Protocol:   "TCP",
				},
			},
		},
	}
	return service
}

func (r *ReconcileJVB) serviceForJvbTcp(namespace string, tcpName string, name string) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tcpName,
			Namespace: namespace,
			Labels: map[string]string{
				JvbPreInstanceName[0]: tcpName,
				JvbName[0]:            JvbName[1],
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				JvbPreInstanceName[0]: name,
				JvbName[0]:            JvbName[1],
			},
			Type:                  corev1.ServiceType("LoadBalancer"),
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyType("Cluster"),
			Ports: []corev1.ServicePort{
				{
					Name:       "jvb-tcp",
					Port:       443,
					TargetPort: intstr.IntOrString{IntVal: 443},
					Protocol:   "TCP",
				},
			},
		},
	}
	return service
}
