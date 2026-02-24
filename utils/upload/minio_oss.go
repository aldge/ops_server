package upload

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

var MinioClient *Minio // 优化性能，但是不支持动态配置

type Minio struct {
	Client *minio.Client
	bucket string
}

func GetMinio(endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*Minio, error) {
	if MinioClient != nil {
		return MinioClient, nil
	}
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL, // Set to true if using https
	})
	if err != nil {
		return nil, err
	}
	// 尝试创建bucket
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			// log.Printf("We already own %s\n", bucketName)
		} else {
			return nil, err
		}
	}
	MinioClient = &Minio{Client: minioClient, bucket: bucketName}
	return MinioClient, nil
}

func (m *Minio) UploadFile(file *multipart.FileHeader) (filePathres, key string, uploadErr error) {
	f, openError := file.Open()
	// mutipart.File to os.File
	if openError != nil {
		global.GVA_LOG.Error("function file.Open() Failed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Failed, err:" + openError.Error())
	}

	filecontent := bytes.Buffer{}
	_, err := io.Copy(&filecontent, f)
	if err != nil {
		global.GVA_LOG.Error("读取文件失败", zap.Any("err", err.Error()))
		return "", "", errors.New("读取文件失败, err:" + err.Error())
	}
	f.Close() // 创建文件 defer 关闭

	// 对文件名进行加密存储
	ext := filepath.Ext(file.Filename)
	filename := utils.MD5V([]byte(strings.TrimSuffix(file.Filename, ext))) + ext
	if global.GVA_CONFIG.Minio.BasePath == "" {
		filePathres = "uploads" + "/" + time.Now().Format("2006-01-02") + "/" + filename
	} else {
		filePathres = global.GVA_CONFIG.Minio.BasePath + "/" + time.Now().Format("2006-01-02") + "/" + filename
	}

	// 根据文件扩展名检测 MIME 类型
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 设置超时10分钟
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	// Upload the file with PutObject   大文件自动切换为分片上传
	info, err := m.Client.PutObject(ctx, global.GVA_CONFIG.Minio.BucketName, filePathres, &filecontent, file.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		global.GVA_LOG.Error("上传文件到minio失败", zap.Any("err", err.Error()))
		return "", "", errors.New("上传文件到minio失败, err:" + err.Error())
	}
	return global.GVA_CONFIG.Minio.BucketUrl + "/" + info.Key, filePathres, nil
}

func (m *Minio) DeleteFile(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Delete the object from MinIO
	err := m.Client.RemoveObject(ctx, m.bucket, key, minio.RemoveObjectOptions{})
	return err
}

// UploadFileFromPath 从本地文件路径上传文件到 Minio
// @param: localFilePath string, fileName string
// @return: filePathres string, key string, uploadErr error
func (m *Minio) UploadFileFromPath(localFilePath string, fileName string, videoID string) (filePathres, key string, uploadErr error) {
	// 打开本地文件
	file, err := os.Open(localFilePath)
	if err != nil {
		global.GVA_LOG.Error("打开本地文件失败", zap.Any("err", err.Error()))
		return "", "", errors.New("打开本地文件失败, err:" + err.Error())
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		global.GVA_LOG.Error("获取文件信息失败", zap.Any("err", err.Error()))
		return "", "", errors.New("获取文件信息失败, err:" + err.Error())
	}

	// 读取文件内容
	filecontent := bytes.Buffer{}
	_, err = io.Copy(&filecontent, file)
	if err != nil {
		global.GVA_LOG.Error("读取文件失败", zap.Any("err", err.Error()))
		return "", "", errors.New("读取文件失败, err:" + err.Error())
	}

	// 对文件名进行加密存储
	ext := filepath.Ext(fileName)
	if videoID != "" {
		// 当 videoID 不为空时，使用特殊路径规则
		videoIDLower := strings.ToLower(videoID)
		// 取前 4 个字符作为两个文件夹名（每个文件夹 2 个字符）
		if len(videoIDLower) >= 4 {
			folder1 := videoIDLower[0:2]
			folder2 := videoIDLower[2:4]
			// 使用原始文件名（不加密）
			if global.GVA_CONFIG.Minio.BasePath == "" {
				filePathres = folder1 + "/" + folder2 + "/" + videoID
			} else {
				filePathres = global.GVA_CONFIG.Minio.BasePath + "/" + folder1 + "/" + folder2 + "/" + videoID
			}
		} else {
			// videoID 长度不足 4 位，使用原始逻辑
			filename := utils.MD5V([]byte(strings.TrimSuffix(fileName, ext))) + ext
			if global.GVA_CONFIG.Minio.BasePath == "" {
				filePathres = "uploads" + "/" + time.Now().Format("2006-01-02") + "/" + filename
			} else {
				filePathres = global.GVA_CONFIG.Minio.BasePath + "/" + time.Now().Format("2006-01-02") + "/" + filename
			}
		}
	} else {
		// videoID 为空，使用原始逻辑
		filename := utils.MD5V([]byte(strings.TrimSuffix(fileName, ext))) + ext
		if global.GVA_CONFIG.Minio.BasePath == "" {
			filePathres = "uploads" + "/" + time.Now().Format("2006-01-02") + "/" + filename
		} else {
			filePathres = global.GVA_CONFIG.Minio.BasePath + "/" + time.Now().Format("2006-01-02") + "/" + filename
		}
	}

	// 根据文件扩展名检测 MIME 类型
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 设置超时10分钟
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	// Upload the file with PutObject   大文件自动切换为分片上传
	info, err := m.Client.PutObject(ctx, global.GVA_CONFIG.Minio.BucketName, filePathres, &filecontent, fileInfo.Size(), minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		global.GVA_LOG.Error("上传文件到minio失败", zap.Any("err", err.Error()))
		return "", "", errors.New("上传文件到minio失败, err:" + err.Error())
	}
	return global.GVA_CONFIG.Minio.BucketUrl + "/" + info.Key, filePathres, nil
}
