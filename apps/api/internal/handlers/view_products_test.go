package handlers_test

import (
	sqlStore "github.com/hallgren/eventsourcing/eventstore/sql"
	"github.com/joe-reed/meal-planner/apps/api/internal/application"
	"github.com/joe-reed/meal-planner/apps/api/internal/database"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/category"
	"github.com/joe-reed/meal-planner/apps/api/internal/domain/product"
	"github.com/joe-reed/meal-planner/apps/api/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViewingProducts(t *testing.T) {
	p1 := product.NewProductBuilder().WithName("Chicken").WithCategory(category.Meat).WithId("8a378ac5-e0d5-405a-8cb5-f03cc1d92d8b").Build()
	p2 := product.NewProductBuilder().WithName("Ice Cream").WithCategory(category.Frozen).WithId("ad11289f-b2b0-4195-ba89-63d92ccc64d7").Build()
	p3 := product.NewProductBuilder().WithName("Carrot").WithCategory(category.Vegetables).WithId("57b5d842-eda3-4368-a496-21956de5e254").Build()

	repo := product.NewFakeProductRepository()
	err := repo.Add(p1)
	assert.NoError(t, err)
	err = repo.Add(p2)
	assert.NoError(t, err)
	err = repo.Add(p3)
	assert.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest("", "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.ProductHandler{Application: application.NewProductApplication(repo)}

	if assert.NoError(t, h.GetProducts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `[{"id":"57b5d842-eda3-4368-a496-21956de5e254","name":"Carrot","category":"Vegetables"},{"id":"8a378ac5-e0d5-405a-8cb5-f03cc1d92d8b","name":"Chicken","category":"Meat"},{"id":"ad11289f-b2b0-4195-ba89-63d92ccc64d7","name":"Ice Cream","category":"Frozen"}]`+"\n", rec.Body.String())
	}
}

func TestViewingProductsGrouped(t *testing.T) {
	db, err := database.CreateDatabase(":memory:")
	assert.NoError(t, err)

	es := sqlStore.Open(db)

	p1 := product.NewProductBuilder().WithName("Chicken").WithCategory(category.Meat).WithId("8a378ac5-e0d5-405a-8cb5-f03cc1d92d8b").Build()
	p2 := product.NewProductBuilder().WithName("Ice Cream").WithCategory(category.Frozen).WithId("ad11289f-b2b0-4195-ba89-63d92ccc64d7").Build()
	p3 := product.NewProductBuilder().WithName("Carrot").WithCategory(category.Vegetables).WithId("57b5d842-eda3-4368-a496-21956de5e254").Build()
	p4 := product.NewProductBuilder().WithName("Broccoli").WithCategory(category.Vegetables).WithId("c1d5d842-eda3-4368-a496-21956de5e254").Build()

	repo, err := product.NewSqliteProductRepository(db)

	assert.NoError(t, err)

	err = repo.Add(p1)
	assert.NoError(t, err)
	err = repo.Add(p2)
	assert.NoError(t, err)
	err = repo.Add(p3)
	assert.NoError(t, err)
	err = repo.Add(p4)

	e := echo.New()
	req := httptest.NewRequest("", "/?grouped=true", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &handlers.ProductHandler{Application: application.NewProductApplication(repo), EventStore: es}

	if assert.NoError(t, h.GetProducts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"Frozen":[{"id":"ad11289f-b2b0-4195-ba89-63d92ccc64d7","name":"Ice Cream","category":"Frozen"}],"Meat":[{"id":"8a378ac5-e0d5-405a-8cb5-f03cc1d92d8b","name":"Chicken","category":"Meat"}],"Vegetables":[{"id":"57b5d842-eda3-4368-a496-21956de5e254","name":"Carrot","category":"Vegetables"},{"id":"c1d5d842-eda3-4368-a496-21956de5e254","name":"Broccoli","category":"Vegetables"}]}`+"\n", rec.Body.String())
	}
}
