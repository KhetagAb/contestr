package regatta

import (
	"contestr/internal/generated/server"
	"contestr/internal/integrations/ejudge"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ContestHandle struct {
	// TODO переделать на абстракцию
	ejudgeXmlParser *ejudge.ContestXMLFetcher
}

func NewContestHandle(
	ejudgeXmlParser *ejudge.ContestXMLFetcher,
) *ContestHandle {
	return &ContestHandle{
		ejudgeXmlParser: ejudgeXmlParser,
	}
}

func (s *ContestHandle) GetContest(ectx echo.Context, contestId int) error {
	ctx := ectx.Request().Context()

	runLog, err := s.ejudgeXmlParser.FetchAndParseXML(ctx, contestId)
	if err != nil {
		return ectx.JSON(http.StatusInternalServerError, err.Error())
	}

	response := server.RegattaContestStandings{
		ContestId:   &contestId,
		ContestName: &runLog.Name,
		Rows:        &[]server.RegattaContestRow{},
	}

	return ectx.JSON(http.StatusOK, response)
}
