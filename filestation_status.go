package filestation

import "fmt"

type FileStationStatus int

const (
	WFM2_FAIL                        FileStationStatus = 0  // UNKNOW ERROR
	WFM2_DONE                        FileStationStatus = 1  // SUCCESS
	WFM2_SUCCESS                     FileStationStatus = 1  // SUCCESS
	WFM2_FILE_EXIST                  FileStationStatus = 2  // FILE EXIST
	WFM2_AUTH_FAIL                   FileStationStatus = 3  // Authentication Failure,認 證 失 敗
	WFM2_PERMISSION_DENY             FileStationStatus = 4  // Permission Denied,存 取 拒 絕
	WFM2_FILE_NO_EXIST               FileStationStatus = 5  // FILE/FOLDER NOT EXIST,檔 案 不 存 在
	WFM2_SRC_FILE_NO_EXIST           FileStationStatus = 5  // FILE/FOLDER NOT EXIST,檔 案 不 存 在
	WFM2_EXTRACTING                  FileStationStatus = 6  // FILE EXTRACTING,檔 案 解 壓 縮 中
	WFM2_OPEN_FILE_FAIL              FileStationStatus = 7  // FILE IO ERROR,檔 案 寫 入 時 發 生 錯 誤
	WFM2_DISABLE                     FileStationStatus = 8  // Web File Manager is not enabled.,Web File Manager尚 未 啟 用
	WFM2_QUOTA_ERROR                 FileStationStatus = 9  // You have reached the disk quota limit.,您 的 磁 碟 容 量 配 額 已 滿
	WFM2_SRC_PERMISSION_DENY         FileStationStatus = 10 // You do not have permission to perform this action.,您 沒 有 權 限 進 行 此 項 操 作
	WFM2_DES_PERMISSION_DENY         FileStationStatus = 11 // You do not have permission to perform this action.,您 沒 有 權 限 進 行 此 項 操 作
	WFM2_ILLEGAL_NAME                FileStationStatus = 12 // 名 稱 不 合 法。因 為 其 中 含 有 以 下 字 元：n " + = / \ ： | * ? < > ; [ ] % , ` ' 字 元 或 特 殊 字 首 "_sn_" 和 "_sn_bk"。
	WFM2_EXCEED_ISO_MAX              FileStationStatus = 13 // The maximum number of allowed ISO shares is 256. Please unmount an ISO share first.//最 大 支 援 的 映 像 檔 資 料 夾 是256。請 先 卸 載 一 個 映 像 檔 資 料 夾。
	WFM2_EXCEED_SHARE_MAX            FileStationStatus = 14 // The maximum number of shares is going to be exceeded.分 享 的 數 目 已 到 達 最 大 的 分 享 數 目 的 限 制
	WFM2_NEED_CHECK                  FileStationStatus = 15
	WFM2_RECYCLE_BIN_NOT_ENABLE      FileStationStatus = 16
	WFM2_CHECK_PASSWORD_FAIL         FileStationStatus = 17 // Enter password,請 輸 入 密 碼
	WFM2_VIDEO_TCS_DISABLE           FileStationStatus = 18 // 媒 體 櫃 未 啟 動
	WFM2_DB_FAIL                     FileStationStatus = 19 // The system is currently busy. Please try again later.系 統 忙 碌 中，請 再 試 一 次。
	WFM2_DB_QUERY_FAIL               FileStationStatus = 19 // The system is currently busy. Please try again later.系 統 忙 碌 中，請 再 試 一 次。
	WFM2_PARAMETER_ERROR             FileStationStatus = 20 // There were input errors. Please try again later.
	WFM2_DEMO_SITE                   FileStationStatus = 21 // Your files are now being transcoded.
	WFM2_TRANSCODE_ONGOING           FileStationStatus = 22 // Your files are now being transcoded.您 的 檔 案 正 在 轉 檔 中。
	WFM2_SRC_VOLUME_ERROR            FileStationStatus = 23 // An error occurred in the source file. Please check and try again later.資 料 來 源 讀 取 異 常，請 檢 查 資 料 來 源 後 再 試 一 次。
	WFM2_DES_VOLUME_ERROR            FileStationStatus = 24 // A write error has occurred at the target destination. Please check and try again later.目 的 地 寫 入 異 常，請 檢 查 後 再 試 一 次。
	WFM2_DES_FILE_NO_EXIST           FileStationStatus = 25 // The target destination is unavailable. Please check and try again later.目 的 地 路 徑 不 存 在，請 檢 查 後 再 試 一 次。
	WFM2_FILE_NAME_TOO_LONG          FileStationStatus = 26 // The file name is too long. Please use a shorter one (maximum: 255 characters). Note that this length is for English characters. For non-English file names, please keep them shorter than the length above.名 稱 長 度 超 過 限 制，請 將 長 度 控 制 在255字 元 之 內。請 注 意，此 長 度 為 英 文 字 元 長 度。故 針 對 非 英 語 語 系 的 檔 案 名 稱，請 注 意 勿 超 過 此 長 度。
	WFM2_FOLDER_ENCRYPTION           FileStationStatus = 27 // This folder has been encrypted. Please decrypt it and try again.資 料 夾 已 加 密，請 先 解 密。
	WFM2_PREPARE                     FileStationStatus = 28 // Processing now, please wait.任 務 進 行 中，請 稍 等。
	WFM2_NO_SUPPORT_MEDIA            FileStationStatus = 29 // This file format is not supported.不 支 援 開 啟 這 類 型 的 格 式。
	WFM2_DLNA_QDMS_DISABLE           FileStationStatus = 30 // Please enable the <qtag>DLNA Media Server</qtag>.請 先 啟 動 <qtag>DLNA Media Server</qtag> 。
	WFM2_RENDER_NOT_FOUND            FileStationStatus = 31 // Cannot find any available DLNA devices.目 前 找 不 到 任 何 可 用 的 播 放 裝 置。
	WFM2_CLOUD_SERVER_ERROR          FileStationStatus = 32 // The SmartLink service is currently busy. Please try again later.SmartLink服 務 忙 碌 中，請 再 試 一 次。
	WFM2_NAME_DUP                    FileStationStatus = 33 // That folder or file name already exists. Please use another name.
	WFM2_EXCEED_SEARCH_MAX           FileStationStatus = 34 // 搜 尋 結 果 超 過1000筆
	WFM2_MEMORY_ERROR                FileStationStatus = 35
	WFM2_COMPRESSING                 FileStationStatus = 36
	WFM2_EXCEED_DAV_MAX              FileStationStatus = 37
	WFM2_UMOUNT_FAIL                 FileStationStatus = 38
	WFM2_MOUNT_FAIL                  FileStationStatus = 39
	WFM2_WEBDAV_ACCOUNT_PASSWD_ERROR FileStationStatus = 40
	WFM2_WEBDAV_SSL_ERROR            FileStationStatus = 41
	WFM2_WEBDAV_REMOUNT_ERROR        FileStationStatus = 42
	WFM2_WEBDAV_HOST_ERROR           FileStationStatus = 43
	WFM2_WEBDAV_TIMEOUT_ERROR        FileStationStatus = 44
	WFM2_WEBDAV_CONF_ERROR           FileStationStatus = 45
	WFM2_WEBDAV_BASE_ERROR           FileStationStatus = 46
)

func (s FileStationStatus) Error() string {
	switch s {
	case WFM2_FAIL:
		return "WFM2_FAIL"
	/*case WFM2_DONE:
	return "WFM2_DONE"*/
	case WFM2_SUCCESS:
		return "WFM2_SUCCESS"
	case WFM2_FILE_EXIST:
		return "WFM2_FILE_EXIST"
	case WFM2_AUTH_FAIL:
		return "WFM2_AUTH_FAIL"
	case WFM2_PERMISSION_DENY:
		return "WFM2_PERMISSION_DENY"
	case WFM2_FILE_NO_EXIST:
		return "WFM2_FILE_NO_EXIST"
	/*case WFM2_SRC_FILE_NO_EXIST:
	return "WFM2_SRC_FILE_NO_EXIST"*/
	case WFM2_EXTRACTING:
		return "WFM2_EXTRACTING"
	case WFM2_OPEN_FILE_FAIL:
		return "WFM2_OPEN_FILE_FAIL"
	case WFM2_DISABLE:
		return "WFM2_DISABLE"
	case WFM2_QUOTA_ERROR:
		return "WFM2_QUOTA_ERROR"
	case WFM2_SRC_PERMISSION_DENY:
		return "WFM2_SRC_PERMISSION_DENY"
	case WFM2_DES_PERMISSION_DENY:
		return "WFM2_DES_PERMISSION_DENY"
	case WFM2_ILLEGAL_NAME:
		return "WFM2_ILLEGAL_NAME"
	case WFM2_EXCEED_ISO_MAX:
		return "WFM2_EXCEED_ISO_MAX"
	case WFM2_EXCEED_SHARE_MAX:
		return "WFM2_EXCEED_SHARE_MAX"
	case WFM2_NEED_CHECK:
		return "WFM2_NEED_CHECK"
	case WFM2_RECYCLE_BIN_NOT_ENABLE:
		return "WFM2_RECYCLE_BIN_NOT_ENABLE"
	case WFM2_CHECK_PASSWORD_FAIL:
		return "WFM2_CHECK_PASSWORD_FAIL"
	case WFM2_VIDEO_TCS_DISABLE:
		return "WFM2_VIDEO_TCS_DISABLE"
	case WFM2_DB_FAIL:
		return "WFM2_DB_FAIL"
	/*case WFM2_DB_QUERY_FAIL:
	return "WFM2_DB_QUERY_FAIL"*/
	case WFM2_PARAMETER_ERROR:
		return "WFM2_PARAMETER_ERROR"
	case WFM2_DEMO_SITE:
		return "WFM2_DEMO_SITE"
	case WFM2_TRANSCODE_ONGOING:
		return "WFM2_TRANSCODE_ONGOING"
	case WFM2_SRC_VOLUME_ERROR:
		return "WFM2_SRC_VOLUME_ERROR"
	case WFM2_DES_VOLUME_ERROR:
		return "WFM2_DES_VOLUME_ERROR"
	case WFM2_DES_FILE_NO_EXIST:
		return "WFM2_DES_FILE_NO_EXIST"
	case WFM2_FILE_NAME_TOO_LONG:
		return "WFM2_FILE_NAME_TOO_LONG"
	case WFM2_FOLDER_ENCRYPTION:
		return "WFM2_FOLDER_ENCRYPTION"
	case WFM2_PREPARE:
		return "WFM2_PREPARE"
	case WFM2_NO_SUPPORT_MEDIA:
		return "WFM2_NO_SUPPORT_MEDIA"
	case WFM2_DLNA_QDMS_DISABLE:
		return "WFM2_DLNA_QDMS_DISABLE"
	case WFM2_RENDER_NOT_FOUND:
		return "WFM2_RENDER_NOT_FOUND"
	case WFM2_CLOUD_SERVER_ERROR:
		return "WFM2_CLOUD_SERVER_ERROR"
	case WFM2_NAME_DUP:
		return "WFM2_NAME_DUP"
	case WFM2_EXCEED_SEARCH_MAX:
		return "WFM2_EXCEED_SEARCH_MAX"
	case WFM2_MEMORY_ERROR:
		return "WFM2_MEMORY_ERROR"
	case WFM2_COMPRESSING:
		return "WFM2_COMPRESSING"
	case WFM2_EXCEED_DAV_MAX:
		return "WFM2_EXCEED_DAV_MAX"
	case WFM2_UMOUNT_FAIL:
		return "WFM2_UMOUNT_FAIL"
	case WFM2_MOUNT_FAIL:
		return "WFM2_MOUNT_FAIL"
	case WFM2_WEBDAV_ACCOUNT_PASSWD_ERROR:
		return "WFM2_WEBDAV_ACCOUNT_PASSWD_ERROR"
	case WFM2_WEBDAV_SSL_ERROR:
		return "WFM2_WEBDAV_SSL_ERROR"
	case WFM2_WEBDAV_REMOUNT_ERROR:
		return "WFM2_WEBDAV_REMOUNT_ERROR"
	case WFM2_WEBDAV_HOST_ERROR:
		return "WFM2_WEBDAV_HOST_ERROR"
	case WFM2_WEBDAV_TIMEOUT_ERROR:
		return "WFM2_WEBDAV_TIMEOUT_ERROR"
	case WFM2_WEBDAV_CONF_ERROR:
		return "WFM2_WEBDAV_CONF_ERROR"
	case WFM2_WEBDAV_BASE_ERROR:
		return "WFM2_WEBDAV_BASE_ERROR"
	}

	return fmt.Sprintf("WMF2_UNKNOWN:%v", int(s))
}
