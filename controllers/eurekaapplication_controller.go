/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	discoveryv1 "github.com/eurek8s/controller/api/v1"
	eurekahandler "github.com/eurek8s/controller/internal/eureka/handler"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	eventType           = v1.EventTypeWarning
	eventReasonNotFound = "NotFound"
)

// EurekaApplicationReconciler reconciles a EurekaApplication object
type EurekaApplicationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	EventRecorder record.EventRecorder
	EurekaHandler *eurekahandler.Handler
}

//+kubebuilder:rbac:groups=discovery.eurek8s.com,resources=eurekaapplications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=discovery.eurek8s.com,resources=eurekaapplications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=discovery.eurek8s.com,resources=eurekaapplications/finalizers,verbs=update

func (r *EurekaApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("eurekaapplication", req.NamespacedName)
	log.Info("Received reconcile event")

	var eurekaApp discoveryv1.EurekaApplication
	if err := r.Get(ctx, req.NamespacedName, &eurekaApp); err != nil {
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.EurekaHandler.Handle(ctx, r.Client, &eurekaApp, req.String()); err != nil {
		if apierrors.IsNotFound(err) {
			switch t := err.(type) {
			case apierrors.APIStatus:
				d := t.Status().Details
				message := fmt.Sprintf("Referenced object not found %s/%s", d.Kind, d.Name)
				r.EventRecorder.Event(&eurekaApp, eventType, eventReasonNotFound, message)
			default:
				r.EventRecorder.Event(&eurekaApp, eventType, eventReasonNotFound, "Referenced object not found")
			}

			log.Error(err, "unable to fetch resource")
		}

		//return ctrl.Result{RequeueAfter: 30 * time.Second}, client.IgnoreNotFound(err)
		r.Log.Info("re-queueing to run after 30s...")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	} else {
		eurekaApp.Status.LastReconcileTime = &metav1.Time{Time: time.Now()}
		log.Info("updating EurekaApplication lastReconcileTime...")
		if err := r.Status().Update(ctx, &eurekaApp); err != nil {
			log.Error(err, "unable to update eureka application status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EurekaApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&discoveryv1.EurekaApplication{}).
		Complete(r)
}
