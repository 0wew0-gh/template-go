package tool

import (
	"fmt"
	"strconv"
	"sync"
)

type OssConfigSync struct {
	m  sync.RWMutex
	si map[int]OssConfig
}
type OssConfig struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	Bucket          string `json:"bucket"`
	Endpoint        string `json:"endpoint"`
	VodRegionID     string `json:"vodRegionID"`
	AccessKeyID     string
	AccessKeySecret string
}

func NewOssConfigSync(si *OssConfigSync) *OssConfigSync {
	if si != nil {
		if len(si.si) == 0 {
			si.si = map[int]OssConfig{}
		}
		return si
	}
	return &OssConfigSync{
		si: map[int]OssConfig{},
	}
}

func (s *OssConfigSync) GetOssConfig() map[int]OssConfig {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.si
}

func (s *OssConfigSync) SetOssConfig(si map[int]OssConfig) {
	s.m.Lock()
	defer s.m.Unlock()
	s.si = si
}

func (s Setting) GetOssConfig() error {
	qd, err := s.SQL.QueryData("*", "site", "", "`id` ASC", "99999")
	if err != nil {
		return nil
	}
	if s.OssConfig == nil {
		s.OssConfig = NewOssConfigSync(s.OssConfig)
	}

	siteMap := s.OssConfig.GetOssConfig()
	for i := 0; i < len(qd); i++ {
		strI := strconv.Itoa(i)
		item := qd[strI]

		id := 0
		tempSite := OssConfig{}
		for k, v := range item {
			switch k {
			case "id":
				tempInt, err := strconv.Atoi(v)
				if err != nil {
					continue
				}
				id = tempInt
				tempSite.ID = tempInt
			case "name":
				tempSite.Name = v
			case "url":
				tempSite.URL = v
			case "bucket":
				tempSite.Bucket = v
				for _, ossv := range s.Config.OSS {
					if v == ossv.Bucket {
						tempSite.AccessKeyID = ossv.AccessKeyID
						tempSite.AccessKeySecret = ossv.AccessKeySecret
						break
					}
				}
			case "endpoint":
				tempSite.Endpoint = v
			case "vodRegionID":
				tempSite.VodRegionID = v
			}
		}
		if id == 0 {
			continue
		}
		siteMap[id] = tempSite
	}
	s.OssConfig.SetOssConfig(siteMap)
	return nil
}

// 根据字符串ID获取站点信息
func (s Setting) GetOssConfigForIDstr(idstr string) (OssConfig, error) {
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return OssConfig{}, err
	}
	return s.GetOssConfigForID(id)
}

// 根据ID获取站点信息
func (s Setting) GetOssConfigForID(id int) (OssConfig, error) {
	siteMap := s.OssConfig.GetOssConfig()
	ossConf, ok := siteMap[id]
	if !ok {
		err := s.GetOssConfig()
		if err != nil {
			return OssConfig{}, err
		}
		siteMap = s.OssConfig.GetOssConfig()
		ossConf, ok = siteMap[id]
		if !ok {
			return OssConfig{}, fmt.Errorf("site not found")
		}
	}
	return ossConf, nil
}
