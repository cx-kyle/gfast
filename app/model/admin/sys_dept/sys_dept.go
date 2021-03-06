// ============================================================================
// This is auto-generated by gf cli tool only once. Fill this file as you wish.
// ============================================================================

package sys_dept

import (
	"database/sql"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

//数据列表数据结构
type Dept struct {
	SearchValue interface{}            `json:"searchValue"`
	CreateBy    string                 `json:"createBy"  orm:"create_by"`
	CreateTime  *gtime.Time            `json:"createTime" orm:"create_time"`
	UpdateBy    string                 `json:"updateBy"  orm:"update_by"`
	UpdateTime  *gtime.Time            `json:"updateTime" orm:"update_time"`
	Remark      string                 `json:"remark"`
	DataScope   interface{}            `json:"dataScope"`
	Params      map[string]interface{} `json:"params"`
	DeptID      int64                  `json:"deptId" orm:"dept_id"`
	ParentID    int64                  `json:"parentId" orm:"parent_id"`
	Ancestors   string                 `json:"ancestors" orm:"ancestors"`
	DeptName    string                 `json:"deptName" orm:"dept_name"`
	OrderNum    int                    `json:"orderNum" orm:"order_num" `
	Leader      string                 `json:"leader" orm:"leader"`
	Phone       string                 `json:"phone" orm:"phone"`
	Email       string                 `json:"email" orm:"email"`
	Status      string                 `json:"status" orm:"status"`
	DelFlag     string                 `json:"delFlag" orm:"del_flag" `
	ParentName  string                 `json:"parentName"`
	Children    []interface{}          `json:"children"`
}

//查询参数

//文章搜索参数
type SearchParams struct {
	DeptName string `p:"deptName"`
	Status   string `p:"status"`
}

type AddParams struct {
	ParentID   int         `json:"parentId" orm:"parent_id"  p:"parentId"  v:"required#父级不能为空"`
	DeptName   string      `json:"deptName" orm:"dept_name" p:"deptName"  v:"required#部门名称不能为空"`
	OrderNum   int         `json:"orderNum" orm:"order_num"  p:"orderNum"  v:"required#排序不能为空"`
	Leader     string      `json:"leader" orm:"leader" p:"leader"  v:"required#负责人不能为空"`
	Phone      string      `json:"phone" orm:"Phone" p:"phone"  v:"required#电话不能为空"`
	Email      string      `json:"email" orm:"email" p:"email"  v:"required#邮箱不能为空"`
	Status     string      `json:"status" orm:"status" p:"status"  v:"required#状态必须"`
	Ancestors  string      `json:"ancestors" orm:"ancestors"`
	DelFlag    string      `json:"delFlag" orm:"del_flag"`
	CreateBy   string      `json:"createBy"  orm:"create_by"`
	CreateTime *gtime.Time `json:"createTime" orm:"create_time"`
	UpdateBy   string      `json:"updateBy"  orm:"update_by"`
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time"`
}

type EditParams struct {
	DeptID int64 `json:"deptId" orm:"dept_id" p:"id" v:"integer|min:1#ID只能为整数|ID只能为正数"`
	AddParams
}

//获取列表数据
func GetList(searchParams *SearchParams) ([]*Dept, error) {
	model := g.DB().Table(Table)
	if searchParams.DeptName != "" {
		model.Where("dept_name like ?", "%"+searchParams.DeptName+"%")
	}

	if searchParams.Status != "" {
		model.Where("status", searchParams.Status)
	}
	depts := ([]*Dept)(nil)
	if err := model.Structs(&depts); err != nil {
		return nil, err
	}

	for _, v := range depts {
		if v.Children == nil {
			v.Children = []interface{}{}
		}

	}

	return depts, nil
}

//添加
func AddDept(data *AddParams) (sql.Result, error) {
	data.DelFlag = "0"
	data.CreateBy = ""
	data.CreateTime = gtime.Now()
	return Model.Data(data).Insert()
}

//编辑
func EditDept(data *EditParams) error {
	data.UpdateBy = ""
	data.UpdateTime = gtime.Now()

	if _, err := Model.Where("dept_id", data.DeptID).Data(data).Update(); err != nil {
		return err
	}
	return nil
}

//删除失败
func DelDept(id int64) error {

	ids, _ := GetChilderenIds(id)
	_, err := Model.Where("dept_id IN(?)", ids).Delete()
	if err != nil {
		return gerror.New("删除失败")
	}
	return nil
}

//根据部门id获取部门信息
func GetDeptById(id int64) (*Dept, error) {

	dept := (*Dept)(nil)

	if err := Model.Where("dept_id", id).Struct(&dept); err != nil {
		return nil, err
	} else {
		return dept, nil
	}

}

/**
获取排除节点
*/
func Exclude(id int64) ([]*Dept, error) {
	ids, err := GetChilderenIds(id)
	if err != nil {
		return nil, err
	}

	model := g.DB().Table(Table)
	if len(ids) > 0 {
		model.Where("dept_id  NOT IN(?)", ids)
	}

	depts := ([]*Dept)(nil)
	if err := model.Structs(&depts); err != nil {
		return nil, err
	}

	for _, v := range depts {
		if v.Children == nil {
			v.Children = []interface{}{}
		}

	}
	return depts, nil
}

/**
根据id获取子孙节点id集合包含本身
*/
func GetChilderenIds(id int64) ([]int64, error) {

	list, err := GetChildrenById(id)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return []int64{}, nil
	}
	var newResult []int64
	for _, v := range list {
		newResult = append(newResult, v.DeptId)
	}

	return newResult, nil
}

/**
根据id获取所有子孙节点包含本身
*/
func GetChildrenById(id int64) ([]*Entity, error) {
	depts, err := Model.All()
	if err != nil {
		return nil, err
	}
	result := recursion(id, depts, true)
	return result, nil
}

/**
根据id获取所有子孙元素

hasroot true - 包含自身  false - 不含自身
*/
func recursion(id int64, depts []*Entity, hasRoot bool) (result []*Entity) {
	for _, v := range depts {
		if hasRoot == true && v.DeptId == id {
			result = append(result, v)
		}
		if v.ParentId == id {
			data := recursion(v.DeptId, depts, false)
			result = append(result, v)
			result = append(result, data...)
		}
	}

	return
}
