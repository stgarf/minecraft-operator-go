package minecraft

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	interviewv1alpha1 "github.com/stgarf/minecraft-operator-go/pkg/apis/interview/v1alpha1"
)

var log = logf.Log.WithName("controller_minecraft")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Minecraft Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMinecraft{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("minecraft-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Minecraft
	err = c.Watch(&source.Kind{Type: &interviewv1alpha1.Minecraft{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Minecraft
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &interviewv1alpha1.Minecraft{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMinecraft{}

// ReconcileMinecraft reconciles a Minecraft object
type ReconcileMinecraft struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Minecraft object and makes changes based on the state read
// and what is in the Minecraft.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMinecraft) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Minecraft")

	// Fetch the Minecraft instance
	instance := &interviewv1alpha1.Minecraft{}
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

	// Define a new Pod object
	// Define a PVC as well
	pvc := newPersistantVolumeClaimForCR(instance)
	pod := newPodForCR(instance)

	// Set Minecraft instance as the owner and controller for PVC
	if err := controllerutil.SetControllerReference(instance, pvc, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Set Minecraft instance as the owner and controller for pod
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this PVC already exists
	foundPVC := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pvc.Name, Namespace: pvc.Namespace}, foundPVC)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new PersistantVolumeClaim", "Pod.Namespace", pvc.Namespace, "Pod.Name", pvc.Name)
		err = r.client.Create(context.TODO(), pvc)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	} else {
		// Pod already exists - don't requeue
		reqLogger.Info("Skip reconcile: PersistantVolumeClaim already exists", "Pod.Namespace", foundPVC.Namespace, "Pod.Name", foundPVC.Name)
	}

	// Check if this Pod already exists
	foundPod := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, foundPod)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	} else {
		// Pod already exists - don't requeue
		reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", foundPod.Namespace, "Pod.Name", foundPod.Name)
	}

	/* 	// Check if this VolumeMount already exists
	   	foundVolumeMount := &corev1.Pod{}
	   	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, foundVolumeMount)
	   	if err != nil && errors.IsNotFound(err) {
	   		reqLogger.Info("Creating a new VolumeMount", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	   		err = r.client.Create(context.TODO(), pod)
	   		if err != nil {
	   			return reconcile.Result{}, err
	   		}
	   	} else if err != nil {
	   		return reconcile.Result{}, err
	   	} else {
	   		// Pod already exists - don't requeue
	   		reqLogger.Info("Skip reconcile: VolumeMount already exists", "Pod.Namespace", foundVolumeMount.Namespace, "Pod.Name", foundVolumeMount.Name)
	   	}

	   	// Check if this Service already exists
	   	foundService := &corev1.Service{}
	   	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, foundService)
	   	if err != nil && errors.IsNotFound(err) {
	   		reqLogger.Info("Creating a new Service", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	   		err = r.client.Create(context.TODO(), pod)
	   		if err != nil {
	   			return reconcile.Result{}, err
	   		}
	   	} else if err != nil {
	   		return reconcile.Result{}, err
	   	} else {
	   		// Pod already exists - don't requeue
	   		reqLogger.Info("Skip reconcile: Service already exists", "Pod.Namespace", foundService.Namespace, "Pod.Name", foundService.Name)
	   	} */

	return reconcile.Result{}, nil
}

// newPodForCR returns a minecraft pod with the same name/namespace as the cr
// https://godoc.org/k8s.io/api/core/v1#Pod
func newPodForCR(cr *interviewv1alpha1.Minecraft) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}

	pvc := &corev1.PersistentVolumeClaimVolumeSource{
		ClaimName: cr.Name + "-pvc",
	}

	fsGroup := int64(1000)
	runAsNonRoot := true
	runAsUser := int64(1000)

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  cr.Name,
					Image: "us.gcr.io/kubeoperatorstest/minecraft:v" + cr.Spec.Version,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      cr.Name + "-storage",
							MountPath: "/server-data",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: cr.Name + "-storage",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: pvc,
					},
				},
			},
			SecurityContext: &corev1.PodSecurityContext{
				FSGroup:      &fsGroup,
				RunAsNonRoot: &runAsNonRoot,
				RunAsUser:    &runAsUser,
			},
		},
	}
}

// newPersistantVolumeClaimForCR returns a PVC for the name/namespace of the cr
// https://godoc.org/k8s.io/api/core/v1#PersistentVolumeClaim
func newPersistantVolumeClaimForCR(cr *interviewv1alpha1.Minecraft) *corev1.PersistentVolumeClaim {
	resList := map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceStorage: resource.MustParse("50Mi"),
	}
	storageClass := "standard"

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pvc",
			Namespace: cr.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: resList,
			},
			StorageClassName: &storageClass,
		},
	}
}

/* // newVolumeMountForCR returns a VolumeMount
// https://godoc.org/k8s.io/api/core/v1#VolumeMount
func newVolumeMountForCR(cr *interviewv1alpha1.Minecraft) *corev1.VolumeMount {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  cr.Name,
					Image: "us.gcr.io/kubeoperatorstest/minecraft:v1.13.2",
				},
			},
		},
	}
}

// newServiceForCR returns a service with the same name/namespace as the cr to load balance traffic
// https://godoc.org/k8s.io/api/core/v1#Service
func newServiceForCR(cr *interviewv1alpha1.Minecraft) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  cr.Name,
					Image: "us.gcr.io/kubeoperatorstest/minecraft:v1.13.2",
				},
			},
		},
	}
}
*/
