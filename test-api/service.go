package user

import(
		"log"
		"golang.org/x/net/context"
		"github.com/jinzhu/gorm"
)

//data structs
type UserDetails struct {

	UserId 	string   `json:"userid"`
	Site    int 	 `json:"site"`
	UserName string  `json:"userName"`
}

type User struct {
	// User Common Field
	UserId           string `gorm:"column:user_id" json:"user_id"`
	FullName         string `gorm:"column:full_name" json:"full_name"`
	FirstName        string `gorm:"column:first_name" json:"first_name"`
	LastName         string `gorm:"column:last_name" json:"last_name"`
	EmailAddress     string `gorm:"column:email_address" json:"email_address"`
	Status           string `gorm:"column:status" json:"status"`
	SystemDateFormat string `gorm:"column:system_date_format" json:"system_date_format"`
	SystemTimeFormat string `json:"system_time_format"`
	SystemTimezone   string `gorm:"column:system_timezone" json:"system_timezone"`
	LoggedIn         int    `gorm:"column:logged_in" json:"logged_in"`
	IsActive         int    `gorm:"column:is_active" json:"is_active"`
	SearchPreference string `gorm:"column:search_preference" json:"search_preference"`
	// Role related Field
	RoleId          int    `gorm:"column:role_id" json:"role_id"`
	RoleName        string `gorm:"column:role_name" json:"role_name"`
	RoleDescription string `gorm:"column:role_description" json:"role_description"`
	// Page or UI related access fields
	RoleAccessId                int    `gorm:"column:role_access_id" json:"role_access_id"`
	MenuId                      int    `gorm:"column:menu_id" json:"menu_id"`
	CanAdd                      int    `gorm:"column:can_add" json:"can_add"`
	CanEdit                     int    `gorm:"column:can_edit" json:"can_edit"`
	CanDelete                   int    `gorm:"column:can_delete" json:"can_delete"`
	CanView                     int    `gorm:"column:can_view" json:"can_view"`
	CanExport                   int    `gorm:"column:can_export" json:"can_export"`
	CanDownload                 int    `gorm:"column:can_download" json:"can_download"`
	CanUpload                   int    `gorm:"column:can_upload" json:"can_upload"`
	CanSearch                   int    `gorm:"column:can_search" json:"can_search"`
	ArchiveComment              int    `gorm:"column:archive_comment" json:"archive_comment"`
	ViewAuditLog                int    `gorm:"column:view_auditlog" json:"view_auditlog"`
	AddComment                  int    `gorm:"column:add_comment" json:"add_comment"`
	CreateDigestBatch           int    `gorm:"column:create_digest_batch" json:"create_digest_batch"`
	ViewQuoteErrors             int    `gorm:"column:view_quote_errors" json:"view_quote_errors"`
	ViewDupQuotes               int    `gorm:"column:view_dup_quotes" json:"view_dup_quotes"`
	ApproveDupQuote             int    `gorm:"column:approve_dup_quote" json:"approve_dup_quote"`
	RejectDupQuote              int    `gorm:"column:reject_dup_quote" json:"reject_dup_quote"`
	PreviewDigestResponse       int    `gorm:"column:preview_digest_response" json:"preview_digest_response"`
	AcceptDigestResponsePreview int    `gorm:"column:accept_digest_response_preview" json:"accept_digest_response_preview"`
	CancelDigestResponsePreview int    `gorm:"column:cancel_digest_response_preview" json:"cancel_digest_response_preview"`
	DownloadBatchDigest         int    `gorm:"column:download_batch_digest" json:"download_batch_digest"`
	CreateDigestWeeklyReport    int    `gorm:"column:create_digest_weekly_report" json:"create_digest_weekly_report"`
	MapPartnameToVendor         int    `gorm:"column:map_partname_to_vendor" json:"map_partname_to_vendor"`
	ViewArchive                 int    `gorm:"column:view_archive" json:"view_archive"`
	ViewDupSKUMapping           int    `gorm:"column:view_dup_skumapping" json:"view_dup_skumapping"`
	AcceptRejectDupSKUMapping   int    `gorm:"column:accept_reject_dup_skumapping" json:"accept_reject_dup_skumapping"`
	PagePropertyAccess          string `gorm:"column:page_property_access" json:"page_property_access"`
	PageTitle                   string `gorm:"column:page_title" json:"page_title"`
	PageURL                     string `gorm:"column:page_url" json:"page_url"`
	MenuLevel                   int    `gorm:"column:menu_level" json:"menu_level"`
	ParentId                    int    `gorm:"column:parent_id" json:"parent_id"`
	SortOrder                   int    `gorm:"column:sort_order" json:"sort_order"`
	Enabled                     int    `gorm:"column:enabled" json:"enabled"`
	Menu                        int    `gorm:"column:menu" json:"menu"`
}


//interface
type Service interface {
	GetUsers (ctx context.Context, user_id string)(userDetails User{}, err error)
	//SaveUser (ctx context.Context,  string) (bool)
}

//data output struct
//type data[serviceName]Service struct {
type dataUserService struct {
	c UserDetails
}

//constructor for our "object"
//func "New[serviceName]Service"
func NewGetUsersService() Service {
	return dataUserService{
		c: UserDetails
	}
}

//do the function
func (c dataUserService) GetUsers(ctx context.Context, user_id string) (userDetails []User{}, err error) {

	config := NewConfig()
	db, err := gorm.Open(config.DBType, config.ConnectionString)
	if err != nil {
		return userDetails, er
	}
	defer db.Close()

	qry := "SELECT *, group_concat(page_property) as page_property_access FROM users a INNER JOIN user_roles b using (role_id, site_id) INNER JOIN role_access c using (role_id) INNER JOIN system_menus d using (menu_id, site_id) LEFT JOIN role_page_properties e using (role_access_id) LEFT JOIN page_properties f using (page_property_id) WHERE user_id = ? and a.site_id= ? and d.enabled = 1 group by c.role_access_id order by sort_order"


	err = db.Raw(qry,user_id).Scan(&userDetails)
	if err != nil {
		log.Println("ERR:", err)
	}

	return userDetails, nil

	
}