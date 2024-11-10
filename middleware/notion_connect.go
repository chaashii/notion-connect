package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"notion-connect/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/jomei/notionapi"
)

type NotionConnectImpl struct {
	NotionConfig config.NotionAPIConfig
}

func NotionConnectInit(notionCfg config.NotionAPIConfig) *NotionConnectImpl {
	return &NotionConnectImpl{
		NotionConfig: config.NotionAPIConfig{
			Key:  notionCfg.Key,
			DbId: notionCfg.DbId,
		},
	}
}

type NotionResponse struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (n *NotionConnectImpl) ConnectNotion(c *gin.Context) {
	// Create a new Notion client
	client := notionapi.NewClient(notionapi.Token(n.NotionConfig.Key))

	// Query the database
	query := &notionapi.DatabaseQueryRequest{
		PageSize: 10, // Adjust this value as needed
	}

	result, err := client.Database.Query(context.Background(), notionapi.DatabaseID(n.NotionConfig.DbId), query)
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}

	c.JSON(http.StatusOK, result)
}

func CallAPIWithHeaders(url, apiKey string) (*NotionResponse, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add headers to the request
	// req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", "Bearer "+apiKey)
	// req.Header.Add("Custom-Header", "custom-value")

	log.Println("resp.Body-", req)
	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the JSON response
	var apiResp NotionResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return &apiResp, nil
}
