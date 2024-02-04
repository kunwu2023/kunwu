package dbpos

type QuickCheckList struct {
	ID               int64  `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt        int64  `gorm:"column:created_at" json:"createdAt"`
	Path             string `gorm:"column:path" json:"path"`
	Results          string `gorm:"column:results" json:"results"`
	ModificationTime int64  `gorm:"column:modificationTime" json:"modificationTime"`
	CloudResultsFlag int64  `gorm:"column:cloud_results_flag;default:1" json:"cloudResultsFlag"` // 1本地检测、2云端正在检测、3云端检测已完成、4云端检测异常
	Size             string `gorm:"column:size" json:"size"`
}

// TableName get sql table name.获取数据库表名
func (m *QuickCheckList) TableName() string {
	return "quick_check_list"
}

// TaskBase [...]
type TaskBase struct {
	ID        int64  `gorm:"primaryKey;column:id" json:"id"`
	CreatedAt int64  `gorm:"column:created_at" json:"createdAt"`
	Model     int64  `gorm:"column:model" json:"model"`     // 1 本地扫描 2 远程扫描
	Status    int64  `gorm:"column:status" json:"status"`   // 1 准备中、 2、扫描中、3、已完成、4、远程扫描异常
	DirPath   string `gorm:"column:dirPath" json:"dirPath"` // 任务名称
	Ip        string `gorm:"column:ip" json:"ip"`
	Port      string `gorm:"column:port" json:"port"`
	UserName  string `gorm:"column:userName" json:"userName"`
	PassWord  string `gorm:"column:passWord" json:"passWord"`
}

// TableName get sql table name.获取数据库表名
func (m *TaskBase) TableName() string {
	return "task_base"
}

// TaskCheckList [...]
type TaskCheckList struct {
	ID               int64  `gorm:"primaryKey;column:id" json:"id"`
	TaskBaseID       int64  `gorm:"column:task_base_id" json:"taskBaseId"`
	Path             string `gorm:"column:path" json:"path"`
	Results          string `gorm:"column:results" json:"results"`
	ModificationTime int64  `gorm:"column:modificationTime" json:"modificationTime"`
	Size             string `gorm:"column:size" json:"size"`
	CloudResultsFlag int64  `gorm:"column:cloud_results_flag;default:1" json:"cloudResultsFlag"` // 1本地检测、2云端正在检测、3云端检测已完成、4云端检测异常
	CreatedAt        int64  `gorm:"column:created_at" json:"createdAt"`
}

// TableName get sql table name.获取数据库表名
func (m *TaskCheckList) TableName() string {
	return "task_check_list"
}
