package controller

import (
	"fmt"
	"ginEssential/config"
	"ginEssential/model"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func AuthTemp(ctx *gin.Context) {
	//构建一个阿里云客户端, 用于发起请求。
	//设置调用者（RAM用户或RAM角色）的AccessKey ID和AccessKey Secret。
	//第一个参数就是bucket所在位置，可查看oss对象储存控制台的概况获取
	//第二个参数就是步骤一获取的AccessKey ID
	//第三个参数就是步骤一获取的AccessKey Secret
	fmt.Printf("config.Instance.Uploader.AliyunOss.Endpoint is %#v\n", config.Instance.Uploader.AliyunOss.Endpoint)
	client, err := sts.NewClientWithAccessKey(config.Instance.Uploader.AliyunOss.RegionId,
		config.Instance.Uploader.AliyunOss.AccessId, config.Instance.Uploader.AliyunOss.AccessSecret)

	//构建请求对象。
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	//设置参数。关于参数含义和设置方法，请参见《API参考》。
	request.RoleArn = config.Instance.Uploader.AliyunOss.RoleArn                 //步骤三获取的角色ARN
	request.RoleSessionName = config.Instance.Uploader.AliyunOss.RoleSessionName //步骤三中的RAM角色名称

	//发起请求，并得到响应。
	response, err := client.AssumeRole(request)
	if err != nil {
		fmt.Print(err.Error())
	}

	// 创建用户
	tempOssVO := model.TempOssVO{
		Endpoint:        config.Instance.Uploader.AliyunOss.Endpoint,
		AccessKeyId:     response.Credentials.AccessKeyId,
		AccessKeySecret: response.Credentials.AccessKeySecret,
		BucketName:      config.Instance.Uploader.AliyunOss.Bucket,
		UploadUrl:       config.Instance.Uploader.AliyunOss.RemotePath,
		Token:           response.Credentials.SecurityToken,
	}

	fmt.Printf("response is %#v\n", response)

	model.Success(ctx, tempOssVO, "请求成功")
}

func Copy(ctx *gin.Context) {

	var fileTemp2FormalDTO = model.FileTemp2FormalDTO{
		RelativePath: ctx.Request.FormValue("relativePath"),
	}

	// 创建OSSClient实例。
	client, err := oss.New(config.Instance.Uploader.AliyunOss.Endpoint,
		config.Instance.Uploader.AliyunOss.AccessId, config.Instance.Uploader.AliyunOss.AccessSecret)
	if err != nil {
		fmt.Println("Error:", err)
	}

	srcBucketName := config.Instance.Uploader.AliyunOss.Bucket
	destinationBucketName := config.Instance.Uploader.AliyunOss.RealBucket

	var fileName = fileTemp2FormalDTO.RelativePath

	destinationObjectName := config.Instance.Uploader.AliyunOss.RealPath + fileName[strings.LastIndex(fileName, "/"):len(fileName)]
	fileName = fileName[1:len(fileName)]

	// 获取存储空间。
	bucket, err := client.Bucket(destinationBucketName)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 获取存储空间。
	bucketSrc, err := client.Bucket(srcBucketName)
	if err != nil {
		fmt.Println("Error:", err)
	}

	meta, err := bucketSrc.GetObjectMeta(fileName)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("meta:", meta)
	contentLength, _ := strconv.Atoi(meta.Get("Content-Length"))

	// 设置分片大小为50MB。
	var partSize = 1024 * 1024 * 50

	// 计算分片总数。
	var partCount = (contentLength / partSize)
	if contentLength%partSize != 0 {
		partCount++
	}

	// 步骤1：初始化一个分片上传事件。
	imur, err := bucket.InitiateMultipartUpload(destinationObjectName)
	// 步骤2：上传分片。
	var parts []oss.UploadPart
	for i := 0; i < partCount; i++ {
		var skipBytes = partSize * i
		// 对每个分片调用UploadPart方法上传。
		part, err := bucket.UploadPartCopy(imur, config.Instance.Uploader.AliyunOss.Bucket, fileName, int64(skipBytes), int64(partSize), i+1)
		if err != nil {
			fmt.Println("Error:", err)
		}
		parts = append(parts, part)
	}
	// 步骤3：完成分片上传。
	cmur, err := bucket.CompleteMultipartUpload(imur, parts)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("cmur:", cmur)

	model.Success(ctx, "ok", "请求成功")
}
