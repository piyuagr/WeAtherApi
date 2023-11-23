/*
Copyright 2023.

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

package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	weatherappv1 "weatherapi.k8s.io/api/v1"
)

// WeatherapiReconciler reconciles a Weatherapi object
type WeatherapiReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=weatherapp.weather.api,resources=weatherapis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=weatherapp.weather.api,resources=weatherapis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=weatherapp.weather.api,resources=weatherapis/finalizers,verbs=update

var (
	jobOwnerKey = ".metadata.controller"
	apiGVStr    = weatherappv1.GroupVersion.String()
)

type coordtype struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type weathertype struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type maintype struct {
	Temp       float32 `json:"temp"`
	Feels_like float32 `json:"feels_like"`
	Temp_min   float32 `json:"temp_min"`
	Temp_max   float32 `json:"temp_max"`
	Pressure   float32 `json:"pressure"`
	Humidity   int     `json:"humidity"`
	Sea_level  int     `json:"sea_level"`
	Grnd_level int     `json:"grnd_level"`
}

type windtype struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust,omitempty"`
}

type cloudtype struct {
	All int `'json:"all"`
}

type systype struct {
	Country string `json:"country"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"`
}

type data struct {
	Coord      coordtype     `json:"coord"`
	Weather    []weathertype `json:"weather"`
	Base       string        `json:"base"`
	Main       maintype      `json:"main"`
	Visibility int           `json:"visibility"`
	Wind       windtype      `json:"wind"`
	Clouds     cloudtype     `json:"clouds"`
	Dt         int           `json:"dt"`
	Sys        systype       `json:"sys"`
	Timezone   int           `json:"timezone"`
	Id         int           `json:"id"`
	Name       string        `json:"name"`
	Cod        int           `json:"cod"`
}

type SlackRequestBody struct {
	Text string `json:"text"`
}

func SendSlackNotification(webhookUrl string, msg string) error {

	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Weatherapi object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *WeatherapiReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here
	var weatherApi weatherappv1.Weatherapi
	fmt.Println("The reconciler is working")
	if err := r.Get(ctx, req.NamespacedName, &weatherApi); err != nil {
		log.Error(err, "unable to fetch Weatherapi")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	response, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + weatherApi.Spec.Location + "&APPID=f6794170a7188ba717037634b2b12706")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var APIdata data
	err1 := json.Unmarshal(responseData, &APIdata) // Could have searched existing go repos for struct for this api. But now done
	if err1 != nil {
		fmt.Println(err1)
	}

	if APIdata.Name == "" {
		log.Error(err, "Incorrect location entered") // If incorrect location entered bring errors
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Convert the Api temprature from kelvin to degree celcius
	APIdata.Main.Temp -= 273.15 // fmt.Println("Struct is:", APIdata)
	weather_report := "Temperature of the " + APIdata.Name + " is: " + fmt.Sprintf("%v", APIdata.Main.Temp)
	fmt.Println(weather_report)

	conditional_min_temp := float32(weatherApi.Spec.Tempmin.MilliValue())
	conditional_min_temp /= 1000

	conditional_max_temp := float32(weatherApi.Spec.Tempmax.MilliValue())
	conditional_max_temp /= 1000

	fmt.Println(conditional_min_temp)
	fmt.Println(conditional_max_temp)

	if (conditional_min_temp <= APIdata.Main.Temp) && (APIdata.Main.Temp <= conditional_max_temp) {
		// Send notification Invoke something that the required condition is achieved.
		fmt.Println("Temperature Condition achieved, Sending Slack Notification ...")
		webhookUrl := "https://hooks.slack.com/services/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"
		err := SendSlackNotification(webhookUrl, "Your Weather Report is given below: "+"\n"+weather_report)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Sent Slack Notification.")
		// stop check after sucessful slack notification
	} else {
		fmt.Println("Condition not satisfied")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WeatherapiReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&weatherappv1.Weatherapi{}).
		Complete(r)
}
