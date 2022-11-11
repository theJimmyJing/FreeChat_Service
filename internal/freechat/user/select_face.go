package freechat_user

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	fc_mysql_model "Open_IM/pkg/common/db/mysql_model/freechat_mysql_model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SelectFace(c *gin.Context) {
	faces, err := fc_mysql_model.GetFacesURL()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrDB, "errMsg": "internal err"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": convert(faces)})
}

func convert(faces []db.UserFaces) map[string][]string {
	var large []string
	var small []string
	var f = make(map[string][]string)
	for _, v := range faces {
		if v.FaceType == db.Large {
			large = append(large, v.FaceURL)
		} else if v.FaceType == db.Small {
			small = append(small, v.FaceURL)
		}
	}
	f["large"] = large
	f["small"] = small
	return f
}
