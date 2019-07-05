// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clusterfeature

import (
	"errors"

	"github.com/goph/logur"
)

// FeatureManagerRegistry operations related to the set of features supported by the application
type FeatureManagerRegistry interface {
	// RegisterFeatureManager registers the feature manager for the given feature name
	RegisterFeatureManager(featureName string, featureManager FeatureManager) error

	//GetFeatureManager retrieves a feature manager by the feature name
	GetFeatureManager(featureName string) (FeatureManager, error)
}

type featureManagerRegistry struct {
	logger   logur.Logger
	registry map[string]FeatureManager
}

func (fr *featureManagerRegistry) RegisterFeatureManager(featureName string, featureManager FeatureManager) error {
	log := logur.WithFields(fr.logger, map[string]interface{}{"feature": featureName})
	log.Info("registering feature ...")

	if _, ok := fr.registry[featureName]; ok {
		log.Debug("feature already registered")
		return errors.New("feature already registered")
	}

	fr.registry[featureName] = featureManager

	log.Info("feature registered")
	return nil
}

func (fr *featureManagerRegistry) GetFeatureManager(featureName string) (FeatureManager, error) {
	log := logur.WithFields(fr.logger, map[string]interface{}{"feature": featureName})
	log.Info("retrieving feature ...")

	if fm, ok := fr.registry[featureName]; ok {
		log.Info("feature retrieved")
		return fm, nil
	}

	log.Debug("feature is not supported")
	return nil, errors.New("feature already registered")

}

func NewFeatureManagerRegistry(logger logur.Logger) FeatureManagerRegistry {
	return &featureManagerRegistry{
		logger:   logur.WithFields(logger, map[string]interface{}{"component": "featureManagerRegistry"}),
		registry: make(map[string]FeatureManager),
	}
}