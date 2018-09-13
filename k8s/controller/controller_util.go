package controller

type ControllerExpectations struct {
	cache.Store
}

type Expectations interface {
	Fulfilled() bool
}

type ControlleeExpectations struct {
	add       int64
	del       int64
	key       string
	timestamp time.Time
}

func (e *ControlleeExpectations) Fulfilled() bool {
	return atomic.LoadInt64(&e.add) <= 0 && atomic.LoadInt64(&e.del) <= 0
}

type PodControlInterface interface {
	CreatePods(namespace string, template *v1.PodTemplateSpec, object runtime.Object) error
	CreatePodsOnNode(nodeName, namespace string, template *v1.PodTemplateSpec, object runtime.Object, controllerRef *metav1.OwnerReference) error
	CreatePodsWithControllerRef(namespace string, template *v1.PodTemplateSpec, object runtime.Object, controllerRef *metav1.OwnerReference) error
	DeletePod(namespace string, podID string, object runtime.Object) error
	PatchPod(namespace, name string, data []byte) error
}

type RealPodControl struct {
	KubeClient clientset.Interface
	Recorder   record.EventRecorder
}

func (r RealPodControl) createPods(nodeName, namespace string, template *v1.PodTemplateSpec, object runtime.Object, controllerRef *metav1.OwnerReference) error {
	pod, _ := GetPodFromTemplate(template, object, controllerRef)
	if len(nodeName) != 0 {
		pod.Spec.NodeName = nodeName
	}
	if labels.Set(pod.Labels).AsSelectorPreValidated().Empty() {
		return fmt.Errorf("unable to create pods, no labels")
	}
	if newPod, err := r.KubeClient.CoreV1().Pods(namespace).Create(pod); err != nil {
		r.Recorder.Eventf(object, v1.EventTypeWarning, FailedCreatePodReason, "Error creating: %v", err)
		return err
	} else {
		accessor, err := meta.Accessor(object)
		if err != nil {
			return nil
		}
		r.Recorder.Eventf(object, v1.EventTypeNormal, SuccessfulCreatePodReason, "Created pod: %v", newPod.Name)
	}
	return nil
}
