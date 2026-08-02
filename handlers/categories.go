package handlers

import (
	"log"
	"net/http"

	"real-time-forum/db"
	"real-time-forum/utils"
)

type categoryItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	SubGroup *string `json:"sub_group"`
}

type categoryGroupResponse struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Categories []categoryItem `json:"categories"`
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	rows, err := db.DB.Query(`
		SELECT
			g.id,
			g.name,
			c.id,
			c.name,
			c.sub_group
		FROM category_groups g
		LEFT JOIN categories c ON c.group_id = g.id
		ORDER BY g.id ASC, c.id ASC
	`)
	if err != nil {
		log.Printf("[ERROR] [GetCategories Query]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not fetch categories")
		return
	}
	defer rows.Close()

	groups := []categoryGroupResponse{}
	groupIndex := make(map[string]int)

	for rows.Next() {
		var groupID, groupName string
		var catID, catName *string
		var subGroup *string

		if err := rows.Scan(&groupID, &groupName, &catID, &catName, &subGroup); err != nil {
			log.Printf("[ERROR] [GetCategories Scan]: %v", err)
			utils.WriteError(w, http.StatusInternalServerError, "error processing categories")
			return
		}

		if _, exists := groupIndex[groupID]; !exists {
			groupIndex[groupID] = len(groups)
			groups = append(groups, categoryGroupResponse{
				ID:         groupID,
				Name:       groupName,
				Categories: []categoryItem{},
			})
		}

		if catID != nil && catName != nil {
			idx := groupIndex[groupID]
			groups[idx].Categories = append(groups[idx].Categories, categoryItem{
				ID:       *catID,
				Name:     *catName,
				SubGroup: subGroup,
			})
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] [GetCategories Rows Loop]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "error reading categories")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, groups); err != nil {
		log.Printf("[ERROR] [GetCategories Response JSON]: %v", err)
	}
}
