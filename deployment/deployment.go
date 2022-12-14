package deployment

import (
	"context"
	"fmt"
	"github.com/Superm4n97/Client-Go/others"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	_ "k8s.io/client-go/util/retry"
	"time"
)

//emptyDeploymentTemplate is a template that contains the basic configuration of a deployment.
//it takes the minimum parameter of a deployment and returns a deployment
func emptyDeploymentTemplate(deploymentName string, replicas int32, podLabel1 string, podLabel2 string, containerName string, image string, containerPortNumber int32) appsv1.Deployment {
	//func emptyDeploymentTemplate() appsv1.Deployment {

	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: others.Int32ToPointer(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					podLabel1: podLabel2,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						podLabel1: podLabel2,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  containerName,
							Image: image,
							Ports: []v1.ContainerPort{
								{
									ContainerPort: containerPortNumber,
								},
							},
						},
					},
				},
			},
		},
	}
}

//UpdateDeployment can only update the replicas of that deployment. It takes name,
//replicas, clientset and update that deployment
func UpdateDeployment(deploymentName string, replicas int32, clientset kubernetes.Interface) {

	fmt.Println("entered into without sleep with replicas ", replicas)

	depClientset := clientset.AppsV1().Deployments(v1.NamespaceDefault)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		dep := GetDeployment(deploymentName, clientset)
		dep.Spec.Replicas = others.Int32ToPointer(replicas)

		time.Sleep(time.Minute)

		_, updateErr := depClientset.Update(context.TODO(), &dep, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		fmt.Println("failed to update")
		panic(retryErr)
	}

	fmt.Println("successfully updated.")

}

func UpdateDeploymentWithSleep(deploymentName string, replicas int32, clientset kubernetes.Interface) {

	fmt.Println("entered with sleep with replicas ", replicas)

	dep := GetDeployment(deploymentName, clientset)
	dep.Spec.Replicas = others.Int32ToPointer(replicas)

	time.Sleep(time.Second * 10)

	_, err := clientset.AppsV1().Deployments(v1.NamespaceDefault).Update(context.TODO(), &dep, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("failed to update deployment with sleep ", deploymentName)
		panic(err)
	}

	fmt.Println("updated replicas with sleep ")
}

//GetDeploymentList gets the list of all the running deployments
func GetDeploymentList(clientset kubernetes.Interface) {
	dep_lst, err := clientset.AppsV1().Deployments(v1.NamespaceDefault).List(context.TODO(), v12.ListOptions{})

	if err != nil {
		fmt.Println("failed to create deployment list")
		panic(err)
	}

	fmt.Println("Deployments:")
	for _, item := range dep_lst.Items {
		fmt.Println(item.Name)
	}
}

// GetDeployment takes name, clientset of a particular deployment and return a deployment
func GetDeployment(name string, clientset kubernetes.Interface) appsv1.Deployment {
	dep, err := clientset.AppsV1().Deployments(v1.NamespaceDefault).Get(context.TODO(), name, v12.GetOptions{})
	if err != nil {
		fmt.Println("failed to get the deployment")
		panic(err)
	}
	//fmt.Println(dep.Name)

	return *dep
}

//CreateDeployment creates a deployment returns nothing
func CreateDeployment(depName string, clientset kubernetes.Interface) {

	dep := emptyDeploymentTemplate(depName, 1, "app", "book-server", "server-container", "superm4n/book-api-server:v0.1.3", 8080)

	res, err := clientset.AppsV1().Deployments(v1.NamespaceDefault).Create(context.TODO(), &dep, v12.CreateOptions{})
	if err != nil {
		fmt.Println("failed to create deployment")
		panic(err)
	}

	fmt.Println(res.Name, " deployment created")
}

func DeleteDeployment(depName string, clientset kubernetes.Interface) {

	err := clientset.AppsV1().Deployments(v1.NamespaceDefault).Delete(context.TODO(), depName, v12.DeleteOptions{})
	if err != nil {
		fmt.Println("failed to create deployment")
		panic(err)
	}

	fmt.Println(depName, " deployment deleted")
}
