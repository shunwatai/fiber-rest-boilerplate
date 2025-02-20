package groupUser

import (
	"encoding/json"
	"errors"
	"golang-api-starter/internal/database"
	"golang-api-starter/internal/helper"
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"golang-api-starter/internal/helper/utils"
	"golang-api-starter/internal/modules/groupResourceAcl"
	"log"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

type GroupDto struct {
	MongoId     helper.Optional[string]                               `json:"_id,omitempty"`
	Id          helper.Optional[helper.FlexInt]                       `json:"id" example:"2"`
	Name        helper.Optional[string]                               `json:"name" example:"emma"`
	Type        helper.Optional[string]                               `json:"type,omitempty"`
	Users       helper.Optional[[]*User]                              `json:"users"`
	Permissions helper.Optional[[]*groupResourceAcl.GroupResourceAcl] `json:"permissions"`
	Disabled    helper.Optional[bool]                                 `json:"disabled"`
	CreatedAt   helper.Optional[helper.CustomDatetime]                `json:"createdAt"`
	UpdatedAt   helper.Optional[helper.CustomDatetime]                `json:"updatedAt"`
}
type GroupsDto []*GroupDto

func (groupDto *GroupDto) Validate(action string) error {
	var ignoreRequiredCheckIfUpdate = func(presented bool) bool {
		if action == "update" {
			return true
		}
		return presented
	}
	var validateErrs []error
	var validations = map[string]map[bool]map[string]map[string]any{
		"name": {
			ignoreRequiredCheckIfUpdate(groupDto.Name.Presented): {
				"omitempty,alphanumspace": {"only allow A-Z1-9, no space": groupDto.Name.Value},
			},
		},
		"type": {
			ignoreRequiredCheckIfUpdate(groupDto.Type.Presented): {
				"omitempty,alphanum": {"only allow A-Z1-9, no space": groupDto.Type.Value},
			},
		},
	}
	for key, presentedRuleValue := range validations {
		for presented, ruleErrmsgValue := range presentedRuleValue {
			if !presented {
				validateErrs = append(validateErrs, errors.New(key+" is required"))
			} else {
				for rule, errmsgValue := range ruleErrmsgValue {
					for errMsg, value := range errmsgValue {
						if err := helper.Validate.Var(value, rule); err != nil {
							validateErrs = append(validateErrs, errors.New(key+" "+errMsg))
						}
					}
				}
			}
		}
	}

	return errors.Join(validateErrs...)
}

func (groupDto *GroupDto) GetId() string {
	if cfg.DbConf.Driver == "mongodb" && groupDto.MongoId.Presented {
		return *groupDto.MongoId.Value
	} else if groupDto.Id.Presented {
		return strconv.Itoa(int(*groupDto.Id.Value))
	} else {
		return ""
	}
}

func (gd *GroupDto) MapToGroup(group *Group) {
	if gd.MongoId.Presented {
		group.MongoId = gd.MongoId.Value
	}
	if gd.Id.Presented {
		group.Id = gd.Id.Value
	}
	if gd.Name.Presented {
		group.Name = *gd.Name.Value
	}
	if gd.Type.Presented {
		group.Type = *gd.Type.Value
	}
	if gd.Users.Presented {
		logger.Debugf("gu>>>>>>>>>>>>>>>????????? %+v", gd.Users)
		group.Users = gd.Users.Value
		logger.Debugf("gu>>>>>>>>>>>>>>>?????????2 %+v", group.Users)
	}
	if gd.Permissions.Presented {
		group.Permissions = gd.Permissions.Value
	}
	if gd.Disabled.Presented {
		group.Disabled = *gd.Disabled.Value
	}
	if gd.CreatedAt.Presented {
		group.CreatedAt = gd.CreatedAt.Value
	}
	if gd.UpdatedAt.Presented {
		group.UpdatedAt = gd.UpdatedAt.Value
	}
}

type Group struct {
	MongoId     *string                               `json:"_id,omitempty" bson:"_id,omitempty" validate:"omitempty,id_custom_validation"`
	Id          *helper.FlexInt                       `json:"id" db:"id" bson:"id,omitempty" example:"2" validate:"omitempty,id_custom_validation"`
	Name        string                                `json:"name" db:"name" bson:"name,omitempty" validate:"required"`
	Type        string                                `json:"type,omitempty" db:"type" bson:"type,omitempty"`
	Users       *[]*User                              `json:"users" validate:"required"`
	Permissions *[]*groupResourceAcl.GroupResourceAcl `json:"permissions" validate:"required"`
	Disabled    bool                                  `json:"disabled" db:"disabled" bson:"disabled,omitempty" validate:"boolean"`
	CreatedAt   *helper.CustomDatetime                `json:"createdAt" db:"created_at" bson:"created_at,omitempty"`
	UpdatedAt   *helper.CustomDatetime                `json:"updatedAt" db:"updated_at" bson:"updated_at,omitempty"`
}
type Groups []*Group

func (g *Group) GetId() string {
	if cfg.DbConf.Driver == "mongodb" {
		return *g.MongoId
	} else {
		return strconv.Itoa(int(*g.Id))
	}
}

func (gs Groups) StructToMap() []map[string]interface{} {
	mapsResults := []map[string]interface{}{}
	for _, g := range gs {
		tmp := map[string]interface{}{}
		result := map[string]interface{}{}
		data, _ := json.Marshal(g)
		json.Unmarshal(data, &tmp)
		for k, v := range tmp {
			result[strcase.ToSnake(k)] = v
		}
		mapsResults = append(mapsResults, result)
	}

	return mapsResults
}

func (gs Groups) RowsToStruct(rows database.Rows) []*Group {
	defer rows.Close()

	records := make([]*Group, 0)
	for rows.Next() {
		var g Group
		err := rows.StructScan(&g)
		if err != nil {
			log.Fatalf("Scan: %v", err)
		}
		records = append(records, &g)
	}

	return records
}

func (gs Groups) GetTags(key ...string) []string {
	if len(gs) == 0 {
		return []string{}
	}

	return gs[0].getTags(key...)
}

func (gs *Groups) PrintValue() {
	for _, v := range *gs {
		if v.Id != nil {
			logger.Debugf("existing --> id: %+v, v: %+v\n", v.GetId(), *v)
		} else {
			logger.Debugf("new --> v: %+v\n", *v)
		}
	}
}

// get the tags by key(json / db / bson) name from the struct
// ref: https://stackoverflow.com/a/40865028
func (g Group) getTags(key ...string) []string {
	var tag string
	if len(key) == 1 {
		tag = key[0]
	} else if cfg.DbConf.Driver == "mongodb" {
		tag = "bson"
	} else {
		tag = "db"
	}

	cols := []string{}
	val := reflect.ValueOf(g)
	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)
		fieldName := t.Name

		switch jsonTag := t.Tag.Get(tag); jsonTag {
		case "-":
		case "":
			// fmt.Println(fieldName)
		default:
			parts := strings.Split(jsonTag, ",")
			name := parts[0]
			if name == "" {
				name = fieldName
			}
			// fmt.Println(name)
			if !slices.Contains(*database.IgnrCols, name) {
				cols = append(cols, name)
			}
		}
	}
	return cols
}

func (groups Groups) SetUsers() {
	if len(groups) == 0 {
		return
	}

	groupIds := make([]string, 0, len(groups))
	for _, group := range groups {
		groupIds = append(groupIds, group.GetId())
	}

	condition := database.GetIdsMapCondition(utils.ToPtr("group_id"), groupIds)
	groupUsers, _ := Repo.Get(condition)

	groupUsersMap := Repo.GetGroupIdMap(groupUsers)

	// map users into group
	for _, group := range groups {
		// if no users, assign empty slice for response json "users": [] instead of "users": null
		group.Users = new([]*User)
		*group.Users = make([]*User, 0)
		// take out the groupUsers by groupId in map and assign
		if gus, haveUsers := groupUsersMap[group.GetId()]; haveUsers {
			for _, gu := range gus {
				gu.User.Groups = nil
				*group.Users = append(*group.Users, gu.User)
			}
		}
	}
}

func (groups Groups) SetPermissions() {
	if len(groups) == 0 {
		return
	}

	groupIds := make([]string, 0, len(groups))
	for _, group := range groups {
		groupIds = append(groupIds, group.GetId())
	}

	condition := database.GetIdsMapCondition(utils.ToPtr("group_id"), groupIds)
	groupResourceAcls, _ := groupResourceAcl.Repo.Get(condition)
	groupAclsMap := groupResourceAcl.Repo.GetGroupIdMap(groupResourceAcls)

	// map permission into group
	for _, group := range groups {
		// if no permissions, assign empty slice for response json "permissions": [] instead of "permissions": null
		group.Permissions = new([]*groupResourceAcl.GroupResourceAcl)
		*group.Permissions = make([]*groupResourceAcl.GroupResourceAcl, 0)
		// take out the groupUsers by groupId in map and assign
		if gas, haveAcls := groupAclsMap[group.GetId()]; haveAcls {
			*group.Permissions = append(*group.Permissions, gas...)
		}
	}
}
