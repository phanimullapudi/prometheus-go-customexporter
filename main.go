package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	masterURL  string
	kubeconfig string
	addr       = flag.String("listen-address", ":8080", "The address to listen for an HTTP requests.")
)

func main() {

	//kubeCfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	kubconfig := flag.String("kubconfig", "/Users/phanimullapudi/Downloads/kubeconfig-home.yaml", "kubeconfig-home.yaml")
	log.Info("Build config from flags..")

	kubeCfg, err := clientcmd.BuildConfigFromFlags("", *kubconfig)
	if err != nil {
		panic(err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(kubeCfg)
	if err != nil {
		log.Fatalf("Error building watcher clientset: %s", err.Error())
	}

	// Prometheus HTTP handler
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(*addr, nil))
	}()

	// Prometheus Metrics
	cns_current_pod_count := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cns_current_pod_count",
		Help: "Get the counts of the pods in the current cluster",
	})

	for {
		// Fetch pods in the cluster
		podList, err := kubeClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Errorf("Error fetching pods")
		} else {
			cns_current_pod_count.Set(float64(len(podList.Items)))
		}

		time.Sleep(time.Second * 10)
	}

}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kuberenetes API server. OVerrides any value in the kubeconfig")
	flag.Parse()
}
