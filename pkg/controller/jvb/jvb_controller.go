package jvb

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	jitsiv1alpha1 "jitsi-operator/pkg/apis/jitsi/v1alpha1"
	"jitsi-operator/pkg/controller/jvb/essentials"
	appsv1 "k8s.io/api/apps/v1"
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

var operatorInstanceNameLabel string = "name"
var operatorNameLabel string = "group"
var operatorNameLabelValue string = "jvb"

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
	err = c.Watch(&source.Kind{Type: &jitsiv1alpha1.Jitsi{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner JVB
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &jitsiv1alpha1.Jitsi{},
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
	instance := &jitsiv1alpha1.Jitsi{}
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

	r.proofInits(request, reqLogger)

	size := int(instance.Spec.Size)
	reqLogger.Info("JVB Size: " + fmt.Sprint(size))

	deployments := &appsv1.DeploymentList{}
	opts := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace),
		client.MatchingLabels{operatorNameLabel: operatorNameLabelValue},
	}
	ctx := context.TODO()
	err = r.client.List(ctx, deployments, opts...)

	if len(deployments.Items) < size {
		r.createDeploymentsAndServicesForJVB(size-len(deployments.Items), request, instance, reqLogger)
	} else if len(deployments.Items) > size {
		r.deleteDeploymentsAndServicesForJVB(len(deployments.Items)-size, request, reqLogger)
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileJVB) getDeployments(request reconcile.Request) *appsv1.DeploymentList {
	deployments := &appsv1.DeploymentList{}
	opts := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace),
	}
	ctx := context.TODO()
	_ = r.client.List(ctx, deployments, opts...)

	return deployments
}

func (r *ReconcileJVB) getPVCs(request reconcile.Request) *corev1.PersistentVolumeClaimList {
	pvcs := &corev1.PersistentVolumeClaimList{}
	opts := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace),
	}
	ctx := context.TODO()
	_ = r.client.List(ctx, pvcs, opts...)

	return pvcs
}

func (r *ReconcileJVB) getServices(request reconcile.Request) *corev1.ServiceList {
	services := &corev1.ServiceList{}
	opts := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace),
	}
	ctx := context.TODO()
	_ = r.client.List(ctx, services, opts...)

	return services
}

func (r *ReconcileJVB) getServiceAccounts(request reconcile.Request) *corev1.ServiceAccountList {
	serviceAccounts := &corev1.ServiceAccountList{}
	opts := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace),
	}
	ctx := context.TODO()
	_ = r.client.List(ctx, serviceAccounts, opts...)

	return serviceAccounts
}

func (r *ReconcileJVB) getConfigMaps(request reconcile.Request) *corev1.ConfigMapList {
	configMaps := &corev1.ConfigMapList{}
	opts := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace),
	}
	ctx := context.TODO()
	_ = r.client.List(ctx, configMaps, opts...)

	return configMaps
}

func (r *ReconcileJVB) getSecrets(request reconcile.Request) *corev1.SecretList {
	secrets := &corev1.SecretList{}
	opts := []client.ListOption{
		client.InNamespace(request.NamespacedName.Namespace),
	}
	ctx := context.TODO()
	_ = r.client.List(ctx, secrets, opts...)

	return secrets
}

func (r *ReconcileJVB) proofInits(request reconcile.Request, reqLogger logr.Logger) {
	prosodyDeployment := true
	jicofoDeployment := true
	prosodyService := true
	prosodyPvc := true
	prosodyServiceAccount := true
	jitsiConfigMap := true
	jitsiWebConfigMap := true
	jitsiSecret := true
	jicofoServiceAccount := true
	webServiceAccount := true
	webService := true
	webDeployment := true
	jvbServiceAccount := true

	//secrets
	secrets := r.getSecrets(request)
	for _, secret := range secrets.Items {
		if secret.ObjectMeta.Name == "jitsi" {
			jitsiSecret = false
		}
	}
	if jitsiSecret {
		ctx := context.TODO()
		_ = r.client.Create(ctx, essentials.JitsiSecret(request.NamespacedName.Namespace))
		reqLogger.Info("Jitsi Secret created")
	}

	//configMap
	configMaps := r.getConfigMaps(request)
	for _, configMap := range configMaps.Items {
		if configMap.ObjectMeta.Name == "jitsi" {
			jitsiConfigMap = false
		}
		if configMap.ObjectMeta.Name == "jitsi-web" {
			jitsiWebConfigMap = false
		}
	}
	if jitsiConfigMap {
		ctx := context.TODO()
		_ = r.client.Create(ctx, essentials.ConfigMap(request.NamespacedName.Namespace))
		reqLogger.Info("Jitsi ConfigMap created")
	}
	if jitsiWebConfigMap {
		ctx := context.TODO()
		_ = r.client.Create(ctx, essentials.ConfigMapWeb(request.NamespacedName.Namespace))
		reqLogger.Info("JitsiWeb ConfigMap created")
	}

	//pvc
	pvcs := r.getPVCs(request)
	for _, pvc := range pvcs.Items {
		if pvc.ObjectMeta.Name == "prosody" {
			prosodyPvc = false
		}
	}
	if prosodyPvc {
		ctx := context.TODO()
		_ = r.client.Create(ctx, essentials.PvcForProsody(request.NamespacedName.Namespace))
		reqLogger.Info("Prosody PVC created")
	}

	// serviceAccount
	serviceAccounts := r.getServiceAccounts(request)
	for _, serviceAccount := range serviceAccounts.Items {
		if serviceAccount.ObjectMeta.Name == "prosody-service-account" {
			prosodyServiceAccount = false
		}
		if serviceAccount.ObjectMeta.Name == "jicofo-service-account" {
			jicofoServiceAccount = false
		}
		if serviceAccount.ObjectMeta.Name == "web-service-account" {
			webServiceAccount = false
		}
		if serviceAccount.ObjectMeta.Name == "jvb-service-account" {
			jvbServiceAccount = false
		}
	}
	if prosodyServiceAccount {
		ctx := context.TODO()
		_ = r.client.Create(ctx, essentials.ServiceAccountForProsody(request.NamespacedName.Namespace))
		reqLogger.Info("Prosody ServiceAccount created")
	}
	if jicofoServiceAccount {
		ctx := context.TODO()
		_ = r.client.Create(ctx, essentials.ServiceAccountForJicofo(request.NamespacedName.Namespace))
		reqLogger.Info("Jicofo ServiceAccount created")
	}
	if webServiceAccount {
		ctx := context.TODO()
		_ = r.client.Create(ctx, essentials.ServiceAccountForWeb(request.NamespacedName.Namespace))
		reqLogger.Info("Web ServiceAccount created")
	}
	if jvbServiceAccount {
		ctx := context.TODO()
		_ = r.client.Create(ctx, essentials.ServiceAccountForJvb(request.NamespacedName.Namespace))
		reqLogger.Info("JVB ServiceAccount created")
	}

	// deployments
	deployments := r.getDeployments(request)
	for _, deployment := range deployments.Items {
		if deployment.ObjectMeta.Name == "prosody" {
			prosodyDeployment = false
		}
		if deployment.ObjectMeta.Name == "jicofo" {
			jicofoDeployment = false
		}
		if deployment.ObjectMeta.Name == "web" {
			webDeployment = false
		}
	}
	if prosodyDeployment {
		ctx := context.TODO()
		_ = r.client.Create(ctx, essentials.DeploymentForProsody(request.NamespacedName.Namespace))
		reqLogger.Info("Prosody Deployment created")
	}
	if jicofoDeployment {
		ctx := context.TODO()
		_ = r.client.Create(ctx, essentials.DeploymentForJicofo(request.NamespacedName.Namespace))
		reqLogger.Info("Jicofo Deployment created")
	}
	if webDeployment {
		ctx := context.TODO()
		_ = r.client.Create(ctx, essentials.DeploymentForWeb(request.NamespacedName.Namespace))
		reqLogger.Info("Web Deployment created")
	}

	//services
	services := r.getServices(request)
	for _, service := range services.Items {
		if service.ObjectMeta.Name == "prosody" {
			prosodyService = false
		}
		if service.ObjectMeta.Name == "web" {
			webService = false
		}
	}
	if prosodyService {
		ctx := context.TODO()
		_ = r.client.Create(ctx, essentials.ServiceForProsody(request.NamespacedName.Namespace))
		reqLogger.Info("Prosody Service created")
	}
	if webService {
		ctx := context.TODO()
		_ = r.client.Create(ctx, essentials.ServiceForWeb(request.NamespacedName.Namespace))
		reqLogger.Info("Web Service created")
	}
}

func (r *ReconcileJVB) createDeploymentsAndServicesForJVB(size int, request reconcile.Request, instance *jitsiv1alpha1.Jitsi, reqLogger logr.Logger) {
	for i := 0; i < size; i++ {
		random := essentials.RandomString(10)

		name := "jvb"

		var jvbServices []corev1.Service
		jvbServices = append(
			jvbServices,
			*essentials.ServiceForJvb(request.NamespacedName.Namespace),
			*essentials.ServiceForJvbHttp(request.NamespacedName.Namespace),
			*essentials.ServiceForJvbTcp(request.NamespacedName.Namespace),
		)

		r.createServicesForJVB(request.NamespacedName.Namespace, name, random, jvbServices, reqLogger)
		r.createDeploymentForJVB(request.NamespacedName.Namespace, random, *essentials.DeploymentForJvb(request.NamespacedName.Namespace), reqLogger)
	}
}

func (r *ReconcileJVB) createDeploymentForJVB(namespace string, random string, originalDeployment appsv1.Deployment, reqLogger logr.Logger) {
	deployment := originalDeployment.DeepCopy()
	instanceName := deployment.Name + "-" + random

	deployment.Namespace = namespace
	deployment.Name = instanceName

	deployment.Labels[operatorInstanceNameLabel] = instanceName
	deployment.Labels[operatorNameLabel] = operatorNameLabelValue

	deployment.Spec.Selector.MatchLabels[operatorInstanceNameLabel] = instanceName
	deployment.Spec.Selector.MatchLabels[operatorNameLabel] = operatorNameLabelValue

	deployment.Spec.Template.Labels[operatorInstanceNameLabel] = instanceName
	deployment.Spec.Template.Labels[operatorNameLabel] = operatorNameLabelValue

	for i, env := range deployment.Spec.Template.Spec.Containers[0].Env {
		if env.Name == "DOCKER_HOST_ADDRESS" {
			services := &corev1.ServiceList{}
			opts := []client.ListOption{
				client.InNamespace(namespace),
				client.MatchingLabels{operatorInstanceNameLabel: instanceName},
			}
			ctx := context.TODO()
			_ = r.client.List(ctx, services, opts...)

			deployment.Spec.Template.Spec.Containers[0].Env[i].Value = services.Items[0].Status.LoadBalancer.Ingress[0].IP
		}
	}

	ctx := context.TODO()
	_ = r.client.Create(ctx, deployment)
	reqLogger.Info(instanceName + " Deployment created")
}

func (r *ReconcileJVB) createServicesForJVB(namespace string, name string, random string, services []corev1.Service, reqLogger logr.Logger) {
	for _, originalService := range services {
		service := originalService.DeepCopy()
		instanceName := service.Name + "-" + random

		service.Namespace = namespace
		service.Name = service.Name + "-" + random

		service.Labels[operatorInstanceNameLabel] = instanceName
		service.Labels[operatorNameLabel] = operatorNameLabelValue

		service.Spec.Selector[operatorInstanceNameLabel] = name + "-" + random
		service.Spec.Selector[operatorNameLabel] = operatorNameLabelValue

		if _, ok := service.Annotations["metallb.universe.tf/allow-shared-ip"]; ok {
			service.Annotations["metallb.universe.tf/allow-shared-ip"] = name + "-" + random
		}

		ctx := context.TODO()
		_ = r.client.Create(ctx, service)
		reqLogger.Info(instanceName + " Service created")
	}
}

func (r *ReconcileJVB) deleteDeploymentsAndServicesForJVB(size int, request reconcile.Request, reqLogger logr.Logger) {
	for i := 0; i < size; i++ {
		services := &corev1.ServiceList{}
		deployments := &appsv1.DeploymentList{}
		opts := []client.ListOption{
			client.InNamespace(request.NamespacedName.Namespace),
			client.MatchingLabels{operatorNameLabel: operatorNameLabelValue},
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
				_ = r.client.Delete(ctx, &currentService, client.GracePeriodSeconds(0))
				reqLogger.Info(currentService.Name + " Service deleted")
			}
		}

		_ = r.client.Delete(ctx, &dep, client.GracePeriodSeconds(0))
		reqLogger.Info(dep.Name + " Deployment deleted")
	}
}
