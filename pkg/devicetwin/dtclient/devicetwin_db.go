package dtclient

import (
	"kubeedge/beehive/pkg/common/log"
	"kubeedge/pkg/common/dbm"
)

//DeviceTwin the struct of device twin
type DeviceTwin struct {
	ID              int64  `orm:"column(id);size(64);auto;pk"`
	DeviceID        string `orm:"column(deviceid); null; type(text)"`
	Name            string `orm:"column(name);null;type(text)"`
	Description     string `orm:"column(description);null;type(text)"`
	Expected        string `orm:"column(expected);null;type(text)"`
	Actual          string `orm:"column(actual);null;type(text)"`
	ExpectedMeta    string `orm:"column(expected_meta);null;type(text)"`
	ActualMeta      string `orm:"column(actual_meta);null;type(text)"`
	ExpectedVersion string `orm:"column(expected_version);null;type(text)"`
	ActualVersion   string `orm:"column(actual_version);null;type(text)"`
	Optional        bool   `orm:"column(optional);null;type(integer)"`
	AttrType        string `orm:"column(attr_type);null;type(text)"`
	Metadata        string `orm:"column(metadata);null;type(text)"`
}

//SaveDeviceTwin save device twin
func SaveDeviceTwin(doc *DeviceTwin) error {
	num, err := dbm.DBAccess.Insert(doc)
	log.LOGGER.Debugf("Insert affected Num: %d, %s", num, err)
	return err
}

//DeleteDeviceTwinByDeviceID delete device twin
func DeleteDeviceTwinByDeviceID(deviceID string) error {
	num, err := dbm.DBAccess.QueryTable(DeviceTwinTableName).Filter("deviceid", deviceID).Delete()
	if err != nil {
		log.LOGGER.Errorf("Something wrong when deleting data: %v", err)
		return err
	}
	log.LOGGER.Debugf("Delete affected Num: %d, %s", num)
	return nil
}

//DeleteDeviceTwin delete device twin
func DeleteDeviceTwin(deviceID string, name string) error {
	num, err := dbm.DBAccess.QueryTable(DeviceTwinTableName).Filter("deviceid", deviceID).Filter("name", name).Delete()
	if err != nil {
		log.LOGGER.Errorf("Something wrong when deleting data: %v", err)
		return err
	}
	log.LOGGER.Debugf("Delete affected Num: %d, %s", num)
	return nil
}

// UpdateDeviceTwinField update special field
func UpdateDeviceTwinField(deviceID string, name string, col string, value interface{}) error {
	num, err := dbm.DBAccess.QueryTable(DeviceTwinTableName).Filter("deviceid", deviceID).Filter("name", name).Update(map[string]interface{}{col: value})
	log.LOGGER.Debugf("Update affected Num: %d, %s", num, err)
	return err
}

// UpdateDeviceTwinFields update special fields
func UpdateDeviceTwinFields(deviceID string, name string, cols map[string]interface{}) error {
	num, err := dbm.DBAccess.QueryTable(DeviceTwinTableName).Filter("deviceid", deviceID).Filter("name", name).Update(cols)
	log.LOGGER.Debugf("Update affected Num: %d, %v", num, err)
	return err
}

// QueryDeviceTwin query Device
func QueryDeviceTwin(key string, condition string) (*[]DeviceTwin, error) {
	twin := new([]DeviceTwin)
	_, err := dbm.DBAccess.QueryTable(DeviceTwinTableName).Filter(key, condition).All(twin)
	if err != nil {
		return nil, err
	}
	return twin, nil

}

//DeviceTwinUpdate the struct for updating device twin
type DeviceTwinUpdate struct {
	DeviceID string
	Name     string
	Cols     map[string]interface{}
}

//UpdateDeviceTwinMulti update device twin multi
func UpdateDeviceTwinMulti(updates []DeviceTwinUpdate) error {
	var err error
	for _, update := range updates {
		err = UpdateDeviceTwinFields(update.DeviceID, update.Name, update.Cols)
		if err != nil {
			return err
		}
	}
	return nil
}

//DeviceTwinTrans transaction of device twin
func DeviceTwinTrans(adds []DeviceTwin, deletes []DeviceDelete, updates []DeviceTwinUpdate) error {
	var err error
	obm := dbm.DBAccess
	obm.Begin()
	for _, add := range adds {
		err = SaveDeviceTwin(&add)
		if err != nil {
			obm.Rollback()
			return err
		}
	}

	for _, delete := range deletes {
		err = DeleteDeviceTwin(delete.DeviceID, delete.Name)
		if err != nil {
			obm.Rollback()
			return err
		}
	}

	for _, update := range updates {
		err = UpdateDeviceTwinFields(update.DeviceID, update.Name, update.Cols)
		if err != nil {
			obm.Rollback()
			return err
		}
	}
	obm.Commit()
	return nil
}
