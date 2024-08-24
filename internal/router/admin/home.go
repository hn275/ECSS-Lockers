package admin

import (
	"encoding/hex"
	"net/http"
	stdtime "time"

	"github.com/zvdv/ECSS-Lockers/internal/database"
	"github.com/zvdv/ECSS-Lockers/internal/httputil"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
	"github.com/zvdv/ECSS-Lockers/internal/time"
	"lukechampine.com/blake3"
)

func Home(w http.ResponseWriter, r *http.Request) {
	data := struct {
		HasData       bool
		Registrations []registration
	}{}

	var err error
	data.Registrations, err = queryAllRegistrations()
	data.HasData = len(data.Registrations) != 0

	if err != nil {
		logger.Error.Println(err)
		httputil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	httputil.WriteTemplatePage(
		w,
		data,
		"templates/nav.html",
		"templates/admin/index.html",
		"templates/admin/lockertable.html")
}

type registration struct {
	ID        string
	RowIndex  uint16
	Locker    string
	Name      string
	Email     string
	Expiry    string
	EmailSent bool
}

func queryAllRegistrations() ([]registration, error) {
	db, lock := database.Lock()
	defer lock.Unlock()

	stmt, err := db.Prepare(`
        SELECT locker, user, name, expiry, expiryEmailSent
        FROM registration;`)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	lockers := make([]registration, 0, 200)
	rowIndex := uint16(1)

	var expiredAt stdtime.Time
	for ; rows.Next(); rowIndex++ {
		reg := registration{
			RowIndex: rowIndex,
		}

		err := rows.Scan(
			&reg.Locker, &reg.Email, &reg.Name,
			&expiredAt, &reg.EmailSent)

		if err != nil {
			return nil, err
		}

		reg.Expiry = expiredAt.Format(time.TimeFormatLayout)

		// generate an id, for UI only
		idHash := blake3.New(16, nil)
		if _, err := idHash.Write([]byte(reg.Locker)); err != nil {
			return nil, err
		}

		reg.ID = hex.EncodeToString(idHash.Sum(nil))

		lockers = append(lockers, reg)
	}

	return lockers, nil
}
