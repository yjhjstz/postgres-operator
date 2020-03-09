package backrest

/*
 Copyright 2019 Crunchy Data Solutions, Inc.
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

import (
	"bytes"
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"time"
	"fmt"

	crv1 "github.com/crunchydata/postgres-operator/apis/cr/v1"
	"github.com/crunchydata/postgres-operator/config"
	"github.com/crunchydata/postgres-operator/events"
	"github.com/crunchydata/postgres-operator/kubeapi"
	"github.com/crunchydata/postgres-operator/operator"
	log "github.com/sirupsen/logrus"
	v1batch "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"
)

type backrestJobTemplateFields struct {
	JobName                       string
	Name                          string
	ClusterName                   string
	Command                       string
	CommandOpts                   string
	PITRTarget                    string
	PodName                       string
	PGOImagePrefix                string
	PGOImageTag                   string
	SecurityContext               string
	PgbackrestStanza              string
	PgbackrestDBPath              string
	PgbackrestRepoPath            string
	PgbackrestRepoType            string
	BackrestLocalAndS3Storage     bool
	PgbackrestRestoreVolumes      string
	PgbackrestRestoreVolumeMounts string
}


var backrestPgHostRegex = regexp.MustCompile("--db-host|--pg1-host")
var backrestPgPathRegex = regexp.MustCompile("--db-path|--pg1-path")

// Backrest ...
func Backrest(namespace string, clientset *kubernetes.Clientset, task *crv1.Pgtask) {

	//create the Job to run the backrest command

	cmd := task.Spec.Parameters[config.LABEL_BACKREST_COMMAND]

	jobFields := backrestJobTemplateFields{
		JobName:                       task.Spec.Parameters[config.LABEL_JOB_NAME],
		ClusterName:                   task.Spec.Parameters[config.LABEL_PG_CLUSTER],
		PodName:                       task.Spec.Parameters[config.LABEL_POD_NAME],
		SecurityContext:               "",
		Command:                       cmd,
		CommandOpts:                   task.Spec.Parameters[config.LABEL_BACKREST_OPTS],
		PITRTarget:                    "",
		PGOImagePrefix:                operator.Pgo.Pgo.PGOImagePrefix,
		PGOImageTag:                   operator.Pgo.Pgo.PGOImageTag,
		PgbackrestStanza:              task.Spec.Parameters[config.LABEL_PGBACKREST_STANZA],
		PgbackrestDBPath:              task.Spec.Parameters[config.LABEL_PGBACKREST_DB_PATH],
		PgbackrestRepoPath:            task.Spec.Parameters[config.LABEL_PGBACKREST_REPO_PATH],
		PgbackrestRestoreVolumes:      "",
		PgbackrestRestoreVolumeMounts: "",
		PgbackrestRepoType:            operator.GetRepoType(task.Spec.Parameters[config.LABEL_BACKREST_STORAGE_TYPE]),
		BackrestLocalAndS3Storage:     operator.IsLocalAndS3Storage(task.Spec.Parameters[config.LABEL_BACKREST_STORAGE_TYPE]),
	}

	podCommandOpts, err := getCommandOptsFromPod(clientset, task, namespace)
	if err != nil {
		log.Error(err.Error())
		return
	}
	jobFields.CommandOpts = jobFields.CommandOpts + " " + podCommandOpts

	var doc2 bytes.Buffer
	err = config.BackrestjobTemplate.Execute(&doc2, jobFields)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if operator.CRUNCHY_DEBUG {
		config.BackrestjobTemplate.Execute(os.Stdout, jobFields)
	}

	newjob := v1batch.Job{}
	err = json.Unmarshal(doc2.Bytes(), &newjob)
	if err != nil {
		log.Error("error unmarshalling json into Job " + err.Error())
		return
	}

	newjob.ObjectMeta.Labels[config.LABEL_PGOUSER] = task.ObjectMeta.Labels[config.LABEL_PGOUSER]
	newjob.ObjectMeta.Labels[config.LABEL_PG_CLUSTER_IDENTIFIER] = task.ObjectMeta.Labels[config.LABEL_PG_CLUSTER_IDENTIFIER]

	kubeapi.CreateJob(clientset, &newjob, namespace)

	//publish backrest backup event
	if cmd == "backup" {
		topics := make([]string, 1)
		topics[0] = events.EventTopicBackup

		f := events.EventCreateBackupFormat{
			EventHeader: events.EventHeader{
				Namespace: namespace,
				Username:  task.ObjectMeta.Labels[config.LABEL_PGOUSER],
				Topic:     topics,
				Timestamp: time.Now(),
				EventType: events.EventCreateBackup,
			},
			Clustername:       jobFields.ClusterName,
			BackupType:        "pgbackrest",
		}

		err := events.Publish(f)
		if err != nil {
			log.Error(err.Error())
		}
	}

}


// getCommandOptsFromPod adds command line options from the primary pod to a backrest job.
// If not already specified in the command options provided in the pgtask, add the IP of the
// primary pod as the value for the "--db-host" parameter.  This will ensure direct
// communication between the repo pod and the primary via the primary's IP, instead of going
// through the primary pod's service (which could be unreliable). also if not already specified
// in the command options provided in the pgtask, then lookup the primary pod for the cluster
// and add the PGDATA dir of the pod as the value for the "--db-path" parameter
func getCommandOptsFromPod(clientset *kubernetes.Clientset, task *crv1.Pgtask,
	namespace string) (commandOpts string, err error) {

	// lookup the primary pod in order to determine the IP of the primary and the PGDATA directory for
	// the current primaty
	selector := config.LABEL_SERVICE_NAME + "=" + task.Spec.Parameters[config.LABEL_PG_CLUSTER] + "," + config.LABEL_DEPLOYMENT_NAME
	pods, err := kubeapi.GetPods(clientset, selector, namespace)
	if err != nil {
		return
	} else if len(pods.Items) > 1 {
		err = fmt.Errorf("More than one primary found when creating backrest job %s",
			task.Spec.Parameters[config.LABEL_JOB_NAME])
		return
	} else if len(pods.Items) == 0 {
		err = fmt.Errorf("Unable to find primary when creating backrest job %s",
			task.Spec.Parameters[config.LABEL_JOB_NAME])
		return
	}
	pod := pods.Items[0]

	var cmdOpts []string

	if !backrestPgHostRegex.MatchString(task.Spec.Parameters[config.LABEL_BACKREST_OPTS]) {
		cmdOpts = append(cmdOpts, fmt.Sprintf("--db-host=%s", pod.Spec.NodeName))
		log.Debug("Backrest primary host " + pod.Spec.NodeName)
	}
	if !backrestPgPathRegex.MatchString(task.Spec.Parameters[config.LABEL_BACKREST_OPTS]) {
		var podDbPath string
		for _, envVar := range pod.Spec.Containers[0].Env {
			if envVar.Name == "PGBACKREST_DB_PATH" {
				podDbPath = envVar.Value
				break
			}
		}
		if podDbPath != "" {
			cmdOpts = append(cmdOpts, fmt.Sprintf("--db-path=%s", podDbPath))
			log.Debug("Backrest podDbPath " + podDbPath)
		} else {
			log.Errorf("Unable to find PGBACKREST_DB_PATH on primary pod %s for backrest job %s",
				pod.Name, task.Spec.Parameters[config.LABEL_JOB_NAME])
			return
		}
	}
	// join options using a space
	commandOpts = strings.Join(cmdOpts, " ")
	return
}