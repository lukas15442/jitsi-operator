package jvb

import (
	"context"
	"fmt"
	fbiv1alpha1 "jitsi-operator/pkg/apis/fbi/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

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

var JvbGroupLabel = []string{"app.kubernetes.io/group-name", "jvb-test"}

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
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
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
		client.MatchingLabels{JvbGroupLabel[0]: JvbGroupLabel[1]},
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

func (r *ReconcileJVB) deploymentForJvb(namespace string, name string) *appsv1.Deployment {
	replicas := int32(1)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				JvbGroupLabel[0]:         JvbGroupLabel[1],
				"app.kubernetes.io/name": name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app.kubernetes.io/name": name,},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app.kubernetes.io/name": name,},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "twalter/openshift-nginx",
						Name:  "redirect-nginx-image",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Protocol:      "TCP",
						}, {
							ContainerPort: 8081,
							Protocol:      "TCP",
						}},
					}},
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
				JvbGroupLabel[0]:         JvbGroupLabel[1],
				"app.kubernetes.io/name": name,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app.kubernetes.io/name": name,
			},
			Ports: []corev1.ServicePort{{
				Port:       8081,
				TargetPort: intstr.IntOrString{IntVal: 8081},
				Protocol:   "TCP",
			}},
		},
	}
	return service
}

func (r *ReconcileJVB) createDeploymentsAndServices(size int, request reconcile.Request) {
	for i := 0; i < size; i++ {
		name := JvbGroupLabel[1] + "-" + RandomString(10)

		dep := r.deploymentForJvb(request.NamespacedName.Namespace, name)
		service := r.serviceForJvb(request.NamespacedName.Namespace, name)
		ctx := context.TODO()
		_ = r.client.Create(ctx, dep)
		_ = r.client.Create(ctx, service)
	}
}

func (r *ReconcileJVB) deleteDeploymentsAndServices(size int, request reconcile.Request) {
	for i := 0; i < size; i++ {
		services := &corev1.ServiceList{}
		deployments := &appsv1.DeploymentList{}
		opts := []client.ListOption{
			client.InNamespace(request.NamespacedName.Namespace),
			client.MatchingLabels{JvbGroupLabel[0]: JvbGroupLabel[1]},
		}
		ctx := context.TODO()
		_ = r.client.List(ctx, services, opts...)
		_ = r.client.List(ctx, deployments, opts...)

		dep := deployments.Items[0]
		name := dep.ObjectMeta.Name

		var service corev1.Service
		for _, currentService := range services.Items {
			if currentService.Name == name {
				service = currentService
			}
		}

		_ = r.client.Delete(ctx, &service, client.GracePeriodSeconds(2))
		_ = r.client.Delete(ctx, &dep, client.GracePeriodSeconds(2))
	}
}
