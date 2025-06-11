package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

func Connect() *pgx.Conn {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error loading .env file")
	}
	var user string = os.Getenv("DB_username")
	var password string = os.Getenv("DB_password")
	var address string = os.Getenv("DB_address")
	var db string = os.Getenv("DB_db")
	var port string = os.Getenv("DB_port")
	connect_string := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, password, address, port, db)
	connect, err := pgx.Connect(context.Background(), connect_string)

	if err != nil {
		fmt.Printf("could not connect to databse as %v", err)
	}
	return connect
}

func InsertUser(email, password string) error {
	HashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hasing as %v", err)
	}

	connect := Connect()
	var id int

	_, err = connect.Exec(context.Background(),
		"INSERT INTO logininfo VALUES (DEFAULT, $1, $2)", email, string(HashedPass))

	rows := connect.QueryRow(context.Background(), "SELECT id FROM logininfo WHERE email=$1", email)
	rows.Scan(&id)

	if err != nil {
		return fmt.Errorf("probably account already exists as %v", err)
	} else {
		_, _ = connect.Exec(context.Background(), "INSERT INTO userstatus VALUES($1, DEFAULT, DEFAULT)", id)
		return nil
	}
}

func VerifyUser(email, password string) (error, interface{}) {
	connect := Connect()
	defer connect.Close(context.Background())
	var id int

	rows := connect.QueryRow(context.Background(), "SELECT password FROM logininfo WHERE email=$1", email)

	var HashedPass string

	err := rows.Scan(&HashedPass)
	if err != nil {
		return fmt.Errorf("no_user\n"), nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(HashedPass), []byte(password))
	if err != nil {
		return fmt.Errorf("wrong_pass\n"), nil
	}
	rowsNew := connect.QueryRow(context.Background(), "SELECT id FROM logininfo WHERE email=$1", email)
	rowsNew.Scan(&id)

	return nil, id
}

func InsertClient(client Client) error {
	connect := Connect()

	_, err := connect.Exec(
		context.Background(),
		`INSERT INTO client (
            name, phone, address, message, callweek,
            highvalue, bookingdone, bookingdate, signed,
            signedprice, workdate, leaduser
        ) VALUES (
            $1, $2, $3, $4, $5,
            $6, $7, $8, $9,
            $10, $11, $12
        )`,
		client.Name,
		client.Phone,
		client.Address,
		client.Message,
		client.CallWeek,
		client.HighValue,
		client.BookingDone,
		client.BookingDate,
		client.Signed,
		client.SignedPrice,
		client.WorkDate,
		client.LeadUser,
	)
	if err != nil {
		fmt.Printf("ran into error as %v", err)
	}

	err = SaveToGoogleSheets(client)
	if err != nil {
		fmt.Printf("ran into sheets copy error as %v", err)
	}

	return nil
}

func GetEveryClientDB() ([]ClientAllInfo, error) {
	var losClient []ClientAllInfo
	conn := Connect()

	rows, err := conn.Query(context.Background(), "SELECT * FROM client")
	if err != nil {
		fmt.Printf("Could not get all clients from database: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var LocalData ClientAllInfo

		err := rows.Scan(
			&LocalData.Id,
			&LocalData.Name,
			&LocalData.Phone,
			&LocalData.Address,
			&LocalData.Message,
			&LocalData.CallWeek,
			&LocalData.HighValue,
			&LocalData.BookingDone,
			&LocalData.BookingDate,
			&LocalData.Signed,
			&LocalData.SignedPrice,
			&LocalData.WorkDate,
			&LocalData.LeadUser,
		)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			return nil, err
		}

		losClient = append(losClient, LocalData)
	}
	return losClient, nil
}

func FetchUserData(id int) (UserInformation, error) {

	var userinfo UserInformation
	conn := Connect()
	rows := conn.QueryRow(context.Background(), "SELECT logininfo.email, logininfo.id,userstatus.isadmin,userstatus.cut FROM logininfo JOIN userstatus ON logininfo.id = userstatus.id WHERE logininfo.id=$1;", id)
	if err := rows.Scan(&userinfo.Email, &userinfo.Id, &userinfo.Isadmin, &userinfo.Cut); err != nil {
		return userinfo, err
	}
	return userinfo, nil
}

func FetchCalendarInfo() ([]CalenderInfo, error) {
	conn := Connect()
	var CalendarData []CalenderInfo

	if rows, err := conn.Query(context.Background(), "SELECT id, starttime, endtime, clientrelation, message FROM calendarinfo"); err != nil {
		return nil, fmt.Errorf("could not fetch, database err as %v", err)
	} else {

		for rows.Next() {
			var LocalData CalenderInfo
			err := rows.Scan(&LocalData.Id, &LocalData.StartTime, &LocalData.EndTime, &LocalData.ClientRelation, &LocalData.Message)
			if err != nil {
				fmt.Printf("Error scanning row: %v\n", err)
				continue
			}
			CalendarData = append(CalendarData, LocalData)
		}
	}
	return CalendarData, nil
}

func AddACalendarInfo(data CalenderInfoReceived) error {
	conn := Connect()

	fmt.Println(data.StartTime)

	if _, err := conn.Exec(context.Background(),
		"INSERT INTO calendarinfo(starttime, endtime, clientrelation, message) VALUES($1, $2, $3, $4)",
		data.StartTime, data.EndTime, data.ClientRelation, data.Message); err != nil {
		return fmt.Errorf("could not insert data, database error: %v", err)
	} else {
		return nil
	}
}

func DeleteACalendarElementDB(id int) error {
	conn := Connect()

	if _, err := conn.Exec(context.Background(), "DELETE FROM calendarinfo WHERE id = $1", id); err != nil {
		fmt.Errorf("error db -", err)
	}
	return nil
}

func ReturnDataSingleClient(id int) (ClientAllInfo, error) {
	conn := Connect()

	row := conn.QueryRow(context.Background(), "SELECT * FROM client WHERE id = $1", id)

	var LocalData ClientAllInfo

	if err := row.Scan(
		&LocalData.Id,
		&LocalData.Name,
		&LocalData.Phone,
		&LocalData.Address,
		&LocalData.Message,
		&LocalData.CallWeek,
		&LocalData.HighValue,
		&LocalData.BookingDone,
		&LocalData.BookingDate,
		&LocalData.Signed,
		&LocalData.SignedPrice,
		&LocalData.WorkDate,
		&LocalData.LeadUser,
	); err != nil {
		return LocalData, err
	} else {
		return LocalData, nil
	}
}

func FetchAllNotesDB() ([]NotesDetail, error) {
	conn := Connect()
	var Data []NotesDetail

	if rows, err := conn.Query(context.Background(), "SELECT id, title, detail, clientrelation FROM notes;"); err != nil {
		return nil, err
	} else {
		for rows.Next() {
			var LocalData NotesDetail
			err := rows.Scan(&LocalData.Id, &LocalData.Title, &LocalData.Detail, &LocalData.ClientRelation)
			if err != nil {
				fmt.Printf("Error scanning row: %v\n", err)
				continue
			}
			Data = append(Data, LocalData)
		}

	}
	return Data, nil
}

func DeleteASpecificNoteDB(id int) error {
	conn := Connect()

	if _, err := conn.Exec(context.Background(), "DELETE FROM notes WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}

func UploadANoteDB(note NotesDetailSender, hasClientRelation bool) error {
	conn := Connect()

	if hasClientRelation {

		if _, err := conn.Exec(context.Background(), "INSERT INTO notes(title,detail,clientrelation) VALUES($1,$2,$3);", note.Title, note.Detail, note.ClientRelation); err != nil {
			return err
		}
	} else {
		if _, err := conn.Exec(context.Background(), "INSERT INTO notes(title,detail) VALUES($1,$2);", note.Title, note.Detail); err != nil {
			return err
		}

	}
	return nil
}

func UpdateANoteDB(note NotesDetailUpdater) error {
	conn := Connect()
	defer conn.Close(context.Background())

	if note.ClientRelation == 0 { //no client relation
		_, err := conn.Exec(
			context.Background(),
			"UPDATE notes SET title=$1, detail=$2,clientrelation=$3 WHERE id=$4",
			note.Title, note.Detail, nil, note.Id,
		)
		if err != nil {
			return err
		}
		return nil

	} else {
		_, err := conn.Exec(
			context.Background(),
			"UPDATE notes SET title=$1, detail=$2, clientrelation=$3 WHERE id=$4",
			note.Title, note.Detail, note.ClientRelation, note.Id,
		)
		if err != nil {
			return err
		}
		return nil

	}

}

func FetchAllFinanceDetailDB() ([]FinanceDetailFull, error) {
	conn := Connect()
	defer conn.Close(context.Background())
	var allDetails []FinanceDetailFull

	if rows, err := conn.Query(context.Background(), "SELECT id,userid,cost,context,hasbeentaken FROM financedetail"); err != nil {
		return nil, err
	} else {
		for rows.Next() {
			var LocalData FinanceDetailFull
			rows.Scan(&LocalData.Id, &LocalData.UserID, &LocalData.Cost, &LocalData.Context, &LocalData.HasBeenTaken)
			allDetails = append(allDetails, LocalData)
		}
	}
	return allDetails, nil
}

func InsertAFinanceDetailDB(data FinanceDetailSend) error {
	conn := Connect()
	defer conn.Close(context.Background())

	if _, err := conn.Exec(context.Background(), "INSERT INTO financedetail (userid,cost,context) VALUES($1,$2,$3)", data.UserID, data.Cost, data.Context); err != nil {
		return err
	}
	return nil
}

func DeleteAFinanceDetailDB(id int) error {
	conn := Connect()
	defer conn.Close(context.Background())

	if _, err := conn.Exec(context.Background(), "DELETE FROM financedetail WHERE id = $1", id); err != nil {
		return err
	}
	return nil
}

func UpdateUtilizedFinanceDetailDB(ids []int) error {
	conn := Connect()
	defer conn.Close(context.Background())
	var deleted_items = 0

	for _, id := range ids {
		if _, err := conn.Exec(context.Background(), "UPDATE financedetail SET hasbeentaken = true WHERE id=$1;", id); err != nil {
			continue
		} else {
			deleted_items += 1
		}
	}

	if deleted_items == len(ids) {
		return nil
	} else {
		return fmt.Errorf("did not delete everything")
	}

}

func UploadAnImageDB(CloudPath, title, description string) error {
	conn := Connect()
	location, _ := time.LoadLocation("America/New_York")
	time_rn := time.Now().In(location)

	_, err := conn.Exec(context.Background(), "INSERT INTO ImageData(url,datetime,title,description) VALUES($1,$2,$3,$4)", CloudPath, time_rn, title, description)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func FetchAllImageDataDB() ([]ImageDataReceived, error) {
	var data []ImageDataReceived
	conn := Connect()
	rows, err := conn.Query(context.Background(), "SELECT id,url,datetime,title,description FROM ImageData")
	if err != nil {
		return data, err
	}
	for rows.Next() {
		var localData ImageDataReceived
		rows.Scan(&localData.Id, &localData.Url, &localData.Datetime, &localData.Title, &localData.Description)
		data = append(data, localData)
	}
	return data, nil
}

func DeleteAImageDataDB(id int) error {
	conn := Connect()

	if _, err := conn.Exec(context.Background(), "DELETE FROM ImageData WHERE id=$1", id); err != nil {
		return err
	}
	return nil
}

func InsertAAuthTokenDB(userid int) (AuthorizationDetailsFetch, error) {
	conn := Connect()
	var data2beSent AuthorizationDetailsFetch

	var currentTime = time.Now()
	var latertime = currentTime.Add(14 * 24 * time.Hour)
	var verification, _ = password.Generate(64, 10, 12, false, true)

	if _, err := conn.Exec(context.Background(), "INSERT INTO authverification(starttime,endtime,verification,userid) VALUES($1,$2,$3,$4)", currentTime, latertime, verification, userid); err != nil {
		return AuthorizationDetailsFetch{}, err
	} else {
		row := conn.QueryRow(context.Background(), "SELECT id, starttime,endtime,verification,userid FROM authverification WHERE verification = $1", verification)
		row.Scan(&data2beSent.Id, &data2beSent.StartTime, &data2beSent.EndTime, &data2beSent.Verification, &data2beSent.UserID)
		return data2beSent, nil

	}

}
func IsTokenValidDB(verification string, userid int) error {
	conn := Connect()

	var data AuthorizationDetailsFetch

	err := conn.QueryRow(context.Background(),
		"SELECT id, starttime, endtime, verification, userid FROM authverification WHERE userid=$1 AND verification=$2",
		userid, verification).Scan(
		&data.Id, &data.StartTime, &data.EndTime, &data.Verification, &data.UserID)

	if err == pgx.ErrNoRows {
		return fmt.Errorf("NoToken")
	} else if err != nil {
		CurrentTime := time.Now()
		if CurrentTime.After(data.StartTime) && CurrentTime.Before(data.EndTime) {
			return nil
		} else {
			return fmt.Errorf("TokenExperied")
		}
	}
	return nil
}
