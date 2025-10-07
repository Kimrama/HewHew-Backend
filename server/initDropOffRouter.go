package server

import (
	_dropOffController "hewhew-backend/pkg/dropOff/controller"
	_dropOffRepository "hewhew-backend/pkg/dropOff/repository"
	_dropOffService "hewhew-backend/pkg/dropOff/service"
)

func (s *fiberServer) initDropOffRouter() {
	dropOffRepository := _dropOffRepository.NewDropOffRepositoryImpl(s.db, s.conf.Supabase)
	dropOffService := _dropOffService.NewDropOffServiceImpl(dropOffRepository)
	dropOffController := _dropOffController.NewDropOffControllerImpl(dropOffService)

	dropOffGroup := s.app.Group("/v1/dropOff")
	dropOffGroup.Post("/", dropOffController.CreateDropOff)
	dropOffGroup.Get("/", dropOffController.GetAllDropOffs)
	dropOffGroup.Get("/:dropOffID", dropOffController.GetDropOffByID)
}
