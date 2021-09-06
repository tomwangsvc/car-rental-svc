package env

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_time "github.com/tomwangsvc/lib-svc/time"
)

type Env struct {
	CloudRun                         bool   `json:"cloud_run"`
	CloudSchedulerPushServiceAccount string `json:"cloud_scheduler_push_service_account"`
	CloudTasksPushServiceAccount     string `json:"cloud_tasks_push_service_account"`
	Debug                            bool   `json:"debug"`
	Id                               string `json:"id"`
	Gae                              bool   `json:"gae"`
	GcpProjectId                     string `json:"gcp_project_id"`
	GcpProjectNumber                 string `json:"gcp_project_number"`
	Localhost                        bool   `json:"localhost"`
	Location                         string `json:"location"`
	MaintenanceMode                  bool   `json:"maintenance_mode"`
	NumCpu                           int    `json:"num_cup"`
	Port                             int    `json:"port"`
	PubsubPushServiceAccount         string `json:"pubsub_push_service_account"`
	RuntimeId                        string `json:"runtime_id"`
	RuntimeLabel                     string `json:"runtime_label"`
	RuntimeService                   string `json:"runtime_service"`
	RuntimeVersion                   string `json:"runtime_version"`
	SpannerDatabaseId                string `json:"spanner_database_id"`
	SpannerInstanceId                string `json:"spanner_instance_id"`
	SvcId                            string `json:"svc_id"`
}

const (
	Dev = "dev"
	Prd = "prd"
	Stg = "stg"
	Uat = "uat"
)

func (e Env) Dev() bool {
	return e.Id == Dev
}

func (e Env) Prd() bool {
	return e.Id == Prd
}

func (e Env) Stg() bool {
	return e.Id == Stg
}

func (e Env) Uat() bool {
	return e.Id == Uat
}

// Production returns true if running in a production environment.
// Currently there is one production environment but there will be others eventually.
func (e Env) Production() bool {
	return e.Prd()
}

func (e Env) KeyRing() string {
	return fmt.Sprintf("projects/%s/locations/global/keyRings/svc/cryptoKeys/svc", e.GcpProjectId)
}

//revive:disable:cyclomatic
func New(svcId string) (*Env, error) {
	log.Println("Initializing env")
	if _, err := time.LoadLocation(lib_time.DefaultTimezoneName); err != nil {
		return nil, lib_errors.Wrap(err, "Failed testing ZONEINFO configuration, failed loading test location")
	}

	envId := os.Getenv("ENV")
	if envId == "" {
		return nil, lib_errors.New("Missing environment variable 'ENV'")
	}
	switch envId {
	default:
		return nil, lib_errors.New("Unrecognized environment variable 'ENV'")
	case Dev, Prd, Stg, Uat:
	}

	port := 8080
	if portEnvValue := os.Getenv("PORT"); portEnvValue == "" {
		log.Printf("Warning: Missing environment variable 'PORT'. Defaulting port to '%d'", port)
	} else {
		if v, err := strconv.Atoi(portEnvValue); err != nil {
			log.Printf("Warning: Unrecognized environment variable 'PORT' %q. Defaulting port to '%d'", portEnvValue, port)
		} else {
			port = v
		}
	}

	maintenanceMode, err := strconv.ParseBool(os.Getenv("MAINTENANCE_MODE"))
	if err != nil {
		maintenanceMode = false
		log.Printf("Warning: Missing environment variable 'MAINTENANCE_MODE'. Defaulting maintenanceMode to %t", maintenanceMode)
	}

	spannerInstanceId := os.Getenv("SPANNER_INSTANCE_ID")
	if spannerInstanceId == "" {
		spannerInstanceId = "grp-svc-1"
		log.Printf("Warning: Missing environment variable 'SPANNER_INSTANCE_ID'. Defaulting spannerInstanceId to %q", spannerInstanceId)
	}

	spannerDatabaseId := os.Getenv("SPANNER_DATABASE_ID")
	if spannerDatabaseId == "" {
		spannerDatabaseId = svcId
		log.Printf("Warning: Missing environment variable 'SPANNER_DATABASE_ID'. Defaulting spannerDatabaseId to %q", spannerDatabaseId)
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		debug = true
		log.Printf("Warning: Missing environment variable 'DEBUG'. Defaulting Debug to %t", debug)
	}

	gaeDeploymentID := os.Getenv("GAE_DEPLOYMENT_ID")
	gaeService := os.Getenv("GAE_SERVICE")
	gaeVersion := os.Getenv("GAE_VERSION")
	gae := gaeDeploymentID != "" && gaeService != "" && gaeVersion != ""

	localhostUser := os.Getenv("USER")
	localhost := localhostUser != ""

	cloudRunRevision := os.Getenv("K_REVISION")
	cloudRunService := os.Getenv("K_SERVICE")
	cloudRun := cloudRunRevision != "" && cloudRunService != ""

	// bool2int := func(v bool) int {
	// 	if v {
	// 		return 1
	// 	}
	// 	return 0
	// }

	// if bool2int(cloudRun)+bool2int(gae)+bool2int(localhost) != 1 {
	// 	return nil, lib_errors.Errorf("Cannot uniquely determine runtime: cloudRun=%t, gae=%t, localhost=%t", cloudRun, gae, localhost)
	// }

	var gcpProjectId, expectedGcpProjectIdEnvVar string
	if localhost {
		expectedGcpProjectIdEnvVar = "GCP_PROJECT_ID"

	} else {
		expectedGcpProjectIdEnvVar = "GOOGLE_CLOUD_PROJECT"
	}
	gcpProjectId = os.Getenv(expectedGcpProjectIdEnvVar)
	if gcpProjectId == "" {
		return nil, lib_errors.Errorf("Missing environment variable %q", expectedGcpProjectIdEnvVar)
	}

	gcpProjectNumber := os.Getenv("GCP_PROJECT_NUMBER")
	if gcpProjectNumber == "" {
		return nil, lib_errors.New("Missing environment variable GCP_PROJECT_NUMBER")
	}

	location := "us-central1"
	if l := os.Getenv("LOCATION"); l != "" {
		location = l
	}

	pubsubPushServiceAccount := fmt.Sprintf("id-pubsub-push-svc@%s.iam.gserviceaccount.com", gcpProjectId)
	if p := os.Getenv("PUBSUB_PUSH_SERVICE_ACCOUNT"); p != "" {
		pubsubPushServiceAccount = p
	}

	cloudSchedulerPushServiceAccount := fmt.Sprintf("id-cloud-scheduler-push-svc@%s.iam.gserviceaccount.com", gcpProjectId)
	if p := os.Getenv("CLOUD_SCHEDULER_PUSH_SERVICE_ACCOUNT"); p != "" {
		cloudSchedulerPushServiceAccount = p
	}

	cloudTasksPushServiceAccount := fmt.Sprintf("id-cloud-tasks-push-svc@%s.iam.gserviceaccount.com", gcpProjectId)
	if p := os.Getenv("CLOUD_TASKS_PUSH_SERVICE_ACCOUNT"); p != "" {
		cloudTasksPushServiceAccount = p
	}

	var runtimeLabel, runtimeID, runtimeVersion, runtimeService string
	if cloudRun {
		if cloudRunService != svcId {
			log.Printf("Warning: Configured service ID %q does not match Cloud Run environment variable K_SERVICE %q", svcId, cloudRunService)
		}
		runtimeLabel = "CLOUD_RUN"
		runtimeID = cloudRunRevision
		runtimeService = cloudRunService
		runtimeVersion = generateRuntimeVersion(cloudRunRevision)

	} else if gae {
		if gaeService != svcId {
			log.Printf("Warning: Configured service ID %q does not match GAE environment variable GAE_SERVICE %q", svcId, gaeService)
		}
		runtimeLabel = "GAE"
		runtimeID = gaeDeploymentID
		runtimeService = gaeService
		runtimeVersion = gaeVersion

	} else if localhost {
		runtimeLabel = localhostUser
		runtimeID = uuid.New().String()
		runtimeService = fmt.Sprintf("%s-%s", svcId, localhostUser)
		runtimeVersion = generateRuntimeVersion(localhostUser)

	} else {
		return nil, lib_errors.Errorf("Cannot determine runtime: cloudRun=%t, gae=%t, localhost=%t", cloudRun, gae, localhost)
	}

	// There are restrictions on the format of runtime labels and IDs when used as metadata for GCP clients so we just define our metadata to match GCP's requiremnets
	runtimeLabel = strings.ToLower(runtimeLabel)
	runtimeID = strings.Replace(strings.ToLower(runtimeID), ".", "-", -1)

	env := Env{
		CloudRun:                         cloudRun,
		CloudSchedulerPushServiceAccount: cloudSchedulerPushServiceAccount,
		CloudTasksPushServiceAccount:     cloudTasksPushServiceAccount,
		Debug:                            debug,
		Gae:                              gae,
		GcpProjectId:                     gcpProjectId,
		GcpProjectNumber:                 gcpProjectNumber,
		Id:                               envId,
		Localhost:                        localhost,
		Location:                         location,
		MaintenanceMode:                  maintenanceMode,
		NumCpu:                           runtime.NumCPU(),
		Port:                             port,
		PubsubPushServiceAccount:         pubsubPushServiceAccount,
		RuntimeId:                        runtimeID,
		RuntimeLabel:                     runtimeLabel,
		RuntimeService:                   runtimeService,
		RuntimeVersion:                   runtimeVersion,
		SpannerDatabaseId:                spannerDatabaseId,
		SpannerInstanceId:                spannerInstanceId,
		SvcId:                            svcId,
	}

	log.Println("Initialized env: ", env)
	return &env, nil
	//revive:enable:cyclomatic
}

func generateRuntimeVersion(runtimeID string) string {
	runtimeVerion := strings.Replace(strings.Replace(time.Now().UTC().Format(time.RFC3339), "-", "", -1), ":", "", -1)
	return fmt.Sprintf("%s@%s", strings.ToLower(runtimeVerion[:strings.Index(runtimeVerion, "Z")-1]), runtimeID)
}
