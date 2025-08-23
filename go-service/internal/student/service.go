package student

import (
	"context"
	"fmt"
	"goservice/internal/client"
	"goservice/internal/models"
	"io"
	"net/http"

	"github.com/jung-kurt/gofpdf"
)

type Service interface {
	GetStudent(ctx context.Context, id int, authCookies []*http.Cookie) (*models.Student, error)
	GenerateReport(ctx context.Context, id int, authCookies []*http.Cookie) (ReportWriter, error)
	Login(ctx context.Context, username, password string) ([]*http.Cookie, error)
}

type service struct {
	backend client.IBackend
}

type ReportWriter interface {
	Output(w io.Writer) error
}

func NewService(b client.IBackend) Service {
	return &service{backend: b}
}

func (s *service) Login(ctx context.Context, username, password string) ([]*http.Cookie, error) {
	return s.backend.Login(ctx, username, password)
}

func (s *service) GetStudent(ctx context.Context, id int, authCookies []*http.Cookie) (*models.Student, error) {
	return s.backend.GetStudentByID(ctx, id, authCookies)
}

func (s *service) GenerateReport(ctx context.Context, id int, authCookies []*http.Cookie) (ReportWriter, error) {
	student, err := s.GetStudent(ctx, id, authCookies)
	if err != nil {
		return nil, err
	}

	return generatePDF(student), nil
}

func generatePDF(student *models.Student) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Student Report")
	pdf.Ln(15)

	// Set Font for content
	pdf.SetFont("Arial", "", 12)

	// Add Student Data
	addLine := func(label, value string) {
		pdf.CellFormat(50, 8, label, "0", 0, "", false, 0, "")
		pdf.CellFormat(100, 8, value, "0", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	addLine("ID:", fmt.Sprintf("%d", student.ID))
	addLine("Name:", student.Name)
	addLine("Email:", student.Email)
	addLine("System Access:", boolToString(student.SystemAccess))
	addLine("Phone:", student.Phone)
	addLine("Gender:", student.Gender)
	addLine("DOB:", student.DOB.Format("2006-01-02"))
	addLine("Class:", student.Class)
	addLine("Section:", student.Section)
	addLine("Roll:", fmt.Sprintf("%d", student.Roll))
	addLine("Father Name:", student.FatherName)
	addLine("Father Phone:", student.FatherPhone)
	addLine("Mother Name:", student.MotherName)
	addLine("Mother Phone:", student.MotherPhone)
	addLine("Guardian Name:", student.GuardianName)
	addLine("Guardian Phone:", student.GuardianPhone)
	addLine("Relation Of Guardian:", student.RelationOfGuardian)
	addLine("Current Address:", student.CurrentAddress)
	addLine("Permanent Address:", student.PermanentAddress)
	addLine("Admission Date:", student.AdmissionDate.Format("2006-01-02"))
	addLine("Reporter Name:", student.ReporterName)

	return pdf
}

func boolToString(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}
