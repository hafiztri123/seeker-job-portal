// internal/handlers/search.go
package handlers

import (
   "github.com/gofiber/fiber/v2"
   "github.com/hafiztri123/internal/core/ports"
)

type SearchHandler struct {
   searchService ports.SearchService
}

func NewSearchHandler(searchService ports.SearchService) *SearchHandler {
   return &SearchHandler{searchService: searchService}
}

type searchRequest struct {
   Keyword       string   `json:"keyword" query:"keyword"`
   LocationType  string   `json:"location_type" query:"location_type"`
   Skills        []string `json:"skills" query:"skills"`
   SalaryMin     *int     `json:"salary_min" query:"salary_min"`
   SalaryMax     *int     `json:"salary_max" query:"salary_max"`
   Page          int      `json:"page" query:"page"`
   PageSize      int      `json:"page_size" query:"page_size"`
}

func (h *SearchHandler) SearchJobs(c *fiber.Ctx) error {
   var req searchRequest
   if err := c.QueryParser(&req); err != nil {
       return err
   }

   if req.Page <= 0 {
       req.Page = 1
   }
   if req.PageSize <= 0 {
       req.PageSize = 10
   }

   results, err := h.searchService.SearchJobs(c.Context(), ports.JobSearchQuery{
       Keyword:      req.Keyword,
       LocationType: req.LocationType,
       Skills:       req.Skills,
       SalaryMin:    req.SalaryMin,
       SalaryMax:    req.SalaryMax,
       Page:         req.Page,
       PageSize:     req.PageSize,
   })
   if err != nil {
       return err
   }

   return c.JSON(fiber.Map{
       "data": results,
       "pagination": fiber.Map{
           "page": req.Page,
           "page_size": req.PageSize,
       },
   })
}