package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/skuid/helm-value-store/store"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type applyRequest struct {
	UUID string `json:"uuid"`
}

type applyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type upsertConfig struct {
	location string
	timeout  int64
}

func upsertRelease(r *store.Release, conf upsertConfig) error {
	_, err := r.Get()

	if err != nil && !strings.Contains(err.Error(), "not found") {
		return err
	}

	if err != nil && strings.Contains(err.Error(), "not found") {
		_, err = r.Install(conf.location, false, conf.timeout)
	} else if err == nil {
		_, err = r.Upgrade(conf.location, false, conf.timeout)
	}

	return err
}

// ApplyChart applies a chart to a tiller server
func (c ApiController) ApplyChart(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Fields for the autit log
	var auditFields []zapcore.Field
	for _, a := range c.authorizers {
		auditFields = append(auditFields, a.LoggingClosure(r)...)
	}

	defer func() {
		successful := err == nil
		auditFields = append(
			auditFields,
			zap.String("controller", "apply"),
			zap.Bool("successful", successful),
		)
		zap.L().Info("Audit Log", auditFields...)
	}()

	applyReq := &applyRequest{}

	if err = json.NewDecoder(r.Body).Decode(applyReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		zap.L().Error("Error decoding request", zap.Error(err))
		return
	}
	auditFields = append(auditFields, zap.String("uuid", applyReq.UUID))

	release := &store.Release{}
	release, err = c.releaseStore.Get(r.Context(), applyReq.UUID)

	applyResp := &applyResponse{}

	if err != nil {
		zap.L().Error("Error getting release", zap.Error(err))

		applyResp.Status = "error"
		applyResp.Message = "Error getting release"
		err = json.NewEncoder(w).Encode(applyResp)
		if err != nil {
			zap.L().Error("Error marshaling response", zap.Error(err))
		}
		return
	}
	auditFields = append(auditFields,
		zap.String("chart", release.Chart),
		zap.String("release", release.Name),
		zap.String("version", release.Version),
		zap.String("namespace", release.Namespace),
	)

	var location string
	location, err = release.Download()

	if err != nil {
		zap.L().Error("Error downloading release", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		applyResp.Status = "error"
		applyResp.Message = "Error downloading release"
		err = json.NewEncoder(w).Encode(applyResp)

		if err != nil {
			zap.L().Error("Error marshaling response", zap.Error(err))
		}
		return
	}

	err = upsertRelease(release, upsertConfig{
		location: location,
		timeout:  c.timeout,
	})

	if err != nil {
		zap.L().Error("Error applying release", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		applyResp.Status = "error"
		applyResp.Message = "Error applying release"
		err = json.NewEncoder(w).Encode(applyResp)
		if err != nil {
			zap.L().Error("Error marshaling response", zap.Error(err))
		}
		return
	}

	applyResp.Status = "success"
	applyResp.Message = fmt.Sprintf("Successfully installed %s", release.Name)
	err = json.NewEncoder(w).Encode(applyResp)
	if err != nil {
		zap.L().Error("Error marshaling response", zap.Error(err))
	}
}
