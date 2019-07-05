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

package clusterfeatureadapter

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"emperror.dev/emperror"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cast"

	"github.com/banzaicloud/pipeline/internal/clusterfeature"
)

// TableName constants
const (
	clusterFeatureTableName = "cluster_features"
)

type featureSpec map[string]interface{}

func (fs *featureSpec) Scan(src interface{}) error {
	value, err := cast.ToStringE(src)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(value), fs)
}

func (fs featureSpec) Value() (driver.Value, error) {
	v, err := json.Marshal(fs)
	if err != nil {
		return "", err
	}
	return v, nil
}

// clusterFeatureModel describes the cluster group model.
type clusterFeatureModel struct {
	// injecting timestamp fields
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name      string
	Status    string
	ClusterId uint
	Spec      featureSpec `gorm:"type:text"`
	CreatedBy uint
}

// TableName changes the default table name.
func (cfm clusterFeatureModel) TableName() string {
	return clusterFeatureTableName
}

// String method prints formatted cluster fields.
func (cfm clusterFeatureModel) String() string {
	return fmt.Sprintf("Id: %d, Creation date: %s, Name: %s", cfm.ID, cfm.CreatedAt, cfm.Name)
}

// GormFeatureRepository component in charge for executing persistence operation on Features.
// TODO: write integration tests
type GormFeatureRepository struct {
	db *gorm.DB
}

// NewGormFeatureRepository returns a feature repository persisting feature state into database using Gorm.
func NewGormFeatureRepository(db *gorm.DB) *GormFeatureRepository {
	return &GormFeatureRepository{db: db}
}

func (r *GormFeatureRepository) SaveFeature(ctx context.Context, clusterId uint, feature clusterfeature.Feature) (uint, error) {
	cfModel := clusterFeatureModel{
		Name:      feature.Name,
		Spec:      feature.Spec,
		ClusterId: clusterId,
		Status:    string(clusterfeature.FeatureStatusPending),
	}

	err := r.db.Save(&cfModel).Error
	if err != nil {
		return 0, emperror.WrapWith(err, "failed to persist feature", "feature", feature.Name)
	}

	return cfModel.ID, nil
}

func (r *GormFeatureRepository) GetFeature(ctx context.Context, clusterId uint, featureName string) (*clusterfeature.Feature, error) {
	fm := clusterFeatureModel{}

	err := r.db.First(&fm, map[string]interface{}{"Name": featureName, "cluster_id": clusterId}).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, emperror.WrapWith(err, "cluster feature not found", "feature", featureName)
	} else if err != nil {
		return nil, emperror.Wrap(err, "could not retrieve feature")
	}

	return r.modelToFeature(&fm)
}

// UpdateFeatureStatus updates the status of the feature
func (r *GormFeatureRepository) UpdateFeatureStatus(ctx context.Context, clusterId uint, featureName string, status string) (*clusterfeature.Feature, error) {
	fm := clusterFeatureModel{
		ClusterId: clusterId,
		Name:      featureName,
	}

	if err := r.db.Find(&fm, fm).Updates(clusterFeatureModel{Status: status}).Error; err != nil {
		return nil, emperror.Wrap(err, "could not update feature status")
	}

	return r.modelToFeature(&fm)
}

// UpdateFeatureStatus updates the status of the feature
func (r *GormFeatureRepository) UpdateFeatureSpec(ctx context.Context, clusterId uint, featureName string, spec map[string]interface{}) (*clusterfeature.Feature, error) {

	fm := clusterFeatureModel{ClusterId: clusterId, Name: featureName}

	if err := r.db.Find(&fm, fm).Updates(clusterFeatureModel{Spec: spec}).Error; err != nil {
		return nil, emperror.Wrap(err, "could not update feature spec")
	}

	return r.modelToFeature(&fm)
}

func (r *GormFeatureRepository) modelToFeature(cfm *clusterFeatureModel) (*clusterfeature.Feature, error) {
	f := clusterfeature.Feature{
		Name:   cfm.Name,
		Status: cfm.Status,
		Spec:   cfm.Spec,
	}

	return &f, nil
}

// DeleteFeature permanently deletes the feature record
func (r *GormFeatureRepository) DeleteFeature(ctx context.Context, clusterId uint, featureName string) error {

	fm := clusterFeatureModel{ClusterId: clusterId, Name: featureName}

	if err := r.db.Delete(&fm, fm).Error; err != nil {
		return emperror.Wrap(err, "could not delete status")
	}

	return nil

}

func (r *GormFeatureRepository) ListFeatures(ctx context.Context, clusterId uint) ([]clusterfeature.Feature, error) {
	var (
		featureModels []clusterFeatureModel
		featureList   []clusterfeature.Feature
	)

	if err := r.db.Find(&featureModels, clusterFeatureModel{ClusterId: clusterId}).Error; err != nil {
		return nil, emperror.WrapWith(err, "could not retrieve features", "clusterID", clusterId)
	}

	// model  --> domain
	for _, feature := range featureModels {
		f, e := r.modelToFeature(&feature)
		if e != nil {
			// skip erroneous model
			continue
		}

		featureList = append(featureList, *f)
	}

	return featureList, nil

}