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
	}
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": convert(faces)})
}

func convert(faces []db.UserFaces) map[string][]string {
	var male []string
	var female []string
	var f = make(map[string][]string)
	for _, v := range faces {
		if v.Gender == db.Male {
			male = append(male, v.FaceURL)
		} else if v.Gender == db.Female {
			female = append(female, v.FaceURL)
		}
	}
	f["male"] = male
	f["female"] = female
	return f
}
