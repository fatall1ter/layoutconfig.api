package infra

import (
	"net/http"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
)

// apiReports docs
// @Summary Gets early created reports
// @Description Gets early created pdf/xlsx reports
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description is_sent filter for reports, 1/true = report has been sent
// @Produce  json
// @Tags reports
// @Param layout_id query string false "digit/uuid format"
// @Param is_sent query string false "boolean format false/true/0/1"
// @Success 200 {object} domain.Reports
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 403 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/reports [get]
func (s *Server) apiReports(c echo.Context) error {
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	isSent := c.QueryParam("is_sent")
	reports, err := repo.FindReports(c.Request().Context(), layoutID, isSent)
	if err != nil {
		s.log.Errorf("repo.FindReports error %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, reports)
}

// ReportFilesResponse http wrapper with metadata
type ReportFilesResponse struct {
	Data domain.ReportFiles `json:"data"`
	Metadata
}

// apiReportFiles docs
// @Summary Gets early created reports files/items
// @Description gets early created reports files/items
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Description is_sent filter for reports, 1/true = report has been sent
// @Produce  json
// @Tags reports
// @Param layout_id query string false "digit/uuid format"
// @Param report_id path string true "digit format"
// @Param is_sent query string false "boolean format false/true/0/1"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Success 200 {object} infra.ReportFilesResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 403 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/reports/{report_id}/files [get]
func (s *Server) apiReportFiles(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	repo, layoutID, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	reportID := c.Param("report_id")
	isSent := c.QueryParam("is_sent")
	files, total, err := repo.FindReportFiles(c.Request().Context(), layoutID, reportID, isSent, offset, limit)
	if err != nil {
		s.log.Errorf("repo.FindReportFiles error %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	response := ReportFilesResponse{
		Data: files,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(files)),
				Offset: offset,
				Limit:  limit,
				Total:  total,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

// apiReportFileContent docs
// @Summary Gets report file content
// @Description gets report file content
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Tags reports
// @Param layout_id query string false "digit/uuid format"
// @Param report_id path string true "digit format"
// @Param file_id path string true "digit format"
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 403 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/reports/{report_id}/files/{file_id} [get]
func (s *Server) apiReportFileContent(c echo.Context) error {
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	fileID := c.Param("file_id")
	file, err := repo.FindReportFileByID(c.Request().Context(), fileID)
	if err != nil {
		s.log.Errorf("repo.FindReportFileByID error %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.Blob(http.StatusOK, file.MimeType, file.Content)
}
