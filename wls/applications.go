package wls

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Application is the root structure for a response from an AdminServer.  An Application instance on an AdminServer will provide details about an application, including it's Health,
// running State, the Type of Application (war, ear, jar, etc.), as well as some more detailed pieces of information, including the Targets it was deployed to, any associated
// WorkManagers, and other pertinent deployed details.
type Application struct {
	Name                  string
	AppType               string `json:"type"`
	State                 string
	Health                string
	TargetStates          []TargetState          `json:"targetStates,omitempty"`
	DataSources           []AppDataSource        `json:"dataSources,omitempty"`
	WorkManagers          []WorkManager          `json:"workManagers,omitempty"`
	MinThreadsConstraints []MinThreadsConstraint `json:"minThreadsConstraints,omitempty"`
	MaxThreadsConstraints []MaxThreadsConstraint `json:"maxThreadsConstraints,omitempty"`
	RequestClasses        []RequestClass         `json:"requestClasses,omitempty"`
}

// TargetState is the state of a Target.  In WebLogic, this could be Running, Prepared
type TargetState struct {
	Target string
	State  string
}

// AppDataSource displays information about an Application's referenced DataSources.  See the DataSources resource for more information on
// what you can gather from a WebLogic DataSource
type AppDataSource struct {
	Name   string
	Server string
	State  string
}

// WorkManager is the struct type explaining an Application's statistics for the Work Managers that are configured for an application.
type WorkManager struct {
	Name              string
	Server            string
	PendingRequests   int
	CompletedRequests int
}

// MinThreadsConstraint provides statistics for the minimum thread constraints that are configured for an application.
type MinThreadsConstraint struct {
	Name                     string
	Server                   string
	PendingRequests          int
	CompletedRequests        int
	ExecutingRequests        int
	OutOfOrderExecutionCount int
	MustRunCount             int
	MaxWaitTime              int
	CurrentWaitTime          int
}

// MaxThreadsConstraint provides statistics for maximum thread constraints that are configured for an application.
type MaxThreadsConstraint struct {
	Name              string
	Server            string
	ExecutingRequests int
	DeferredRequests  int
}

// RequestClass provides statistics for the request classes that are configured for an application.
type RequestClass struct {
	Name                 string
	Server               string
	RequestClassType     string
	CompletedCount       int
	TotalThreadUse       int
	PendingRequestCount  int
	VirtualTimeIncrement int
}

func (a *Application) GoString() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Name: %v|AppType: %v|State: %v|Health: %v\n", a.Name, a.AppType, a.State, a.Health))
	if len(a.TargetStates) > 0 {
		buffer.WriteString("Target States\n")
		for i := range a.TargetStates {
			buffer.WriteString("\t")
			buffer.WriteString(fmt.Sprintf("Target: %v|State: %v\n", a.TargetStates[i].Target, a.TargetStates[i].State))
		}
	}
	if len(a.DataSources) > 0 {
		buffer.WriteString("Data Sources\n")
		for i := range a.DataSources {
			buffer.WriteString("\t")
			buffer.WriteString(fmt.Sprintf("Name: %v|Server: %v|State: %v\n", a.DataSources[i].Name, a.DataSources[i].Server, a.DataSources[i].State))
		}
	}
	if len(a.WorkManagers) > 0 {
		buffer.WriteString("Work Managers\n")
		for i := range a.WorkManagers {
			buffer.WriteString("\t")
			buffer.WriteString(fmt.Sprintf("Name: %v|Server: %v|Pending Requests: %v|Completed Requests: %v\n", a.WorkManagers[i].Name,
				a.WorkManagers[i].Server, a.WorkManagers[i].PendingRequests, a.WorkManagers[i].CompletedRequests))
		}
	}
	if len(a.MinThreadsConstraints) > 0 {
		buffer.WriteString("Minimum Threads Constraints\n")
		for i := range a.MinThreadsConstraints {
			buffer.WriteString("\t")
			buffer.WriteString(fmt.Sprintf("Name: %v|Server: %v|Pending Req: %v|Completed Req: %v\n", a.MinThreadsConstraints[i].Name,
				a.MinThreadsConstraints[i].Server, a.MinThreadsConstraints[i].PendingRequests, a.MinThreadsConstraints[i].CompletedRequests))
			buffer.WriteString("\t")
			buffer.WriteString(fmt.Sprintf("Executing Reqs: %v|OutOfOrder Exec. Cnt: %v|Must Run Count: %v\n", a.MinThreadsConstraints[i].ExecutingRequests,
				a.MinThreadsConstraints[i].OutOfOrderExecutionCount, a.MinThreadsConstraints[i].MustRunCount))
			buffer.WriteString("\t")
			buffer.WriteString(fmt.Sprintf("Max Wait Time: %v|Current Wait Time: %v\n", a.MinThreadsConstraints[i].MaxWaitTime,
				a.MinThreadsConstraints[i].CurrentWaitTime))
		}
	}
	if len(a.MaxThreadsConstraints) > 0 {
		buffer.WriteString("Max Thread Constraints\n")
		for i := range a.MaxThreadsConstraints {
			buffer.WriteString("\t")
			buffer.WriteString(fmt.Sprintf("Name: %v|Server: %v|Executing Reqs: %v|Deferred Reqs: %v\n",
				a.MaxThreadsConstraints[i].Name, a.MaxThreadsConstraints[i].Server, a.MaxThreadsConstraints[i].ExecutingRequests,
				a.MaxThreadsConstraints[i].DeferredRequests))
		}
	}
	if len(a.RequestClasses) > 0 {
		buffer.WriteString("Request Classes\n")
		for i := range a.RequestClasses {
			buffer.WriteString("\t")
			buffer.WriteString(fmt.Sprintf("Name: %v|Server %v|Request Class Type: %v|Completed Count: %v\n",
				a.RequestClasses[i].Name, a.RequestClasses[i].Server, a.RequestClasses[i].RequestClassType,
				a.RequestClasses[i].CompletedCount))
			buffer.WriteString("\t")
			buffer.WriteString(fmt.Sprintf("Total Thread Use: %v|Pending Request Cnt: %v|Virtual Time Increment: %v\n",
				a.RequestClasses[i].TotalThreadUse, a.RequestClasses[i].PendingRequestCount, a.RequestClasses[i].VirtualTimeIncrement))
		}
	}

	return buffer.String()
}

// Applications returns all applications deployed in the domain and their run-time information, including the application type and their state and health.
// - isfullFormat specifies whether to request the FULL format for an Application.  Much more data is brought back for
//   each of the subytpes within an Application.  By default, this is false.
// This function returns a listing of []Application's on the specified AdminServer, or an error denoting any issues
// making the callout.
func (s *AdminServer) Applications(isFullFormat bool) ([]Application, error) {
	url := fmt.Sprintf("%v%v/applications", s.AdminURL, MonitorPath)
	if isFullFormat {
		url = url + "?format=full"
	}
	w, err := requestAndUnmarshal(url, s)
	if err != nil {
		return nil, err
	}
	var applications []Application
	if err := json.Unmarshal(w.Body.Items, &applications); err != nil {
		return nil, err
	}
	return applications, nil
}

// Application returns the run-time information of a specified application, including statistics for entity beans, application-scoped work managers, and data sources.
// on how to get all of the []Application's on the server.
// This will always return a full format, including all of the details in the underlying struct types.
// It may also return an error if there were any issues calling out to the AdminServer
func (s *AdminServer) Application(app string) (*Application, error) {
	url := fmt.Sprintf("%v%v/applications/%v", s.AdminURL, MonitorPath, app)
	w, err := requestAndUnmarshal(url, s)
	if err != nil {
		return nil, err
	}
	var application Application
	if err := json.Unmarshal(w.Body.Item, &application); err != nil {
		return nil, err
	}
	return &application, nil
}
