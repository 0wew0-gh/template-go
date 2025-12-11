package tool

import (
	"mime/multipart"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vod"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type UploadAuthDTO struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
}
type UploadAddressDTO struct {
	Endpoint string
	Bucket   string
	FileName string
}
type ProgressListener struct {
	Current           int64
	Progress          int64
	DoneProgress      int64
	Total             int64
	FileProgressToken string
	Hash              []string
	VideoID           []string
	Done              []bool
}

func (s Setting) InitVodClient(vodRegionID string, accessKeyId string, accessKeySecret string) (client *vod.Client, err error) {
	// 创建授权对象
	credential := &credentials.AccessKeyCredential{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	// 自定义config
	config := sdk.NewConfig()
	config.AutoRetry = true     // 失败是否自动重试
	config.MaxRetryTime = 3     // 最大重试次数
	config.Timeout = 3000000000 // 连接超时，单位：纳秒；默认为3秒
	// 创建vodClient实例
	return vod.NewClientWithOptions(vodRegionID, config, credential)
}

func MyCreateUploadVideo(client *vod.Client, title string, name string) (response *vod.CreateUploadVideoResponse, err error) {
	request := vod.CreateCreateUploadVideoRequest()
	request.Title = title
	request.FileName = name
	//request.CateId = "-1"
	//Cover URL示例：http://example.alicdn.com/tps/TB1qnJ1PVXXXXXCXXXXXXXXXXXX-700-****.png
	// request.CoverURL = "<your CoverURL>"
	request.AcceptFormat = "JSON"
	return client.CreateUploadVideo(request)
}

func MyGetPlayInfo(client *vod.Client, videoID string) (response *vod.GetPlayInfoResponse, err error) {
	request := vod.CreateGetPlayInfoRequest()
	request.VideoId = videoID
	request.AcceptFormat = "JSON"
	return client.GetPlayInfo(request)
}

func InitOssClient(uploadAuthDTO UploadAuthDTO, uploadAddressDTO UploadAddressDTO) (*oss.Client, error) {
	client, err := oss.New(uploadAddressDTO.Endpoint,
		uploadAuthDTO.AccessKeyId,
		uploadAuthDTO.AccessKeySecret,
		oss.SecurityToken(uploadAuthDTO.SecurityToken),
		oss.Timeout(86400*7, 86400*7))
	return client, err
}

func UploadLocalFile(client *oss.Client, uploadAddressDTO UploadAddressDTO, upf multipart.File, pl *ProgressListener) error {
	// 获取存储空间。
	bucket, err := client.Bucket(uploadAddressDTO.Bucket)
	if err != nil {
		return err
	}
	// 上传本地文件。
	return bucket.PutObject(uploadAddressDTO.FileName, upf, oss.Progress(pl))
}

func (p ProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		// fmt.Println("传输已启动")
	case oss.TransferDataEvent:
		// fmt.Println("传输中")
		// p.Current = event.ConsumedBytes
		// // 计算总进度
		// p.Progress = p.DoneProgress + p.Current
		// upMediaProgress[p.FileProgressToken].progress = p.Progress
	case oss.TransferCompletedEvent:
	case oss.TransferFailedEvent:
		// fmt.Println("传输失败")
	}
}
