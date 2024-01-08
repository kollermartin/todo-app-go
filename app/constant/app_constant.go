package constant

const (
	GetTodosLogEventErrorKey   string = "todo_get_all_fail"
	GetTodosLogEventKey        string = "todo_get_all"
	CreateTodoLogEventKey      string = "todo_create"
	CreateTodoLogEventErrorKey string = "todo_create_fail"
	GetTodoLogEventErrorKey    string = "todo_get_fail"
	GetTodoLogEventKey         string = "todo_get"
	UpdateTodoLogEventKey      string = "todo_update"
	UpdateTodoLogEventErrorKey string = "todo_update_fail"
	DeleteTodoLogEventKey      string = "todo_delete"
	DeleteTodoLogEventErrorKey string = "todo_delete_fail"
	DbIdNotFoundMsg            string = "Id not found"
	DbQueryFailMsg             string = "Failed to query database"
	DbExecFailMsg              string = "Failed to execute database query"
	DbRowsAffectedFailMsg      string = "Failed to get rows affected"
	DbScanFailMsg              string = "Failed to scan database row"
	ErrMsgInternalServer       string = "Internal server error"
	ConfigLoadLogEventErrorKey string = "config_load_fail"
	DbInitErrorEventKey        string = "db_init_fail"
)
