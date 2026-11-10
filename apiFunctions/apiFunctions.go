package apiFunctions

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const OSU_SBGI = "003210"

var db *sql.DB

type Course struct {
	CourseCode string `json:"course_code"`
	Title      string `json:"title"`
}

type DetailedCourse struct {
	Course
	SgbiCode string `json:"sgbi_code"`
	SubjCode string `json:"subj_code"`
	School   string `json:"school"`
}

type SearchResult struct {
	Course
	SubjCode string `json:"subj_code"`
}

type Courses []Course
type DetailedCourses []DetailedCourse
type SearchResults []SearchResult

func InitDb() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	database := os.Getenv("DATABASE")
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
	db, err = sql.Open("mysql", dataSource)
	if err != nil {
		return err
	}
	return db.Ping()
}

func GetTransferCourses(subj, code string) (DetailedCourses, error) {
	query := `
        SELECT c.course_code,
               c.subj_code,
               c.sgbi_code,
               c.title,
               s.school
        FROM courses c
        JOIN schools s ON s.sgbi_code = c.sgbi_code
        WHERE c.osu_eqv_sgbi_code   = ?
          AND c.osu_eqv_subj_code   = ?
          AND c.osu_eqv_course_code = ?
          AND c.sgbi_code <> ?
    `
	data := DetailedCourses{}                                    // Will hold returned rows data
	rows, err := db.Query(query, OSU_SBGI, subj, code, OSU_SBGI) // query databse
	if err != nil {                                              // error handle query
		return nil, err
	}
	defer rows.Close() // close rows after done with function
	var course_code, subj_code, sgbi_code, title, school string
	for rows.Next() { // extract data from rows
		err := rows.Scan(&course_code, &subj_code, &sgbi_code, &title, &school)
		if err != nil { // error handle rows
			return nil, err
		}
		// add course to result
		data = append(data, DetailedCourse{Course{course_code, title}, sgbi_code, subj_code, school})
	}
	if err := rows.Err(); err != nil { // Final error check
		return nil, err
	}
	return data, nil
}

func GetCourses(sgbi, subj string) (Courses, error) {
	query := `
		SELECT course_code,
			   title
		FROM COURSES
		WHERE sgbi_code = ?
		AND subj_code = ?
	`
	data := Courses{}
	rows, err := db.Query(query, sgbi, subj)
	if err != nil {
		return nil, err
	}
	var course_code, title string
	for rows.Next() {
		err := rows.Scan(&course_code, &title)
		if err != nil {
			return nil, err
		}
		data = append(data, Course{course_code, title})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func GetSubjects(sgbi string) ([]string, error) {
	query := `
		SELECT subj_code
		FROM subjects
		WHERE sgbi_code = ?`
	var data []string
	rows, err := db.Query(query, sgbi)
	if err != nil {
		return nil, err
	}
	var subject string
	for rows.Next() {
		err := rows.Scan(&subject)
		if err != nil {
			return nil, err
		}
		data = append(data, subject)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func SearchCourses(sgbi, search string) (SearchResults, error) {
	query := `
		SELECT course_code,
		subj_code,
		title
		FROM courses 
		WHERE sgbi_code = ?
		AND title REGEXP ?
		ORDER BY title ASC
		LIMIT 25;
	`
	var data SearchResults
	rows, err := db.Query(query, sgbi, "^"+search)
	if err != nil {
		return nil, err
	}
	var course_code, subj_code, title string
	for rows.Next() {
		err := rows.Scan(&course_code, &subj_code, &title)
		if err != nil {
			return nil, err
		}
		data = append(data, SearchResult{Course{course_code, title}, subj_code})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
}
