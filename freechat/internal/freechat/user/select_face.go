package freechat_user

import (
	"Open_IM/pkg/common/constant"
	"encoding/json"
	"freechat/pkg/common/db"
	fc_mysql_model "freechat/pkg/common/db/mysql_model/freechat_mysql_model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FaceURL struct {
	ID    int    `json:"id"`
	Large string `json:"large"`
	Small string `json:"small"`
}

func SelectFace(c *gin.Context) {
	faces, err := fc_mysql_model.GetFacesURL()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrDB, "errMsg": "internal err"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": convert(faces)})
}

func convert(faces []db.UserFaces) []FaceURL {
	var f = []FaceURL{}
	for _, v := range faces {
		var u = FaceURL{}
		err := json.Unmarshal([]byte(v.FaceURL), &u)
		if err != nil {
			return f
		}
		u.ID = v.ID
		f = append(f, u)
	}
	return f
}
